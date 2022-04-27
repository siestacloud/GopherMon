package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/siestacloud/service-monitoring/internal/mtrx"
)

var (
	cms runtime.MemStats
	mp  *mtrx.MetricsPool
	err error
)

func main() {

	ctx, cansel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cansel()
	//Задаем интервал сбора метрик
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(20) * time.Second
	go takeMetrics(ctx, pollInterval)
	go postMetrics(ctx, reportInterval)

	<-ctx.Done()
	time.Sleep(time.Second)
	os.Exit(0)
}

func takeMetrics(ctx context.Context, pollInterval time.Duration) {

	var ic int64
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(pollInterval):
			ic++
			// Получаем все метрики
			runtime.ReadMemStats(&cms)
			// Берем только нужные
			mp, err = mtrxMotion(ic, &cms)
			if err != nil {
				log.Println(err)
			}
			// fmt.Printf("%v\n%v\n\n", cms.HeapReleased, cmp.M["HeapReleased"])
			// Just encode to json and print
			// b, _ := json.Marshal(cms)
			// fmt.Println(string(b))
		}
	}
}

func postMetrics(ctx context.Context, reportInterval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(reportInterval):
			url()
		}
	}
}

//
func url() {
	// конструируем клиент
	client := &http.Client{}
	for _, metric := range mp.M {
		var buf bytes.Buffer
		err := metric.MarshalMetricsinJSON(&buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("JSON request agent", buf.String())
		// конструируем запрос
		request, err := http.NewRequest("POST", "http://127.0.0.1:8080/update/", &buf)
		if err != nil {
			fmt.Printf("Request %s\n\n", err)
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "application/json")
		// отправляем запрос
		resp, err := client.Do(request)
		if err != nil {
			fmt.Printf("Do %s\n\n", err)
			continue
		}
		if resp != nil {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("%v\n", string(b))
			resp.Body.Close()
			continue
		}

	}

}

//Формирую метрики по заданию, заполняю общий пул метрик
func mtrxMotion(c int64, cms *runtime.MemStats) (*mtrx.MetricsPool, error) {
	mtrxPool := mtrx.NewMetricsPool()

	//Создаю метрику PollCount
	pollCount := mtrx.NewMetric()
	if err := pollCount.SetID("PollCount"); err != nil {
		return nil, err
	}
	if err := pollCount.SetType("counter"); err != nil {
		return nil, err
	}
	if err := pollCount.SetValue(c); err != nil {
		return nil, err
	}
	if !mtrxPool.Add(pollCount.ID, *pollCount) {
		return nil, errors.New("unable add PollCount mtrx into MetricsPool: " + pollCount.GetID() + pollCount.GetType())
	}
	//Создаю метрику RandomValue
	rand.Seed(time.Now().UTC().UnixNano())

	randomValue := mtrx.NewMetric()
	if err := randomValue.SetID("RandomValue"); err != nil {
		return nil, err
	}
	if err := randomValue.SetType("gauge"); err != nil {
		return nil, err
	}
	if err := randomValue.SetValue(rand.Float64()); err != nil {
		return nil, err
	}

	if !mtrxPool.Add(randomValue.ID, *randomValue) {
		return nil, errors.New("unable add PollCount mtrx into MetricsPool: " + randomValue.GetID() + randomValue.GetType())
	}
	//Создаю метрики runtime
	val := reflect.ValueOf(cms).Elem()
	for i := 0; i < val.NumField(); i++ {
		id := val.Type().Field(i).Name
		v := fmt.Sprint(val.FieldByName(val.Type().Field(i).Name))

		m := mtrx.NewMetric()
		if err := m.SetID(id); err != nil {
			return nil, err
		}
		if err := m.SetType("gauge"); err != nil {
			return nil, err
		}
		if err := m.SetValue(v); err != nil {

			continue
		}
		if !mtrxPool.Add(m.ID, *m) {
			return nil, errors.New("unable add runtime mtrx into MetricsPool: " + m.GetID() + "  " + m.GetType())
		}
	}
	return mtrxPool, nil
}
