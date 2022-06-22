package utils

import (
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type Gauge float64
type Counter int64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

var Counters = []string{
	"PollCount",
}

type JsonMetrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Metrics struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
	Jsons    []JsonMetrics
}

func NewMetricsStorage() *Metrics {
	m := &Metrics{}
	m.Init()
	return m
}

func (m *Metrics) Init() {
	m.Gauges = map[string]Gauge{}
	m.Counters = map[string]Counter{}
	m.Jsons = make([]JsonMetrics, 100)
}

func (m *Metrics) Poll() {
	m.Counters["PollCount"] += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		switch f.Kind() {
		case reflect.Uint64, reflect.Uint32:
			m.Gauges[mtrx.Type().Field(i).Name] = Gauge(f.Uint())
		case reflect.Float32, reflect.Float64:
			m.Gauges[mtrx.Type().Field(i).Name] = Gauge(f.Float())
		case reflect.Int32, reflect.Int64:
			m.Gauges[mtrx.Type().Field(i).Name] = Gauge(f.Int())
		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	m.Gauges["RandomValue"] = Gauge(r.Float64())
	log.Println("Poll metrics")
}
