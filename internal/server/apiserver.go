package apiserver

import (
	"net/http"
)

type Config struct {
	bind_addr string
}

type APIServer struct {
	config *Config
}

func New(config *Config) *APIServer {
	return &APIServer{config: config}
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func (s *APIServer) Start() error {
	http.HandleFunc("/update", UpdateMetrics)
	return http.ListenAndServe(s.config.bind_addr, nil)
}
