package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

var MemStorage memstorage.MemStorage

var Chi chi.Router

func Initialize() {
	Chi = chi.NewRouter()

	if MemStorage.GaugeStorage == nil || MemStorage.CounterStorage == nil {
		MemStorage.GaugeStorage = make(map[string]memstorage.Gauge)
		MemStorage.CounterStorage = make(map[string]memstorage.Counter)

		MemStorage.New()
	}
}

func NoMetricsTypeHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no type of metrics ", http.StatusBadRequest)
	loggingResponse.WriteHeader(http.StatusBadRequest)
}

func NoMetricsHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no metrics in request ", http.StatusNotFound)
	loggingResponse.WriteHeader(http.StatusNotFound)
}

func NoMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "no value of metrics ", http.StatusBadRequest)
	loggingResponse.WriteHeader(http.StatusBadRequest)

}

func MetricsHandler(res http.ResponseWriter, req *http.Request) {
	sentMetricType := strings.ToLower(chi.URLParam(req, "type"))

	if !servermetrics.IfHasCorrectType(sentMetricType) {
		http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
		loggingResponse.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: Делать проверку на нахождении в словаре
	sentMetric := strings.ToLower(chi.URLParam(req, "metric"))
	sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))

	if sentMetricType == servermetrics.GaugeType {
		gaugeNamesDictionary := servermetrics.GetGaugeMetricsResponseNames()
		gaugeNameInDictionary, ok := gaugeNamesDictionary[sentMetric]

		val, e := servermetrics.StringToGauge(sentMetricValue, 64)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		if ok {
			MemStorage.Save(gaugeNameInDictionary, val)
		} else {
			MemStorage.Save(sentMetric, val)
		}

		loggingResponse.WriteHeader(http.StatusOK)
	} else if sentMetricType == servermetrics.CounterType {
		counterNamesDictionary := servermetrics.GetCounterMetricsResponseNames()
		counterNameInDictionary, ok := counterNamesDictionary[sentMetric]

		val, e := servermetrics.StringToCounter(sentMetricValue)
		if e != nil {
			http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		if ok {
			currentValue, _ := MemStorage.Get(counterNameInDictionary, servermetrics.CounterType)
			currentCounterValue, _ := servermetrics.StringToCounter(currentValue)

			counterValue := currentCounterValue + val

			MemStorage.Save(counterNameInDictionary, counterValue)
			loggingResponse.WriteHeader(http.StatusOK)
		} else {
			currentValue, _ := MemStorage.Get(sentMetric, servermetrics.CounterType)
			currentCounterValue, _ := servermetrics.StringToCounter(currentValue)

			counterValue := currentCounterValue + val

			MemStorage.Save(sentMetric, counterValue)
			loggingResponse.WriteHeader(http.StatusOK)
		}
	}
}

func JSONUpdateMetricsHandler(res http.ResponseWriter, req *http.Request) {
	metric := servermetrics.Metrics{
		ID:    "0",
		MType: "",
		Delta: nil,
		Value: nil,
	}

	var sMetrics memstorage.Metrics
	var buf bytes.Buffer

	res.Header().Add("Content-Type", "application/json")
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Connection", "Keep-Alive")

	if req.Method == http.MethodPost {
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &sMetrics); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			loggingResponse.WriteHeader(http.StatusBadRequest)
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
			return
		}

		if sentMetricType == servermetrics.GaugeType {
			if sentMetricsGaugeValue == nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
				loggingResponse.WriteHeader(http.StatusBadRequest)
				return
			} else {
				storageValue := servermetrics.Gauge(*sentMetricsGaugeValue)
				MemStorage.Save(sentMetricName, storageValue)
			}
		} else if sentMetricType == servermetrics.CounterType {
			if sentMetricsCounterValue == nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
				loggingResponse.WriteHeader(http.StatusBadRequest)
				return
			}
			//Получили значение в сторадже
			currentValue, _ := MemStorage.Get(sentMetricName, servermetrics.CounterType)
			currentValueCounter, _ := servermetrics.StringToCounter(currentValue)

			sentValue := servermetrics.Counter(*sentMetricsCounterValue)

			newCounterValue := currentValueCounter + sentValue

			//Сохранили в мемсторадже
			MemStorage.Save(sentMetricName, newCounterValue)
		}

		metric.ID = sentMetricName
		metric.MType = sentMetricType

		//Выдаем полученное в ответе
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := MemStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := servermetrics.StringToCounter(metricDelta)

			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
		} else if sentMetricType == servermetrics.GaugeType {
			metricValue, _ := MemStorage.Get(sentMetricName, sentMetricType)
			gaugeMetricValue, _ := servermetrics.StringToGauge(metricValue, 64)
			floatMetricValue := float64(gaugeMetricValue)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf(string(resp))

		res.Header().Set("Content-Type", "application/json")
		res.Write(resp)
		loggingResponse.WriteStatusCode(http.StatusOK)
	}
}

func JSONGetMetricsHandler(res http.ResponseWriter, req *http.Request) {
	metric := servermetrics.Metrics{
		ID:    "0",
		MType: "",
		Delta: nil,
		Value: nil,
	}

	res.Header().Add("Content-Type", "application/json")
	res.Header().Set("Connection", "Keep-Alive")

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
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		sentMetricType := sMetrics.MType
		sentMetricName := sMetrics.ID

		metric.ID = sentMetricName
		metric.MType = sentMetricType
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := MemStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := strconv.Atoi(metricDelta)
			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
		} else if sentMetricType == servermetrics.GaugeType {
			metricValue, _ := MemStorage.Get(sentMetricName, sentMetricType)
			floatMetricValue, _ := strconv.ParseFloat(metricValue, 64)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			loggingResponse.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf(string(resp))

		res.Header().Set("Content-Type", "application/json")
		res.Write(resp)

		loggingResponse.WriteStatusCode(http.StatusOK)
		return
	}
}
