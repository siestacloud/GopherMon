package config

import "time"

//Config Server configuration struct
type ServerConfig struct {
	Server struct {
		LogLevel string `mapstructure:"log_level"`
		LogFile  string `mapstructure:"log_file"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Timeout  struct {
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
	} `mapstructure:"server"`
}
