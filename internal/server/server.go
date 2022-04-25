package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/siestacloud/service-monitoring/internal/handlers"
	"github.com/siestacloud/service-monitoring/internal/server/config"

	"github.com/sirupsen/logrus"
)

//APIServer main server struct
type APIServer struct {
	config   *config.ServerConfig
	logger   *logrus.Logger
	e        *echo.Echo
	handlers *handlers.MyHandler
}

//New return point to new server
func New(config *config.ServerConfig) *APIServer {
	return &APIServer{
		config:   config,
		logger:   logrus.New(),
		e:        echo.New(),
		handlers: handlers.New(),
	}
}

//Start - method start server
func (s *APIServer) Start() error {

	// Set up a channel to listen to for interrupt signals.
	var runChan = make(chan os.Signal, 1)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.config.Server.Timeout.Server,
	)

	defer cancel()

	server := &http.Server{
		Addr:         s.config.Server.Host + ":" + s.config.Server.Port,
		ReadTimeout:  s.config.Server.Timeout.Read * time.Second,
		WriteTimeout: s.config.Server.Timeout.Write * time.Second,
		IdleTimeout:  s.config.Server.Timeout.Idle * time.Second,
	}

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureEchoRouter()

	// Run the server on a new goroutine
	go func() {
		if err := s.e.StartServer(server); err != nil {
			s.logger.Info("Fail starting http server: ", err)
			log.Fatal()
		}
	}()

	// Block on this let know, why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	s.logger.Infof("Server is shutting down due to %+v\n", interrupt)

	if err := server.Shutdown(ctx); err != nil {
		s.logger.Errorf("Server was unable to gracefully shutdown due to err: %+v", err)
		return err
	}
	s.logger.Info("Server was gracefully shutdown")
	return nil
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.Server.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	if s.config.Server.LogLevel == "debug" {
		s.logger.SetReportCaller(true)
		s.logger.SetFormatter(&logrus.TextFormatter{})
	}
	return nil
}

//configureRouter Set handlers for URL path's
func (s *APIServer) configureEchoRouter() {
	s.e.POST("/update/:type/:name/:value", s.handlers.Update())
	s.e.POST("/update/", s.handlers.UpdateJson())
	s.e.GET("/value/:type/:name", s.handlers.ShowMetric())
	s.e.GET("/value/", s.handlers.ShowMetricJSON())

	s.e.GET("/", s.handlers.ShowAllMetrics())
	// Prometheus endpoint
	// s.e.Use(middleware.Logger())
	// s.e.Use(middleware.Recover())
	// s.e.GET("/", s.handlers.Update())
	// a.e.GET("/success", s.handleSuccess)
	// a.e.POST("/pull", s.handleUploadPost)
	// a.e.GET("/metrics", Handler(promhttp.Handler())
	// a.e.GET("/css/*", s.staticHandle)
	// a.e.GET("/js/*", s.staticHandle)
	// a.e.GET("/img/*", s.staticHandle)

}
