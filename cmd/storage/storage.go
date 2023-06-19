package storage

import servermetrics "github.com/NTsareva/go-metrics-tpl.git/internal/server/metrics"

type Gauge servermetrics.Gauge
type Counter servermetrics.Counter

type MetricType interface {
	Gauge | Counter
}

type Storage interface {
	SaveGauge(metrics string, value Gauge) error
	SaveCounter(metrics string, value Counter) error
	Remove(metric string) error
	PrintAll() (string, error)
}
