package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	serverMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

func (serverHandlers *SeverHandlers) AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	body, _ := serverHandlers.MemStorage.PrintAll()
	io.WriteString(res, body)
}

func (serverHandlers *SeverHandlers) MetricHandler(res http.ResponseWriter, req *http.Request) {

	metricType := strings.ToLower(chi.URLParam(req, "type"))
	if metricType != serverMetrics.CounterType && metricType != serverMetrics.GaugeType {
		http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		loggingResponse.WriteHeader(http.StatusNotFound)
	} else {
		metric := strings.ToLower(chi.URLParam(req, "metric"))
		if metricType == serverMetrics.CounterType {
			metricValue, ok := serverHandlers.MemStorage.MetricValueIfExists(metric, metricType)
			if ok {
				io.WriteString(res, metricValue+"   ")
				loggingResponse.WriteHeader(http.StatusOK)
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
				loggingResponse.WriteHeader(http.StatusNotFound)
			}
		}
		if metricType == serverMetrics.GaugeType {
			metricValue, ok := serverHandlers.MemStorage.MetricValueIfExists(metric, metricType)
			if ok {
				loggingResponse.Write([]byte(metricValue))
				loggingResponse.WriteHeader(http.StatusOK)
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
				loggingResponse.WriteHeader(http.StatusNotFound)
			}
		}
	}

}
