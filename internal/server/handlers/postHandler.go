package handlers

import (
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func NoMetricsTypeHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no type of metrics ", http.StatusBadRequest)
	return
}

func NoMetricsHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no metrics in request ", http.StatusNotFound)
	return
}

func NoMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no value of metrics ", http.StatusBadRequest)
	return
}

func MetricsValueHandler(res http.ResponseWriter, req *http.Request) {
	//Type check
	sentMetricType := strings.ToLower(chi.URLParam(req, "type"))

	if !servermetrics.IfHasCorrestType(sentMetricType) {
		http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
	}
	//Проверяем, что метрика попадает в список
	//sentMetric := strings.ToLower(chi.URLParam(req, "metric"))

	//Проверяем, что значение соответствует типу
	sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))
	if sentMetricType == servermetrics.Gauge {
		_, e := servermetrics.StringToGauge(sentMetricValue, 64)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			return
		}
	}

	//TODO: когда дойдем до каунтеров, сделать проверку, что тип Counter

	//Сохраняем метрику

	return
}
