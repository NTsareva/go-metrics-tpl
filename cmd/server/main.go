package main

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/handlers"
	"net/http"
)

type gauge float64
type counter int64

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.PostHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
