package apperr

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

type Code string

const (
	CodeInvalidArgument  Code = "INVALID_ARGUMENT"
	CodeUnauthenticated  Code = "UNAUTHENTICATED"
	CodePermissionDenied Code = "PERMISSION_DENIED"
	CodeNotFound         Code = "NOT_FOUND"
	CodeConflict         Code = "CONFLICT"
	CodeInternal         Code = "INTERNAL"
)

type Error struct {
	code    Code
	message string
	cause   error
}

func New(code Code, message string) error {
	return &Error{
		code:    code,
		message: message,
	}
}

func Wrap(code Code, message string, err error) error {
	return &Error{
		code:    code,
		message: message,
		cause:   err,
	}
}

func (e *Error) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return string(e.code)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func IsCode(err error, code Code) bool {
	return GetCode(err) == code
}

func GetCode(err error) Code {
	if err == nil {
		return CodeInternal
	}
	var appErr *Error
	if errors.As(err, &appErr) && appErr.code != "" {
		return appErr.code
	}
	return CodeInternal
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	var appErr *Error
	if errors.As(err, &appErr) && appErr.message != "" {
		return appErr.message
	}
	return err.Error()
}

func HTTPStatus(err error) int {
	switch GetCode(err) {
	case CodeInvalidArgument:
		return http.StatusBadRequest
	case CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodePermissionDenied:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func GRPCStatusCode(err error) codes.Code {
	switch GetCode(err) {
	case CodeInvalidArgument:
		return codes.InvalidArgument
	case CodeUnauthenticated:
		return codes.Unauthenticated
	case CodePermissionDenied:
		return codes.PermissionDenied
	case CodeNotFound:
		return codes.NotFound
	case CodeConflict:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
