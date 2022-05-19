package server

type Config struct {
	BindAddr       string `json:"bind_addr"`
	PollInterval   int64  `json:"poll_interval"`
	ReportInterval int64  `json:"report_interval"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:       "127.0.0.1:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
}
