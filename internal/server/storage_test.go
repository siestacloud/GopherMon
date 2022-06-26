package server

import (
	"errors"
	"testing"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestDB_Set(t *testing.T) {

	init_db := NewDB()
	init_db.Metrics = utils.NewMetricsStorage()
	m := utils.NewMetrics("TestGauge", "gauge")
	*m.Value = 123.123
	init_db.Metrics["TestGauge"] = *m
	m = utils.NewMetrics("TestCounter", "counter")
	*m.Delta = 123
	init_db.Metrics["TestCounter"] = *m

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
			db:   init_db,
			args: args{t: "gauge", name: "Mymetric", val: "1.329184"},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   init_db,
			args: args{t: "gauge", name: "Mymetric", val: "1,FSD29184"},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "TestCounter",
			db:   init_db,
			args: args{t: "Counter", name: "Mymetric", val: "14563"},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   init_db,
			args: args{t: "Counter", name: "Mymetric", val: "1.329184"},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "TestInvalidType",
			db:   init_db,
			args: args{t: "Mytype", name: "Mymetric", val: "1sdfgsd4"},
			want: want{
				err: errors.New("invalid type"),
			},
		}}
	for i, test := range tests {
		db := NewDB()
		db.Metrics = utils.NewMetricsStorage()
		m := utils.NewMetrics("TestGauge", "gauge")
		*m.Value = 123.123
		db.Metrics["TestGauge"] = *m
		m = utils.NewMetrics("TestCounter", "counter")
		*m.Delta = 123
		db.Metrics["TestCounter"] = *m
		db.Set(test.args.t, test.args.name, test.args.val)
		tests[i].want.db = db
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.db.Set(tt.args.t, tt.args.name, tt.args.val)
		})
	}
}

func TestDB_Get(t *testing.T) {
	init_db := NewDB()
	init_db.Set("gauge", "TestGauge", "123.123")
	init_db.Set("counter", "TestCounter", "123")
	init_db.Set("gauge", "Mygauge", "14.563")
	init_db.Set("counter", "Mycounter", "14563")

	type args struct {
		t    string
		name string
	}
	type want struct {
		res utils.Metrics
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
			db:   init_db,
			args: args{t: "gauge", name: "Mygauge"},
			want: want{res: init_db.Metrics["Mygauge"],
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   init_db,
			args: args{t: "gauge", name: "Unknown"},
			want: want{res: *utils.NewMetrics("Mymetric", "untype"),
				err: errors.New("unknown metric"),
			},
		},
		{
			name: "TestCounter",
			db:   init_db,
			args: args{t: "counter", name: "Mycounter"},
			want: want{res: init_db.Metrics["Mycounter"],
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   init_db,
			args: args{t: "counter", name: "Unknown"},
			want: want{res: *utils.NewMetrics("Mymetric", "untype"),
				err: errors.New("unknown metric"),
			},
		},
		{
			name: "TestInvalidType",
			db:   init_db,
			args: args{t: "untype", name: "Mymetric"},
			want: want{res: *utils.NewMetrics("Mymetric", "untype"),
				err: errors.New("invalid type"),
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Get(tt.args.t, tt.args.name)
			if err != nil {
				assert.Equal(t, err.Error(), tt.want.err.Error())
				return
			}
			assert.Equal(t, tt.want.res, *got)
		})
	}
}
