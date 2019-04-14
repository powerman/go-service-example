package api

import (
	"errors"
	"net/http"

	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

func (svc *service) authenticate(apiKey string) (*app.Auth, error) {
	if apiKey == "anonymous" {
		return nil, errors.New("account blocked, go away")
	}
	return &app.Auth{UserID: apiKey}, nil
}

func (svc *service) authorize(r *http.Request, principal interface{}) error {
	auth := principal.(*app.Auth)
	if r.Method != "GET" && auth.UserID != "admin" {
		return errors.New("only admin may make changes")
	}
	return nil
}
