package server

import (
	"context"
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
func New(config *config.ServerConfig) (*APIServer, error) {
	sf, err := storage.NewStorage(config.StoreFile)
	if err != nil {
		return nil, err
	}
	if config.StoreFile != "" { //Если передан путь до файла и указан флаг восстановить, метрики будут прочитаны из хранилища
		if config.Restore {
			err := sf.ReadStorage()
			if err != nil {

				return nil, err
			}
		}
	}

	return &APIServer{
		s: sf,
		l: logrus.New(),
		e: echo.New(),
		c: config,
	}, nil
}

//Start - method start server
func (s *APIServer) Start() error {
	s.l.Warn("cfg: ", s.c)
	s.l.Warn("mtrx from storage: ", s.s.Mp)

	// var err error
	// Set up a channel to listen to for interrupt signals.
	var runChan = make(chan os.Signal, 1)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	// defer cancel()
	defer func() {
		s.l.Warn("Server saving metrics pool...")
		if s.c.StoreFile != "" {
			if err := s.s.WriteStorage(); err != nil {
				s.l.Error("failed save metrics pool: ", err)
				cancel()
			}
		}

		s.l.Info("Server was gracefully shutdown")
		cancel()
	}()

	server := &http.Server{
		Addr: s.c.Address,
		// ReadTimeout:  s.c.Server.Timeout.Read,
		// WriteTimeout: s.c.Server.Timeout.Write,
		// IdleTimeout:  s.c.Server.Timeout.Idle,
	}

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureEchoRouter()

	// Run the server on a new goroutine
	go func() {
		if err := s.e.StartServer(server); err != nil {
			s.l.Info(err)
		}
	}()

	if s.c.StoreFile != "" {
		if s.c.StoreInterval != 0 {

			go s.StoreInterval()
		}
	}

	// Block on this let know, why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	s.l.Infof("Server is shutting down due to %+v\n", interrupt)

	if err := server.Shutdown(ctx); err != nil {
		s.l.Errorf("Server was unable to gracefully shutdown due to err: %+v", err)
		return err
	}

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
}

func (s *APIServer) StoreInterval() {
	for {
		time.Sleep(s.c.StoreInterval)
		if err := s.s.WriteStorage(); err != nil {
			s.l.Error("error store interval: ", err)
		}
		s.l.Info("Storage update")
	}
}
