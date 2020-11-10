package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type FieldErrors map[string]interface{}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type AppError struct {
	FieldErrors FieldErrors `json:"fieldErrors,omitempty"`
	Message     string      `json:"error,omitempty"`
	Code        int         `json:"code,omitempty"`
}

func (appErr *AppError) Error() string {
	return appErr.Message
}

type MultipleError interface {
	WrappedErrors() []error
	Length() int
	Add(error)
}

type MError = multierror.Error
type MultiErrors struct {
	*MError
	httpStatus int
}

func NewMultiErrors() *MultiErrors {
	return &MultiErrors{
		MError: &multierror.Error{},
	}
}

func (e *MultiErrors) Cause() error {
	return e.MError.ErrorOrNil()
}

func (e *MultiErrors) StackTrace() errors.StackTrace {
	err := e.ErrorOrNil()

	if err, ok := err.(StackTracer); ok {
		return err.StackTrace()
	}
	return nil
}

func (e *MultiErrors) SetStatusCode(httpStatus int) {
	e.httpStatus = httpStatus
}

func (e *MultiErrors) GetStatusCode() int {
	return e.httpStatus
}

func (e *MultiErrors) Add(err error) {
	if err != nil {
		e.MError = multierror.Append(e.MError, err)
	}
}

func RequestReadError() *AppError {
	return &AppError{Message: "unable to read request", Code: http.StatusInternalServerError}
}

func ServerError(msg string, fieldErrors map[string]interface{}) *AppError {
	return &AppError{Message: msg, FieldErrors: fieldErrors, Code: http.StatusInternalServerError}
}

func ReqContextError() *AppError {
	return ServerError("failed to get request context", nil)
}

func InvalidRequestError(fieldErrors map[string]interface{}) *AppError {
	return &AppError{Message: "request is invalid", FieldErrors: fieldErrors, Code: http.StatusBadRequest}
}

func RequestBodyDecodeError(err error) *AppError {
	var fieldErrs map[string]interface{}
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		// convert type error to JSON type for readability
		jsonType := reflectKindToJSONType(typeErr.Type.Kind())
		fieldErrs = map[string]interface{}{
			typeErr.Field: fmt.Sprintf("expected %s", jsonType)}
	}

	return &AppError{
		Message:     "Error decoding request body",
		Code:        http.StatusBadRequest,
		FieldErrors: fieldErrs}
}

func reflectKindToJSONType(t reflect.Kind) string {
	switch t {
	case reflect.Bool:
		fallthrough
	case reflect.String:
		return t.String()

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		return "int"

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return "float"

	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return "array"

	case reflect.Map:
		fallthrough
	case reflect.Struct:
		return "object"
	}

	return "invalid"
}
