package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/siestacloud/service-monitoring/internal/handlers"
	"github.com/siestacloud/service-monitoring/internal/server/config"

	"github.com/sirupsen/logrus"
)

//APIServer main server struct
type APIServer struct {
	config   *config.ServerConfig
	logger   *logrus.Logger
	mux      *http.ServeMux
	handlers *handlers.MyHandler
}

//New return point to new server
func New(config *config.ServerConfig) *APIServer {
	return &APIServer{
		config:   config,
		logger:   logrus.New(),
		mux:      http.NewServeMux(),
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
		Handler:      s.mux,
		ReadTimeout:  s.config.Server.Timeout.Read * time.Second,
		WriteTimeout: s.config.Server.Timeout.Write * time.Second,
		IdleTimeout:  s.config.Server.Timeout.Idle * time.Second,
	}

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("Start API server on ", s.config.Server.Port)

	// Run the server on a new goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
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
func (s *APIServer) configureRouter() {
	s.mux.HandleFunc("/update/", s.handlers.Update())
	// Prometheus endpoint

}
