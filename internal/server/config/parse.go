package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/spf13/viper"
)

// CLI .
type CLI struct {
	ConfigPath string `help:"Config path" type:"path" default:"config.yaml"`
}

var ()

//
func Parse(c *CLI, cfg *ServerConfig) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(c.ConfigPath)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.apiserver.yaml") // call multiple times to add many search paths
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("Fatal error config file: %w \n ", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.DurationVar(&cfg.StoreInterval, "i", 300, "Duration time of saving. Possible values: int")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Path to store file. Possible values: /path/to/file")
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Address for server. Possible values: localhost:8080")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore metrics pool. Possible values: true false")

	err = env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	return nil
}
