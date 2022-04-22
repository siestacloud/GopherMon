package metricscustom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricPool(t *testing.T) {
	tmp := MetricsPool{M: map[string]Metric{}}
	test := []struct {
		name string
		want MetricsPool
	}{
		{
			name: "Test #1",
			want: tmp,
		},
	}
	for _, tt := range test {
		mp := NewMetricsPool()
		assert.Equal(t, tt.want, *mp)
	}
}

func TestNewMetrics(t *testing.T) {

	test := []struct {
		name   string
		values []string
		want   Metric
	}{
		{
			name:   "Test #1",
			values: []string{"testMetric", "counter", "1"},
			want:   Metric{Name: "testMetric", Value: 1, Types: "counter"},
		},
		{
			name:   "Test #2",
			values: []string{"metric2", "gauge", "123"},
			want:   Metric{Name: "metric2", Value: 123, Types: "gauge"},
		},
		{
			name:   "Test #3",
			values: []string{"metrics3", "counter", "111"},
			want:   Metric{Name: "metrics3", Value: 111, Types: "counter"},
		},
	}
	for _, tt := range test {
		m, status := NewMetric(tt.values[1], tt.values[0], tt.values[2])
		if status != "" {
			fmt.Println(status)
		}
		assert.Equal(t, tt.want, *m)
	}
}

func TestConvert(t *testing.T) {
	tmp := MetricsPool{M: map[string]Metric{}}
	test := []struct {
		name string

		want MetricsPool
	}{
		{
			name: "Test #1",
			want: tmp,
		},
	}
	for _, tt := range test {
		mp := NewMetricsPool()
		assert.Equal(t, tt.want, *mp)
	}
}
