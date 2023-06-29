package memstorage

import (
	"errors"
	servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"
)

type Gauge servermetrics.Gauge
type Counter servermetrics.Counter

const (
	GaugeType   string = "gauge"
	CounterType string = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Storage interface {
	New()
	Save(metrics string, value interface{}) error
	SaveGauge(metrics string, value Gauge) error
	SaveCounter(metrics string, value Counter) error
	Remove(metric string) error
	PrintAll() (string, error)
	Get(metrics string, metricType string) (string, error)
}

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
		memStorage.CounterStorage[k] = Counter(v)
	}
}

func (memStorage *MemStorage) Save(metrics string, value interface{}) error {
	switch i := value.(type) {
	case Gauge:
		memStorage.GaugeStorage[metrics] = i
		return nil
	case Counter:
		memStorage.CounterStorage[metrics] = i
		return nil
	default:
		return errors.New("no such type")
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

func (memStorage *MemStorage) Get(metric string, metricType string) (string, error) {
	metricsGauge := memStorage.GaugeStorage
	metricsCounter := memStorage.CounterStorage

	if metricType == GaugeType {
		return servermetrics.GaugeToString(servermetrics.Gauge(metricsGauge[metric])), nil
	} else if metricType == CounterType {
		return servermetrics.CounterToString(servermetrics.Counter(metricsCounter[metric])), nil
	} else {
		return "", errors.New("Incorrect type, should be Gauge or Counter")
	}
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
