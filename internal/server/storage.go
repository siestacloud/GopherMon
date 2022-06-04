package server

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/MustCo/Mon_go/internal/utils"
)

type Storage interface {
	Init()
	Get(t, name string) (utils.Gauge, error)
	Set(t, name, value string) error
}

type DB struct {
	Metrics *utils.Metrics
}

func (db *DB) Init() {
	db.Metrics = new(utils.Metrics)
	db.Metrics.Init()
}

func (db *DB) Set(t, name, val string) error {
	log.Println(val)
	switch t {
	case "Gauge":

		g, err := strconv.ParseFloat(val, 64)
		log.Printf("Gauge %s %s = %e", name, val, g)
		db.Metrics.Gauges[name] = utils.Gauge(g)
		return err
	case "Counter":
		ctr, err := strconv.Atoi(val)
		db.Metrics.Counters[name] = utils.Counter(ctr)
		return err
	default:
		return errors.New("invalid type")
	}
}

func (db *DB) Get(t, name string) (utils.Gauge, error) {
	switch t {
	case "gauge":
		return db.Metrics.Gauges[name], nil
	case "counter":
		return utils.Gauge(db.Metrics.Counters[name]), nil
	}
	return 0, errors.New("invalid type")
}

func (db *DB) String() string {
	result := fmt.Sprintf("Counters:%v\nGauges:%v\n", db.Metrics.Counters, db.Metrics.Gauges)
	return result

}
