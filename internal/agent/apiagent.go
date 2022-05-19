package agent

import (
	"fmt"
	"time"

	"github.com/MustCo/Mon_go/internal/utils"
)

type APIAgent struct {
	config *Config
}

func New(config *Config) *APIAgent {
	return &APIAgent{config: config}
}

func (c *APIAgent) Report(*utils.Metrics) error {
	return nil
}

func (c *APIAgent) Start() error {
	m := new(utils.Metrics)
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
