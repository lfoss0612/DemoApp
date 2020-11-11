package context

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lfoss0612/DemoApp/constants"
	"github.com/lfoss0612/DemoApp/env"
	"github.com/lfoss0612/DemoApp/logger"
	"github.com/lfoss0612/DemoApp/server"
)

const requestContextKey string = "context-key"

type ContextKey string

const (
	HOSTNAME      ContextKey = "HostName"
	AMZNTRACEID   ContextKey = "AmznTraceId"
	LOGGER        ContextKey = "Logger"
	HOST          ContextKey = "Host"
	URL           ContextKey = "URL"
	METHOD        ContextKey = "Method"
	BODY          ContextKey = "Body"
	HEADER        ContextKey = "Header"
	SCHEME        ContextKey = "Scheme"
	PATTERN       ContextKey = "Pattern"
	RESPONSE_BODY ContextKey = "responseBody"
	STATUS_CODE   ContextKey = "statusCode"
)

const (
	TransactionIDHeader = "tx-Correlation-id"
)

type Context struct {
	context context.Context
}

type valueOnlyContext struct {
	context.Context
}

func (valueOnlyContext) Deadline() (deadline time.Time, ok bool) { return }
func (valueOnlyContext) Done() <-chan struct{}                   { return nil }
func (valueOnlyContext) Err() error                              { return nil }

func (ctx *Context) Clone() *Context {
	clonedContext := &Context{
		context: valueOnlyContext{ctx.context},
	}
	return clonedContext
}

func (ctx *Context) NewGoRoutine(name string) *Context {
	clonedContext := ctx.Clone()
	return clonedContext
}

func BuildContextFromRequest(r *http.Request) *Context {

	ctx := context.WithValue(r.Context(), HOST, r.Host)
	ctx = context.WithValue(ctx, URL, r.URL)
	ctx = context.WithValue(ctx, METHOD, r.Method)
	ctx = context.WithValue(ctx, SCHEME, getCurrentScheme(r))
	ctx = context.WithValue(ctx, BODY, readBody(r))
	ctx = context.WithValue(ctx, HEADER, r.Header)

	ctx = context.WithValue(ctx, HOSTNAME, constants.MachineID)

	log := logger.NewLogger()
	ctx = context.WithValue(ctx, LOGGER, log)
	pattern, routeErr := server.GetRoutePattern(r)

	ctx = context.WithValue(ctx, PATTERN, pattern)

	requestCtx := &Context{
		context: ctx,
	}

	requestCtx.AddLogField(logger.Hostname, constants.MachineID)
	requestCtx.AddLogField(logger.Application, env.EnvVar.BuildInfo.App)
	environment := "local"
	if env.EnvVar != nil {
		environment = env.EnvVar.Env
	}
	requestCtx.AddLogField(logger.Environment, environment)

	// logging the route error after adding the additional log fields
	if routeErr != nil {
		log.Error("unable to determine route's pattern")
	}

	return requestCtx
}

func GetContextFromRequest(r *http.Request) (*Context, error) {
	return GetContext(r.Context())
}

func GetContext(ctx context.Context) (*Context, error) {
	requestContext, ok := ctx.Value(requestContextKey).(*Context)

	if !ok {
		requestContext = &Context{
			context: context.Background(),
		}
		return requestContext, errors.New("error retrieving request context")
	}

	return requestContext, nil
}

func (ctx *Context) AddRequestDataToLog() {
	if ctx.GetURL() != nil {
		url := fmt.Sprintf("%s://%s%s", ctx.GetScheme(), ctx.GetHost(), ctx.GetURL().Path)
		ctx.AddLogField(logger.Url, url)
	}

	if ctx.GetMethod() != "" {
		ctx.AddLogField(logger.Method, ctx.GetMethod())
	}

	if len(ctx.GetBody()) > 0 {
		ctx.AddLogField(logger.Body, ctx.GetBodyAsString())
	}
}

func (ctx *Context) AddToRequest(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), requestContextKey, ctx))
}

func (ctx *Context) GetHost() string {
	if ctx.context.Value(HOST) != nil {
		return ctx.context.Value(HOST).(string)
	}
	return ""
}

func (ctx *Context) GetURL() *url.URL {
	if ctx.context.Value(URL) != nil {
		return ctx.context.Value(URL).(*url.URL)
	}
	return nil
}

func (ctx *Context) GetMethod() string {
	if ctx.context.Value(METHOD) != nil {
		return ctx.context.Value(METHOD).(string)
	}
	return ""
}

func (ctx *Context) GetBody() []byte {
	if ctx.context.Value(BODY) != nil {
		return ctx.context.Value(BODY).([]byte)
	}
	return nil
}

func (ctx *Context) GetBodyAsString() string {
	return string(ctx.GetBody())
}

func (ctx *Context) GetPattern() string {
	if ctx.context.Value(PATTERN) != nil {
		return ctx.context.Value(PATTERN).(string)
	}
	return ""
}

func (ctx *Context) SetPattern(pattern string) {
	ctx.context = context.WithValue(ctx.context, PATTERN, pattern)
}

func (ctx *Context) GetAmznTraceId() string {
	if ctx.context.Value(AMZNTRACEID) != nil {
		return ctx.context.Value(AMZNTRACEID).(string)
	}
	return ""
}

func (ctx *Context) SetAmznTraceId(amznTraceId string) {
	ctx.context = context.WithValue(ctx.context, AMZNTRACEID, amznTraceId)
}

func (ctx *Context) GetResponseBody() interface{} {
	if ctx.context.Value(RESPONSE_BODY) != nil {
		return ctx.context.Value(RESPONSE_BODY)
	}
	return ""
}

func (ctx *Context) SetResponseBody(responseBody interface{}) {
	ctx.context = context.WithValue(ctx.context, RESPONSE_BODY, responseBody)
}

func (ctx *Context) GetStatusCode() int {
	if ctx.context.Value(STATUS_CODE) != nil {
		return ctx.context.Value(STATUS_CODE).(int)
	}
	return 0
}

func (ctx *Context) SetStatusCode(statusCode int) {
	ctx.context = context.WithValue(ctx.context, STATUS_CODE, statusCode)
}

func (ctx *Context) Logger() *logger.Logger {
	if ctx.context.Value(LOGGER) != nil {
		return ctx.context.Value(LOGGER).(*logger.Logger)
	}
	return nil
}

func (ctx *Context) GetHeader() http.Header {
	if ctx.context.Value(HEADER) != nil {
		return ctx.context.Value(HEADER).(http.Header)
	}
	return nil
}

func (ctx *Context) GetScheme() string {
	if ctx.context.Value(SCHEME) != nil {
		return ctx.context.Value(SCHEME).(string)
	}
	return ""
}

func (ctx *Context) GetContext() context.Context {
	return ctx.context
}

func (ctx *Context) AddLogField(key logger.LogField, value interface{}) {
	ctx.Logger().AddConstantField(key, value)
}

func (ctx *Context) GetCustom(key interface{}) interface{} {
	return ctx.context.Value(key)
}

func (ctx *Context) SetCustom(key interface{}, value interface{}) {
	ctx.context = context.WithValue(ctx.context, key, value)
}
