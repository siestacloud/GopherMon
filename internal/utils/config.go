package utils

type Config struct {
	Address        string `env:"ADDRESS,required" envDefault:"127.0.0.1:8080" json:"ADDRESS"`
	PollInterval   int64  `env:"POLL_INTERVAL" envDefault:"2" json:"POLL_INTERVAL"`
	ReportInterval int64  `env:"REPORT_INTERVAL"  envDefault:"10" json:"REPORT_INTERVAL"`
}

func NewConfig() *Config {
	return &Config{}
}
