package request

import (
	"net/http"

	"github.com/gorilla/schema"

	democtx "github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
)

func readAndValidateRequest(r *http.Request, requestValue request.Value, ctx *democtx.Context) *demoerrors.AppError {
	if requestValue != nil {
		decoder := schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)

		if readPathErr := readPathParams(decoder, r, requestValue); readPathErr != nil {
			return &demoerrors.AppError{Message: readPathErr.Error(), Code: http.StatusBadRequest}
		}

		if readQueryErr := readQueryParams(decoder, r, requestValue); readQueryErr != nil {
			return &demoerrors.AppError{Message: readQueryErr.Error(), Code: http.StatusBadRequest}
		}

		if readBodyErr := readBody(r, requestValue, ctx); readBodyErr != nil {
			return readBodyErr
		}

		validationErrors := requestValue.Validate()
		if validationErrors != nil {
			appErr := &demoerrors.AppError{Message: validationErrors.Error(), Code: http.StatusBadRequest}
			if valErrs, ok := validationErrors.(*demoerrors.AppError); ok {
				appErr.FieldErrors = valErrs.FieldErrors
			}
			return appErr
		}

	}
	return nil
}

func readPathParams(decoder *schema.Decoder, r *http.Request, requestValue request.Value) error {
	pathParams := mux.Vars(r)
	if len(pathParams) > 0 {
		values := url.Values{}
		for key, val := range pathParams {
			values.Set(key, val)
		}

		decodeErr := decoder.Decode(requestValue, values)
		if decodeErr != nil {
			return &demoerrors.AppError{Message: decodeErr.Error(), Code: http.StatusBadRequest}
		}
	}

	return nil
}

func readQueryParams(decoder *schema.Decoder, r *http.Request, requestValue request.Value) error {
	if (r.Method == http.MethodGet || r.Method == http.MethodDelete) && len(r.URL.Query()) > 0 {
		decodeErr := decoder.Decode(requestValue, r.URL.Query())
		if decodeErr != nil {
			return &demoerrors.AppError{Message: decodeErr.Error(), Code: http.StatusBadRequest}
		}
	}
	return nil
}

func readBody(r *http.Request, requestValue request.Value, ctx *democtx.Context) *demoerrors.AppError {
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if err := request.UnmarshalJSONRequest(r, requestValue); err != nil {
			return err
		}
	}
	return nil
}
