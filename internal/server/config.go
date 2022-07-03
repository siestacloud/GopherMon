package server

import "os"

type Config struct {
	BindAddr string `json:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "127.0.0.1:8080",
	}
}

func EnvConfig() *Config {
	var address string

	if address = os.Getenv("ADDRESS"); address == "" {
		address = "127.0.0.1:8080"
	}
	return &Config{
		BindAddr: address,
	}
}
