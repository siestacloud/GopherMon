package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/MustCo/Mon_go/internal/agent"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiagent.json", "path to config file")
}

func main() {
	flag.Parse()
	config := agent.NewConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, config)
	agent := agent.New(config)
	err = agent.Start()
	if err != nil {
		log.Fatal(err)
	}

}
