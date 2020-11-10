package server

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	democtx "github.com/lfoss0612/DemoApp/context"
	"github.com/lfoss0612/DemoApp/env"
	demoerrors "github.com/lfoss0612/DemoApp/errors"
	"github.com/lfoss0612/DemoApp/logger"
	"github.com/lfoss0612/DemoApp/middleware"
	"github.com/lfoss0612/DemoApp/request"
	"github.com/lfoss0612/DemoApp/response"
)

// Server wires up routes and starts api server
type Server struct {
	shutdownReq chan bool
	router      *Router
}

// New provides a new Server
func New(routes []*Route, middlewares []MiddlewareFunc) *Server {
	s := &Server{}
	router := NewRouter()
	router.StrictSlash(true)
	router.addRoutes(routes)
	router.addMiddlewares(middlewares)
	s.router = router
}

//WaitShutdown wait for server shutdown
func (s *Server) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)
	serverLog := logger.NewLogger()

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		serverLog.Infof("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		serverLog.Infof("Shutdown request (/shutdown %v)", sig)
	}

	serverLog.Infof("Stopping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		serverLog.Infof("Shutdown request error: %v", err)
	}
}

//Shutdown shutdown server call
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}
