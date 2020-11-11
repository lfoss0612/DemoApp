package handlers

import (
	"net/http"
	"os"

	"github.com/lfoss0612/DemoApp/context"
	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/response"
)

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

func HealthCheck(ctx *context.Context, w http.ResponseWriter, requestValue request.Value) {

	// do some additional test of health here. For now, respond 200
	// TODO(ny): CORS middleware with more restrictive settings.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Methods", "Content-Type, api_key, Authorization")

	if ctx.GetMethod() == "HEAD" {
		ctx.SetStatusCode(http.StatusOK)
		w.WriteHeader(http.StatusOK)
	} else {
		ctx.SetStatusCode(http.StatusOK)
		status := getStatus()
		ctx.SetResponseBody(status)
		response.WriteJSON(w, http.StatusOK, status)
	}
}

// WriteSuccess writes env to json as response.
func getStatus() Status {
	return Status{
		Status:         "UP",
		Color:          os.Getenv("ACTIVE_COLOR"),
		Version:        os.Getenv("APP_VERSION"),
		ClusterName:    os.Getenv("CLUSTER_NAME"),
		BuildDatetime:  os.Getenv("BUILD_TS"),
		LastDeployment: os.Getenv("DEPLOY_TS"),
		Commit:         os.Getenv("COMMIT"),
		RefName:        os.Getenv("REF_NAME"),
		DeployedBy:     os.Getenv("DEPLOYED_BY"),
	}
}
