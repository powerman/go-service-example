package openapi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	oapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/client"
	"github.com/powerman/go-service-goswagger-clean-example/api/openapi/model"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
	"github.com/powerman/go-service-goswagger-clean-example/internal/srv/openapi"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/def"
	"github.com/powerman/go-service-goswagger-clean-example/pkg/netx"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"
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

func testNewServer(t *check.C) (cleanup func(), c *client.AddressBook, url string, mockAppl *app.MockAppl) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockAppl = app.NewMockAppl(ctrl)

	server, err := openapi.NewServer(mockAppl, openapi.Config{
		APIKeyAdmin: "admin",
		Addr:        netx.NewAddr("localhost", 0),
	})
	t.Must(t.Nil(err, "NewServer"))
	t.Must(t.Nil(server.Listen(), "server.Listen"))
	errc := make(chan error, 1)
	go func() { errc <- server.Serve() }()

	cleanup = func() {
		t.Helper()
		t.Nil(server.Shutdown(), "server.Shutdown")
		t.Nil(<-errc, "server.Serve")
		ctrl.Finish()
	}

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

	return cleanup, c, url, mockAppl
}
