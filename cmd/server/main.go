package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/MustCo/Mon_go/internal/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.json", "path to config json file")
}

func main() {
	flag.Parse()
	config := server.NewConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, config)

	server := server.New(config)
	server.Start()

}
