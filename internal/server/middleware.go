package server

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// Process is the middleware function.
func (s *APIServer) ShowStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip") { //Заголовка нет
			return next(c)
		}

		c.Response().Header().Add("Content-Encoding", "gzip")
		// s.mutex.Lock()
		// defer s.mutex.Unlock()
		// s.RequestCount++
		h := c.Request().Header
		s.l.Info("Middleware: request have Accept-Encoding header ", h)
		return next(c)
	}
}
