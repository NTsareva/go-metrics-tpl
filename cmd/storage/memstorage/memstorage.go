package memstorage

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
	"strconv"
)

type gauge storage.Gauge
type counter storage.Counter

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

func StringToGauge(s string, bitSize int) (gauge, error) {
	v, e := strconv.ParseFloat(s, bitSize)
	if e != nil {
		return 0.0, e
	}
	return gauge(v), nil
}

func GaugeToString(gv gauge) string {
	value := strconv.FormatFloat(float64(gv), 'f', 3, 64)
	if value[len(value)-1] == '0' {
		value = strconv.FormatFloat(float64(gv), 'f', 2, 64)
	}

	if value[len(value)-1] == '0' && value[len(value)-2] == '0' {
		value = strconv.FormatFloat(float64(gv), 'f', 1, 64)
	}

	return value
}

func CounterToString(cv counter) string {
	value := strconv.Itoa(int(cv))
	return value
}

func StringToCounter(s string) (counter, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return counter(value), nil
}

// Хотела тут использовать дженерики, но запуталась
type MemStorage struct {
	GaugeStorage   map[string]gauge
	CounterStorage map[string]counter
}

func (ms *MemStorage) New() {
	ms.GaugeStorage = make(map[string]gauge)
	ms.CounterStorage = make(map[string]counter)

	metricsGauge := servermetrics.MetricsGauge{}
	metricsGauge.New()

	for k, v := range metricsGauge.RuntimeMetrics {
		ms.GaugeStorage[k] = gauge(v)
	}

	metricsCounter := servermetrics.MetricsCount{}
	metricsCounter.New()

	for k, v := range metricsCounter.RuntimeMetrics {
		ms.GaugeStorage[k] = gauge(v)
	}
}

func (ms *MemStorage) SaveGauge(metrics string, value gauge) error {
	ms.GaugeStorage[metrics] = value
	return nil
}

func (ms *MemStorage) SaveCounter(metrics string, value counter) error {
	ms.CounterStorage[metrics] = value
	return nil
}

func (ms *MemStorage) Remove(metrics string) error {
	_, okGauge := ms.GaugeStorage[metrics]
	_, okCounter := ms.CounterStorage[metrics]

	if okGauge {
		delete(ms.GaugeStorage, metrics)
	}

	if okCounter {
		delete(ms.CounterStorage, metrics)
	}

	return nil
}

func (ms *MemStorage) IfExists(metrics string, value counter) bool {
	_, okGauge := ms.GaugeStorage[metrics]
	_, okCounter := ms.CounterStorage[metrics]

	return okGauge || okCounter
}

func (ms *MemStorage) PrintAll() (string, error) {
	outerString := ""

	metricsGauge := ms.GaugeStorage
	for k, v := range metricsGauge {
		outerString = outerString + k + ": " + GaugeToString(v) + "\n"
	}

	metricsCounter := ms.CounterStorage
	for k, v := range metricsCounter {
		outerString = outerString + k + ": " + CounterToString(v) + "\n"
	}

	return outerString, nil
}
