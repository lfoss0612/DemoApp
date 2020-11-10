package middleware

import (
	"errors"
	"net/http"

	democtx "github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/response"
)

type PanicReport struct {
	PanicMsg interface{} `json:"panic-msg"`
	Stack    string      `json:"stack"`
}

// HandlePanic is a Middleware Handler that recovers the request from a panic
func getPanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if requestContext, ctxErr := democtx.GetContextFromRequest(r); ctxErr == nil {
					var theErr error
					if panicErr, ok := err.(error); ok {
						theErr = errors.New("unexpected system error: " + panicErr.Error())
					} else {
						theErr = errors.New("unexpected system error")
					}
					response.WriteError(w, &demoerrors.AppError{Message: theErr.Error(), Code: http.StatusInternalServerError}, requestContext)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
