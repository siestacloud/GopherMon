package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/config"
	"github.com/siestacloud/service-monitoring/internal/server/repository"
	"github.com/siestacloud/service-monitoring/internal/server/service"
	"github.com/siestacloud/service-monitoring/internal/server/transport/rest"
	"github.com/siestacloud/service-monitoring/internal/server/transport/rest/handler"
	"github.com/sirupsen/logrus"
)

var (
	cfg config.Cfg
)

func main() {
	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	mp := core.NewMetricsPool()
	repos := repository.NewRepository(mp)
	services := service.NewService(repos)
	handlers := handler.NewHandler(&cfg, services)
	s, err := rest.NewServer(&cfg, handlers)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.StoreFile != "" {
		if cfg.StoreInterval != 0 {
			go func() {
				for {
					time.Sleep(cfg.StoreInterval)
					if err := mp.WLS(cfg.StoreFile); err != nil {
						logrus.Error("error store interval: ", err)
					}
					logrus.Info("Storage update")
				}
			}()
		}
	}

	if err = s.Start(); err != nil {
		if cfg.StoreFile != "" {
			file, err := json.MarshalIndent(mp, "", " ")
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(cfg.StoreFile, file, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}
	os.Exit(0)
}
