package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/siestacloud/service-monitoring/internal/database"
	"github.com/siestacloud/service-monitoring/internal/metricscustom"
	"github.com/siestacloud/service-monitoring/internal/storage"
)

type MyHandler struct {
	s *storage.Storage
}

func New() *MyHandler {
	return &MyHandler{
		s: storage.New(),
	}
}

//test NOT USEING
func (h *MyHandler) NotUSing() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Only text/plain requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		m := strings.ReplaceAll(r.URL.Path, "/update/", "")
		ms := strings.Split(m, "/")
		if len(ms) != 3 {
			http.Error(w, "", http.StatusNotFound)
		}

		file, err := database.New("metrics.json")
		if err != nil {
			log.Printf("ERR storage %s", err)
			http.Error(w, "Unable open file for metric", http.StatusMethodNotAllowed)
			return
		}
		defer file.Close()

		buf, err := file.ReadMetrics()
		if err != nil {
			log.Printf("ERR read storage %s", err)
			return
		}

		var mp = metricscustom.NewMetricsPool()

		err = json.Unmarshal(buf, &mp.M)
		if err != nil {
			log.Println(err)
			return
		}

		js, err := json.MarshalIndent(mp.M, "", " ")
		if err != nil {
			log.Printf("ERR %s", err)
			return
		}
		_, err = file.Write(js)
		if err != nil {
			log.Printf("ERR %s", err)
			return
		}

	}
}

//Update POST update/:type/:name/:value
func (h *MyHandler) Update() echo.HandlerFunc {

	return func(c echo.Context) error {
		fmt.Println("New request on: ", c.Request().URL.Path)
		if c.Request().Method != http.MethodPost {
			return c.HTML(http.StatusMethodNotAllowed, `"{"message":"Method Not Allowed"}"`)
		}
		defer c.Request().Body.Close()

		t := c.Param("type")
		n := c.Param("name")
		v := c.Param("value")

		s, status := metricscustom.NewMetric(t, n, v)
		if status != "" {
			switch status {
			case "unknown metric type":
				return c.HTML(http.StatusNotImplemented, `"{"message":"Unknown Metric Type"}"`)
			case "incorrect value":
				return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect Metric Value"}"`)
			default:
				return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect Metric"}"`)
			}
		}

		fmt.Println("New metric: ", s)
		h.s.Update(s)
		fmt.Println("In Storage: ")
		for k, v := range h.s.Mp.M {
			fmt.Printf("	Metric:  %s\n	    Name:%s\n	    Value:%v\n		Delta:%v\n	    Type:%s\n\n", k, v.ID, v.Value, v.Delta, v.MType)
		}
		return c.HTML(http.StatusOK, `"{"message":"Successful Metric Add/Update"}"`)
	}
}

//   POST /update/
func (h *MyHandler) UpdateJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		fmt.Println("New request on: ", c.Request().URL.Path)
		c.Response().Header().Add("Content-Type", "application/json")
		if c.Request().Method != http.MethodPost {
			return c.HTML(http.StatusMethodNotAllowed, `"{"message":"Method Not Allowed"}"`)
		}
		defer c.Request().Body.Close()
		m := metricscustom.Metric{}
		if err := json.NewDecoder(c.Request().Body).Decode(&m); err != nil {
			log.Println("Unable decode JSON", err)
			return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect metric"}"`)
		}

		status := m.Check()
		if status != "" {
			switch status {
			case "unknown metric type":
				return c.HTML(http.StatusNotImplemented, `"{"message":"Unknown Metric Type"}"`)
			case "incorrect value":
				return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect Metric Value"}"`)
			case "empty value":
				return c.HTML(http.StatusBadRequest, `"{"message":"empty value"}"`)
			case "empty delta":
				return c.HTML(http.StatusBadRequest, `"{"message":"empty delta"}"`)
			default:
				return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect Metric"}"`)
			}
		}

		fmt.Println("Metric from request: ", m)
		h.s.Update(&m)
		fmt.Println("Metric from storage: ", h.s.Mp.M[m.ID])
		return c.HTML(http.StatusOK, `"{"message":"Successful Metric Add/Update json"}"`)
	}
}

// GET /value/:type/:name
func (h *MyHandler) ShowMetric() echo.HandlerFunc {

	return func(c echo.Context) error {
		fmt.Println("New request on: ", c.Request().URL.Path)

		defer c.Request().Body.Close()
		t := c.Param("type")
		n := c.Param("name")
		metric := h.s.Take(t, n)
		if metric == nil {
			return c.HTML(http.StatusNotFound, `"{"message":"Metric Not Found"}"`)
		}
		if t == "counter" {
			return c.HTML(http.StatusOK, fmt.Sprintf("%v", metric.Delta))
		}
		return c.HTML(http.StatusOK, fmt.Sprintf("%v", metric.Value))
	}
}

// GET /
func (h *MyHandler) ShowAllMetrics() echo.HandlerFunc {

	return func(c echo.Context) error {
		fmt.Println("New request on: ", c.Request().URL.Path)

		defer c.Request().Body.Close()

		mp, err := h.s.TakeAll()
		if err != nil || mp == nil {
			return c.HTML(http.StatusNotFound, "")
		}

		return c.HTML(http.StatusOK, string(mp))
	}
}

//  POST /value/
func (h *MyHandler) ShowMetricJSON() echo.HandlerFunc {

	return func(c echo.Context) error {
		log.Println("New request on: ", c.Request().URL.Path)
		c.Response().Header().Add("Content-Type", "application/json")
		if c.Request().Method != http.MethodPost {
			return c.HTML(http.StatusMethodNotAllowed, `"{"message":"Method Not Allowed"}"`)
		}
		defer c.Request().Body.Close()
		m := metricscustom.Metric{}
		// if err := json.NewDecoder(c.Request().Body).Decode(&m); err != nil {
		// 	log.Println("Unable decode JSON", err)
		// 	return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect metric"}"`)
		// }

		message, _ := ioutil.ReadAll(c.Request().Body)
		log.Println("In request body: ", string(message))
		err := json.Unmarshal([]byte(string(message)), &m)
		if err != nil {
			log.Println("Unable decode JSON", err)
			return c.HTML(http.StatusBadRequest, `"{"message":"Incorrect metric"}"`)
		}
		defer c.Request().Body.Close()

		metric := h.s.Take(m.MType, m.ID)
		if metric == nil {
			return c.HTML(http.StatusOK, `"{"message":"Metric Not Found"}"`)
		}
		var buf bytes.Buffer
		err = metric.MarshalMetricsinJSON(&buf)
		if err != nil {
			log.Panicln(err)
			return c.HTML(http.StatusOK, `"{"message":"Unable marshal metric"}"`)

		}
		return c.HTML(http.StatusOK, buf.String())
	}
}
