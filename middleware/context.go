package middleware

import (
	"net/http"

	democtx "github.com/lfoss0612/DemoApp/context"
	"github.com/lfoss0612/DemoApp/logger"
)

func getContextHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := democtx.BuildContextFromRequest(r)

		if txID := r.Header.Get(democtx.TransactionIDHeader); txID != "" {
			ctx.AddLogField(logger.TransactionId, txID)
		}

		ctx.AddRequestDataToLog()
		r = ctx.AddToRequest(r)

		next.ServeHTTP(w, r)
	})
}
