package server

import (
	"net/http"
)

type APIServer struct {
	config *Config
}

func New(config *Config) *APIServer {
	return &APIServer{config: config}
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *APIServer) Start() error {
	http.HandleFunc("/update", UpdateMetrics)
	return http.ListenAndServe(s.config.BindAddr, nil)
}
