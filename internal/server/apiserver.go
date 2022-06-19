package server

import (
	"context"
	"fmt"
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
	var html string
	html = "<html>\n"
	metrics := handler.DB.GetAll()
	for k, v := range metrics {
		html += fmt.Sprintf("<p>%s: %s</p>\n", k, v)
	}
	html += "</html>"
	return c.HTML(http.StatusOK, html)
}

func (handler *UpdateHandler) getMetric(c echo.Context) error {
	t := c.Param("type")
	name := c.Param("name")
	val, err := handler.DB.Get(t, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "text/plain")

	return c.HTML(http.StatusOK, val.String())
}
func (handler *UpdateHandler) postMetric(c echo.Context) error {

	t := c.Param("type")
	name := c.Param("name")
	val := c.Param("value")
	err := handler.DB.Set(t, name, val)
	switch {
	case err == nil:
		log.Println(handler.DB)
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
