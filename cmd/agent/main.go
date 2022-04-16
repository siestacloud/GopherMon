package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/siestacloud/service-monitoring/internal/agent/metricscustom"
)

var (
	cms *metricscustom.CustomMetrics
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	//Задаем интервал сбора метрик
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(10) * time.Second
	go takeMetrics(pollInterval)
	go postMetrics(reportInterval)

	func() {
		for {
			sig := <-sigs
			switch sig {
			case os.Interrupt:
				HandleSignal(sig)
				os.Exit(0)
			case syscall.SIGTERM:
				HandleSignal(sig)
				os.Exit(0)
			case syscall.SIGINT:
				HandleSignal(sig)
				os.Exit(0)
			case syscall.SIGQUIT:
				HandleSignal(sig)
				os.Exit(0)
			default:
				// fmt.Println("Ignoring: ", sig)
			}
		}
	}()
}

func takeMetrics(pollInterval time.Duration) {
	//обьект обертка над runtime.MemStats
	var cmemstats metricscustom.CustomMemStats
	var intervalcounter int64
	for {
		select {
		// case <-time.After(postinterval):
		// 	fmt.Println("10 SECONDS!!!!")
		case <-time.After(pollInterval):
			intervalcounter++
			// Получаем все метрики
			cmemstats.ParseAllMetrics()
			// Берем только нужные
			cms = cmemstats.Convert(intervalcounter)
			// Just encode to json and print
			// b, _ := json.Marshal(cms)
			// fmt.Println(string(b))
		}
	}
}

func postMetrics(reportInterval time.Duration) {
	for {
		select {
		case <-time.After(reportInterval):
			url()
		}
	}
}

func url() {
	for _, v := range cms.G {
		url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", v.Types, v.Name, int(v.Value))

		// конструируем запрос
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Request %s\n\n", err)
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "text/plain")
		// конструируем клиент
		client := &http.Client{}
		// отправляем запрос
		_, err = client.Do(request)
		if err != nil {
			fmt.Printf("Do %s\n\n", err)
		}
		// defer resp.Body.Close()
		// fmt.Printf("%v", resp)
	}
	for _, v := range cms.C {
		url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", v.Types, v.Name, strconv.FormatInt(v.Value, 10))

		// конструируем запрос
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("req %s\n\n", err)
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "text/plain")
		// конструируем клиент
		client := &http.Client{}
		// отправляем запрос
		_, err = client.Do(request)
		if err != nil {
			fmt.Printf("do %s\n\n", err)
		}
		// defer resp.Body.Close()
		// fmt.Printf("%v", resp)
	}
}

func HandleSignal(signal os.Signal) {
	fmt.Println("HandleSignal() Received:", signal)
}
