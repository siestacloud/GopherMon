package server

import (
	"log"
	"net/http"
	"strings"
)

type APIServer struct {
	config *Config
}

type Handler struct {
	DB Storage
}

func New(config *Config) *APIServer {
	return &APIServer{config: config}
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	API := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if API[0] != "update" {
		log.Println("Bad Request", r.URL.Path, API)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(API) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := handler.DB.Set(API[1], API[2], API[3])
	if err != nil {
		log.Println("Bad Request", r.URL.Path, API)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(handler.DB)
	w.WriteHeader(http.StatusOK)
}

func (s *APIServer) Start() error {
	updater := new(Handler)
	updater.DB = new(DB)
	updater.DB.Init()
	http.Handle("/update/", updater)
	return http.ListenAndServe(s.config.BindAddr, nil)
}
