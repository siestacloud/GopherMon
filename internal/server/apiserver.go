package server

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIServer struct {
	config *Config
}

type UpdateHandler struct {
	DB Storage
}

func New(config *Config) *APIServer {
	return &APIServer{config: config}
}

func NewUpdateHandler() *UpdateHandler {
	updater := new(UpdateHandler)
	updater.DB = NewDB()
	return updater
}

func (handler *UpdateHandler) getAllMetrics(c echo.Context) error {
	metrics := handler.DB.GetAll()
	resp := c.Response()
	resp.Header().Set("Content-Type", "application/json")

	return c.JSON(http.StatusOK, metrics)
}

func (handler *UpdateHandler) getMetric(c echo.Context) error {
	t := c.Param("type")
	name := c.Param("name")
	log.Printf("Get Metric type:%s name:%s", t, name)
	val, err := handler.DB.Get(t, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, val)
}
func (handler *UpdateHandler) postMetric(c echo.Context) error {

	t := c.Param("type")
	name := c.Param("name")
	val := c.Param("value")
	log.Printf("Post Metric type:%s name:%s value:%s", t, name, val)
	err := handler.DB.Set(t, name, val)
	switch {
	case err == nil:
		return nil
	case err.Error() == "invalid type":
		return echo.NewHTTPError(http.StatusNotImplemented, err.Error())
	default:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

}

func (s *APIServer) Start(ctx context.Context) error {
	updater := NewUpdateHandler()
	e := echo.New()
	e.GET("/", updater.getAllMetrics)
	e.GET("/value/:type/:name", updater.getMetric)
	e.POST("/update/:type/:name/:value", updater.postMetric)
	return e.Start(s.config.BindAddr)
}
