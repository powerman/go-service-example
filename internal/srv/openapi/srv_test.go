package openapi_test

import (
	"context"
	"net/http"
	"path"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/powerman/check"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/restapi"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/def"
)

func TestServeSwagger(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	_, tsURL, _ := testNewServer(t)
	c := &http.Client{}
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	t.Nil(err)
	basePath := swaggerSpec.BasePath()

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
		ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
		req, err := http.NewRequestWithContext(ctx, "GET", tsURL+tc.path, nil)
		t.Nil(err)
		resp, err := c.Do(req)
		t.Nil(err, tc.path)
		t.Equal(resp.StatusCode, tc.want, tc.path)
		cancel()
	}
}
