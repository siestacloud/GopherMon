package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/viper"
)

// Setting up configuration parametrs
func Parse(cfg *Cfg) error {

	// Читаю конфиг
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.mon.yaml")
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error config file: %w \n ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	// Читаю флаги, переопределяю параметры, если флаги заданы
	flag.DurationVar(&cfg.StoreInterval, "i", 300000000000, "Duration time of saving. Possible values: 20s 3s")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Path to store file. Possible values: /path/to/file")
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Address for server. Possible values: localhost:8080")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore metrics pool. Possible values: true false")
	flag.Parse()

	// Читаю переменные окружения, переопределяю параметры, если пер окр заданы
	err = env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
