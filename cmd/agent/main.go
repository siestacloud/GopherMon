package main

import (
	"encoding/json"
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

func NewMonitor(duration int) {
	//обьект обертка над runtime.MemStats
	var cmemstats metricscustom.CustomMemStats
	var intervalcounter int64
	//Задаем интервал сбора метрик
	var interval = time.Duration(duration) * time.Second
	var postinterval = time.Duration(4) * time.Second
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go postMetrics(postinterval)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case os.Interrupt:
				HandleSignal(sig)
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
				fmt.Println("Ignoring: ", sig)

			}
		}
	}()
	for {
		select {
		// case <-time.After(postinterval):
		// 	fmt.Println("10 SECONDS!!!!")
		case <-time.After(interval):
			intervalcounter++
			// Получаем все метрики
			cmemstats.ParseAllMetrics()
			// Берем только нужные
			cms = cmemstats.Convert(intervalcounter)

			// Just encode to json and print
			b, _ := json.Marshal(cms)
			fmt.Println(string(b))
		}
	}
}

func postMetrics(interval time.Duration) {

	for {
		select {
		case <-time.After(interval):
			url()
		}
	}
}
func main() {

	NewMonitor(2)

}

func url() {
	for _, v := range cms.G {
		url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", v.Types, v.Name, v.Value)

		// конструируем запрос
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			// обработаем ошибку
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "text/plain")
		// конструируем клиент
		client := &http.Client{}
		// отправляем запрос
		resp, err := client.Do(request)
		if err != nil {
			// log.Fatal(err)
		}
		fmt.Printf("%v", resp)
	}
	for _, v := range cms.C {
		url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", v.Types, v.Name, strconv.FormatInt(v.Value, 10))

		// конструируем запрос
		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			// обработаем ошибку
		}
		// устанавливаем заголовки
		request.Header.Add("Content-Type", "text/plain")
		// конструируем клиент
		client := &http.Client{}
		// отправляем запрос
		resp, err := client.Do(request)
		if err != nil {
			// обработаем ошибку
		}
		fmt.Printf("%v", resp)
	}
}

func HandleSignal(signal os.Signal) {
	fmt.Println("HandleSignal() Received:", signal)
}
