package handlers

import (
	"encoding/json"
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

		m := strings.ReplaceAll(r.URL.Path, "/update/", "")
		ms := strings.Split(m, "/")
		if len(ms) != 3 {
			http.Error(w, "", http.StatusNotFound)
		}
		v, err := strconv.ParseUint(ms[2], 10, 64)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
		}
		s := metricscustom.Metric{Name: ms[1], Types: ms[0], Value: v}

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
			if v.Name == s.Name {
				if v.Types == "counter" {
					v.Value += s.Value
				}
				if v.Types == "gauge" {
					v.Value = s.Value
				}
				check = true
			}
		}
		if !check {
			mp.M[s.Name] = s
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
		defer r.Body.Close()
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	http.Error(w, "Only text/plain requests are allowed!", http.StatusMethodNotAllowed)
		// 	return
		// }
		m := strings.ReplaceAll(r.URL.Path, "/update/", "")
		ms := strings.Split(m, "/")
		fmt.Println(ms)
		if len(ms) != 3 {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		v, err := strconv.ParseUint(ms[2], 10, 64)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !checkType(ms[0]) {
			http.Error(w, "", http.StatusNotImplemented)
			return
		}
		s := metricscustom.Metric{Name: ms[1], Types: ms[0], Value: v}

		fmt.Printf("New metric: %s %v %s\n\n\n", s.Name, s.Value, s.Types)
		h.s.Update(&s)
		fmt.Println("In Storage: ")
		for k, v := range h.s.Mp.M {
			fmt.Printf("	Metric:  %s\n	    Name:%s\n	    Value:%v\n	    Type:%s\n\n", k, v.Name, v.Value, v.Types)
		}

	}
}

func checkType(s string) bool {
	var types []string = []string{"gauge", "counter"}
	for _, v := range types {
		if s == v {
			return true
		}
	}
	return false
}
