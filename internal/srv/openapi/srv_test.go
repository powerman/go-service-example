package openapi_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/structlog"

	"github.com/powerman/go-service-example/api/openapi/restapi"
	"github.com/powerman/go-service-example/api/openapi/restapi/op"
	"github.com/powerman/go-service-example/internal/apix"
	"github.com/powerman/go-service-example/internal/srv/openapi"
	"github.com/powerman/go-service-example/pkg/def"
)

type Ctx = context.Context

var healthCheckEndpoint = new(op.HealthCheckURL).String()

func fetch(t *check.C, url string, headers ...string) *http.Response {
	t.Helper()
	c := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	t.Cleanup(cancel)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	t.Must(t.Nil(err))
	for i := 0; i < len(headers); i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}
	resp, err := c.Do(req)
	t.Must(t.Nil(err))
	return resp
}

func TestServeNoCache(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, _ := testNewServer(t, openapi.Config{})

	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)

	resp := fetch(t, tsURL+healthCheckEndpoint)
	t.Equal(resp.StatusCode, 200)
	t.Equal(resp.Header.Get("Expires"), "0")
	t.Equal(resp.Header.Get("Cache-Control"), "no-cache, no-store, must-revalidate")
	t.Equal(resp.Header.Get("Pragma"), "no-cache")
}

func TestServeXFF(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, _ := testNewServer(t, openapi.Config{})

	var remoteIP string
	mockAppl.EXPECT().HealthCheck(gomock.Any()).DoAndReturn(func(ctx Ctx) (interface{}, error) {
		remoteIP = apix.FromContext(ctx)
		return nil, nil
	}).Times(2)

	resp := fetch(t, tsURL+healthCheckEndpoint)
	t.Equal(resp.StatusCode, 200)
	t.Equal(remoteIP, "127.0.0.1")
	resp = fetch(t, tsURL+healthCheckEndpoint, "X-Forwarded-For", "192.168.1.1, 1.2.3.4, 4.3.2.1")
	t.Equal(resp.StatusCode, 200)
	t.Equal(remoteIP, "1.2.3.4")
}

func TestServeLogger(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	for _, basePath := range []string{"", "/", "/base/path", "/base/path/"} {
		basePath := basePath
		t.Run("BasePath="+basePath, func(tt *testing.T) {
			t := check.T(tt)
			t.Parallel()
			_, tsURL, mockAppl, _ := testNewServer(t, openapi.Config{BasePath: basePath})

			var log *structlog.Logger
			mockAppl.EXPECT().HealthCheck(gomock.Any()).DoAndReturn(func(ctx Ctx) (interface{}, error) {
				log = structlog.FromContext(ctx, nil)
				return nil, nil
			})

			t.Log(tsURL + path.Join(basePath, healthCheckEndpoint))
			resp := fetch(t, tsURL+path.Join(basePath, healthCheckEndpoint))
			t.Equal(resp.StatusCode, 200)
			var buf bytes.Buffer
			log.SetOutput(&buf)
			log.Info("test")
			t.Match(buf.String(), ": 127.0.0.1 +GET +"+healthCheckEndpoint+": `test`")
		})
	}
}

func TestServeRecover(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, _ := testNewServer(t, openapi.Config{})

	mockAppl.EXPECT().HealthCheck(gomock.Any()).DoAndReturn(func(ctx Ctx) (interface{}, error) {
		panic("boom")
	})

	resp := fetch(t, tsURL+healthCheckEndpoint)
	t.Equal(resp.StatusCode, 500)
}

func TestServeRecoverNetError(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, logc := testNewServer(t, openapi.Config{})

	c := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	cancelled := make(chan struct{})

	mockAppl.EXPECT().HealthCheck(gomock.Any()).DoAndReturn(func(ctx Ctx) (interface{}, error) {
		cancel()
		<-cancelled
		return make([]byte, 10000), nil // Overflow 8KB output buffer when sending response.
	})

	req, err := http.NewRequestWithContext(ctx, "GET", tsURL+healthCheckEndpoint, nil)
	t.Must(t.Nil(err))
	resp, err := c.Do(req)
	t.Must(t.Err(err, context.Canceled))
	t.Nil(resp)
	cancelled <- struct{}{}

	select {
	case <-time.After(def.TestTimeout):
		t.Error("nothing logged after request")
	case msg := <-logc:
		t.Match(msg, "`recovered`.*broken pipe")
	}
}

func TestServeAccessLog(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, logc := testNewServer(t, openapi.Config{})
	reAccessLog := regexp.MustCompile("`handled`|`failed to handle`|`panic`")

	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)
	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, io.EOF)
	mockAppl.EXPECT().HealthCheck(gomock.Any()).DoAndReturn(func(ctx Ctx) (interface{}, error) {
		panic("boom")
	})

	tests := []struct {
		want    int
		wantMsg string
	}{
		{200, "200 GET +" + healthCheckEndpoint + ": `handled`"},
		{500, "500 GET +" + healthCheckEndpoint + ": `failed to handle`"},
		{500, "500 GET +" + healthCheckEndpoint + ": `panic`"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			resp := fetch(t, tsURL+healthCheckEndpoint)
			t.Equal(resp.StatusCode, tc.want)
			var msg string
			for !reAccessLog.MatchString(msg) {
				select {
				case <-time.After(def.TestTimeout):
					t.Error("timeout waiting for access log message")
				case msg = <-logc:
				}
			}
			t.Match(msg, tc.wantMsg)
		})
	}
}

func TestServeSwagger(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	for _, basePath := range []string{"", "/", "/base/path"} {
		basePath := basePath
		t.Run("BasePath="+basePath, func(tt *testing.T) {
			t := check.T(tt)
			t.Parallel()
			_, tsURL, _, _ := testNewServer(t, openapi.Config{BasePath: basePath})

			swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
			t.Nil(err)
			if basePath == "" {
				basePath = swaggerSpec.BasePath()
			}

			testCases := []struct {
				path string
				want int
			}{
				{"/", 404},
				{"/swagger.yml", 404},
				{"/swagger.yaml", 404},
				{"/swagger.json", 200},
				{basePath, 404},
				{path.Join(basePath, "docs"), 200},
				{path.Join(basePath, "swagger.json"), 200},
			}
			for _, tc := range testCases {
				resp := fetch(t, tsURL+tc.path)
				t.Equal(resp.StatusCode, tc.want, tc.path)
			}
		})
	}
}

func TestServeCORS(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, mockAppl, _ := testNewServer(t, openapi.Config{})

	mockAppl.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)

	resp := fetch(t, tsURL+healthCheckEndpoint, "Origin", "google.com")
	t.Equal(resp.StatusCode, 200)
	t.Equal(resp.Header.Get("Access-Control-Allow-Origin"), "*")
}
