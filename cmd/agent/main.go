package main

import (
	agentMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"time"
)

const URL = "http://localhost:8080"
const pollInterval = 2
const reportInterval = 10

func main() {
	var mg agentMetrics.MetricsGauge
	var mc agentMetrics.MetricsCount
	mg.New()
	mc.New()

	//Костыль, запуталась в параллелизме
	for {
		for i := 1; i < reportInterval/pollInterval; i++ {
			metricsRenew(mg, mc)
			time.Sleep(pollInterval * time.Second)
		}
		metricsRenew(mg, mc)
		SendRuntimeMetrics(&mg, &mc)
	}

}

func metricsRenew(mg agentMetrics.MetricsGauge, mc agentMetrics.MetricsCount) {
	log.Println(1)
	mg.Renew()
	mc.Renew()
}

func SendRuntimeMetrics(m *agentMetrics.MetricsGauge, cm *agentMetrics.MetricsCount) {
	client := resty.New()

	client.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)

	client.
		SetHeader("Content-Type", "plain/text").
		SetHeader("Accept", "plain/text")

	for k, v := range m.RuntimeMetrics {
		url := URL + "/update/gauge/" + k + "/" + strconv.FormatFloat(float64(v), 'f', 64, 64)

		response, err := client.R().
			Post(url)

		if err != nil {
			panic(err)
		}

		log.Println(response)
		log.Println(url)
	}

	for k, v := range cm.RuntimeMetrics {
		url := URL + "/update/counter/" + k + "/" + strconv.Itoa(int(v))

		response, err := client.R().
			Post(url)

		if err != nil {
			panic(err)
		}

		log.Println(response)
		log.Println(url)
	}
}
