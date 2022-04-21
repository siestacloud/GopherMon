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
	Value uint64
}

func NewMetric(t, n, val string) (*Metric, error) {
	v, err := strconv.ParseUint(val, 0, 64)
	if err != nil {
		return nil, err
	}
	return &Metric{
		Name:  n,
		Types: t,
		Value: v,
	}, nil
}

//Convert Возвращает обьект с нужными метриками
func (m *MetricsPool) Convert(counter int64, cms *runtime.MemStats) {
	rand.Seed(time.Now().UTC().UnixNano())

	m.M["Alloc"] = Metric{Name: "alloc", Value: cms.Alloc, Types: "gauge"}
	m.M["BuckHashSys"] = Metric{Name: "buckhashsys", Value: cms.BuckHashSys, Types: "gauge"}
	m.M["Frees"] = Metric{Name: "frees", Value: cms.Frees, Types: "gauge"}
	m.M["GCCPUFraction"] = Metric{Name: "gccpufraction", Value: uint64(cms.GCCPUFraction), Types: "gauge"}
	m.M["GCSys"] = Metric{Name: "gcsys", Value: cms.GCSys, Types: "gauge"}
	m.M["HeapAlloc"] = Metric{Name: "heapalloc", Value: cms.HeapAlloc, Types: "gauge"}
	m.M["HeapIdle"] = Metric{Name: "heapidle", Value: cms.HeapIdle, Types: "gauge"}
	m.M["HeapInuse"] = Metric{Name: "heapinuse", Value: cms.HeapInuse, Types: "gauge"}
	m.M["HeapObjects"] = Metric{Name: "heapobjects", Value: cms.HeapObjects, Types: "gauge"}
	m.M["HeapReleased"] = Metric{Name: "heapreleased", Value: cms.HeapReleased, Types: "gauge"}
	m.M["HeapSys"] = Metric{Name: "heapsys", Value: cms.HeapSys, Types: "gauge"}
	m.M["LastGC"] = Metric{Name: "lastgc", Value: cms.LastGC, Types: "gauge"}
	m.M["Lookups"] = Metric{Name: "lookups", Value: cms.Lookups, Types: "gauge"}
	m.M["MCacheInuse"] = Metric{Name: "mcacheinuse", Value: cms.MCacheInuse, Types: "gauge"}
	m.M["MCacheSys"] = Metric{Name: "mcachesys", Value: cms.MCacheSys, Types: "gauge"}
	m.M["MSpanInuse"] = Metric{Name: "mspaninuse", Value: cms.MSpanInuse, Types: "gauge"}
	m.M["MSpanSys"] = Metric{Name: "mspansys", Value: cms.MSpanSys, Types: "gauge"}
	m.M["Mallocs"] = Metric{Name: "mallocs", Value: cms.Mallocs, Types: "gauge"}
	m.M["NextGC"] = Metric{Name: "nextgc", Value: cms.NextGC, Types: "gauge"}
	m.M["NumForcedGC"] = Metric{Name: "numforcedgc", Value: uint64(cms.NumForcedGC), Types: "gauge"}
	m.M["NumGC"] = Metric{Name: "numgc", Value: uint64(cms.NumGC), Types: "gauge"}
	m.M["OtherSys"] = Metric{Name: "othersys", Value: cms.OtherSys, Types: "gauge"}
	m.M["PauseTotalNs"] = Metric{Name: "pausetotalns", Value: cms.PauseTotalNs, Types: "gauge"}
	m.M["StackInuse"] = Metric{Name: "stackinuse", Value: cms.StackInuse, Types: "gauge"}
	m.M["StackSys"] = Metric{Name: "stacksys", Value: cms.StackSys, Types: "gauge"}
	m.M["Sys"] = Metric{Name: "sys", Value: cms.Sys, Types: "gauge"}
	m.M["TotalAlloc"] = Metric{Name: "totalalloc", Value: cms.TotalAlloc, Types: "gauge"}
	m.M["RandomValue"] = Metric{Name: "randomvalue", Value: uint64(rand.Intn(100)), Types: "gauge"}
	m.M["PollCount"] = Metric{Name: "totalalloc", Value: cms.TotalAlloc, Types: "counter"}
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
