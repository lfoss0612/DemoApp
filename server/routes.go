package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/lfoss0612/DemoApp/handlers"
)

type HandlerFunc = func(ctx *request.Context, w http.ResponseWriter, requestValue request.Value)

// Route Interface
type Route interface {
	GetName() string
	GetMethod() string
	GetPattern() string
	GetRequestFactory() request.Factory
	GetFunction() HandlerFunc
}

type router struct {
	mux.Router
}

func (r *router) addRoutes(routes []*Routes) {
	for _, route := range routes {
		router.Handle(route.GetPattern(), http.HandlerFunc(buildHandler(route.GetFunction(), route.GetRequestFactory()))).Methods(route.GetMethod())
	}
}

func (r *router) addMiddleware(middlewares []MiddlewareFunc) {
	for _, middleware := range middlewares {
		r.router.Use(middlewareFunc)
	}
}

func buildHandler(theHandler func(ctx *democtx.Context, w http.ResponseWriter, requestValue request.Value), requestFactory request.Factory) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestValue := requestFactory.NewInstance()

		ctx, err := democtx.GetContextFromRequest(r)

		if err != nil {
			response.WriteError(w, &demoerrors.AppError{Message: err.Error(), Code: http.StatusInternalServerError}, ctx)
			return
		}

		if readErr := readAndValidateRequest(r, requestValue, ctx); readErr != nil {
			response.WriteError(w, readErr, ctx)
			return
		}

		theHandler(ctx, w, requestValue)
	})
}
