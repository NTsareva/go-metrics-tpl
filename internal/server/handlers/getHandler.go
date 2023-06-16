package handlers

import (
	"io"
	"net/http"

	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
)

func AllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	storage := memstorage.MemStorage{}
	body, _ := storage.PrintAll()
	body += "wow"
	io.WriteString(res, body)
}
