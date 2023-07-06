package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/filestorage"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

var memStorage memstorage.MemStorage
var Chi chi.Router
var Consumer *filestorage.Consumer
var Producer *filestorage.Producer

func Initialize(isRestore bool, filePath string) {
	Chi = chi.NewRouter()

	if memStorage.GaugeStorage == nil || memStorage.CounterStorage == nil {
		memStorage.GaugeStorage = make(map[string]memstorage.Gauge)
		memStorage.CounterStorage = make(map[string]memstorage.Counter)

	}

	if isRestore && filePath != "" {

		file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)

		if err != nil {
			log.Println(err)
		}

		defer file.Close()
		var metric *servermetrics.Metrics

		scanner := bufio.NewScanner(file)
		fmt.Println(scanner.Text())

		for scanner.Scan() {

			line := scanner.Text()

			err := json.Unmarshal([]byte(line), &metric)

			if err != nil {
				log.Println("Error of read")
				continue
			}

			if metric.MType == memstorage.GaugeType {
				memStorage.Save(metric.ID, metric.Value)
			} else if metric.MType == memstorage.CounterType {
				memStorage.Save(metric.ID, metric.Delta)
			}:q
			log.Println(memStorage.PrintAll())
		}

	} else {
		memStorage.New()
	}

	var err error
	time.Sleep(60)
	Producer, err = filestorage.NewProducer(filePath)

	if err != nil {
		log.Println(err)

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
			memStorage.Save(gaugeNameInDictionary, val)
		} else {
			memStorage.Save(sentMetric, val)
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
			currentValue, _ := memStorage.Get(counterNameInDictionary, servermetrics.CounterType)
			currentCounterValue, _ := servermetrics.StringToCounter(currentValue)

			counterValue := currentCounterValue + val

			memStorage.Save(counterNameInDictionary, counterValue)
			loggingResponse.WriteHeader(http.StatusOK)
		} else {
			currentValue, _ := memStorage.Get(sentMetric, servermetrics.CounterType)
			currentCounterValue, _ := servermetrics.StringToCounter(currentValue)

			counterValue := currentCounterValue + val

			memStorage.Save(sentMetric, counterValue)
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
				memStorage.Save(sentMetricName, storageValue)
			}
		} else if sentMetricType == servermetrics.CounterType {
			if sentMetricsCounterValue == nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
				loggingResponse.WriteHeader(http.StatusBadRequest)
				return
			}
			//Получили значение в сторадже
			currentValue, _ := memStorage.Get(sentMetricName, servermetrics.CounterType)
			currentValueCounter, _ := servermetrics.StringToCounter(currentValue)

			sentValue := servermetrics.Counter(*sentMetricsCounterValue)

			newCounterValue := currentValueCounter + sentValue

			//Сохранили в мемсторадже
			memStorage.Save(sentMetricName, newCounterValue)
		}

		metric.ID = sentMetricName
		metric.MType = sentMetricType

		//Выдаем полученное в ответе
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := memStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := servermetrics.StringToCounter(metricDelta)

			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
		} else if sentMetricType == servermetrics.GaugeType {
			metricValue, _ := memStorage.Get(sentMetricName, sentMetricType)
			gaugeMetricValue, _ := servermetrics.StringToGauge(metricValue, 64)
			floatMetricValue := float64(gaugeMetricValue)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		reqHeader := req.Header.Get("Content-Type")

		if res.Header().Get("Content-Type") == "" {
			res.Header().Add("Content-Type", reqHeader)
		}

		res.Header().Set("Content-Type", reqHeader)

		res.Write(resp)
		loggingResponse.WriteStatusCode(http.StatusOK)
	}

	reqHeader := req.Header.Get("Content-Type")

	if res.Header().Get("Content-Type") == "" {
		res.Header().Add("Content-Type", reqHeader)
	}

	res.Header().Set("Content-Type", reqHeader)

}

func JSONGetMetricsHandler(res http.ResponseWriter, req *http.Request) {
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
			loggingResponse.WriteHeader(http.StatusBadRequest)
			return
		}

		sentMetricType := sMetrics.MType
		sentMetricName := sMetrics.ID

		metric.ID = sentMetricName
		metric.MType = sentMetricType
		if sentMetricType == servermetrics.CounterType {
			metricDelta, _ := memStorage.Get(sentMetricName, sentMetricType)
			intMetricDelta, _ := strconv.Atoi(metricDelta)
			int64MetricDelta := int64(intMetricDelta)
			metric.Delta = &int64MetricDelta
		} else if sentMetricType == servermetrics.GaugeType {
			metricValue, _ := memStorage.Get(sentMetricName, sentMetricType)
			floatMetricValue, _ := strconv.ParseFloat(metricValue, 64)
			metric.Value = &floatMetricValue
		}

		resp, err := json.Marshal(metric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			loggingResponse.WriteHeader(http.StatusInternalServerError)
			return
		}

		reqHeader := req.Header.Get("Content-Type")

		if res.Header().Get("Content-Type") == "" {
			res.Header().Add("Content-Type", reqHeader)
		}

		res.Header().Set("Content-Type", reqHeader)
		res.Write(resp)

		loggingResponse.WriteStatusCode(http.StatusOK)
		return
	}
}
