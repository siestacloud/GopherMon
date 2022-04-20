package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

//handleUpload upload file /upload
func (h *MyHandler) HandleUpdate() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Only text/plain requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		m, err := metricURI(r.URL.Path)
		if err != nil {
			log.Println(err)
			return
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
		var check bool
		for _, v := range mp.M {
			if v.Name == m.Name {
				if v.Types == "counter" {
					v.Value += m.Value
				}
				if v.Types == "gauge" {
					v.Value = m.Value
				}
				check = true
			}
		}
		if !check {
			mp.M[m.Name] = *m
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
		// if err = mp.WriteMetricJSON(file.F); err != nil {
		// 	http.Error(w, "Unable write metrics", http.StatusMethodNotAllowed)
		// 	return
		// }
		// fmt.Println(file.F.Name())

		// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
	}
}

//Update upload file /upload
func (h *MyHandler) Update() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Only text/plain requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		m, err := metricURI(r.URL.Path)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("New metric: %s %v %s\n\n\n", m.Name, m.Value, m.Types)
		h.s.Update(m)
		fmt.Println("In Storage: ")
		for k, v := range h.s.Mp.M {
			fmt.Printf("	Metric:  %s\n	    Name:%s\n	    Value:%v\n	    Type:%s\n\n", k, v.Name, v.Value, v.Types)
		}

	}
}

//makeMetric take metric from URI
func metricURI(uri string) (*metricscustom.Metric, error) {
	m := strings.ReplaceAll(uri, "/update/", "")
	ms := strings.Split(m, "/")
	if len(ms) != 3 {
		return nil, errors.New("Incorrect len URI")
	}
	v, err := strconv.ParseUint(ms[2], 10, 64)
	if err != nil {
		return nil, err
	}
	s := metricscustom.Metric{Name: ms[1], Types: ms[0], Value: v}

	return &s, nil
}
