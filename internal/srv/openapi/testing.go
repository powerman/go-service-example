package openapi

import (
	"errors"

	"github.com/go-openapi/swag"

	"github.com/powerman/go-service-example/api/openapi/model"
)

// APIError returns model.Error with given code and msg.
func APIError(code int32, msg string) *model.Error {
	return &model.Error{
		Code:    swag.Int32(code),
		Message: swag.String(msg),
	}
}

// ErrPayload returns err.Payload or *model.Error(nil) or err for unknown errors.
func ErrPayload(err error) interface{} {
	var errDefault interface{ GetPayload() *model.Error }
	switch ok := errors.As(err, &errDefault); true {
	case ok:
		return errDefault.GetPayload()
	case err == nil:
		return (*model.Error)(nil)
	default:
		return err
	}
}
