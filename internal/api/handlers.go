package api

import (
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

func (svc *service) getContacts(params op.GetContactsParams, auth *app.Auth) op.GetContactsResponder {
	ctx, log := fromRequest(params.HTTPRequest, auth)
	cs, err := svc.app.Contacts(ctx, log, *auth)
	if err != nil {
		return defError(err, op.NewGetContactsDefault(0)).(op.GetContactsResponder)
	}
	return op.NewGetContactsOK().WithPayload(apiContacts(cs))
}

func (svc *service) postContacts(params op.PostContactsParams, auth *app.Auth) op.PostContactsResponder {
	ctx, log := fromRequest(params.HTTPRequest, auth)
	c := appContact(params.Contact)
	log.Debug("calling AddContact")
	err := svc.app.AddContact(ctx, log, *auth, &c)
	if err != nil {
		return op.NewPostContactsDefault(500).WithPayload(apiError(err))
	}
	return op.NewPostContactsCreated().WithPayload(apiContact(c))
}
