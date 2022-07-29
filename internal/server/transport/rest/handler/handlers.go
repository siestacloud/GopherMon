package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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
		// При подключенной postgres
		if h.cfg.URLPostgres != "" {
			_, err := h.services.MtrxList.Add(mtrx)
			if err != nil {
				return errResponse(c, http.StatusBadRequest, "invalid mtrx type: "+err.Error())
				// errPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String(), err)
			}
			infoPrint("send client Ok", "request: "+h.cfg.Address+c.Request().URL.String())
			return c.JSON(http.StatusOK, statusResponse{"mtrx save in postgres db"})
		}
		// При отключенной postgres
		// метрика сохр в оперативной памяти
		err := h.services.RAM.Add(n, mtrx)
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "invalid mtrx type: "+err.Error())
		}

		// При подключенном local storage
		if h.cfg.StoreFile != "" { //Если не указан путь до файла метрика не сохранится на диск
			if h.cfg.StoreInterval == 0 { //Если интервал сохранения равен нулю, новая метрика незамедлительно сохранится на диск
				if err := h.services.RAM.WriteLocalStorage(h.cfg.StoreFile); err != nil {
					logrus.Error("error save metric pool after request: ", err)
				}
			}
		}

		infoPrint("send client Ok", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.JSON(http.StatusOK, statusResponse{"mtrx save"})
	}
}

//   POST /update/
func (h *Handler) UpdateJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, "unable read data from body request: "+err.Error())
		}
		infoPrint("in tune", "	mtrx in request: "+string(body))
		defer c.Request().Body.Close()

		//Формирую новую метрику из полученного запроса
		mtrx := core.NewMetric()
		if err := mtrx.UnmarshalMetricJSON(body); err != nil {
			return errResponse(c, http.StatusBadRequest, "unable read data from body request: "+err.Error())
		}
		// проверяю целостность данных и подпись
		err = h.services.CheckHash(h.cfg.Key, mtrx)
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "unable compare hash: "+err.Error())
		}
		infoPrint("in tune", "	success compared hash")

		// При подключенной postgres
		if h.cfg.URLPostgres != "" {
			_, err := h.services.MtrxList.Add(mtrx)
			if err != nil {
				return errResponse(c, http.StatusBadRequest, "invalid mtrx type: "+err.Error())
			}
			infoPrint("send client Ok", "request: "+h.cfg.Address+c.Request().URL.String())
			return c.JSON(http.StatusOK, statusResponse{"mtrx save in postgres db"})
		}

		// При отключенной postgres
		// метрика сохр в оперативной памяти
		err = h.services.RAM.Add(mtrx.GetID(), mtrx)
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "unable read data from body request: "+err.Error())
		}

		// При подключенном local storage
		if h.cfg.StoreFile != "" { //Если не указан путь до файла метрика не сохранится на диск
			if h.cfg.StoreInterval == 0 { //Если интервал сохранения равен нулю, новая метрика незамедлительно сохранится на диск
				if err := h.services.RAM.WriteLocalStorage(h.cfg.StoreFile); err != nil {
					logrus.Error("error save metric pool after request: ", err)
				}
			}
		}
		infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.String(http.StatusOK, string(body))
	}
}

// GET /value/:type/:name
func (h *Handler) ShowMetric() echo.HandlerFunc {

	return func(c echo.Context) error {
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		defer c.Request().Body.Close()
		var err error

		t := c.Param("type")
		n := c.Param("name")

		//Формирую новую метрику из полученного запроса
		mtrx := core.NewMetric()
		sMtrx := core.NewMetric()
		if err := mtrx.SetID(n); err != nil {
			return errResponse(c, http.StatusNotImplemented, "invalid mtrx name: "+err.Error())
		}
		if err := mtrx.SetType(t); err != nil {
			return errResponse(c, http.StatusNotImplemented, "invalid mtrx type: "+err.Error())
		}
		// При подключенной postgres
		if h.cfg.URLPostgres != "" {
			sMtrx, err = h.services.MtrxList.Get(mtrx.GetID())
			if err != nil {
				return errResponse(c, http.StatusNotFound, "mtrx not found in postges")
			}
		} else {
			// При отключенной postgres
			// Ищу метрику в RAM
			sMtrx = h.services.LookUP(mtrx.GetID())
			if sMtrx == nil {
				return errResponse(c, http.StatusNotFound, "mtrx not found")
			}
		}

		if mtrx.GetType() != sMtrx.MType {
			return errResponse(c, http.StatusNotFound, "mtrx found but types not equal")
		}

		if sMtrx.GetType() == "counter" {
			d, err := sMtrx.GetDelta()
			if err != nil {
				return errResponse(c, http.StatusNotFound, "mtrx from storage has empty value "+err.Error())
			}
			infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String())
			return c.HTML(http.StatusOK, fmt.Sprintf("%v", d))
		}
		v, err := sMtrx.GetValue()
		if err != nil {
			return errResponse(c, http.StatusNotFound, "mtrx from storage has empty value "+err.Error())
		}
		infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.HTML(http.StatusOK, fmt.Sprintf("%v", v))
	}
}

