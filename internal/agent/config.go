package agent

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ReportAddr     string `json:"report_addr"`
	PollInterval   int64  `json:"poll_interval"`
	ReportInterval int64  `json:"report_interval"`
}

func NewConfig() *Config {
	return &Config{
		ReportAddr:     "127.0.0.1:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
}
func EnvConfig() *Config {
	var address string
	var poll, report int64
	var err error

	if os.Getenv("ADDRESS") != "" {
		address = os.Getenv("ADDRESS")
	} else {
		address = "127.0.0.1:8080"
	}
	if reportS := os.Getenv("REPORT_INTERVAL"); reportS != "" {
		report, err = strconv.ParseInt(reportS, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		report = 10
	}
	if pollS := os.Getenv("REPORT_INTERVAL"); pollS != "" {
		poll, err = strconv.ParseInt(pollS, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		poll = 2
	}

	return &Config{
		ReportAddr:     address,
		PollInterval:   poll,
		ReportInterval: report,
	}
}
