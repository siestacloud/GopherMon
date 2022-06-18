package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/go-resty/resty/v2"
)

type APIAgent struct {
	config *Config
	client *resty.Client
}

func New(config *Config) *APIAgent {
	return &APIAgent{config: config}
}

func (c *APIAgent) Report(ms *utils.Metrics) error {
	for name, v := range ms.Counters {
		err := c.sendMetric(name, v)
		if err != nil {
			return err
		}
	}
	for name, v := range ms.Gauges {
		err := c.sendMetric(name, v)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *APIAgent) sendMetric(name string, m interface{}) error {

	var val, t string
	switch m := m.(type) {
	case utils.Gauge:
		t = "gauge"
		val = strconv.FormatFloat(float64(m), 'e', 2, 64)
	case utils.Counter:
		t = "counter"
		val = strconv.FormatInt(int64(m), 10)
	}
	log.Printf("Send metric: %s : %v", name, val)
	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"host": c.config.ReportAddr,
			"type": t,
			"name": name,
			"val":  val,
		}).
		Post("http://{host}/update/{type}/{name}/{val}")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		fmt.Println("  Status Code:", resp.StatusCode())
		fmt.Println("  Status     :", resp.Status())
		fmt.Println("  Proto      :", resp.Proto())
		fmt.Println("  Time       :", resp.Time())
		fmt.Println("  Received At:", resp.ReceivedAt())
		fmt.Println("  Body       :\n", resp)
		return errors.New("invalid status code")
	}
	return nil
}

func (c *APIAgent) Start(ctx context.Context) error {
	m := new(utils.Metrics)
	c.client = resty.New()
	c.client.R().
		SetHeader("Content-Type", "text/plain")

	reports := time.NewTicker(time.Duration(c.config.ReportInterval) * time.Second)
	polls := time.NewTicker(time.Duration(c.config.PollInterval) * time.Second)
	m.Init()
	for {
		select {
		case <-reports.C:
			err := c.Report(m)
			if err != nil {
				log.Println("Error", err)
			}
		case <-polls.C:
			m.Poll()
		case <-ctx.Done():
			log.Println("Exit by context")
			return nil
		}
	}

}
