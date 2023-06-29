package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	serverMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

func AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	body, _ := memStorage.PrintAll()
	io.WriteString(res, body)
}

func MetricHandler(res http.ResponseWriter, req *http.Request) {
	metricType := strings.ToLower(chi.URLParam(req, "type"))
	if metricType != serverMetrics.CounterType && metricType != serverMetrics.GaugeType {
		http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		loggingResponse.WriteHeader(http.StatusNotFound)
	} else {
		metric := strings.ToLower(chi.URLParam(req, "metric"))
		if metricType == serverMetrics.CounterType {
			metricValue, ok := memStorage.MetricValueIfExists(metric, metricType)
			if ok {
				io.WriteString(res, metricValue+"   ")
				loggingResponse.WriteHeader(http.StatusOK)
				return
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
				loggingResponse.WriteHeader(http.StatusNotFound)
				return
			}
		}

		if metricType == serverMetrics.GaugeType {
			metricValue, ok := memStorage.MetricValueIfExists(metric, metricType)
			if ok {
				loggingResponse.Write([]byte(metricValue))
				loggingResponse.WriteHeader(http.StatusOK)
				return
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
				loggingResponse.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}
}
