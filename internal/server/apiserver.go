package server

import (
	"log"
	"net/http"
	"strings"
)

type APIServer struct {
	config *Config
}

type UpdateHandler struct {
	DB Storage
}

func New(config *Config) *APIServer {
	return &APIServer{config: config}
}

func NewUpdateHandler() *UpdateHandler {
	updater := new(UpdateHandler)
	updater.DB = NewDB()
	return updater
}

func (handler *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	API := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if API[0] != "update" {
		log.Println("Bad Request", r.URL.Path, API)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if len(API) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := handler.DB.Set(API[1], API[2], API[3])
	switch {
	case err == nil:
	case err.Error() == "invalid type":
		log.Println("Bad Request", r.URL.Path, API)
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return

	default:
		log.Println("Bad Request", r.URL.Path, API)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(handler.DB)
}

func (s *APIServer) Start() error {
	updater := NewUpdateHandler()
	http.Handle("/update/", updater)
	return http.ListenAndServe(s.config.BindAddr, nil)
}
