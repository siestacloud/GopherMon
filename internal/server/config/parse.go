package config

import (
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/spf13/viper"
)

// CLI .
type CLI struct {
	ConfigPath string `help:"Config path" type:"path" default:"config.yaml"`
	// Add more here
}

var ()

func Parse(path string, out *ServerConfig) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.apiserver.yaml") // call multiple times to add many search paths
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("Fatal error config file: %w \n ", err)
	}

	if err := viper.Unmarshal(&out); err != nil {
		log.Fatal(err)
	}
	err = env.Parse(out)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
