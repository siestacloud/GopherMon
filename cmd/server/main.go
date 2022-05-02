package main

import (
	"log"
	"os"

	"github.com/siestacloud/service-monitoring/internal/server"
	"github.com/siestacloud/service-monitoring/internal/server/config"
)

var (
	cli config.CLI
	cfg config.ServerConfig
)

func main() {

	err := config.Parse(&cli, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	s, err := server.New(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Start(); err != nil {
		os.Exit(0)
	}
	os.Exit(0)
}
