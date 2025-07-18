package ownerr

import (
	"github.com/pkg/errors"
)

var (
	ErrUnknown            = errors.New("")
	ErrBadRequest         = errors.New("Controller: Bad Request")
	ErrNotFound           = errors.New("Controller: Err Not Found")
	ErrAlreadyExist       = errors.New("Controller: Err Already Exist")
	ErrUnauthorized       = errors.New("Controller: Unauthorized")
	ErrForbidden          = errors.New("Controller: Forbidden")
	ErrInternal           = errors.New("Controller: Internal Error")
	ErrMovedPermanently   = errors.New("Controller: Moved Permanently")
	ErrServiceUnavailable = errors.New("Controller: Service Unavailable")
)

type Error struct {
	ControllerErr error
	ServiceErr    error
}

func NewError(controllerErr error, serviceErr error) *Error {
	return &Error{
		ControllerErr: controllerErr,
		ServiceErr:    serviceErr,
	}
}

func (e Error) Error() string {
	if e.ControllerErr == nil {
		return ""
	}
	return e.ControllerErr.Error()
}

func (e Error) Unwrap() []error {
	return []error{e.ControllerErr, e.ServiceErr}
}
