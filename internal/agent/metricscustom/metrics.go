package metricscustom

import (
	"math/rand"
	"runtime"
	"time"
)

//CustomMemStats Обертка над runtime.MemStats
type CustomMemStats struct {
	runtime.MemStats
}

//CustomMetrics Необходимые метрики
type CustomMetrics struct {
	G map[string]Gaudge
	C map[string]Counter
}

type Gaudge struct {
	Name  string
	Types string
	Value float64
}

type Counter struct {
	Name  string
	Types string
	Value int64
}

//newCustomMetrics создание обьекта и формирование необходимых метрик
func newCustomMetrics(c *CustomMemStats, i int64) *CustomMetrics {
	rand.Seed(time.Now().UnixNano())
	var cms CustomMetrics
	cms.G = map[string]Gaudge{}
	cms.C = map[string]Counter{}
	cms.G["Alloc"] = Gaudge{Name: "alloc", Value: float64(c.Alloc), Types: "gaudge"}
	cms.G["BuckHashSys"] = Gaudge{Name: "buckhashsys", Value: float64(c.BuckHashSys), Types: "gaudge"}
	cms.G["Frees"] = Gaudge{Name: "frees", Value: float64(c.Frees), Types: "gaudge"}
	cms.G["GCCPUFraction"] = Gaudge{Name: "gccpufraction", Value: float64(c.GCCPUFraction), Types: "gaudge"}
	cms.G["GCSys"] = Gaudge{Name: "gcsys", Value: float64(c.GCSys), Types: "gaudge"}
	cms.G["HeapAlloc"] = Gaudge{Name: "heapalloc", Value: float64(c.HeapAlloc), Types: "gaudge"}
	cms.G["HeapIdle"] = Gaudge{Name: "heapidle", Value: float64(c.HeapIdle), Types: "gaudge"}
	cms.G["HeapInuse"] = Gaudge{Name: "heapinuse", Value: float64(c.HeapInuse), Types: "gaudge"}
	cms.G["HeapObjects"] = Gaudge{Name: "heapobjects", Value: float64(c.HeapObjects), Types: "gaudge"}
	cms.G["HeapReleased"] = Gaudge{Name: "heapreleased", Value: float64(c.HeapReleased), Types: "gaudge"}
	cms.G["HeapSys"] = Gaudge{Name: "heapsys", Value: float64(c.HeapSys), Types: "gaudge"}
	cms.G["LastGC"] = Gaudge{Name: "lastgc", Value: float64(c.LastGC), Types: "gaudge"}
	cms.G["Lookups"] = Gaudge{Name: "lookups", Value: float64(c.Lookups), Types: "gaudge"}
	cms.G["MCacheInuse"] = Gaudge{Name: "mcacheinuse", Value: float64(c.MCacheInuse), Types: "gaudge"}
	cms.G["MCacheSys"] = Gaudge{Name: "mcachesys", Value: float64(c.MCacheSys), Types: "gaudge"}
	cms.G["MSpanInuse"] = Gaudge{Name: "mspaninuse", Value: float64(c.MSpanInuse), Types: "gaudge"}
	cms.G["MSpanSys"] = Gaudge{Name: "mspansys", Value: float64(c.MSpanSys), Types: "gaudge"}
	cms.G["Mallocs"] = Gaudge{Name: "mallocs", Value: float64(c.Mallocs), Types: "gaudge"}
	cms.G["NextGC"] = Gaudge{Name: "nextgc", Value: float64(c.NextGC), Types: "gaudge"}
	cms.G["NumForcedGC"] = Gaudge{Name: "numforcedgc", Value: float64(c.NumForcedGC), Types: "gaudge"}
	cms.G["NumGC"] = Gaudge{Name: "numgc", Value: float64(c.NumGC), Types: "gaudge"}
	cms.G["OtherSys"] = Gaudge{Name: "othersys", Value: float64(c.OtherSys), Types: "gaudge"}
	cms.G["PauseTotalNs"] = Gaudge{Name: "pausetotalns", Value: float64(c.PauseTotalNs), Types: "gaudge"}
	cms.G["StackInuse"] = Gaudge{Name: "stackinuse", Value: float64(c.StackInuse), Types: "gaudge"}
	cms.G["StackSys"] = Gaudge{Name: "stacksys", Value: float64(c.StackSys), Types: "gaudge"}
	cms.G["Sys"] = Gaudge{Name: "sys", Value: float64(c.Sys), Types: "gaudge"}
	cms.G["TotalAlloc"] = Gaudge{Name: "totalalloc", Value: float64(c.TotalAlloc), Types: "gaudge"}
	cms.G["RandomValue"] = Gaudge{Name: "randomvalue", Value: float64(rand.Intn(100)), Types: "gaudge"}
	cms.C["PollCount"] = Counter{Name: "totalalloc", Value: int64(c.TotalAlloc), Types: "counter"}

	return &cms
}

//ParseAllMetrics метод обертка считывает все метрики
func (c *CustomMemStats) ParseAllMetrics() {
	runtime.ReadMemStats(&c.MemStats)
}

//Convert Возвращает обьект с нужными метриками
func (c *CustomMemStats) Convert(counter int64) *CustomMetrics {
	cms := newCustomMetrics(c, counter)
	// Number of goroutines

	return cms
}
