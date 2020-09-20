package openapi

import (
	"errors"

	"github.com/powerman/go-service-example/api/openapi/restapi/op"
	"github.com/powerman/go-service-example/internal/app"
)

func (srv *server) HealthCheck(params op.HealthCheckParams) op.HealthCheckResponder {
	ctx, log := fromRequest(params.HTTPRequest, nil)
	status, err := srv.app.HealthCheck(ctx)
	switch {
	default:
		return errHealthCheck(log, err, codeInternal)
	case err == nil:
		return op.NewHealthCheckOK().WithPayload(status)
	}
}

func (srv *server) ListContacts(params op.ListContactsParams, auth *app.Auth) op.ListContactsResponder {
	ctx, log := fromRequest(params.HTTPRequest, auth)
	cs, err := srv.app.Contacts(ctx, *auth, appSeekPage(params.Args.SeekPagination))
	switch {
	default:
		return errListContacts(log, err, codeInternal)
	case err == nil:
		return op.NewListContactsOK().WithPayload(apiContacts(cs))
	}
}

func (srv *server) AddContact(params op.AddContactParams, auth *app.Auth) op.AddContactResponder {
	ctx, log := fromRequest(params.HTTPRequest, auth)
	log.Debug("calling AddContact")
	c, err := srv.app.AddContact(ctx, *auth, *params.Args.Name)
	switch {
	default:
		return errAddContact(log, err, codeInternal)
	case errors.Is(err, app.ErrContactExists):
		return errAddContact(log, err, codeContactExists)
	case err == nil:
		return op.NewAddContactCreated().WithPayload(apiContact(*c))
	}
}
