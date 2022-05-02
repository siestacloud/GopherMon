package config

import "time"

//Config Server configuration struct
type ServerConfig struct {
	Server `mapstructure:"server"`
}

type Server struct {
	LogLevel      string        `mapstructure:"log_level"`
	LogFile       string        `mapstructure:"log_file"`
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080" mapstructure:"address"`
	Restore       bool          `env:"RESTORE" envDefault:"true" mapstructure:"restore"`
	StoreInterval time.Duration `env:"STORE_INT" envDefault:"300s" mapstructure:"storeinterval"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json" mapstructure:"storefile"`
	Timeout       struct {
		// Server is the general server timeout to use
		// for graceful shutdowns
		Server int `mapstructure:"server"`

		// Write is the amount of time to wait until an HTTP server
		// write opperation is cancelled
		Write int `mapstructure:"write"`

		// Read is the amount of time to wait until an HTTP server
		// read operation is cancelled
		Read int `mapstructure:"read"`

		// Read is the amount of time to wait
		// until an IDLE HTTP session is closed
		Idle int `mapstructure:"idle"`
	} `mapstructure:"timeout"`
}
