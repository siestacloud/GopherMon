package config

import "time"

//Config Server configuration struct
type ServerConfig struct {
	Server `mapstructure:"server"`
}

type Server struct {
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`
	Address  string `env:"ADDRESS" envDefault:"localhost:8080" mapstructure:"address"`
	// StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json" mapstructure:"storefile"`
	Restore bool `env:"RESTORE" envDefault:"true" mapstructure:"restore"`
	// StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300" mapstructure:"storeinterval"`
	Timeout struct {
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
