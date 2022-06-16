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
	logrus.WithFields(logrus.Fields{"tag": "layer: transport"}).Info(message)
	return c.JSON(statusCode, errorResponse{message})
}

func infoPrint(status, message string) {
	logrus.WithFields(logrus.Fields{"Layer": "transport", "Status": status}).Info(message)
}
