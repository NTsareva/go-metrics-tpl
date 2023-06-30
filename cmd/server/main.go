package main

import (
	"flag"
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
	handlers.Initialize()

	r := handlers.Chi

	r.Use(handlers.WithLogging)

	r.Post("/update", handlers.NoMetricsTypeHandler)                  //Done
	r.Post("/update/", handlers.JSONUpdateMetricsHandler)             //Done
	r.Post("/update/{type}", handlers.NoMetricsHandler)               //Done
	r.Post("/update/{type}/", handlers.NoMetricsHandler)              //Done
	r.Post("/update/{type}/{metric}", handlers.NoMetricValueHandler)  //Done
	r.Post("/update/{type}/{metric}/", handlers.NoMetricValueHandler) //Done
	r.Post("/update/{type}/{metric}/{value}", handlers.MetricsHandler)

	r.Get("/", handlers.AllMetricsHandler)
	r.Get("/value/{type}/{metric}", handlers.MetricHandler)
	r.Post("/value/", handlers.JSONGetMetricsHandler)

	return r
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Print(err, "#1")
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
		sugar.Fatalf(err.Error(), "event", "start server")
	}
}
