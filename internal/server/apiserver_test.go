package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		db *DB
		sc int
	}
	tests := []struct {
		name    string
		handler *UpdateHandler
		args    args
		want    want
	}{
		{
			name:    "Test1",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodGet, "/status", nil)},
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{}, Counters: map[string]utils.Counter{}}}, sc: http.StatusMethodNotAllowed},
		},
		{
			name:    "Test2",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123", nil)},
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{}}}, sc: http.StatusOK},
		},
		{
			name:    "Test3",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/status", nil)},
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{}, Counters: map[string]utils.Counter{}}}, sc: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.handler.ServeHTTP(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			assert.Equal(t, tt.want.sc, res.StatusCode)
			assert.Equal(t, tt.want.db, tt.handler.DB)
		})
	}
}
