package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/siestacloud/service-monitoring/internal/server/config"
	"github.com/siestacloud/service-monitoring/internal/server/transport/rest/handler"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

//APIServer main server struct.
type Server struct {
	e *echo.Echo
	c *config.Cfg
	h *handler.Handler
}

//New return point to new server.
func NewServer(config *config.Cfg, h *handler.Handler) (*Server, error) {
	return &Server{
		e: echo.New(),
		c: config,
		h: h,
	}, nil
}

//Start - method start server
func (s *Server) Start() error {
	if err := s.cfgLogRus(); err != nil {
		return err
	}
	cfgjson, _ := json.MarshalIndent(s.c, "  ", " ")
	logrus.Info(string(cfgjson))
	// var err error
	// Set up a channel to listen to for interrupt signals.
	var runChan = make(chan os.Signal, 1)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		100*time.Second,
	)
	defer cancel()

	server := &http.Server{
		Addr: s.c.Address,
	}

	s.cfgRouter()

	// Run the server on a new goroutine
	go func() {
		if err := s.e.StartServer(server); err != nil {
			logrus.Info(err)
		}
	}()

	// Block on this let know, why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	logrus.Infof("Server is shutting down due to %+v\n", interrupt)

	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	logrus.Info("Server was gracefully shutdown")

	return nil
}

func (s *Server) cfgLogRus() error {
	level, err := logrus.ParseLevel("info")
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	if viper.GetString("server.logrus.level") == "debug" {
		logrus.SetReportCaller(true)
	}
	if viper.GetBool("server.logrus.json") {

		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}

//configureRouter Set handlers for URL path's
func (s *Server) cfgRouter() {
	// s.e.Use(s.ShowStatus)
	s.e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	s.e.POST("/update/:type/:name/:value", s.h.UpdateParam())
	s.e.GET("/value/:type/:name", s.h.ShowMetric())
	s.e.POST("/update", s.h.UpdateJSON())
	s.e.POST("/update/", s.h.UpdateJSON())
	s.e.POST("/value/", s.h.ShowMetricJSON())
	s.e.GET("/", s.h.ShowAllMetrics())
	s.e.GET("/ping", s.h.CheckDB())

}
