package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "positive test #1",
			path: "/update/gauge/alloc/12.0",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "plain/text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Дописать все проверки
			request := httptest.NewRequest(http.MethodPost, tt.path, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			PostHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, tt.want.code)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
		})
	}
}

func Test_checkMetricsValue(t *testing.T) {
	type args struct {
		res         http.ResponseWriter
		metricValue string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkMetricsValue(tt.args.res, tt.args.metricValue)
		})
	}
}

func Test_checkType(t *testing.T) {
	type args struct {
		res         http.ResponseWriter
		metricsType string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkType(tt.args.res, tt.args.metricsType)
		})
	}
}

func Test_methodTypeHandler(t *testing.T) {
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			methodTypeHandler(tt.args.res, tt.args.req)
		})
	}
}
