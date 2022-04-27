package storage

import (
	"encoding/json"

	"github.com/siestacloud/service-monitoring/internal/mtrx"
)

type Storage struct {
	Mp *mtrx.MetricsPool
}

func New() *Storage {
	return &Storage{
		Mp: mtrx.NewMetricsPool(),
	}
}

func (s *Storage) TakeAll() ([]byte, error) {

	js, err := json.MarshalIndent(s.Mp, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}
