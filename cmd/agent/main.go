package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, config)
	agent := agent.New(config)
	err = agent.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
