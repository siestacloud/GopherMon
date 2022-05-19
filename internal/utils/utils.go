package utils

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
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
		"PollCount": 0,
	}
}

func (m *Metrics) Poll() {
	m.Counters["PollCount"] += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		switch {
		case f.CanUint():
			m.Gauges_runtime[mtrx.Type().Field(i).Name] = gauge(f.Uint())
		case f.CanFloat():
			m.Gauges_runtime[mtrx.Type().Field(i).Name] = gauge(f.Float())
		case f.CanInt():
			m.Gauges_runtime[mtrx.Type().Field(i).Name] = gauge(f.Int())
		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	m.Gauges_my["RandomValue"] = gauge(r.Float64())
}
