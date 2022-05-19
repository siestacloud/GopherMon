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
	GaugesRuntime map[string]gauge
	GaugesMy      map[string]gauge
	Counters      map[string]counter
}

func (m *Metrics) Init() {
	m.GaugesRuntime = map[string]gauge{}
	m.GaugesMy = map[string]gauge{
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
		switch f.Kind() {
		case reflect.Uint64, reflect.Uint32:
			m.GaugesRuntime[mtrx.Type().Field(i).Name] = gauge(f.Uint())
		case reflect.Float32, reflect.Float64:
			m.GaugesRuntime[mtrx.Type().Field(i).Name] = gauge(f.Float())
		case reflect.Int32, reflect.Int64:
			m.GaugesRuntime[mtrx.Type().Field(i).Name] = gauge(f.Int())
		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	m.GaugesMy["RandomValue"] = gauge(r.Float64())
}
