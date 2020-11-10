package middleware

import (
	"net/http"

	"github.com/lfoss0612/DemoApp/logger"

	"time"

	democtx "github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/response"
)

func getLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()
		next.ServeHTTP(w, r)

		timeTaken := int64(time.Now().Sub(startTime) / time.Millisecond)
		ctx, err := democtx.GetContextFromRequest(r)
		if err != nil {
			response.WriteError(w, &demoerrors.AppError{Message: err.Error(), Code: http.StatusInternalServerError}, ctx)
			return
		}
		if "/api/v1/health" != ctx.GetPattern() {
			logMap := make(map[logger.LogField]interface{})
			logMap[logger.Method] = ctx.GetMethod()
			logMap[logger.Pattern] = ctx.GetPattern()
			logMap[logger.Uri] = ctx.GetURL().RequestURI()
			logMap[logger.Body] = ctx.GetBodyAsString()
			logMap[logger.StatusCode] = ctx.GetStatusCode()
			logMap[logger.TimeTaken] = timeTaken

			status := demoerrors.PROCESSING_ERROR
			if ctx.GetStatusCode()/100 == 2 {
				status = demoerrors.PROCESSING_SUCCESS
			} else {
				logMap[logger.ResponseBody] = ctx.GetResponseBody()
			}

			ctx.Logger().WithLogFields(logMap).Info(status)
		}
	})
}
