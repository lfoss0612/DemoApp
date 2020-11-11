package handlers

import (
	"fmt"
	"net/http"

	"github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/response"
)

type TimeRequest struct {
	Duration int `json:"duration"`
}

func (t *TimeRequest) Validate() error {
	return nil
}

func (t *TimeRequest) NewInstance() request.Value {
	return &TimeRequest{}
}

func Timer(ctx *context.Context, w http.ResponseWriter, requestValue request.Value) {
	req, ok := requestValue.(*TimeRequest)
	if !ok {
		response.WriteError(w, &demoerrors.AppError{Message: "unable to read request", Code: http.StatusInternalServerError})
	}
	ctx.SetStatusCode(http.StatusOK)

	fmt.Println("Duration:", req.Duration)
	response.WriteJSON(w, http.StatusOK, req.Duration)

}
