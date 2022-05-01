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

//main entry point
func main() {

	err := config.Parse(cli.ConfigPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)
	s := server.New(&cfg)
	if err = s.Start(); err != nil {
		os.Exit(0)
	}
	os.Exit(0)
}
