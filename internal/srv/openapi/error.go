//go:generate genny -in=$GOFILE -out=gen.$GOFILE gen "HealthCheck=ListContacts,AddContact"

package openapi

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/powerman/go-service-example/api/openapi/model"
	"github.com/powerman/go-service-example/api/openapi/restapi/op"
	"github.com/powerman/go-service-example/pkg/def"
)

func errHealthCheck(log Log, err error, code errCode) op.HealthCheckResponder {
	if code.status < http.StatusInternalServerError {
		log.Info("client error", def.LogHTTPStatus, code.status, "code", code.extra, "err", err)
	} else {
		log.PrintErr("server error", def.LogHTTPStatus, code.status, "code", code.extra, "err", err)
	}

	msg := err.Error()
	if code.status == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = "internal error" //nolint:goconst // Duplicated by go:generate.
	}

	return op.NewHealthCheckDefault(code.status).WithPayload(&model.Error{
		Code:    swag.Int32(code.extra),
		Message: swag.String(msg),
	})
}
