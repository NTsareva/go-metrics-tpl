package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
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
	loggingResponse.WriteHeader(http.StatusBadRequest)
}

func (serverHandlers *SeverHandlers) NoMetricsHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no metrics in request ", http.StatusNotFound)
	loggingResponse.WriteHeader(http.StatusNotFound)
}

func (serverHandlers *SeverHandlers) NoMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no value of metrics ", http.StatusBadRequest)
	loggingResponse.WriteHeader(http.StatusBadRequest)

}

func (serverHandlers *SeverHandlers) MetricsHandler(res http.ResponseWriter, req *http.Request) {
	sentMetricType := strings.ToLower(chi.URLParam(req, "type"))

	if !servermetrics.IfHasCorrectType(sentMetricType) {
		http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
		loggingResponse.WriteHeader(http.StatusBadRequest)
	}

	sentMetric := strings.ToLower(chi.URLParam(req, "metric"))

	sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))
	if sentMetricType == servermetrics.GaugeType {

		val, e := servermetrics.StringToGauge(sentMetricValue, 64)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
		}

		serverHandlers.MemStorage.Save(sentMetric, val)
		loggingResponse.WriteHeader(http.StatusOK)
	}

	if sentMetricType == servermetrics.CounterType {
		val, e := servermetrics.StringToCounter(sentMetricValue)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
		}

		currentValue, _ := serverHandlers.MemStorage.Get(sentMetric, servermetrics.CounterType)
		currentCounterValue, _ := servermetrics.StringToCounter(currentValue)

		counterValue := currentCounterValue + val

		serverHandlers.MemStorage.Save(sentMetric, counterValue)
		loggingResponse.WriteHeader(http.StatusOK)
	}
}

func (serverHandlers *SeverHandlers) JSONUpdateMetricsHandler(res http.ResponseWriter, req *http.Request) {
	metric := servermetrics.Metrics{
		ID:    "0",
		MType: "",
		Delta: nil,
		Value: nil,
	}

	var sMetrics memstorage.Metrics
	var buf bytes.Buffer

	if req.Method == http.MethodPost {
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &sMetrics); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		sentMetricType := sMetrics.MType
		sentMetricsGaugeValue := sMetrics.Value
		sentMetricsCounterValue := sMetrics.Delta
		sentMetricName := sMetrics.ID

		//Распарсили то, что получили от JSON
		//Проверили тип
		//Если тип GAUGE, просто сохранили
		//Если тип Counter, скастовали, взяли адрес, сохранили
		//В теле ответа запросили из стораджа, скастовали

		if !servermetrics.IfHasCorrectType(sentMetricType) {
			http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
		}

		if sentMetricType == servermetrics.GaugeType {
			if sentMetricsGaugeValue == nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
				loggingResponse.WriteHeader(http.StatusBadRequest)
			} else {
				storageValue := servermetrics.Gauge(*sentMetricsGaugeValue)
				serverHandlers.MemStorage.Save(sentMetricName, storageValue)
				loggingResponse.WriteHeader(http.StatusOK)
			}
		}

		if sentMetricType == servermetrics.CounterType {
			if sentMetricsCounterValue == nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
				loggingResponse.WriteHeader(http.StatusBadRequest)
			}
			//Получили значение в сторадже
			currentValue, _ := serverHandlers.MemStorage.Get(sentMetricName, servermetrics.CounterType)
			currentValueCounter, _ := servermetrics.StringToCounter(currentValue)
			sentValue := servermetrics.Counter(*sentMetricsCounterValue)

			newCounterValue := currentValueCounter + sentValue

			//Сохранили в мемсторадже
			serverHandlers.MemStorage.Save(sentMetricName, newCounterValue)
			loggingResponse.WriteHeader(http.StatusOK)
		}

		metric.ID = sentMetricName
		metric.MType = sentMetricType

		//Выдаем полученное в ответе
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := serverHandlers.MemStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := servermetrics.StringToCounter(metricDelta)

			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
			metricValue := 0.0
			metric.Value = &metricValue
		} else if sentMetricType == servermetrics.GaugeType {
			metricDelta := int64(0)
			metric.Delta = &metricDelta
			metricValue, _ := serverHandlers.MemStorage.Get(sentMetricName, sentMetricType)
			gaugeMetricValue, _ := servermetrics.StringToGauge(metricValue, 64)
			floatMetricValue := float64(gaugeMetricValue)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
	}
}

func (serverHandlers *SeverHandlers) JSONGetMetricsHandler(res http.ResponseWriter, req *http.Request) {
	metric := servermetrics.Metrics{
		ID:    "0",
		MType: "",
		Delta: nil,
		Value: nil,
	}

	var sMetrics memstorage.Metrics
	var buf bytes.Buffer

	if req.Method == http.MethodPost {

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &sMetrics); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		sentMetricType := sMetrics.MType
		sentMetricName := sMetrics.ID

		metric.ID = sentMetricName
		metric.MType = sentMetricType
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := serverHandlers.MemStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := strconv.Atoi(metricDelta)
			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
			metricValue := 0.0
			metric.Value = &metricValue
		} else if sentMetricType == servermetrics.GaugeType {
			metricDelta := int64(0)
			metric.Delta = &metricDelta
			metricValue, _ := serverHandlers.MemStorage.Get(sentMetricName, sentMetricType)
			floatMetricValue, _ := strconv.ParseFloat(metricValue, 64)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(resp)
		res.WriteHeader(http.StatusOK)
	}
}
