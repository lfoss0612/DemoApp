package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	demoerrors "github.com/lfoss0612/DemoApp/errors"
)

// WriteJSON marshals anything to JSON and writes it as a response.
func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(w, appErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(w, appErr)
		return
	}
}

// WriteError forms an httpError and writes it as a response.
func WriteError(w http.ResponseWriter, err *demoerrors.AppError) {
	type httpError struct {
		FieldErrors map[string]interface{} `json:"fieldErrors" description:"field errors if any"`
		Error       string                 `json:"error" description:"error message"`
		Code        int                    `json:"code" description:"HTTP status code, same as on response"`
	}

	if err.Code == 0 {
		err.Code = http.StatusInternalServerError
	}

	e := httpError{
		FieldErrors: err.FieldErrors,
		Error:       err.Message,
		Code:        err.Code,
	}

	WriteJSON(w, err.Code, e)
}
