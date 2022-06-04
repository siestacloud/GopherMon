package server

type Config struct {
	BindAddr string `json:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "127.0.0.1:8080",
	}
}
