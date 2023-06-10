package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	var gm MetricsGauge
	gm.New()

	for {
		gm.Renew()
		time.Sleep(pollInterval * time.Second)
		SendRuntimeMetrics(&gm)
		time.Sleep(reportInterval * time.Second)
	}

}

func SendRuntimeMetrics(m *MetricsGauge) {
	url := "http://localhost:8080"
	for k, v := range m.runtimeMetrics {
		url = url + "/update/gauge/" + k + "/" + strconv.FormatFloat(v, 'f', 1, 64)
		response, err := http.Post(url, "text/plain", nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(response.Status)
		log.Print(response.Status)

	}

}
