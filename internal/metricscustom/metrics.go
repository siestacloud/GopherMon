package metricscustom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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
	Name  string
	Types string
	Value float64
}

type Value struct{}

func NewMetric(t, n, v string) (*Metric, string) {
	if !checkType(t) {
		fmt.Println("check")
		return nil, "unknown metric type"
	}

	V, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, "incorrect value"
	}
	return &Metric{
		Name:  n,
		Types: t,
		Value: V,
	}, ""
}

//Convert Возвращает обьект с нужными метриками
func (m *MetricsPool) Convert(counter int64, cms *runtime.MemStats) {
	rand.Seed(time.Now().UTC().UnixNano())

	m.M["Alloc"] = Metric{Name: "Alloc", Value: float64(cms.Alloc), Types: "gauge"}
	m.M["BuckHashSys"] = Metric{Name: "BuckHashSys", Value: float64(cms.BuckHashSys), Types: "gauge"}
	m.M["Frees"] = Metric{Name: "Frees", Value: float64(cms.Frees), Types: "gauge"}
	m.M["GCCPUFraction"] = Metric{Name: "GCCPUFraction", Value: float64(cms.GCCPUFraction), Types: "gauge"}
	m.M["GCSys"] = Metric{Name: "GCSys", Value: float64(cms.GCSys), Types: "gauge"}
	m.M["HeapAlloc"] = Metric{Name: "HeapAlloc", Value: float64(cms.HeapAlloc), Types: "gauge"}
	m.M["HeapIdle"] = Metric{Name: "HeapIdle", Value: float64(cms.HeapIdle), Types: "gauge"}
	m.M["HeapInuse"] = Metric{Name: "HeapInuse", Value: float64(cms.HeapInuse), Types: "gauge"}
	m.M["HeapObjects"] = Metric{Name: "HeapObjects", Value: float64(cms.HeapObjects), Types: "gauge"}
	m.M["HeapReleased"] = Metric{Name: "HeapReleased", Value: float64(cms.HeapReleased), Types: "gauge"}
	m.M["HeapSys"] = Metric{Name: "HeapSys", Value: float64(cms.HeapSys), Types: "gauge"}
	m.M["LastGC"] = Metric{Name: "LastGC", Value: float64(cms.LastGC), Types: "gauge"}
	m.M["Lookups"] = Metric{Name: "Lookups", Value: float64(cms.Lookups), Types: "gauge"}
	m.M["MCacheInuse"] = Metric{Name: "MCacheInuse", Value: float64(cms.MCacheInuse), Types: "gauge"}
	m.M["MCacheSys"] = Metric{Name: "MCacheSys", Value: float64(cms.MCacheSys), Types: "gauge"}
	m.M["MSpanInuse"] = Metric{Name: "MSpanInuse", Value: float64(cms.MSpanInuse), Types: "gauge"}
	m.M["MSpanSys"] = Metric{Name: "MSpanSys", Value: float64(cms.MSpanSys), Types: "gauge"}
	m.M["Mallocs"] = Metric{Name: "Mallocs", Value: float64(cms.Mallocs), Types: "gauge"}
	m.M["NextGC"] = Metric{Name: "NextGC", Value: float64(cms.NextGC), Types: "gauge"}
	m.M["NumForcedGC"] = Metric{Name: "NumForcedGC", Value: float64(cms.NumForcedGC), Types: "gauge"}
	m.M["NumGC"] = Metric{Name: "NumGC", Value: float64(cms.NumGC), Types: "gauge"}
	m.M["OtherSys"] = Metric{Name: "OtherSys", Value: float64(cms.OtherSys), Types: "gauge"}
	m.M["PauseTotalNs"] = Metric{Name: "PauseTotalNs", Value: float64(cms.PauseTotalNs), Types: "gauge"}
	m.M["StackInuse"] = Metric{Name: "StackInuse", Value: float64(cms.StackInuse), Types: "gauge"}
	m.M["StackSys"] = Metric{Name: "StackSys", Value: float64(cms.StackSys), Types: "gauge"}
	m.M["Sys"] = Metric{Name: "Sys", Value: float64(cms.Sys), Types: "gauge"}
	m.M["TotalAlloc"] = Metric{Name: "TotalAlloc", Value: float64(cms.TotalAlloc), Types: "gauge"}
	m.M["RandomValue"] = Metric{Name: "RandomValue", Value: float64(rand.Intn(100)), Types: "gauge"}
	m.M["PollCount"] = Metric{Name: "PollCount", Value: float64(counter), Types: "counter"}
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
