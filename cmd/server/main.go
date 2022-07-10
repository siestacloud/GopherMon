package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/config"
	"github.com/siestacloud/service-monitoring/internal/server/repository"
	"github.com/siestacloud/service-monitoring/internal/server/service"
	"github.com/siestacloud/service-monitoring/internal/server/transport/rest"
	"github.com/siestacloud/service-monitoring/internal/server/transport/rest/handler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfg config.Cfg
)

func main() {
	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	mp := core.NewMetricsPool()
	repos := repository.NewRepository(mp, db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(&cfg, services)
	s, err := rest.NewServer(&cfg, handlers)
	if err != nil {
		log.Fatal(err)
	}

	// если в cfg задан путь до файла с mtrx и задан интервал
	if cfg.StoreFile != "" {
		if cfg.StoreInterval != 0 {
			if err := mp.RLS(cfg.StoreFile); err != nil {
				log.Fatal(err)
			}
			// пул будет сохранятся на диск с опр интервалом
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
		logrus.Errorf("Server was unable to gracefully shutdown due to err: %+v", err)
		os.Exit(1)
	}

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
