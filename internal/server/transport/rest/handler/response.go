package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func errResponse(c echo.Context, statusCode int, message string) error {
	logrus.WithFields(logrus.Fields{"layer": "transport", "status": statusCode}).Warn(message)
	return c.JSON(statusCode, errorResponse{message})
}

func infoPrint(status, message string) {
	logrus.WithFields(logrus.Fields{"layer": "transport", "status": status}).Info(message)
}
func errPrint(status, message string, err error) {
	logrus.WithFields(logrus.Fields{"layer": "transport", "status": status, "error": err}).Error(message)
}
