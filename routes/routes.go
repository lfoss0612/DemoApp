package routes

import (
	"context"
	"net/http"

	democtx "github.com/lfoss0612/DemoApp/context"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/handlers"
	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/response"
	"github.com/lfoss0612/DemoApp/server"
)

type HandleFuncWithContext = func(ctx *democtx.Context, w http.ResponseWriter, requestValue request.Value)

// Route holds route metadata
type Route struct {
	Name        string
	method      string
	pattern     string
	reqFactory  request.Factory
	handlerFunc server.HandlerFunc
}

func (r *Route) GetName() string {
	return r.Name
}

func (r *Route) GetMethod() string {
	return r.method
}

func (r *Route) GetPattern() string {
	return r.pattern
}
func (r *Route) GetRequestFactory() request.Factory {
	return r.reqFactory
}

func (r *Route) GetFunction() server.HandlerFunc {
	return r.handlerFunc
}

func NewRoute(name, pattern string) *Route {
	return &Route{
		Name:    name,
		pattern: pattern,
		method:  "GET", //default to GET
	}
}

func (r *Route) Method(method string) *Route {
	r.method = method
	return r
}

func (r *Route) HandlerFunc(f server.HandlerFunc) *Route {
	r.handlerFunc = f
	return r
}

func (r *Route) ReqFactory(rf request.Factory) *Route {
	r.reqFactory = rf
	return r
}

func handleWithContext(fn HandleFuncWithContext) server.HandlerFunc {
	return server.HandlerFunc(func(requestContext context.Context, w http.ResponseWriter, requestValue request.Value) {
		ctx, err := democtx.GetContext(requestContext)
		if err != nil {
			response.WriteError(w, &demoerrors.AppError{Message: err.Error(), Code: http.StatusInternalServerError})
			return
		}
		fn(ctx, w, requestValue)
	})
}

// Routes getter
func Routes() []server.Route {
	return []server.Route{
		NewRoute("HealthCheck", "/api/v1/health").HandlerFunc(handleWithContext(handlers.HealthCheck)),
		NewRoute("HealthCheck", "/api/v1/health").HandlerFunc(handleWithContext(handlers.HealthCheck)).Method("HEAD"),
		NewRoute("timer", "/api/v1/timer/{duration}").HandlerFunc(handleWithContext(handlers.Timer)).ReqFactory(&handlers.TimeRequest{}),
	}
}
