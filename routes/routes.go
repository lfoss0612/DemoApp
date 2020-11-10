package routes

import (
	"net/http"

	"github.com/lfoss0612/DemoApp/handlers"
	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/server"
)

// Route holds route metadata
type Route struct {
	Name       string
	ReqMethod  string
	Pattern    string
	ReqFactory request.Factory
	ReqHandler http.HandlerFunc
}

func (r *Route) GetName() string {
	return r.Name
}

func (r *Route) GetMethod() string {
	return r.ReqMethod
}

func (r *Route) GetPattern() string {
	return r.Pattern
}
func (r *Route) GetRequestFactory() request.Factory {
	return r.ReqFactory
}

func (r *Route) GetFunction() HandlerFunc {
	return r.ReqHandler
}

func NewRoute(name, pattern string) *Route {
	return &Route{
		Name:      name,
		Pattern:   pattern,
		ReqMethod: "GET",
	}
}

func (r *Route) Method(method string) *Route {
	r.ReqMethod = method
	return r
}

func (r *Route) Handler(handler http.Handler) *Route {
	r.ReqHandler = handler
	return r
}

func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	r.ReqHandler = http.HandlerFunc(f)
	return r
}

// Routes getter
func Routes() []*server.Route {
	return []*Route{
		NewRoute("HealthShow", "/api/v1/health").HandlerFunc(handlers.HealthHandler),
		NewRoute("HealthShow", "/api/v1/health").HandlerFunc(handlers.HealthHandler).Method("HEAD"),
	}
}
