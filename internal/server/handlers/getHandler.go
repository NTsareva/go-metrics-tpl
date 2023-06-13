package handlers

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/MemStorage"
	"io"
	"net/http"
)

func AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	storage := MemStorage.MemStorage{}
	body, _ := storage.PrintAll()
	body += "wow"
	io.WriteString(res, body)
	return
}
