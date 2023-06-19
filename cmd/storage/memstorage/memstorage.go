package memstorage

import (
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

type Gauge servermetrics.Gauge
type Counter servermetrics.Counter

type MemStorage struct {
	GaugeStorage   map[string]Gauge
	CounterStorage map[string]Counter
}

func (memStorage *MemStorage) New() {
	memStorage.GaugeStorage = make(map[string]Gauge)
	memStorage.CounterStorage = make(map[string]Counter)

	metricsGauge := servermetrics.MetricsGauge{}
	metricsGauge.New()

	for k, v := range metricsGauge.RuntimeMetrics {
		memStorage.GaugeStorage[k] = Gauge(v)
	}

	metricsCounter := servermetrics.MetricsCount{}
	metricsCounter.New()

	for k, v := range metricsCounter.RuntimeMetrics {
		memStorage.GaugeStorage[k] = Gauge(v)
	}
}

func (memStorage *MemStorage) SaveGauge(metrics string, value Gauge) error {
	memStorage.GaugeStorage[metrics] = value
	return nil
}

func (memStorage *MemStorage) SaveCounter(metrics string, value Counter) error {
	memStorage.CounterStorage[metrics] = value
	return nil
}

func (memStorage *MemStorage) Remove(metrics string) error {
	_, okGauge := memStorage.GaugeStorage[metrics]
	_, okCounter := memStorage.CounterStorage[metrics]

	if okGauge {
		delete(memStorage.GaugeStorage, metrics)
	}

	if okCounter {
		delete(memStorage.CounterStorage, metrics)
	}

	return nil
}

func (memStorage *MemStorage) PrintAll() (string, error) {
	outerString := ""

	metricsGauge := memStorage.GaugeStorage
	for k, v := range metricsGauge {
		outerString = outerString + k + ": " + servermetrics.GaugeToString(servermetrics.Gauge(v)) + "\n"
	}

	metricsCounter := memStorage.CounterStorage
	for k, v := range metricsCounter {
		outerString = outerString + k + ": " + servermetrics.CounterToString(servermetrics.Counter(v)) + "\n"
	}

	return outerString, nil
}
