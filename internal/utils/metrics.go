package mymetrics

import (
	"math/rand"
	"reflect"
	"runtime"
)

type gauge float64
type counter int64

type Metrics struct {
	Gauges_runtime map[string]gauge
	Gauges_my      map[string]gauge
	Counters       map[string]counter
}

func (m *Metrics) Init() {
	m.Gauges_runtime = map[string]gauge{}
	m.Gauges_my = map[string]gauge{
		"RandomValue": 0,
	}
	m.Counters = map[string]counter{
		"PoolCount": 0,
	}
}

func (m *Metrics) Poll() {
	m.Counters["PollCount"] += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics)
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		switch f.Kind() {
		case reflect.Float64:
			m.Gauges_runtime[f.Elem().Type().Name()] = gauge(f.Float())
		case reflect.Uint:
			m.Gauges_runtime[f.Elem().Type().Name()] = gauge(f.Uint())
		case reflect.Int:
			m.Gauges_runtime[f.Elem().Type().Name()] = gauge(f.Int())
		}

	}
	m.Gauges_my["RandomValue"] = gauge(rand.Float64())
}
