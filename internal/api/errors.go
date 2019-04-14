package api

import "github.com/powerman/go-service-goswagger-clean-example/internal/api/model"

type defaultResponder interface {
	SetStatusCode(code int)
	SetPayload(payload *model.Error)
}

func defError(err error, r defaultResponder) defaultResponder {
	r.SetStatusCode(500)
	r.SetPayload(&model.Error{
		Code:    500,
		Message: err.Error(),
	})
	return r
}

func apiError(err error) *model.Error {
	return &model.Error{
		Code:    500,
		Message: err.Error(),
	}
}
