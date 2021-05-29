package openapi

import (
	"errors"
	"net/http"

	oapierrors "github.com/go-openapi/errors"

	"github.com/powerman/go-service-example/internal/app"
)

var errRequireAdmin = errors.New("only admin can make changes")

func (srv *server) authenticate(apiKey string) (*app.Auth, error) {
	switch apiKey {
	case "anonymous":
		return nil, oapierrors.Unauthenticated("invalid credentials")
	case srv.cfg.APIKeyAdmin:
		return &app.Auth{UserID: "admin"}, nil
	default:
		return &app.Auth{UserID: "user:" + apiKey}, nil
	}
}

func (srv *server) authorize(r *http.Request, principal interface{}) error {
	auth := principal.(*app.Auth) //nolint:forcetypeassert // Want panic.
	if r.Method != "GET" && auth.UserID != "admin" {
		return errRequireAdmin
	}
	return nil
}