// GET /
func (h *Handler) ShowAllMetrics() echo.HandlerFunc {
	return func(c echo.Context) error {
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		defer c.Request().Body.Close()
		mp, err := h.services.GetAlljson()
		if err != nil || mp == nil {
			return errResponse(c, http.StatusNotFound, "mtrx pool empty "+err.Error())
		}
		infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.HTML(http.StatusOK, string(mp))
	}
}

//  POST /value/
func (h *Handler) ShowMetricJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		defer c.Request().Body.Close()

		var err error

		//Формирую новую метрику из полученного запроса
		mtrx := core.NewMetric()
		sMtrx := core.NewMetric()

		if err := json.NewDecoder(c.Request().Body).Decode(&mtrx); err != nil {
			return errResponse(c, http.StatusNotFound, "unable decode mtrx"+err.Error())
		}
		infoPrint("in tune", "	mtrx in request: "+mtrx.GetID())
		b, _ := json.Marshal(mtrx)
		fmt.Println("IN VALUE", string(b))

		// При подключенной postgres
		if h.cfg.URLPostgres != "" {
			sMtrx, err = h.services.MtrxList.Get(mtrx.GetID())
			if err != nil {
				return errResponse(c, http.StatusNotFound, "mtrx not found in postges")
			}
		} else {
			// При отключенной postgres
			// Ищу метрику в RAM
			sMtrx = h.services.LookUP(mtrx.GetID())
			if sMtrx == nil {
				return errResponse(c, http.StatusNotFound, "mtrx not found")
			}
		}
		infoPrint("in tune", fmt.Sprintf("	mtrx in db: %+v", mtrx))

		// генерирую хеш
		err = sMtrx.SetHash(h.cfg.Key)
		if err != nil {
			return errResponse(c, http.StatusNotFound, "unable set hash")

		}
		var buf bytes.Buffer
		err = sMtrx.MarshalMetricsinJSON(&buf)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, "unable convert mtrx to json format before send"+err.Error())
		}
		infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String()+" response will send: "+buf.String())
		return c.String(http.StatusOK, buf.String())
	}
}

// GET /ping
func (h *Handler) CheckDB() echo.HandlerFunc {
	return func(c echo.Context) error {
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		defer c.Request().Body.Close()

		if err := h.services.MtrxList.TestDB(); err != nil {
			return errResponse(c, http.StatusInternalServerError, "postgres db fail connect")
		}

		infoPrint("200", "request: "+h.cfg.Address+c.Request().URL.String())
		return c.JSON(http.StatusOK, statusResponse{"ok"})
	}
}

func (h *Handler) MultupleMtrxJSON() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Response().Header().Add("Content-Type", "application/json")
		infoPrint("in tune", "request: "+h.cfg.Address+c.Request().URL.String())
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return errResponse(c, http.StatusInternalServerError, "unable read data from body request: "+err.Error())
		}
		infoPrint("in tune", "	mtrx in request: "+string(body))

		defer c.Request().Body.Close()

		mtrxCase, err := core.UnmarshalMetricCaseJSON(body)
		if err != nil {
			return errResponse(c, http.StatusBadRequest, "unable read new client mtrx from body request: "+err.Error())
		}

		b, _ := json.Marshal(mtrxCase)
		fmt.Println(string(b))

		infoPrint("in tune", "	success parse in object mtrx")
		var val int
		// При подключенной postgres
		if h.cfg.URLPostgres != "" {
			val, err = h.services.MtrxList.Flush(mtrxCase)
			if err != nil {
				return errResponse(c, http.StatusNotFound, "unable insert or update mrtxs: "+err.Error())
			}
		}

		return c.JSON(http.StatusOK, statusResponse{strconv.Itoa(val)})
	}
}
