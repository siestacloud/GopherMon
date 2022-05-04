package config

import "time"

//Config Server configuration struct
type ServerConfig struct {
	Server `mapstructure:"server"`
}

type Server struct {
	LogLevel      string        `mapstructure:"log_level"`
	LogFile       string        `mapstructure:"log_file"`
	Address       string        `env:"ADDRESS" mapstructure:"address"`
	Restore       bool          `env:"RESTORE" mapstructure:"restore"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" mapstructure:"storeinterval"`
	StoreFile     string        `env:"STORE_FILE" mapstructure:"storefile"`
	Timeout       struct {
		// Server is the general server timeout to use
		// for graceful shutdowns
		Server time.Duration `mapstructure:"server"`

		// Write is the amount of time to wait until an HTTP server
		// write opperation is cancelled
		Write time.Duration `mapstructure:"write"`

		// Read is the amount of time to wait until an HTTP server
		// read operation is cancelled
		Read time.Duration `mapstructure:"read"`

		// Read is the amount of time to wait
		// until an IDLE HTTP session is closed
		Idle time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeout"`
}
