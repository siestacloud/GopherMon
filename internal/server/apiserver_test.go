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
	test_metrics := utils.NewMetricsStorage()
	test_metrics["TestGauge"] = *utils.NewMetrics("TestGauge", "gauge")
	*test_metrics["TestGauge"].Value = 123.124
	test_metrics["TestCounter"] = *utils.NewMetrics("TestCounter", "counter")
	*test_metrics["TestCounter"].Delta = 123

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
			name:    "Invalid Path",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodGet, "/status", nil)},
			want:    want{db: &DB{Metrics: utils.MetricsStorage{}}, sc: http.StatusMethodNotAllowed, err: echo.NewHTTPError(http.StatusNotImplemented, "invalid type")}},
		{
			name:    "Post_gauge",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123.124", nil)},
			want:    want{db: &DB{Metrics: utils.MetricsStorage{"TestGauge": test_metrics["TestGauge"]}}, sc: http.StatusOK, err: nil},
		},
		{
			name:    "Post_counter",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/counter/TestCounter/123", nil)},
			want:    want{db: &DB{Metrics: utils.MetricsStorage{"TestCounter": test_metrics["TestCounter"]}}, sc: http.StatusOK, err: nil},
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
