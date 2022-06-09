package server

import (
	"errors"
	"testing"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestDB_Set(t *testing.T) {
	type args struct {
		t    string
		name string
		val  string
	}
	type want struct {
		db  *DB
		err error
	}
	tests := []struct {
		name string
		db   *DB
		args args
		want want
	}{
		{
			name: "TestGauge",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
			args: args{t: "gauge", name: "Mymetric", val: "1.329184"},
			want: want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 1.329184}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
			args: args{t: "gauge", name: "Mymetric", val: "1,FSD29184"},
			want: want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
				err: errors.New("error"),
			},
		},
		{
			name: "TestCounter",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
			args: args{t: "Counter", name: "Mymetric", val: "14563"},
			want: want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
			args: args{t: "Counter", name: "Mymetric", val: "1.329184"},
			want: want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
				err: errors.New("error"),
			},
		},
		{
			name: "TestInvalidType",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
			args: args{t: "Mytype", name: "Mymetric", val: "1sdfgsd4"},
			want: want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{"TestCounter": 1234}}},
				err: errors.New("invalid type"),
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.db.Set(tt.args.t, tt.args.name, tt.args.val)
			assert.Equal(t, tt.want.db, tt.db)
		})
	}
}

func TestDB_Get(t *testing.T) {
	type args struct {
		t    string
		name string
	}
	type want struct {
		res utils.Gauge
		err error
	}
	tests := []struct {
		name string
		db   *DB
		args args
		want want
	}{
		{
			name: "TestGauge",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 14.563}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
			args: args{t: "Gauge", name: "Mymetric"},
			want: want{res: 14.563,
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 14.563}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
			args: args{t: "gauge", name: "Unknown"},
			want: want{res: 0,
				err: nil,
			},
		},
		{
			name: "TestCounter",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 14.563}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
			args: args{t: "Counter", name: "Mymetric"},
			want: want{res: 14563,
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 14.563}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
			args: args{t: "counter", name: "Unknown"},
			want: want{res: 0,
				err: nil,
			},
		},
		{
			name: "TestInvalidType",
			db:   &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123, "Mymetric": 14.563}, Counters: map[string]utils.Counter{"TestCounter": 1234, "Mymetric": 14563}}},
			args: args{t: "untype", name: "Mymetric"},
			want: want{res: 0,
				err: errors.New("invalid type"),
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.db.Get(tt.args.t, tt.args.name)
			assert.Equal(t, tt.want.res, got)
		})
	}
}
