package openapi

import (
	"errors"

	"github.com/powerman/go-service-example/api/openapi/restapi/op"
	"github.com/powerman/go-service-example/internal/app"
)

func (srv *server) listContacts(params op.ListContactsParams, auth *app.Auth) op.ListContactsResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, auth)
	cs, err := srv.app.Contacts(ctx, *auth)
	switch {
	default:
		return errListContacts(log, err, codeInternal)
	case err == nil:
		return op.NewListContactsOK().WithPayload(apiContacts(cs))
	}
}

func (srv *server) addContact(params op.AddContactParams, auth *app.Auth) op.AddContactResponder {
	ctx, log, _ := fromRequest(params.HTTPRequest, auth)
	log.Debug("calling AddContact")
	c, err := srv.app.AddContact(ctx, *auth, *params.Contact.Name)
	switch {
	default:
		return errAddContact(log, err, codeInternal)
	case errors.Is(err, app.ErrContactExists):
		return errAddContact(log, err, codeContactExists)
	case err == nil:
		return op.NewAddContactCreated().WithPayload(apiContact(*c))
	}
}
