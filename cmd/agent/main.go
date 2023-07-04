package main

import (
	"flag"
	"fmt"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	agentConfig "github.com/NTsareva/go-metrics-tpl.git/cmd/agent/config"
	agentMetrics "github.com/NTsareva/go-metrics-tpl.git/internal/agent/metrics"
)

type MetricsBody struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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
		SetRetryCount(20).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(180 * time.Second)

	client.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Encoding", "gzip")

	url := fmt.Sprintf("http://%s/update/", agentURL)

	postClient := client.R()

	for k, v := range metricsGauge.RuntimeMetrics {
		deltaValue := int64(0)
		floatGaugeValue, _ := strconv.ParseFloat(agentMetrics.GaugeToString(v), 64)
		requestBody := MetricsBody{
			ID:    k,
			MType: servermetrics.GaugeType,
			Delta: &deltaValue,
			Value: &floatGaugeValue,
		}
		response, err := postClient.SetBody(requestBody).Post(url)

		if err != nil {
			log.Print(err)
			continue
		}

		log.Println(response)
		log.Println(url)

	}

	for k, v := range metricsCount.RuntimeMetrics {
		deltaValue, _ := strconv.Atoi(agentMetrics.CounterToString(v))
		int64DeltaValue := int64(deltaValue)
		floatGaugeValue := 0.0
		requestBody := MetricsBody{
			ID:    k,
			MType: servermetrics.CounterType,
			Delta: &int64DeltaValue,
			Value: &floatGaugeValue,
		}
		response, err := postClient.SetBody(requestBody).Post(url)

		if err != nil {

			log.Print(err)
			continue
		}

		log.Println(response)
		log.Println(url)
	}
}
