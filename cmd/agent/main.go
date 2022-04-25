package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/siestacloud/service-monitoring/internal/metricscustom"
)

var (
	cms runtime.MemStats
	cmp = metricscustom.NewMetricsPool()
)

func main() {

	ctx, cansel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cansel()
	//Задаем интервал сбора метрик
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(10) * time.Second
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
			cmp.AddMetrics(ic, &cms)
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

func url() {
	var buf bytes.Buffer
	for _, metric := range cmp.M {

		err := metric.MarshalMetricsinJSON(&buf)
		if err != nil {
			log.Fatal(err)
		}
		// конструируем запрос
		request, err := http.NewRequest("POST", "http://localhost:8080/update/", &buf)
		if err != nil {
			fmt.Printf("Request %s\n\n", err)
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "text/plain")
		// конструируем клиент
		client := &http.Client{}
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
			continue
		}

	}

}
