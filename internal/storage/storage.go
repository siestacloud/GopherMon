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
		if v.ID == m.ID {
			if v.MType == "counter" {
				fmt.Println("Ok")
				metr, status := metricscustom.NewMetric(m.MType, m.ID, fmt.Sprint(m.Delta))
				if status != "" {
					return false
				}
				fmt.Println(metr)
				s.Mp.M[k] = *metr

			}
			if v.MType == "gauge" {
				s.Mp.M[k] = *m
			}
			return true
		}
	}
	s.Mp.M[m.ID] = *m
	return false
}

//
func (s *Storage) Take(t, n string) *metricscustom.Metric {

	for k, v := range s.Mp.M {
		if v.ID == n {
			if v.MType == t {
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
