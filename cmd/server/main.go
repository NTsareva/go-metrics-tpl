package main

import (
	"flag"
	"fmt"
	"github.com/NTsareva/go-metrics-tpl.git/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var serverParams struct {
	address         string
	storeInterval   int64
	fileStoragePath string
	ifRestore       bool
}

func MetricsRouter() chi.Router {
	handlers.Initialize(serverParams.ifRestore, serverParams.fileStoragePath)

	r := handlers.Chi
	r.Use(handlers.WithLogging)
	r.Use(handlers.WithGzipActions)

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
		log.Println(err)
	}

	addr := "localhost:8080"
	defer logger.Sync()

	sugar := *logger.Sugar()

	sugar.Infow(
		"Starting server",
		"addr", addr,
	)

	flag.StringVar(&serverParams.address, "a", addr, "input address")
	flag.Int64Var(&serverParams.storeInterval, "i", 300, "store interval")
	flag.StringVar(&serverParams.fileStoragePath, "f", "tmp/metrics-db.json", "save file path")
	flag.BoolVar(&serverParams.ifRestore, "r", true, "if should restore")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverParams.address = envRunAddr
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		serverParams.storeInterval, _ = strconv.ParseInt(envStoreInterval, 10, 64)
	}

	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		serverParams.address = envStoragePath
	}

	if envIfRestore := os.Getenv("RESTORE"); envIfRestore != "" {
		serverParams.address = envIfRestore
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			time.Sleep(time.Duration(serverParams.storeInterval) * time.Second)
			handlers.WriteMemstorageToFile(serverParams.fileStoragePath)
		}
	}()

	if err := http.ListenAndServe(serverParams.address, MetricsRouter()); err != nil {
		sugar.Fatalf(err.Error(), "event", "start server")
	}

	sig := <-signalCh
	handlers.WriteMemstorageToFile(serverParams.fileStoragePath)
	fmt.Println("Resieved sig", sig)
	os.Exit(0)

}
