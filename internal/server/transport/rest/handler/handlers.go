package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/config"
	"github.com/siestacloud/service-monitoring/internal/server/service"
)

type Handler struct {
	cfg      *config.Cfg
	services *service.Service
}

func NewHandler(cfg *config.Cfg, services *service.Service) *Handler {
	return &Handler{
		cfg:      cfg,
		services: services,
	}
}

//Update POST update/:type/:name/:value
func (h *Handler) UpdateParam() echo.HandlerFunc {

	return func(c echo.Context) error {
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		defer c.Request().Body.Close()
		//Получение параметров с url
		t := c.Param("type")
		n := c.Param("name")
		v := c.Param("value")

		//Формирую новую метрику из полученного запроса
		mtrx := core.NewMetric()
		if err := mtrx.SetID(n); err != nil {
			return errResponse(c, http.StatusNotImplemented, "invalid mtrx name: "+err.Error())
		}
		if err := mtrx.SetType(t); err != nil {
			return errResponse(c, http.StatusNotImplemented, "invalid mtrx type: "+err.Error())
		}
		if err := mtrx.SetValue(v); err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid mtrx type: "+err.Error())
		}
		err := h.services.Add(n, mtrx)
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid mtrx type: "+err.Error())
		}

		if h.cfg.StoreFile != "" { //Если не указан путь до файла метрика не сохранится на диск
			if h.cfg.StoreInterval == 0 { //Если интервал сохранения равен нулю, новая метрика незамедлительно сохранится на диск
				if err := h.services.RAM.WriteLocalStorage(h.cfg.StoreFile); err != nil {
					logrus.Error("error save metric pool after request: ", err)
				}
			}
		}
		infoPrint("send client Ok", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.JSON(http.StatusOK, statusResponse{"ok"})
	}
}

//   POST /update/
func (h *Handler) UpdateJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		logrus.Info("New request on: ", c.Request().URL.String())
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logrus.Error("check err", err)
		}
		logrus.Info("/update/ mtrx from request", string(body))

		defer c.Request().Body.Close()

		mtrx := core.NewMetric()
		if err := mtrx.UnmarshalMetricJSON(body); err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusBadRequest, "")
		}
		logrus.Info(" /update/ new mtrx object  ", mtrx)

		err = h.services.Add(mtrx.GetID(), mtrx)
		if err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusBadRequest, "")
		}

		if h.cfg.StoreFile != "" { //Если не указан путь до файла метрика не сохранится на диск
			if h.cfg.StoreInterval == 0 { //Если интервал сохранения равен нулю, новая метрика незамедлительно сохранится на диск
				if err := h.services.RAM.WriteLocalStorage(h.cfg.StoreFile); err != nil {
					logrus.Error("error save metric pool after request: ", err)
				}
			}
		}
		return c.String(http.StatusOK, string(body))

	}
}

// GET /value/:type/:name
func (h *Handler) ShowMetric() echo.HandlerFunc {

	return func(c echo.Context) error {
		logrus.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()

		t := c.Param("type")
		n := c.Param("name")

		mtrx := core.NewMetric()
		if err := mtrx.SetID(n); err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		if err := mtrx.SetType(t); err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		sMtrx := h.services.LookUP(mtrx.GetID())
		if sMtrx == nil {
			logrus.Error("metric not found")
			return c.HTML(http.StatusNotFound, "")
		}
		if mtrx.GetType() != sMtrx.MType {
			logrus.Error("Metric found but type not equal")
			return c.HTML(http.StatusNotFound, "")
		}

		if sMtrx.GetType() == "counter" {
			d, err := sMtrx.GetDelta()
			if err != nil {
				logrus.Error("mtrs from storage has empty value", err)
				return c.HTML(http.StatusNotFound, "")
			}
			return c.HTML(http.StatusOK, fmt.Sprintf("%v", d))
		}
		v, err := sMtrx.GetValue()
		if err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusNotFound, "")
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("%v", v))
	}
}

// GET /
func (h *Handler) ShowAllMetrics() echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()
		mp, err := h.services.GetAlljson()
		if err != nil || mp == nil {
			return c.HTML(http.StatusNotFound, "")
		}
		return c.HTML(http.StatusOK, string(mp))
	}
}

//  POST /value/
func (h *Handler) ShowMetricJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		logrus.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()

		// message, _ := bytes.ReadAll(c.Request().Body)
		// s.l.Info(string(message))

		mtrx := core.NewMetric() // Промежуточный обьект, поля которого будут проверены
		if err := json.NewDecoder(c.Request().Body).Decode(&mtrx); err != nil {
			logrus.Error(err)
			return c.HTML(http.StatusNotFound, "")
		}
		logrus.Info(" /value/  from request: ", mtrx)

		//Произвожу поиск метрики в базе
		sMtrx := h.services.LookUP(mtrx.GetID())
		if sMtrx == nil {
			logrus.Error("metric not found")
			return c.HTML(http.StatusNotFound, "")
		}

		var buf bytes.Buffer
		err := sMtrx.MarshalMetricsinJSON(&buf)
		if err != nil {
			logrus.Error("message Unable marshal metric", err)
			return c.HTML(http.StatusOK, "")
		}
		logrus.Info("/value/ response will send: ", buf.String())
		return c.String(http.StatusOK, buf.String())
	}
}

//		message, _ := ioutil.ReadAll(c.Request().Body)
// 		log.Println("In request body: ", string(message))
// 		err := json.Unmarshal([]byte(string(message)), &m)
// 		if err != nil {
// 			log.Println("Unable decode JSON", err)
// 			return c.HTML(http.StatusBadRequest, "")
// }
