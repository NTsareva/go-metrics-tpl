package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

type SeverHandlers struct {
	Chi        chi.Router
	MemStorage memstorage.MemStorage
}

func (serverHandlers *SeverHandlers) New() {
	serverHandlers.Chi = chi.NewRouter()
	serverHandlers.MemStorage.New()
}

func (serverHandlers *SeverHandlers) NoMetricsTypeHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no type of metrics ", http.StatusBadRequest)
}

func (serverHandlers *SeverHandlers) NoMetricsHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no metrics in request ", http.StatusNotFound)
}

func (serverHandlers *SeverHandlers) NoMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no value of metrics ", http.StatusBadRequest)
}

func (serverHandlers *SeverHandlers) MetricsHandler(res http.ResponseWriter, req *http.Request) {
	sentMetricType := strings.ToLower(chi.URLParam(req, "type"))

	if !servermetrics.IfHasCorrestType(sentMetricType) {
		http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
	}

	sentMetric := strings.ToLower(chi.URLParam(req, "metric"))

	sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))
	if sentMetricType == servermetrics.GaugeType {

		val, e := servermetrics.StringToGauge(sentMetricValue, 64)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
		}

		serverHandlers.MemStorage.GaugeStorage[sentMetric] = memstorage.Gauge(val)
	}

	if sentMetricType == servermetrics.CounterType {
		val, e := servermetrics.StringToCounter(sentMetricValue)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
		}
		currentValue := serverHandlers.MemStorage.CounterStorage[sentMetric]

		serverHandlers.MemStorage.CounterStorage[sentMetric] = currentValue + memstorage.Counter(val)
	}

}
