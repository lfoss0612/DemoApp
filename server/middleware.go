package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MiddlewareFunc alias for mux MiddlewareFunc
type MiddlewareFunc func(http.Handler) http.Handler

func (r *router) addMiddlewares(middlewares []MiddlewareFunc) {
	for _, mw := range middlewares {
		r.Router.Use(mux.MiddlewareFunc(mw))
	}
}
