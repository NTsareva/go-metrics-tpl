package main

import (
	"net/http"
)

type gauge float64
type counter int64

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", PostHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
