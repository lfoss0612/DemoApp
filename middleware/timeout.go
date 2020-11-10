package middleware

import (
	"fmt"
	"net/http"
	"time"

	democtx "github.com/lfoss0612/DemoApp/context"

	"github.com/go-http-utils/timeout"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/response"
	"github.com/lfoss0612/DemoApp/env"
)

const timeoutLen = time.Duration(env.EnvVar.ServerReadTimeoutInSeconds) * time.Second

func getTimeoutHandler(next http.Handler) http.Handler {
		return timeout.Handler(next, timeoutLen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requestContext, err := democtx.GetContextFromRequest(r); err == nil {
				timeoutErrorMsg := fmt.Sprintf("request timeout: %s %s", r.Method, r.RequestURI)

				response.WriteError(w, &demoerrors.AppError{Message: timeoutErrorMsg, Code: http.StatusRequestTimeout}, requestContext)
			} else {
				w.WriteHeader(http.StatusRequestTimeout)

				if _, err := w.Write([]byte("Service timeout")); err != nil {
					requestContext.Logger().Error(err)
				}
			}
		})
	}
}
