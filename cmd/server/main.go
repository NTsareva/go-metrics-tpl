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

func init() {
	flag.StringVar(&serverParams.address, "a", "localhost:8080", "input address")
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

		//Проверяем, что значение соответствует типу
		sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))
		if sentMetricType == memstorage.Gauge {

			val, e := memstorage.StringToGauge(sentMetricValue, 64)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}

			memStorage.GaugeStorage[sentMetric] = val
		}

		if sentMetricType == memstorage.Counter {
			val, e := memstorage.StringToCounter(sentMetricValue)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}
			currentValue := memStorage.CounterStorage[sentMetric]

			memStorage.CounterStorage[sentMetric] = currentValue + val
		}

	})

	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		body, _ := memStorage.PrintAll()
		io.WriteString(res, body)
	})

	r.Get("/value/{type}/{metric}", func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		if metricType != memstorage.Counter && metricType != memstorage.Gauge {
			http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		}

		metric := strings.ToLower(chi.URLParam(req, "metric"))
		if metricType == memstorage.Counter {
			metricValue, ok := memStorage.CounterStorage[metric]
			if ok {
				io.WriteString(res, memstorage.CounterToString(metricValue)+"   ")
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
			}
		}
		if metricType == memstorage.Gauge {
			metricValue, ok := memStorage.GaugeStorage[metric]
			if ok {
				res.Write([]byte(memstorage.GaugeToString(metricValue)))
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
			}
		}
	})

	return r
}

func main() {
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}
	log.Fatal(http.ListenAndServe(serverParams.address, MetricsRouter()))
	fmt.Println(serverParams.address)
}
