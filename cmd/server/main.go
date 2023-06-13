package main

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/MemStorage"
	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	r := chi.NewRouter()
	var ms MemStorage.MemStorage
	ms.New()

	r.Post("/update", handlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", handlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/{type}", handlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", handlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", handlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", handlers.NoMetricValueHandler) //Done
	r.Post("/update/{type}/{metric}/{value}", func(res http.ResponseWriter, req *http.Request) {
		//Проверка, что тип корректный
		sentMetricType := strings.ToLower(chi.URLParam(req, "type"))

		if !servermetrics.IfHasCorrestType(sentMetricType) {
			http.Error(res, "incorrect type of metrics "+sentMetricType+" ", http.StatusBadRequest)
		}

		sentMetric := strings.ToLower(chi.URLParam(req, "metric"))

		//Проверяем, что значение соответствует типу
		sentMetricValue := strings.ToLower(chi.URLParam(req, "value"))
		if sentMetricType == MemStorage.Gauge {
			//_, okGauge := ms.GaugeStorage[sentMetric]
			//
			//if !okGauge {
			//	http.Error(res, "unknown type of metrics "+sentMetric, http.StatusNotFound)
			//}

			val, e := MemStorage.StringToGauge(sentMetricValue, 64)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}

			ms.GaugeStorage[sentMetric] = val
		}

		//TODO: когда дойдем до каунтеров, сделать проверку, что тип Counter

		if sentMetricType == MemStorage.Counter {
			val, e := MemStorage.StringToCounter(sentMetricValue)
			if e != nil {
				http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
			}

			ms.CounterStorage[sentMetric] = val
		}
	})

	r.Get("/", handlers.AllMetricsHandler) //GET all metrics
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		body, _ := ms.PrintAll()
		io.WriteString(res, body)
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
