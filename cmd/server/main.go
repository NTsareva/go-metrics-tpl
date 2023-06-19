package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

var serverParams struct {
	address string
}

func MetricsRouter() chi.Router {
	r := chi.NewRouter()
	var memStorage memstorage.MemStorage
	memStorage.New()

	r.Post("/update", handlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", handlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/{type}", handlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", handlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", handlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", handlers.NoMetricValueHandler) //Done
	//Надо бы вынести в отдельную функцию c глобальным мемстораджем. но мне это не нравится
	r.Post("/update/{type}/{metric}/{value}", func(res http.ResponseWriter, req *http.Request) {
		//Проверка, что тип корректный
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

			memStorage.GaugeStorage[sentMetric] = memstorage.Gauge(val)
		}

		if sentMetricType == servermetrics.CounterType {
			val, e := servermetrics.StringToCounter(sentMetricValue)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}
			currentValue := memStorage.CounterStorage[sentMetric]

			memStorage.CounterStorage[sentMetric] = currentValue + memstorage.Counter(val)
		}

	})

	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		body, _ := memStorage.PrintAll()
		io.WriteString(res, body)
	})

	r.Get("/value/{type}/{metric}", func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		if metricType != servermetrics.CounterType && metricType != servermetrics.GaugeType {
			http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		}

		metric := strings.ToLower(chi.URLParam(req, "metric"))
		if metricType == servermetrics.CounterType {
			metricValue, ok := memStorage.CounterStorage[metric]
			if ok {
				io.WriteString(res, servermetrics.CounterToString(servermetrics.Counter(metricValue))+"   ")
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
			}
		}
		if metricType == servermetrics.GaugeType {
			metricValue, ok := memStorage.GaugeStorage[metric]
			if ok {
				res.Write([]byte(servermetrics.GaugeToString(servermetrics.Gauge(metricValue))))
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
			}
		}
	})

	return r
}

func main() {
	flag.StringVar(&serverParams.address, "a", "localhost:8080", "input address")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}
	log.Fatal(http.ListenAndServe(serverParams.address, MetricsRouter()))
	fmt.Println(serverParams.address)
}
