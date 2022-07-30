package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/sirupsen/logrus"
)

// MyApiError — описание ошибки при неверном запросе
type (
	Config struct {
		Address        string        `env:"ADDRESS"`
		PollInterval   time.Duration `env:"POLL_INTERVAL"`
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
		Key            string        `env:"KEY"`
	}

	APIError struct {
		Code      int       `json:"code"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
	}
)

var (
	cfg = Config{}
	cms runtime.MemStats
	mp  *core.MetricsPool
	err error
)

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Address for server. Possible values: localhost:8080")
	flag.DurationVar(&cfg.PollInterval, "p", 2000000000, "Poll interval. Possible values: 1s 12s 1m")
	flag.DurationVar(&cfg.ReportInterval, "r", 10000000000, "Report interval. Possible values: 1s 12s 1m")
	flag.StringVar(&cfg.Key, "k", "", "key for data sign. Possible values: 123qwe123")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfgjson, _ := json.MarshalIndent(cfg, "  ", " ")
	logrus.Info(string(cfgjson))
}
func main() {

	ctx, cansel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cansel()

	go takeMetrics(ctx, cfg.PollInterval)
	go postMetrics(ctx, cfg.ReportInterval)

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
	logger := logrus.New()
	logger.Out = ioutil.Discard

	for _, metric := range mp.M {
		fmt.Println(metric, "   ", metric.Value)
		client := resty.New().SetRetryCount(2).SetLogger(logger).
			SetRetryWaitTime(1 * time.Second).
			SetRetryMaxWaitTime(2 * time.Second)
		var responseErr APIError
		_, err := client.R().
			SetError(&responseErr).SetDoNotParseResponse(false).
			SetBody(metric).
			Post("http://" + cfg.Address + "/update/")
		if err != nil {
			// fmt.Println("resp err:  ", responseErr)
			log.Println("AGENT resp err:: ", err)
		}

		// fmt.Println(metric, "   ", metric.Value)
		// // var buf bytes.Buffer
		// // err := metric.MarshalMetricsinJSON(&buf)
		// // if err != nil {
		// // 	log.Fatal(err)
		// // }
		// // fmt.Println("JSON request agent", buf.String())
		// // // конструируем запрос
		// body, err := json.Marshal(metric)
		// if err != nil {
		// 	fmt.Println("json marshal err: ", err)
		// 	continue
		// }
		// fmt.Println("SHOW METRIC", string(body))
		// // конструируем запрос
		// request, err := http.NewRequest("POST", "http://127.0.0.1:8080/update/", bytes.NewBuffer(body))
		// if err != nil {
		// 	fmt.Printf("Request %s\n\n", err)
		// }
		// // устанавливаем заголовки
		// request.Header.Add("Content-Type", "text/plain")
		// // Close the connection
		// request.Close = true
		// // конструируем клиент
		// client := &http.Client{}
		// // отправляем запрос
		// resp, err := client.Do(request)
		// if err != nil {
		// 	fmt.Printf("Do %s\n\n", err)
		// 	continue
		// }
		// resp.Body.Close()
		// // resp, err := http.Post("http://127.0.0.1:8080/update/", "application/json", bytes.NewBuffer(body))
		// // if err != nil {
		// // 	fmt.Println("DO POST err: ", err)
		// // 	break
		// // }
		// fmt.Printf("Status: %s  \n", resp.Status)
		// continue
	}

}

//Формирую метрики по заданию, заполняю общий пул метрик
func mtrxMotion(c int64, cms *runtime.MemStats) (*core.MetricsPool, error) {
	mtrxPool := core.NewMetricsPool()

	//Создаю метрику PollCount
	pollCount := core.NewMetric()
	if err := pollCount.SetID("PollCount"); err != nil {
		return nil, err
	}
	if err := pollCount.SetType("counter"); err != nil {
		return nil, err
	}
	if err := pollCount.SetValue(c); err != nil {
		return nil, err
	}
	if err := pollCount.SetHash(cfg.Key); err != nil {
		return nil, err
	}
	if err := mtrxPool.Create(pollCount.ID, *pollCount); err != nil {
		return nil, errors.New("unable add PollCount mtrx into MetricsPool: " + pollCount.GetID() + pollCount.GetType())
	}
	//Создаю метрику RandomValue
	rand.Seed(time.Now().UTC().UnixNano())

	randomValue := core.NewMetric()
	if err := randomValue.SetID("RandomValue"); err != nil {
		return nil, err
	}
	if err := randomValue.SetType("gauge"); err != nil {
		return nil, err
	}
	if err := randomValue.SetValue(rand.Float64()); err != nil {
		return nil, err
	}
	if err := randomValue.SetHash(cfg.Key); err != nil {
		return nil, err
	}
	if err := mtrxPool.Create(randomValue.ID, *randomValue); err != nil {
		return nil, errors.New("unable add randomValue mtrx into MetricsPool: " + randomValue.GetID() + randomValue.GetType())
	}

	//Создаю метрики из пакета runtime / у cms тип runtime.MemStats
	val := reflect.ValueOf(cms).Elem()
	//итерируюсь по всем полям cms, runtime.MemStats

	for i := 0; i < val.NumField(); i++ {
		id := val.Type().Field(i).Name                             //достаю имя поля
		v := fmt.Sprint(val.FieldByName(val.Type().Field(i).Name)) // значение в этом поле

		m := core.NewMetric()               // создаю свою метрику
		if err := m.SetID(id); err != nil { // у обьекта метрики определены методы, через которые заполняются поля имя метрики значение и тип
			return nil, err
		}
		if err := m.SetType("gauge"); err != nil {
			return nil, err
		}
		if err := m.SetValue(v); err != nil {
			continue
		}
		if err := m.SetHash(cfg.Key); err != nil {
			return nil, err
		}
		if err := mtrxPool.Create(m.ID, *m); err != nil { // Метрика добавляется в общий пул (мапку)
			return nil, errors.New("unable add runtime mtrx into MetricsPool: " + m.GetID() + "  " + m.GetType())
		}
	}
	memory, _ := mem.VirtualMemory()

	totalMemory := core.NewMetric()
	if err := totalMemory.SetID("TotalMemory"); err != nil {
		return nil, err
	}
	if err := totalMemory.SetType("gauge"); err != nil {
		return nil, err
	}
	if err := totalMemory.SetValue(float64(memory.Total)); err != nil {
		return nil, err
	}
	if err := totalMemory.SetHash(cfg.Key); err != nil {
		return nil, err
	}
	if err := mtrxPool.Create(totalMemory.ID, *totalMemory); err != nil {
		return nil, errors.New("unable add totalMemory mtrx into MetricsPool: " + totalMemory.GetID() + totalMemory.GetType())
	}

	freeMemory := core.NewMetric()
	if err := freeMemory.SetID("FreeMemory"); err != nil {
		return nil, err
	}
	if err := freeMemory.SetType("gauge"); err != nil {
		return nil, err
	}
	if err := freeMemory.SetValue(float64(memory.Free)); err != nil {
		return nil, err
	}
	if err := freeMemory.SetHash(cfg.Key); err != nil {
		return nil, err
	}
	if err := mtrxPool.Create(freeMemory.ID, *freeMemory); err != nil {
		return nil, errors.New("unable add freeMemory mtrx into MetricsPool: " + freeMemory.GetID() + freeMemory.GetType())
	}

	cpuPercent, _ := cpu.Percent(1000000, true)

	for i, cpu := range cpuPercent {

		CPUutilization := core.NewMetric()
		if err := CPUutilization.SetID("CPUutilization" + strconv.Itoa(i+1)); err != nil {
			return nil, err
		}
		if err := CPUutilization.SetType("gauge"); err != nil {
			return nil, err
		}
		if err := CPUutilization.SetValue(cpu); err != nil {
			continue
		}
		if err := CPUutilization.SetHash(cfg.Key); err != nil {
			return nil, err
		}
		if err := mtrxPool.Create(CPUutilization.ID, *CPUutilization); err != nil { // Метрика добавляется в общий пул (мапку)
			return nil, errors.New("unable add CPUutilization into MetricsPool: " + CPUutilization.GetID() + "  " + CPUutilization.GetType())
		}

	}
	return mtrxPool, nil
}
