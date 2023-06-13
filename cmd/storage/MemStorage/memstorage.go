package MemStorage

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"strconv"
)

type gauge storage.Gauge
type counter storage.Counter

func GaugeToStringWithError(gv gauge) (string, error) {
	value := strconv.FormatFloat(float64(gv), 'f', 1, 64)
	return value, nil
}

func GaugeToString(gv gauge) string {
	value := strconv.FormatFloat(float64(gv), 'f', 1, 64)
	return value
}

func CounterToString(cv counter) string {
	value := strconv.Itoa(int(cv))
	return value
}

// Хотела тут использовать дженерики, но запуталась
type MemStorage struct {
	gaugeStorage   map[string]gauge
	counterStorage map[string]counter
}

func (ms *MemStorage) New() {
	ms.gaugeStorage = make(map[string]gauge)
	ms.counterStorage = make(map[string]counter)

	metricsGauge := servermetrics.MetricsGauge{}
	metricsGauge.New()

	for k, v := range metricsGauge.RuntimeMetrics {
		ms.gaugeStorage[k] = gauge(v)
	}

	metricsCounter := servermetrics.MetricsCount{}
	metricsCounter.New()

	for k, v := range metricsCounter.RuntimeMetrics {
		ms.gaugeStorage[k] = gauge(v)
	}
}

func (ms *MemStorage) SaveGauge(metrics string, value gauge) error {
	ms.gaugeStorage[metrics] = value
	return nil
}

func (ms *MemStorage) SaveCounter(metrics string, value counter) error {
	ms.counterStorage[metrics] = value
	return nil
}

func (ms *MemStorage) Remove(metrics string) error {
	_, okGauge := ms.gaugeStorage[metrics]
	_, okCounter := ms.counterStorage[metrics]

	if okGauge {
		delete(ms.gaugeStorage, metrics)
	}

	if okCounter {
		delete(ms.counterStorage, metrics)
	}

	return nil
}

func (ms *MemStorage) IfExists(metrics string, value counter) bool {
	_, okGauge := ms.gaugeStorage[metrics]
	_, okCounter := ms.counterStorage[metrics]

	return okGauge || okCounter
}

func (ms *MemStorage) PrintAll() (string, error) {
	outerString := ""

	metricsGauge := ms.gaugeStorage
	for k, v := range metricsGauge {
		outerString = outerString + k + ": " + GaugeToString(v) + "\n"
	}

	metricsCounter := ms.counterStorage
	for k, v := range metricsCounter {
		outerString = outerString + k + ": " + CounterToString(v) + "\n"
	}

	return outerString, nil
}
