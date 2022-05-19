package agent

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"time"

	"github.com/MustCo/Mon_go/internal/utils"
)

type APIAgent struct {
	config *Config
	client http.Client
}

func New(config *Config) *APIAgent {
	return &APIAgent{config: config}
}

func (c *APIAgent) Report(ms *utils.Metrics) error {
	for name, v := range ms.Counters {
		metric := reflect.ValueOf(v)
		err := c.SendMetric(name, &metric)
		if err != nil {
			return err
		}
	}
	for name, v := range ms.GaugesMy {
		metric := reflect.ValueOf(v)
		err := c.SendMetric(name, &metric)
		if err != nil {
			return err
		}

	}
	for name, v := range ms.GaugesRuntime {
		metric := reflect.ValueOf(v)
		err := c.SendMetric(name, &metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *APIAgent) SendMetric(name string, m *reflect.Value) error {
	var url string
	switch m.Kind() {
	case reflect.Uint64, reflect.Uint32:
		url = fmt.Sprintf("%s/update/%s/%s/%v", c.config.ReportAddr, m.Type().Name(), name, m.Uint())
	case reflect.Float32, reflect.Float64:
		url = fmt.Sprintf("%s/update/%s/%s/%v", c.config.ReportAddr, m.Type().Name(), name, m.Float())
	case reflect.Int32, reflect.Int64:
		url = fmt.Sprintf("%s/update/%s/%s/%v", c.config.ReportAddr, m.Type().Name(), name, m.Int())
	}
	fmt.Println(url)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "text/plain")
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(dump))
	resp, err := c.client.Do(r)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *APIAgent) Start() error {
	m := new(utils.Metrics)
	c.client = http.Client{}
	reports := time.NewTicker(time.Duration(c.config.ReportInterval) * time.Second)
	polls := time.NewTicker(time.Duration(c.config.PollInterval) * time.Second)
	m.Init()
	for {
		select {
		case <-reports.C:
			c.Report(m)
			_, err := fmt.Println("Report metrics")
			if err != nil {
				return err
			}
		case <-polls.C:
			m.Poll()
			fmt.Println("Poll metrics")

		}
	}

}
