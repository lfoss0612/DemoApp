package handlers

import (
	"net/http"

	democtx "github.com/lfoss0612/DemoApp/context"
	"github.com/lfoss0612/DemoApp/response"
)

func HealthCheck(ctx *request.Context, w http.ResponseWriter, requestValue request.Value) {

	// do some additional test of health here. For now, respond 200
	// TODO(ny): CORS middleware with more restrictive settings.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Methods", "Content-Type, api_key, Authorization")

	if r.Method == "HEAD" {
		w.WriteHeader(http.StatusOK)
	} else {
		response.WriteSuccess(w, ctx)
	}
}
