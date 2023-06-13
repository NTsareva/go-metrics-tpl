package main

import (
	"flag"
	"fmt"
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var serverParams struct {
	address string
}

func init() {
	flag.StringVar(&serverParams.address, "a", "localhost:8080", "input address")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}
}

func MetricsRouter() chi.Router {
	r := chi.NewRouter()
	var ms memstorage.MemStorage
	ms.New()

	r.Post("/update", handlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", handlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/{type}", handlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", handlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", handlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", handlers.NoMetricValueHandler) //Done
	//Надо бы вынести в отдельную функцию, но пока не разобралась, как мемсторадж использовать в параметрах
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

			ms.GaugeStorage[sentMetric] = val
		}

		if sentMetricType == memstorage.Counter {
			val, e := memstorage.StringToCounter(sentMetricValue)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}
			currentValue := ms.CounterStorage[sentMetric]

			ms.CounterStorage[sentMetric] = currentValue + val
		}

	})

	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		body, _ := ms.PrintAll()
		io.WriteString(res, body)
	})

	r.Get("/value/{type}/{metric}", func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "type"))
		if metricType != memstorage.Counter && metricType != memstorage.Gauge {
			http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		}

		metric := strings.ToLower(chi.URLParam(req, "metric"))
		if metricType == memstorage.Counter {
			metricValue, ok := ms.CounterStorage[metric]
			if ok {
				io.WriteString(res, memstorage.CounterToString(metricValue)+"   ")
			} else {
				http.Error(res, "no such metric", http.StatusNotFound)
			}
		}
		if metricType == memstorage.Gauge {
			metricValue, ok := ms.GaugeStorage[metric]
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
	log.Fatal(http.ListenAndServe(serverParams.address, MetricsRouter()))
	fmt.Println(serverParams.address)
}
