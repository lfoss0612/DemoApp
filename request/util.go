package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	demoerrors "github.com/lfoss0612/DemoApp/errors"
)

// UnmarshalJSONRequest Read request body bytes into targetStruct
func UnmarshalJSONRequest(r *http.Request, targetStruct interface{}) *demoerrors.AppError {
	// extract body
	if r.Body == nil {
		return &demoerrors.AppError{Message: demoerrors.REQUEST_BODY_MISSING, Code: http.StatusBadRequest}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &demoerrors.AppError{Message: demoerrors.ERROR_READING_REQUEST_BODY, Code: http.StatusBadRequest}
	}

	// decode json
	err = json.Unmarshal(body, targetStruct)
	if err != nil {
		return demoerrors.RequestBodyDecodeError(err)
	}

	return nil
}
