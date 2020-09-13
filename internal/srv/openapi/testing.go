package openapi

import (
	"github.com/go-openapi/swag"
	"github.com/powerman/go-service-example/api/openapi/client/op"
	"github.com/powerman/go-service-example/api/openapi/model"
)

// APIError returns model.Error with given code and msg.
func APIError(code int32, msg string) *model.Error {
	return &model.Error{
		Code:    swag.Int32(code),
		Message: swag.String(msg),
	}
}

// ErrPayload returns err.Payload or err for unknown errors.
func ErrPayload(err error) interface{} {
	switch errDefault := err.(type) {
	default:
		return err
	case *op.ListContactsDefault:
		return errDefault.Payload
	case *op.AddContactDefault:
		return errDefault.Payload
	}
}
