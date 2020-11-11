package middleware

import (
	"fmt"
	"net/http"

	democtx "github.com/lfoss0612/DemoApp/context"

	"github.com/go-http-utils/timeout"
	"github.com/lfoss0612/DemoApp/env"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/response"
)

func getTimeoutHandler(next http.Handler) http.Handler {
	return timeout.Handler(next, env.EnvVar.ServerReadTimeoutInSeconds, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestContext, err := democtx.GetContextFromRequest(r); err == nil {
			timeoutErrorMsg := fmt.Sprintf("request timeout: %s %s", r.Method, r.RequestURI)

			response.WriteError(w, &demoerrors.AppError{Message: timeoutErrorMsg, Code: http.StatusRequestTimeout})
		} else {
			w.WriteHeader(http.StatusRequestTimeout)

			if _, err := w.Write([]byte("Service timeout")); err != nil {
				requestContext.Logger().Error(err)
			}
		}
	}))
}
