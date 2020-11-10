package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/pkg/errors"

	"golang.org/x/net/netutil"

	"github.com/lfoss0612/DemoApp/env"
	"github.com/lfoss0612/DemoApp/logger"
	"github.com/lfoss0612/DemoApp/server"
	"github.com/lfoss0612/DemoApp/routes"
	"github.com/lfoss0612/DemoApp/middleware"
)

// Variables specified at build time.
var (
	AppName      = "Demo"
	BuildVersion = "v1.x.x"
	BuildDate    = "yyyy-mm-dd"
	BuildCommit  = "git"
)

var serverLog *logger.Logger

func init() {
	readEnv()

	appName := fmt.Sprintf("%s(%s)", AppName, env.EnvVar.Env)

	logInfo(appName)
}

func readEnv() {
	env.Init()

	buildInfo := env.Build{App: AppName, Date: BuildDate, Commit: BuildCommit}
	env.EnvVar.BuildInfo = buildInfo

	setupLogger()

	return
}

func setupLogger() {
	logger.Init(env.EnvVar.LogLevel, os.Stdout, logger.CustomFormatter{})
	serverLog = logger.NewLogger().WithLogField(logger.Application, AppName).WithLogField(logger.Environment, env.EnvVar.Env)
}

func logInfo(appName string) {
	// Log server startup

	serverLog.Infof("Starting %s...", appName)
	serverLog.Infof("Build date: %s", BuildDate)
	serverLog.Infof("GIT commit ID: %s", BuildCommit)
	serverLog.Infof("System architecture: %s", runtime.GOARCH)
	serverLog.Infof("System OS: %s", runtime.GOOS)
	serverLog.Infof("Go version: %s", runtime.Version())
	serverLog.Infof("Logging at level %s and higher", env.EnvVar.LogLevel)

	//display Environment Variables on Startup
	serverLog.Info("Loaded Environment Variables")
}
func initServer() (*server.Server, error) {
	s := server.New(
		routes.Routes(),
		middleware.MiddleWares()		
	)

	return s, nil
}

func main() {

	// Setup server routes
	s, err := initServer()
	if err != nil {
		serverLog.WithError(err).Fatalf("Unable to Initialize Server")
	}

	// Setup Middleware chain
	router := s.BuildRouter()

	// Start the web server
	serverLog.Infof("Listening on %s", env.EnvVar.Port)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", env.EnvVar.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 13,
	}

	l, listenerErr := net.Listen("tcp", srv.Addr)

	if listenerErr != nil {
		serverLog.WithError(listenerErr).Fatalf("unable to listen on %s", srv.Addr)
	}

	defer func() {
		_ = l.Close() // nolint
	}()

	serveErr := srv.Serve(netutil.LimitListener(l, env.EnvVar.ServerMaxSimultaneousConnections))

	if serveErr != nil {
		serverLog.Error(errors.Wrap(serveErr, "Error in Starting the Server."))
	}

	s.WaitShutdown()

	serverLog.Infof("Server Stopped!")
}
