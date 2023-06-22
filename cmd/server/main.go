package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"

	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
)

var serverParams struct {
	address string
}

func MetricsRouter() chi.Router {
	var serverHandlers handlers.SeverHandlers
	serverHandlers.New()
	r := serverHandlers.Chi

	r.Post("/update", serverHandlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", serverHandlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/{type}", serverHandlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", serverHandlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", serverHandlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", serverHandlers.NoMetricValueHandler) //Done
	r.Post("/update/{type}/{metric}/{value}", serverHandlers.MetricsHandler)

	r.Get("/", serverHandlers.AllMetricsHandler)
	r.Get("/value/{type}/{metric}", serverHandlers.MetricHandler)

	return r
}

func main() {
	flag.StringVar(&serverParams.address, "a", "localhost:8080", "input address")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}

	log.Print(http.ListenAndServe(serverParams.address, MetricsRouter()))
	fmt.Println(serverParams.address)
}
