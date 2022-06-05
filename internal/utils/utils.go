package utils

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

type Gauge float64
type Counter int64

var Counters = []string{
	"PollCount",
}

type Metrics struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func (m *Metrics) Init() {
	m.Gauges = map[string]Gauge{}
	m.Counters = map[string]Counter{}
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
}