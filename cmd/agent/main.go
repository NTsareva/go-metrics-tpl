package main

import (
	"flag"
	agentMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"time"
)

var agentParams struct {
	address        string
	pollInterval   int
	reportInterval int
}

func main() {
	flag.StringVar(&agentParams.address, "a", "localhost:8080", "input address")
	flag.IntVar(&agentParams.pollInterval, "p", 2, "input poll interval")
	flag.IntVar(&agentParams.reportInterval, "r", 10, "input report interval")
	// делаем разбор командной строки
	flag.Parse()

	var mg agentMetrics.MetricsGauge
	var mc agentMetrics.MetricsCount
	mg.New()
	mc.New()

	reportInterval := agentParams.reportInterval
	pollInterval := agentParams.pollInterval

	tempReportInterval := 0
	tempPollInterval := 0

	//Костыль, запуталась в параллелизме
	for {
		time.Sleep(1 * time.Second)
		tempPollInterval += 1
		tempReportInterval += 1

		if tempPollInterval == pollInterval {
			metricsRenew(mg, mc)
			tempPollInterval = 0
		}
		if tempReportInterval == reportInterval {
			SendRuntimeMetrics(&mg, &mc)
			tempReportInterval = 0
		}
	}
}

func metricsRenew(mg agentMetrics.MetricsGauge, mc agentMetrics.MetricsCount) {
	log.Println(1)
	mg.Renew()
	mc.Renew()
}

func SendRuntimeMetrics(m *agentMetrics.MetricsGauge, cm *agentMetrics.MetricsCount) {
	client := resty.New()

	agentURL := agentParams.address

	client.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)

	client.
		SetHeader("Content-Type", "plain/text").
		SetHeader("Accept", "plain/text")

	for k, v := range m.RuntimeMetrics {
		url := "http://" + agentURL + "/update/gauge/" + k + "/" + strconv.FormatFloat(float64(v), 'f', 64, 64)

		response, err := client.R().
			Post(url)

		if err != nil {
			panic(err)
		}

		log.Println(response)
		log.Println(url)
	}

	for k, v := range cm.RuntimeMetrics {
		url := "http://" + agentURL + "/update/counter/" + k + "/" + strconv.Itoa(int(v))

		response, err := client.R().
			Post(url)

		if err != nil {
			panic(err)
		}

		log.Println(response)
		log.Println(url)
	}
}
