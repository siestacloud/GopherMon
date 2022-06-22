package server

import (
	"errors"
	"strconv"
	"sync"

	"github.com/MustCo/Mon_go/internal/utils"
)

type Storage interface {
	Get(t, name string) (*utils.Metrics, error)
	GetAll() map[string]utils.Metrics
	Set(t, name, val string) error
}

func NewDB() *DB {
	db := new(DB)
	db.Metrics = utils.NewMetricsStorage()
	return db
}

type DB struct {
	mut     sync.Mutex
	Metrics utils.MetricsStorage
}

func (db *DB) Set(t, name, val string) error {
	metrica := utils.NewMetrics(name, t)
	switch metrica.MType {
	case "counter":
		d, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		metrica.Delta = &d
	case "gauge":
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		metrica.Value = &v
	default:
		return errors.New("invalid type")
	}
	db.mut.Lock()
	db.Metrics[metrica.ID] = *metrica
	db.mut.Unlock()
	return nil

}

func (db *DB) Get(t, name string) (*utils.Metrics, error) {
	if _, ok := utils.Types[t]; ok {
		db.mut.Lock()
		if m, ok := db.Metrics[name]; ok {
			db.mut.Unlock()
			return &m, nil
		}
	}
	return nil, errors.New("")
}

func (db *DB) GetAll() map[string]utils.Metrics {
	return db.Metrics
}
