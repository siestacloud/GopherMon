package utils

import (
	"fmt"
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

var Types = map[string]bool{
	"counter": true,
	"gauge":   true,
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) String() string {
	s := fmt.Sprintf("ID:%s\ntype:%s\ndelta:%d\nvalue:%f\n", m.ID, m.MType, *m.Delta, *m.Value)
	return s
}

func NewMetrics(id, mtype string) *Metrics {
	return &Metrics{
		ID:    id,
		MType: mtype,
		Value: new(float64),
		Delta: new(int64),
	}
}

type MetricsStorage map[string]Metrics

func NewMetricsStorage() MetricsStorage {
	m := NewMetrics("PollCount", "counter")
	return MetricsStorage{m.ID: *m}
}

func (m MetricsStorage) Poll() {
	g := "gauge"
	*m["PollCount"].Delta += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		m[mtrx.Type().Field(i).Name] = *NewMetrics(mtrx.Type().Field(i).Name, g)
		switch f.Kind() {
		case reflect.Int32, reflect.Int64:
			*m[mtrx.Type().Field(i).Name].Value = float64(f.Int())

		case reflect.Uint64, reflect.Uint32:
			*m[mtrx.Type().Field(i).Name].Value = float64(f.Uint())

		case reflect.Float32, reflect.Float64:
			*m[mtrx.Type().Field(i).Name].Value = f.Float()
		}

	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	m["RandomValue"] = *NewMetrics("RandomValue", g)
	*m["RandomValue"].Value = r.Float64()
	log.Println("Poll metrics")
}
