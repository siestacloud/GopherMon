package storage

import (
	"encoding/json"
	"fmt"

	"github.com/siestacloud/service-monitoring/internal/metricscustom"
)

type Storage struct {
	Mp *metricscustom.MetricsPool
}

func New() *Storage {
	return &Storage{
		Mp: metricscustom.NewMetricsPool(),
	}
}

func (s *Storage) Update(m *metricscustom.Metric) bool {

	for k, v := range s.Mp.M {
		if v.Name == m.Name {
			if v.Types == "counter" {
				metr, err := metricscustom.NewMetric(m.Types, m.Name, fmt.Sprint(m.Value+v.Value))
				if err != nil {
					return false
				}
				s.Mp.M[k] = *metr

			}
			if v.Types == "gauge" {
				s.Mp.M[k] = *m
			}
			return true
		}
	}
	s.Mp.M[m.Name] = *m
	return false
}

func (s *Storage) Take(t, n string) *metricscustom.Metric {

	for k, v := range s.Mp.M {
		if v.Name == n {
			if v.Types == t {
				m := s.Mp.M[k]
				return &m
			}
			return nil
		}
	}

	return nil
}

func (s *Storage) TakeAll() ([]byte, error) {

	js, err := json.MarshalIndent(s.Mp, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}
