package metricscustom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

//CustomMetrics Необходимые метрики
type MetricsPool struct {
	M map[string]Metric
}

//NewMetricsPool обертка считывает все метрики
func NewMetricsPool() *MetricsPool {
	return &MetricsPool{
		M: map[string]Metric{},
	}
}

//Metric .
type Metric struct {
	ID    string  `json:"id"`    // имя метрики
	MType string  `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta"` // значение метрики в случае передачи counter
	Value float64 `json:"value"` // значение метрики в случае передачи gauge
}

func NewMetric(t, n, v string) (*Metric, string) {
	switch t {
	case "counter":
		V, err := strconv.ParseInt(v, 10, 64)
		if err != nil {

			return nil, "incorrect value"
		}
		return &Metric{
			ID:    n,
			MType: t,
			Delta: V,
		}, ""

	case "gauge":
		V, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, "incorrect value"
		}
		// fmt.Println("V:  ", n, V)
		return &Metric{
			ID:    n,
			MType: t,
			Value: V,
		}, ""
	}
	return nil, "unknown metric type"
}

func (m *Metric) Check() string {
	switch m.MType {
	case "counter":
		if m.Value != 0 {
			return "incorrect value"
		}
		if m.Delta == 0 {
			return "empty delta"
		}
		return ""
	case "gauge":
		if m.Delta != 0 {
			return "incorrect value"
		}
		if m.Value == 0 {
			return "empty value"
		}
		return ""
	}
	return "unknown metric type"
}

func (m *Metric) Find(mp *MetricsPool) string {
	for _, v := range mp.M {
		if v.ID == m.ID {
			if v.MType == m.MType {
				switch m.MType {
				case "counter":
					m.Delta = v.Delta
					return ""
				case "gauge":
					m.Value = v.Value
					return ""
				}
				return "unknown metric type"
			}
			return "metrics found but type incorrect"
		}
	}
	return "metric not found"
}

func (m *MetricsPool) AddMetrics(counter int64, cms *runtime.MemStats) {

	rand.Seed(time.Now().UTC().UnixNano())
	m.M["PollCount"] = Metric{ID: "PollCount", Delta: counter, MType: "counter"}
	m.M["RandomValue"] = Metric{ID: "RandomValue", Value: rand.Float64(), MType: "gauge"}

	// val := reflect.ValueOf(cms).Elem()
	// n := val.Type().Field(0).Name
	// v := fmt.Sprint(val.FieldByName(val.Type().Field(0).Name))
	// M, _ := NewMetric("gauge", n, v)
	// m.M[M.ID] = *M

	val := reflect.ValueOf(cms).Elem()
	for i := 0; i < val.NumField(); i++ {
		t := "gauge"
		n := val.Type().Field(i).Name
		v := fmt.Sprint(val.FieldByName(val.Type().Field(i).Name))
		M, status := NewMetric(t, n, v)
		if status != "" {

			continue
		}
		m.M[val.Type().Field(i).Name] = *M
	}
	fmt.Println(cms.LastGC)
	fmt.Println(m.M["LastGC"])
}

// WriteMetricJSON сериализует структуру Metric в JSON, и если всё отрабатывает
// успешно, то вызывается соответствующий метод Write() из io.Writer.
func (m *Metric) MarshalMetricsinJSON(w io.Writer) error {
	js, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(js)
	return err
}

func (m *Metric) UnmarshalMetricJSON(r io.Reader) error {

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return err
	}
	fmt.Println("ReadMetricJSON", m)
	return nil
}
