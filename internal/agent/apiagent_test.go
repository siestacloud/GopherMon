package agent

import (
	"reflect"
	"testing"
)

func TestAPIAgent_SendMetric(t *testing.T) {
	type args struct {
		name string
		m    reflect.Value
	}
	tests := []struct {
		name    string
		c       *APIAgent
		args    args
		wantErr bool
	}{{
		name:    "Test1",
		c:       &APIAgent{config: &Config{ReportAddr: "http://127.0.0.1:8080", ReportInterval: 10, PollInterval: 2}},
		args:    args{name: "Mymetric", m: reflect.ValueOf(0)},
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.sendMetric(tt.args.name, &tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("APIAgent.SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
