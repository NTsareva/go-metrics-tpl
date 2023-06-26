package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

func (serverHandlers *SeverHandlers) AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	body, _ := serverHandlers.MemStorage.PrintAll()
	io.WriteString(res, body)
}

func (serverHandlers *SeverHandlers) MetricHandler(res http.ResponseWriter, req *http.Request) {

	metricType := strings.ToLower(chi.URLParam(req, "type"))
	if metricType != servermetrics.CounterType && metricType != servermetrics.GaugeType {
		http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		loggingResponse.WriteHeader(http.StatusNotFound)
	}

	metric := strings.ToLower(chi.URLParam(req, "metric"))
	if metricType == servermetrics.CounterType {
		metricValue, ok := serverHandlers.MemStorage.CounterStorage[metric]
		if ok {
			io.WriteString(res, servermetrics.CounterToString(servermetrics.Counter(metricValue))+"   ")
			loggingResponse.WriteHeader(http.StatusOK)
		} else {
			http.Error(res, "no such metric", http.StatusNotFound)
			loggingResponse.WriteHeader(http.StatusNotFound)
		}
	}
	if metricType == servermetrics.GaugeType {
		metricValue, ok := serverHandlers.MemStorage.GaugeStorage[metric]
		if ok {
			loggingResponse.Write([]byte(servermetrics.GaugeToString(servermetrics.Gauge(metricValue))))
			loggingResponse.WriteHeader(http.StatusOK)
		} else {
			http.Error(res, "no such metric", http.StatusNotFound)
			loggingResponse.WriteHeader(http.StatusNotFound)
		}
	}
}
