package error

import (
	"fmt"
	"reflect"
)

type Code int

func (c Code) ToInt() int {
	return int(c)
}

type Error struct {
	Code    Code
	Message string
}

func New(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) HTTPError() string {
	if e.Code == Internal {
		return "internal server error"
	}

	return e.Message
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func (e *Error) Is(target error) bool {
	targetError, ok := target.(*Error)
	if !ok {
		return false
	}

	return reflect.DeepEqual(e, targetError)
}

func NewErrorFromDomain(err error) *Error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *Error:
		return err.(*Error)
	default:
		return New(Internal, err.Error())
	}
}

const (
	OK               Code = 0
	InvalidArgument  Code = 1
	NotFound         Code = 2
	AlreadyExists    Code = 3
	PermissionDenied Code = 4
	Unauthenticated  Code = 5
	DeadlineExceeded Code = 6
	Internal         Code = 7
)
