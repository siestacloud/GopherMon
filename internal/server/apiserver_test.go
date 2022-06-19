package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/labstack/echo/v4"
)

func TestUpdateHandler_postMetric(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		db  *DB
		sc  int
		err error
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
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{}, Counters: map[string]utils.Counter{}}}, sc: http.StatusMethodNotAllowed, err: echo.NewHTTPError(http.StatusNotImplemented, "invalid type")},
		},
		{
			name:    "Test2",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123", nil)},
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{"TestGauge": 123}, Counters: map[string]utils.Counter{}}}, sc: http.StatusOK, err: nil},
		},
		{
			name:    "Test3",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/status", nil)},
			want:    want{db: &DB{Metrics: &utils.Metrics{Gauges: map[string]utils.Gauge{}, Counters: map[string]utils.Counter{}}}, sc: http.StatusBadRequest, err: echo.NewHTTPError(http.StatusNotImplemented, "invalid type")},
		},
	}
	e := echo.New()
	updater := NewUpdateHandler()
	e.GET("/", updater.getAllMetrics)
	e.GET("/value/:type/:name", updater.getMetric)
	e.POST("/update/:type/:name/:value", updater.postMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := e.NewContext(tt.args.r, tt.args.w)
			updater.postMetric(c)
		})
	}
}
