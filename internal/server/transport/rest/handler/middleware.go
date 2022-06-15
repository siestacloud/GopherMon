package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Process is the middleware function.
func (h *Handler) ShowStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		h := c.Request().Header
		logrus.Info("Middleware: request have Accept-Encoding header ", h)

		if !strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip") {
			//Заголовка Accept-Encoding нет
			return next(c)
		}

		c.Response().Header().Add("Content-Encoding", "gzip")
		// s.mutex.Lock()
		// defer s.mutex.Unlock()
		// s.RequestCount++

		return next(c)
	}
}
