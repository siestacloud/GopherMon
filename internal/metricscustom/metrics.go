package metricscustom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Value struct{}

func NewMetric(t, n, v string) (*Metric, string) {
	fmt.Println("TYPE:  ", t)
	if !checkType(t) {
		fmt.Println("check")
		return nil, "unknown metric type"
	}

	V, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, "incorrect value"
	}
	return &Metric{
		ID:    n,
		MType: t,
		Value: &V,
	}, ""
}

func (m *MetricsPool) AddMetrics(counter int64, cms *runtime.MemStats) {
	rand.Seed(time.Now().UTC().UnixNano())
	m.M["PollCount"] = Metric{ID: "PollCount", Delta: &counter, MType: "counter"}

	val := reflect.ValueOf(cms).Elem()
	for i := 0; i < val.NumField(); i++ {

		t := "gauge"
		n := val.Type().Field(i).Name
		v := fmt.Sprint(val.FieldByName(val.Type().Field(i).Name))

		M, status := NewMetric(t, n, v)
		if status != "" {
			log.Println(status)
		}
		m.M[val.Type().Field(i).Name] = *M
	}
}

// WriteMetricJSON сериализует структуру Metric в JSON, и если всё отрабатывает
// успешно, то вызывается соответствующий метод Write() из io.Writer.
func (m *MetricsPool) WriteMetricsJSON(w io.Writer) error {
	js, err := json.MarshalIndent(m.M, "", "	")
	if err != nil {
		return err
	}

	_, err = w.Write(js)
	return err
}

func (m *MetricsPool) ReadMetricsJSON(r io.Reader) error {

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &m.M)
	if err != nil {
		return err
	}
	fmt.Println("ReadMetricJSON", m)
	return nil
}

func checkType(s string) bool {
	var types = []string{"gauge", "counter"}
	for _, v := range types {
		if s == v {
			return true
		}
	}
	return false
}
