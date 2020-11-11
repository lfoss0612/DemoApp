package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/response"
)

// HandlerFunc function to handle route
type HandlerFunc = func(ctx context.Context, w http.ResponseWriter, requestValue request.Value)

// Route Interface
type Route interface {
	GetName() string
	GetMethod() string
	GetPattern() string
	GetRequestFactory() request.Factory
	GetFunction() HandlerFunc
}

type router struct {
	*mux.Router
}

func newRouter() *router {
	return &router{
		Router: mux.NewRouter(),
	}
}

func (r *router) addRoutes(routes []Route) {
	for _, route := range routes {
		r.Router.Handle(route.GetPattern(), http.HandlerFunc(buildHandler(route.GetFunction(), route.GetRequestFactory()))).Methods(route.GetMethod())
	}
}

func buildHandler(theHandler HandlerFunc, requestFactory request.Factory) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestFactory == nil {
			theHandler(r.Context(), w, nil)
		}

		requestValue := requestFactory.NewInstance()

		if readErr := request.ReadAndValidateRequest(r, requestValue); readErr != nil {
			response.WriteError(w, readErr)
			return
		}

		theHandler(r.Context(), w, requestValue)
	})
}

// GetRoutePattern Returns the Pattern of the route used to match
func GetRoutePattern(r *http.Request) (string, error) {
	return mux.CurrentRoute(r).GetPathTemplate()
}
