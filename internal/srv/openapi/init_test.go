package openapi_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	oapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/powerman/go-service-example/api/openapi/client"
	"github.com/powerman/go-service-example/api/openapi/model"
	"github.com/powerman/go-service-example/internal/app"
	"github.com/powerman/go-service-example/internal/srv/openapi"
	"github.com/powerman/go-service-example/pkg/def"
	"github.com/powerman/go-service-example/pkg/netx"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	openapi.InitMetrics(reg, "test")
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	apiError401  = openapi.APIError(401, "unauthenticated for invalid credentials")
	apiError403  = openapi.APIError(403, "only admin can make changes")
	apiError500  = openapi.APIError(500, "internal error")
	apiError1000 = openapi.APIError(1000, app.ErrContactExists.Error())
	apiKeyAdmin  = oapiclient.APIKeyAuth("API-Key", "header", "admin")
	apiKeyUser   = oapiclient.APIKeyAuth("API-Key", "header", "u1")
	authAdmin    = app.Auth{UserID: "admin"}
	authUser     = app.Auth{UserID: "user:u1"}
	appContact1  = app.Contact{ID: 1, Name: "A"}
	appContact2  = app.Contact{ID: 2, Name: "B"}
	apiContact1  = &model.Contact{ID: 1, Name: swag.String("A")}
	apiContact2  = &model.Contact{ID: 2, Name: swag.String("B")}
)

func testNewServer(t *check.C, cfg openapi.Config) (c *client.AddressBook, url string, mockAppl *app.MockAppl, logc <-chan string) {
	cfg.Addr = netx.NewAddr("localhost", 0)
	cfg.APIKeyAdmin = "admin"

	t.Helper()
	ctrl := gomock.NewController(t)

	mockAppl = app.NewMockAppl(ctrl)

	server, err := openapi.NewServer(mockAppl, cfg)
	t.Must(t.Nil(err, "NewServer"))

	piper, pipew := io.Pipe()
	server.SetHandler(interceptLog(pipew, server.GetHandler()))
	logch := make(chan string, 64) // Keep some unread log messages.
	go func() {
		scanner := bufio.NewScanner(piper)
		for scanner.Scan() {
			select {
			default: // Do not hang test because of some unread log messages.
			case logch <- scanner.Text():
			}
		}
		close(logch)
	}()

	t.Must(t.Nil(server.Listen(), "server.Listen"))
	errc := make(chan error, 1)
	go func() { errc <- server.Serve() }()

	t.Cleanup(func() {
		t.Helper()
		t.Nil(server.Shutdown(), "server.Shutdown")
		t.Nil(<-errc, "server.Serve")
		pipew.Close()
	})

	ln, err := server.HTTPListener()
	t.Must(t.Nil(err, "server.HTTPListener"))
	c = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes:  []string{"http"},
		Host:     ln.Addr().String(),
		BasePath: client.DefaultBasePath,
	})
	url = fmt.Sprintf("http://%s", ln.Addr().String())

	// Avoid race between server.Serve() and server.Shutdown().
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	t.Must(t.Nil(err))
	_, err = (&http.Client{}).Do(req)
	t.Must(t.Nil(err, "connect to service"))
	<-logch

	return c, url, mockAppl, logch
}

func interceptLog(out io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := structlog.FromContext(r.Context(), nil)
		log.SetOutput(out)
		r = r.WithContext(structlog.NewContext(r.Context(), log))
		next.ServeHTTP(w, r)
	})
}
