package main

import (
	"log"

	"github.com/siestacloud/service-monitoring/internal/server"
	"github.com/siestacloud/service-monitoring/internal/server/config"
)

var (
	cli config.CLI
	cfg config.ServerConfig
)

func main() {

	err := config.ParseFile(cli.ConfigPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(&cfg)
	if err = s.Start(); err != nil {
		return
	}

}
