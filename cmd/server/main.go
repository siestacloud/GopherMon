package main

import (
	"context"
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
	ctx := context.TODO()
	server := server.New(config)
	log.Fatal(server.Start(ctx))

}
