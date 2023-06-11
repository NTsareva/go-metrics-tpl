package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockServerHandler (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}


func TestSendRuntimeMetrics(t *testing.T) {
	type args struct {
		m *MetricsGauge
	}
	tests := []struct {
		name string

	}{
		{
			name: "positive test #1"
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(mockServerHandler))
			defer server.Close()



		})
	}
}
