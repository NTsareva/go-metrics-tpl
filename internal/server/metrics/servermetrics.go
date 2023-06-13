package servermetrics

import (
	"strconv"
)

const (
	Gauge   string = "gauge"
	Counter        = "counter"
)

type gauge float64
type counter int64

func IfHasCorrestType(s string) bool {
	if s == Gauge || s == Counter {
		return true
	}

	return false
}

func StringToGauge(s string, bitSize int) (gauge, error) {
	v, e := strconv.ParseFloat(s, bitSize)
	if e != nil {
		return 0.0, e
	}
	return gauge(v), nil
}
func GaugeToStringWithError(gv gauge) (string, error) {
	value := strconv.FormatFloat(float64(gv), 'f', 1, 64)
	return value, nil
}

func GaugeToString(gv gauge) string {
	value := strconv.FormatFloat(float64(gv), 'f', 1, 64)
	return value
}

type MetricsGauge struct {
	RuntimeMetrics map[string]gauge
	PollCount      counter
	RandomValue    gauge
}

type MetricsCount struct {
	RuntimeMetrics map[string]counter
}

func (mg *MetricsGauge) New() {
	mg.RuntimeMetrics = make(map[string]gauge)
	mg.RuntimeMetrics["Alloc"] = 0.0
	mg.RuntimeMetrics["BuckHashSys"] = 0.0
	mg.RuntimeMetrics["Frees"] = 0.0
	mg.RuntimeMetrics["GCCPUFraction"] = 0.0
	mg.RuntimeMetrics["GCSys"] = 0.0
	mg.RuntimeMetrics["HeapAlloc"] = 0.0
	mg.RuntimeMetrics["HeapIdle"] = 0.0
	mg.RuntimeMetrics["HeapInuse"] = 0.0
	mg.RuntimeMetrics["HeapObjects"] = 0.0
	mg.RuntimeMetrics["HeapReleased"] = 0.0
	mg.RuntimeMetrics["HeapSys"] = 0.0
	mg.RuntimeMetrics["LastGC"] = 0.0
	mg.RuntimeMetrics["Lookups"] = 0.0
	mg.RuntimeMetrics["MCacheInuse"] = 0.0
	mg.RuntimeMetrics["MCacheSys"] = 0.0
	mg.RuntimeMetrics["MSpanInuse"] = 0.0
	mg.RuntimeMetrics["MSpanSys"] = 0.0
	mg.RuntimeMetrics["Mallocs"] = 0.0
	mg.RuntimeMetrics["NextGC"] = 0.0
	mg.RuntimeMetrics["NumForcedGC"] = 0.0
	mg.RuntimeMetrics["NumGC"] = 0.0
	mg.RuntimeMetrics["OtherSys"] = 0.0
	mg.RuntimeMetrics["PauseTotalNs"] = 0.0
	mg.RuntimeMetrics["StackInuse"] = 0.0
	mg.RuntimeMetrics["StackSys"] = 0.0
	mg.RuntimeMetrics["Sys"] = 0.0
	mg.RuntimeMetrics["TotalAlloc"] = 0.0
}

func (mc *MetricsCount) New() {
	mc.RuntimeMetrics = make(map[string]counter)
	mc.RuntimeMetrics["PollCount"] = 0
}
