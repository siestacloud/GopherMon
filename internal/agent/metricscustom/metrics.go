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
	cms.G["Alloc"] = Gaudge{Name: "Alloc", Value: float64(c.Alloc), Types: "Gaudge"}
	cms.G["BuckHashSys"] = Gaudge{Name: "BuckHashSys", Value: float64(c.BuckHashSys), Types: "Gaudge"}
	cms.G["Frees"] = Gaudge{Name: "Frees", Value: float64(c.Frees), Types: "Gaudge"}
	cms.G["GCCPUFraction"] = Gaudge{Name: "GCCPUFraction", Value: float64(c.GCCPUFraction), Types: "Gaudge"}
	cms.G["GCSys"] = Gaudge{Name: "GCSys", Value: float64(c.GCSys), Types: "Gaudge"}
	cms.G["HeapAlloc"] = Gaudge{Name: "HeapAlloc", Value: float64(c.HeapAlloc), Types: "Gaudge"}
	cms.G["HeapIdle"] = Gaudge{Name: "HeapIdle", Value: float64(c.HeapIdle), Types: "Gaudge"}
	cms.G["HeapInuse"] = Gaudge{Name: "HeapInuse", Value: float64(c.HeapInuse), Types: "Gaudge"}
	cms.G["HeapObjects"] = Gaudge{Name: "HeapObjects", Value: float64(c.HeapObjects), Types: "Gaudge"}
	cms.G["HeapReleased"] = Gaudge{Name: "HeapReleased", Value: float64(c.HeapReleased), Types: "Gaudge"}
	cms.G["HeapSys"] = Gaudge{Name: "HeapSys", Value: float64(c.HeapSys), Types: "Gaudge"}
	cms.G["LastGC"] = Gaudge{Name: "LastGC", Value: float64(c.LastGC), Types: "Gaudge"}
	cms.G["Lookups"] = Gaudge{Name: "Lookups", Value: float64(c.Lookups), Types: "Gaudge"}
	cms.G["MCacheInuse"] = Gaudge{Name: "MCacheInuse", Value: float64(c.MCacheInuse), Types: "Gaudge"}
	cms.G["MCacheSys"] = Gaudge{Name: "MCacheSys", Value: float64(c.MCacheSys), Types: "Gaudge"}
	cms.G["MSpanInuse"] = Gaudge{Name: "MSpanInuse", Value: float64(c.MSpanInuse), Types: "Gaudge"}
	cms.G["MSpanSys"] = Gaudge{Name: "MSpanSys", Value: float64(c.MSpanSys), Types: "Gaudge"}
	cms.G["Mallocs"] = Gaudge{Name: "Mallocs", Value: float64(c.Mallocs), Types: "Gaudge"}
	cms.G["NextGC"] = Gaudge{Name: "NextGC", Value: float64(c.NextGC), Types: "Gaudge"}
	cms.G["NumForcedGC"] = Gaudge{Name: "NumForcedGC", Value: float64(c.NumForcedGC), Types: "Gaudge"}
	cms.G["NumGC"] = Gaudge{Name: "NumGC", Value: float64(c.NumGC), Types: "Gaudge"}
	cms.G["OtherSys"] = Gaudge{Name: "OtherSys", Value: float64(c.OtherSys), Types: "Gaudge"}
	cms.G["PauseTotalNs"] = Gaudge{Name: "PauseTotalNs", Value: float64(c.PauseTotalNs), Types: "Gaudge"}
	cms.G["StackInuse"] = Gaudge{Name: "StackInuse", Value: float64(c.StackInuse), Types: "Gaudge"}
	cms.G["StackSys"] = Gaudge{Name: "StackSys", Value: float64(c.StackSys), Types: "Gaudge"}
	cms.G["Sys"] = Gaudge{Name: "Sys", Value: float64(c.Sys), Types: "Gaudge"}
	cms.G["TotalAlloc"] = Gaudge{Name: "TotalAlloc", Value: float64(c.TotalAlloc), Types: "Gaudge"}
	cms.G["RandomValue"] = Gaudge{Name: "RandomValue", Value: float64(float64(rand.Intn(100))), Types: "Gaudge"}
	cms.C["PollCount"] = Counter{Name: "TotalAlloc", Value: int64(c.TotalAlloc), Types: "Counter"}

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
