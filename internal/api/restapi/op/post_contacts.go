// Code generated by go-swagger; DO NOT EDIT.

package op

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/powerman/go-service-goswagger-clean-example/internal/app"
)

// PostContactsHandlerFunc turns a function with the right signature into a post contacts handler
type PostContactsHandlerFunc func(PostContactsParams, *app.Auth) PostContactsResponder

// Handle executing the request and returning a response
func (fn PostContactsHandlerFunc) Handle(params PostContactsParams, principal *app.Auth) PostContactsResponder {
	return fn(params, principal)
}

// PostContactsHandler interface for that can handle valid post contacts params
type PostContactsHandler interface {
	Handle(PostContactsParams, *app.Auth) PostContactsResponder
}

// NewPostContacts creates a new http.Handler for the post contacts operation
func NewPostContacts(ctx *middleware.Context, handler PostContactsHandler) *PostContacts {
	return &PostContacts{Context: ctx, Handler: handler}
}

/*PostContacts swagger:route POST /contacts postContacts

Add contact

*/
type PostContacts struct {
	Context *middleware.Context
	Handler PostContactsHandler
}

func (o *PostContacts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostContactsParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *app.Auth
	if uprinc != nil {
		principal = uprinc.(*app.Auth) // this is really a app.Auth, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
