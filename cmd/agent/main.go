package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	agentConfig "github.com/NTsareva/go-metrics-tpl.git/cmd/agent/config"
	agentMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/agent/metrics"
)

func main() {
	flag.StringVar(&agentConfig.AgentParams.Address, "a", "localhost:8080", "input address")
	flag.IntVar(&agentConfig.AgentParams.PollInterval, "p", 2, "input poll interval")
	flag.IntVar(&agentConfig.AgentParams.ReportInterval, "r", 10, "input report interval")

	flag.Parse()

	if addressFromEnv := os.Getenv("ADDRESS"); addressFromEnv != "" {
		agentConfig.AgentParams.Address = addressFromEnv
	}

	if reportIntervalFromEnv := os.Getenv("REPORT_INTERVAL"); reportIntervalFromEnv != "" {
		agentConfig.AgentParams.ReportInterval, _ = strconv.Atoi(reportIntervalFromEnv)
	}

	if pollIntervalFromEnv := os.Getenv("POLL_INTERVAL"); pollIntervalFromEnv != "" {
		agentConfig.AgentParams.PollInterval, _ = strconv.Atoi(pollIntervalFromEnv)
	}

	var metricsGauge agentMetrics.MetricsGauge
	var metricsCount agentMetrics.MetricsCount
	metricsGauge.New()
	metricsCount.New()

	reportInterval := agentConfig.AgentParams.ReportInterval
	pollInterval := agentConfig.AgentParams.PollInterval

	tempReportInterval := 0
	tempPollInterval := 0

	for {
		time.Sleep(1 * time.Second)
		tempPollInterval += 1
		tempReportInterval += 1

		if tempPollInterval == pollInterval {
			metricsRenew(metricsGauge, metricsCount)
			tempPollInterval = 0
		}
		if tempReportInterval == reportInterval {
			sendRuntimeMetrics(&metricsGauge, &metricsCount)
			tempReportInterval = 0
		}
	}
}

func metricsRenew(metricsGauge agentMetrics.MetricsGauge, metricsCount agentMetrics.MetricsCount) {
	metricsGauge.Renew()
	metricsCount.Renew()
}

func sendRuntimeMetrics(metricsGauge *agentMetrics.MetricsGauge, metricsCount *agentMetrics.MetricsCount) {
	client := resty.New()
	agentURL := agentConfig.AgentParams.Address

	client.
		SetRetryCount(3).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)

	client.
		SetHeader("Content-Type", "plain/text").
		SetHeader("Accept", "plain/text")

	postClient := client.R()

	for k, v := range metricsGauge.RuntimeMetrics {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%s", agentURL, k, agentMetrics.GaugeToString(v))

		response, err := postClient.Post(url)

		if err != nil {
			log.Print(err)
		}

		log.Println(response)
		log.Println(url)
	}

	for k, v := range metricsCount.RuntimeMetrics {
		url := fmt.Sprintf("http://%s/update/counter/%s/%s", agentURL, k, agentMetrics.CounterToString(v))

		response, err := postClient.Post(url)

		if err != nil {
			log.Print(err)
		}

		log.Println(response)
		log.Println(url)
	}
}
