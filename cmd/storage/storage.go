package storage

type Gauge float64
type Counter int64

type MetricType interface {
	Gauge | Counter
}

type Storage interface {
	SaveGauge(metrics string, value Gauge) error
	SaveCounter(metrics string, value Counter) error
	Remove(metric string) error
	IfExists(metric string) bool
	PrintAll() (string, error)
}
