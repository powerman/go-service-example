package api

import (
	"github.com/powerman/go-service-goswagger-clean-example/internal/api/restapi/op"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

func (svc *service) getContacts(params op.GetContactsParams, auth *app.Auth) op.GetContactsResponder {
	ctx := params.HTTPRequest.Context()
	cs, err := svc.app.Contacts(ctx, svc.log, *auth)
	if err != nil {
		// TODO Maybe --strict should leave default response as middleware.Responder?
		return defError(err, op.NewGetContactsDefault(0)).(op.GetContactsResponder)
	}
	return op.NewGetContactsOK().WithPayload(apiContacts(cs))
}

func (svc *service) postContacts(params op.PostContactsParams, auth *app.Auth) op.PostContactsResponder {
	ctx := params.HTTPRequest.Context()
	c := appContact(params.Contact)
	err := svc.app.AddContact(ctx, svc.log, *auth, &c)
	if err != nil {
		return op.NewPostContactsDefault(500).WithPayload(apiError(err))
	}
	return op.NewPostContactsCreated().WithPayload(apiContact(c))
}
