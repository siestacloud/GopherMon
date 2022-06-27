package server

import (
	"errors"
	"strconv"
	"sync"

	"github.com/MustCo/Mon_go/internal/utils"
)

type Storage interface {
	Get(t, name string) (*string, error)
	GetAll() map[string]utils.Metrics
	Set(t, name, val string) error
	SetMetrica(metrica *utils.Metrics) error
	GetMetrica(t, name string) (*utils.Metrics, error)
}

func NewDB() *DB {
	db := new(DB)
	db.Metrics = utils.NewMetricsStorage()
	db.mut = new(sync.Mutex)
	return db
}

type DB struct {
	mut     *sync.Mutex
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
		*metrica.Delta += d
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

func (db *DB) SetMetrica(metrica *utils.Metrics) error {
	db.mut.Lock()
	if _, ok := db.Metrics[metrica.ID]; !ok {
		db.Metrics[metrica.ID] = *utils.NewMetrics(metrica.ID, metrica.MType)
	}
	db.mut.Unlock()

	switch metrica.MType {
	case "gauge":
		db.mut.Lock()
		*db.Metrics[metrica.ID].Value = *metrica.Value
		db.mut.Unlock()
	case "counter":
		db.mut.Lock()
		*db.Metrics[metrica.ID].Delta += *metrica.Delta
		db.mut.Unlock()

	default:
		return errors.New("invalid type")
	}

	return nil

}

func (db *DB) Get(t, name string) (*string, error) {
	if _, ok := utils.Types[t]; ok {
		db.mut.Lock()
		defer db.mut.Unlock()
		if m, ok := db.Metrics[name]; db.Metrics[name].MType == t && ok {
			s := ""
			switch m.MType {
			case "gauge":
				s = strconv.FormatFloat(*m.Value, 'f', 3, 64)
			case "counter":
				s = strconv.FormatInt(*m.Delta, 10)
			}
			return &s, nil
		}
		return nil, errors.New("unknown metric")
	}
	return nil, errors.New("invalid type")
}

func (db *DB) GetMetrica(t, name string) (*utils.Metrics, error) {
	if _, ok := utils.Types[t]; ok {
		db.mut.Lock()
		defer db.mut.Unlock()
		if m, ok := db.Metrics[name]; db.Metrics[name].MType == t && ok {

			return &m, nil
		}
		return nil, errors.New("unknown metric")
	}
	return nil, errors.New("invalid type")
}

func (db *DB) GetAll() map[string]utils.Metrics {
	return db.Metrics
}
