package handlers

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"io"
	"net/http"
)

func AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	storage := memstorage.MemStorage{}
	body, _ := storage.PrintAll()
	body += "wow"
	io.WriteString(res, body)
}
