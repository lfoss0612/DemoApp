package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/netutil"

	"github.com/lfoss0612/DemoApp/logger"
)

// Server wires up routes and starts api server
type Server struct {
	shutdownReq chan bool
	router      *router
	listener    net.Listener
	Port        string
	Address     string
}

// New provides a new Server
func New(routes []Route, middlewares []MiddlewareFunc, port string, maxSimultaneousConnections int) (*Server, error) {
	s := &Server{}
	router := newRouter()
	router.StrictSlash(true)
	router.addRoutes(routes)
	router.addMiddlewares(middlewares)
	s.router = router

	s.Port = port
	// Start the web server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		MaxHeaderBytes: 1 << 13,
	}
	s.Address = srv.Addr
	listener, listenerErr := net.Listen("tcp", srv.Addr)

	if listenerErr != nil {
		return nil, fmt.Errorf("unable to listen on %s", srv.Addr)
	}

	s.listener = listener

	serveErr := srv.Serve(netutil.LimitListener(listener, maxSimultaneousConnections))

	if serveErr != nil {
		return nil, fmt.Errorf("error in starting server at address %s", srv.Addr)
	}

	return s, nil
}

// Close used to close the listener on the server
func (s *Server) Close() error {
	return s.listener.Close()
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
