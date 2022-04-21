package storage

import (
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
