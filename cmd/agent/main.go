package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	var gm MetricsGauge
	gm.New()

	for {
		go gm.Renew()
		SendRuntimeMetrics(&gm)
	}
}

func SendRuntimeMetrics(m *MetricsGauge) {
	time.Sleep(reportInterval * time.Second)

	url := "http://localhost:8080"
	for k, v := range m.runtimeMetrics {
		url = url + "/update/gauge/" + k + "/" + strconv.FormatFloat(v, 'f', 1, 64)
		response, err := http.Post(url, "text/plain", nil)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		log.Println(response)
	}
}
