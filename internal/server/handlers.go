package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/siestacloud/service-monitoring/internal/mtrx"
)

//Update POST update/:type/:name/:value
func (s *APIServer) UpdateParam() echo.HandlerFunc {

	return func(c echo.Context) error {
		s.l.Infoln("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()
		//Получение параметров с url
		t := c.Param("type")
		n := c.Param("name")
		v := c.Param("value")

		//Формирую новую метрику из полученного запроса
		mtrx := mtrx.NewMetric()
		if err := mtrx.SetID(n); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		if err := mtrx.SetType(t); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		if err := mtrx.SetValue(v); err != nil {
			s.l.Error("incorrect value in request", err)
			return c.HTML(http.StatusBadRequest, "")
		}
		s.l.Info("NEW   ", mtrx)
		s.l.Debug(mtrx)
		if !s.s.Mp.Update(n, *mtrx) {
			s.l.Warn("unable find and update metric in storage")
			s.l.Warn("try add new metric")
			if !s.s.Mp.Add(n, *mtrx) {
				s.l.Error("unable add metric in storage")
				return c.HTML(http.StatusBadRequest, "")
			}
			s.l.Warn("OK")
		}
		s.s.Mp.PrintAll()
		return c.HTML(http.StatusOK, "")
	}
}

//   POST /update/
func (s *APIServer) UpdateJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		s.l.Info("New request on: ", c.Request().URL.String())
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			s.l.Error("check err", err)
		}
		s.l.Info("/update/ mtrx from request", string(body))

		defer c.Request().Body.Close()

		mtrx := mtrx.NewMetric()
		if err := mtrx.UnmarshalMetricJSON(body); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusBadRequest, "")
		}
		s.l.Info(" /update/ new mtrx object  ", mtrx)
		if !s.s.Mp.Update(mtrx.GetID(), *mtrx) {
			s.l.Warn("unable find and update metric in storage")
			s.l.Warn("try add new metric")
			if !s.s.Mp.Add(mtrx.GetID(), *mtrx) {
				s.l.Error("unable add metric in storage")
				return c.HTML(http.StatusBadRequest, "")
			}
			s.l.Warn("OK")
		}

		//Произвожу поиск метрики в базе
		sMtrx := s.s.Mp.LookUP(mtrx.GetID())
		if sMtrx == nil {
			s.l.Error("/update/ metric not found in storage")

		}
		d, _ := sMtrx.GetDelta()
		v, _ := sMtrx.GetValue()
		s.l.Info(" /update/  mtrx object from storage  ", sMtrx, "dalta: ", d, "  value: ", v)
		// s.s.Mp.PrintAll()
		return c.String(http.StatusOK, string(body))

	}
}

// GET /value/:type/:name
func (s *APIServer) ShowMetric() echo.HandlerFunc {

	return func(c echo.Context) error {
		s.l.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()

		t := c.Param("type")
		n := c.Param("name")

		mtrx := mtrx.NewMetric()
		if err := mtrx.SetID(n); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		if err := mtrx.SetType(t); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotImplemented, "")
		}
		sMtrx := s.s.Mp.LookUP(mtrx.GetID())
		if sMtrx == nil {
			s.l.Error("metric not found")
			return c.HTML(http.StatusNotFound, "")
		}
		if mtrx.GetType() != sMtrx.MType {
			s.l.Error("Metric found but type not equal")
			return c.HTML(http.StatusNotFound, "")
		}

		if sMtrx.GetType() == "counter" {
			d, err := sMtrx.GetDelta()
			if err != nil {
				s.l.Error(err)
				return c.HTML(http.StatusNotFound, "")
			}
			return c.HTML(http.StatusOK, fmt.Sprintf("%v", d))
		}
		v, err := sMtrx.GetValue()
		if err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotFound, "")
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("%v", v))
	}
}

// GET /
func (s *APIServer) ShowAllMetrics() echo.HandlerFunc {

	return func(c echo.Context) error {
		s.l.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()
		mp, err := s.s.TakeAll()
		if err != nil || mp == nil {
			return c.HTML(http.StatusNotFound, "")
		}
		return c.HTML(http.StatusOK, string(mp))
	}
}

//  POST /value/
func (s *APIServer) ShowMetricJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		s.l.Info("New request on: ", c.Request().URL.String())
		defer c.Request().Body.Close()

		// message, _ := bytes.ReadAll(c.Request().Body)
		// s.l.Info(string(message))

		mtrx := mtrx.NewMetric() // Промежуточный обьект, поля которого будут проверены
		if err := json.NewDecoder(c.Request().Body).Decode(&mtrx); err != nil {
			s.l.Error(err)
			return c.HTML(http.StatusNotFound, "")
		}
		s.l.Info(" /value/  from request: ", mtrx)

		//Произвожу поиск метрики в базе
		sMtrx := s.s.Mp.LookUP(mtrx.GetID())
		if sMtrx == nil {
			s.l.Error("metric not found")
			return c.HTML(http.StatusNotFound, "")
		}

		var buf bytes.Buffer
		err := sMtrx.MarshalMetricsinJSON(&buf)
		if err != nil {
			s.l.Error("message Unable marshal metric", err)
			return c.HTML(http.StatusOK, "")
		}
		s.l.Info("/value/ response will send: ", buf.String())
		return c.String(http.StatusOK, buf.String())
	}
}

//		message, _ := ioutil.ReadAll(c.Request().Body)
// 		log.Println("In request body: ", string(message))
// 		err := json.Unmarshal([]byte(string(message)), &m)
// 		if err != nil {
// 			log.Println("Unable decode JSON", err)
// 			return c.HTML(http.StatusBadRequest, "")
// 		}
