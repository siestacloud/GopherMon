package server

import (
	"github.com/labstack/echo/v4"
)

// Process is the middleware function.
func (s *APIServer) ShowStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		// s.mutex.Lock()
		// defer s.mutex.Unlock()
		// s.RequestCount++
		status := c.Request().Header
		s.l.Info("Middleware: ", status)
		return nil
	}
}
