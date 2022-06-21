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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
type MetricsStorage map[string]Metrics

func NewMetricsStorage() MetricsStorage {
	m := MetricsStorage{}
	return m
}

func (m MetricsStorage) Poll() {
	*m["PollCount"].Delta += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		switch f.Kind() {
		case reflect.Uint64, reflect.Uint32:
			*m[mtrx.Type().Field(i).Name].Delta = int64(f.Uint())
		case reflect.Int32, reflect.Int64:
			*m[mtrx.Type().Field(i).Name].Delta = int64(f.Int())
		case reflect.Float32, reflect.Float64:
			*m[mtrx.Type().Field(i).Name].Value = float64(f.Float())

		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	*m["RandomValue"].Value = r.Float64()
	log.Println("Poll metrics")
}
