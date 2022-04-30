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
	"github.com/siestacloud/service-monitoring/internal/server/config"
	"github.com/siestacloud/service-monitoring/internal/storage"

	"github.com/sirupsen/logrus"
)

//APIServer main server struct.
type APIServer struct {
	c *config.ServerConfig
	s *storage.Storage
	l *logrus.Logger
	e *echo.Echo
}

//New return point to new server.
func New(config *config.ServerConfig) *APIServer {
	return &APIServer{
		s: storage.New(),
		l: logrus.New(),
		e: echo.New(),
		c: config,
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
		s.c.Server.Timeout.Server,
	)

	defer cancel()

	server := &http.Server{
		Addr:         s.c.Address,
		ReadTimeout:  s.c.Server.Timeout.Read * time.Second,
		WriteTimeout: s.c.Server.Timeout.Write * time.Second,
		IdleTimeout:  s.c.Server.Timeout.Idle * time.Second,
	}

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureEchoRouter()

	// Run the server on a new goroutine
	go func() {
		if err := s.e.StartServer(server); err != nil {
			s.l.Info("Fail starting http server: ", err)
			log.Fatal()
		}
	}()

	// Block on this let know, why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	s.l.Infof("Server is shutting down due to %+v\n", interrupt)

	if err := server.Shutdown(ctx); err != nil {
		s.l.Errorf("Server was unable to gracefully shutdown due to err: %+v", err)
		return err
	}
	s.l.Info("Server was gracefully shutdown")
	return nil
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.c.Server.LogLevel)
	if err != nil {
		return err
	}
	s.l.SetLevel(level)
	if s.c.Server.LogLevel == "debug" {
		s.l.SetReportCaller(true)
		s.l.SetFormatter(&logrus.TextFormatter{})
	}
	return nil
}

//configureRouter Set handlers for URL path's
func (s *APIServer) configureEchoRouter() {
	s.e.POST("/update/:type/:name/:value", s.UpdateParam())
	s.e.POST("/update", s.UpdateJSON())
	s.e.POST("/update/", s.UpdateJSON())
	s.e.GET("/value/:type/:name", s.ShowMetric())
	s.e.POST("/value/", s.ShowMetricJSON())
	s.e.GET("/", s.ShowAllMetrics())

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
