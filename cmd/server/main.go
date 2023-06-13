package main

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/MemStorage"
	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	var storage MemStorage.MemStorage
	storage.New()

	r.Post("/update", handlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/", handlers.NoMetricsTypeHandler)                //Done
	r.Post("/update/{type}", handlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", handlers.NoMetricValueHandler) //Done
	r.Post("/update/{type}/{metric}/{value}", handlers.MetricsValueHandler)

	r.Get("/", handlers.AllMetricsHandler) //GET all metrics
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		body, _ := storage.PrintAll()
		io.WriteString(res, body)
		return
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
