package config

import (
	"log"

	"github.com/spf13/viper"
)

// CLI .
type CLI struct {
	ConfigPath string `help:"Config path" type:"path" default:"config.yaml"`
	// Add more here

}

func ParseFile(path string, out interface{}) error {
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
	return nil
}
