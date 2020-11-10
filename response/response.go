package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	democtx "github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
)

// WriteJSON marshals anything to JSON and writes it as a response.
func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}, r *democtx.Context) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(w, appErr, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(w, appErr, r)
		return
	}
	r.SetStatusCode(statusCode)
	r.SetResponseBody(string(b))
}

// WriteSuccess writes env to json as response.
func WriteSuccess(w http.ResponseWriter, r *democtx.Context) {
	activeColor := os.Getenv("ACTIVE_COLOR")
	appVersion := os.Getenv("APP_VERSION")
	clusterName := os.Getenv("CLUSTER_NAME")
	buildTimestamp := os.Getenv("BUILD_TS")
	deployTimestamp := os.Getenv("DEPLOY_TS")
	commit := os.Getenv("COMMIT")
	refName := os.Getenv("REF_NAME")
	deployedBy := os.Getenv("DEPLOYED_BY")

	writeJson(w, http.StatusOK, Status{Status: "UP", Color: activeColor, Version: appVersion,
		ClusterName: clusterName, BuildDatetime: buildTimestamp, LastDeployment: deployTimestamp,
		Commit: commit, RefName: refName, DeployedBy: deployedBy}, r)
}

// WriteError forms an httpError and writes it as a response.
func WriteError(w http.ResponseWriter, err *demoerrors.AppError, r *democtx.Context) {
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

	writeJson(w, err.Code, e, r)
}

func writeJson(writer http.ResponseWriter, statusCode int, message interface{}, r *democtx.Context) {
	writer.Header().Set("Content-Type", "application/json")
	msg, err := json.Marshal(message)
	if err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(writer, appErr, r)
		return
	}

	writer.WriteHeader(statusCode)

	if _, err := writer.Write(msg); err != nil {
		appErr := &demoerrors.AppError{FieldErrors: nil, Message: fmt.Sprintf("ERROR: %s", err.Error()), Code: http.StatusInternalServerError}
		WriteError(writer, appErr, r)
		return
	}

	r.SetStatusCode(statusCode)
	r.SetResponseBody(string(msg))
}

type Status struct {
	Status         string `json:"status"`
	Color          string `json:"color"`
	Version        string `json:"version"`
	ClusterName    string `json:"cluster"`
	Commit         string `json:"commit"`
	RefName        string `json:"refName"`
	BuildDatetime  string `json:"buildDatetime"`
	LastDeployment string `json:"lastDeployment"`
	DeployedBy     string `json:"deployedBy"`
}
