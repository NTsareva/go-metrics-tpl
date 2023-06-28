package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

var serverParams struct {
	address string
}

func MetricsRouter() chi.Router {
	var serverHandlers handlers.SeverHandlers
	serverHandlers.New()
	r := serverHandlers.Chi

	r.Use(handlers.WithLogging)

	r.Post("/update", serverHandlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", serverHandlers.NoMetricsTypeHandler)                 //Done
	r.Post("/update/", serverHandlers.JSONUpdateMetricsHandler)             //Done
	r.Post("/update/{type}", serverHandlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", serverHandlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", serverHandlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", serverHandlers.NoMetricValueHandler) //Done
	r.Post("/update/{type}/{metric}/{value}", serverHandlers.MetricsHandler)

	r.Get("/", serverHandlers.AllMetricsHandler)
	r.Get("/value/{type}/{metric}", serverHandlers.MetricHandler)
	r.Post("/value/", serverHandlers.JSONGetMetricsHandler)

	return r
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Print(err)
	}

	addr := "localhost:8080"
	defer logger.Sync()

	sugar := *logger.Sugar()

	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	flag.StringVar(&serverParams.address, "a", addr, "input address")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}

	if err := http.ListenAndServe(serverParams.address, MetricsRouter()); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

	fmt.Println(serverParams.address)
}
