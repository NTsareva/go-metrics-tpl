package agentmetrics

import (
	"runtime"
	"time"
)

type gauge float64
type counter int64

const pollInterval = 2
const ReportInterval = 10

type MetricsGauge struct {
	RuntimeMetrics map[string]gauge
	PollCount      counter
	RandomValue    gauge
}

func (m *MetricsGauge) New() {
	m.RuntimeMetrics = make(map[string]gauge)
	m.RuntimeMetrics["Alloc"] = 0
	m.RuntimeMetrics["BuckHashSys"] = 0
	m.RuntimeMetrics["Frees"] = 0
	m.RuntimeMetrics["GCCPUFraction"] = 0
	m.RuntimeMetrics["GCSys"] = 0
	m.RuntimeMetrics["HeapAlloc"] = 0
	m.RuntimeMetrics["HeapIdle"] = 0
	m.RuntimeMetrics["HeapInuse"] = 0
	m.RuntimeMetrics["HeapObjects"] = 0
	m.RuntimeMetrics["HeapReleased"] = 0
	m.RuntimeMetrics["HeapSys"] = 0
	m.RuntimeMetrics["LastGC"] = 0
	m.RuntimeMetrics["Lookups"] = 0
	m.RuntimeMetrics["MCacheInuse"] = 0
	m.RuntimeMetrics["MCacheSys"] = 0
	m.RuntimeMetrics["MSpanInuse"] = 0
	m.RuntimeMetrics["MSpanSys"] = 0
	m.RuntimeMetrics["Mallocs"] = 0
	m.RuntimeMetrics["NextGC"] = 0
	m.RuntimeMetrics["NumForcedGC"] = 0
	m.RuntimeMetrics["NumGC"] = 0
	m.RuntimeMetrics["OtherSys"] = 0
	m.RuntimeMetrics["PauseTotalNs"] = 0
	m.RuntimeMetrics["StackInuse"] = 0
	m.RuntimeMetrics["StackSys"] = 0
	m.RuntimeMetrics["Sys"] = 0
	m.RuntimeMetrics["TotalAlloc"] = 0
}

func (m *MetricsGauge) Renew() {
	time.Sleep(pollInterval * time.Second)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.RuntimeMetrics["Alloc"] = gauge(ms.Alloc)
	m.RuntimeMetrics["BuckHashSys"] = gauge(ms.BuckHashSys)
	m.RuntimeMetrics["Frees"] = gauge(ms.Frees)
	m.RuntimeMetrics["GCCPUFraction"] = gauge(ms.GCCPUFraction)
	m.RuntimeMetrics["GCSys"] = gauge(ms.GCSys)
	m.RuntimeMetrics["HeapAlloc"] = gauge(ms.HeapAlloc)
	m.RuntimeMetrics["HeapIdle"] = gauge(ms.HeapIdle)
	m.RuntimeMetrics["HeapInuse"] = gauge(ms.HeapInuse)
	m.RuntimeMetrics["HeapObjects"] = gauge(ms.HeapObjects)
	m.RuntimeMetrics["HeapReleased"] = gauge(ms.HeapReleased)
	m.RuntimeMetrics["HeapSys"] = gauge(ms.HeapSys)
	m.RuntimeMetrics["LastGC"] = gauge(ms.LastGC)
	m.RuntimeMetrics["Lookups"] = gauge(ms.Lookups)
	m.RuntimeMetrics["MCacheInuse"] = gauge(ms.MCacheInuse)
	m.RuntimeMetrics["MCacheSys"] = gauge(ms.MCacheSys)
	m.RuntimeMetrics["MSpanInuse"] = gauge(ms.MSpanInuse)
	m.RuntimeMetrics["MSpanSys"] = gauge(ms.MSpanSys)
	m.RuntimeMetrics["Mallocs"] = gauge(ms.Mallocs)
	m.RuntimeMetrics["NextGC"] = gauge(ms.NextGC)
	m.RuntimeMetrics["NumForcedGC"] = gauge(ms.NumForcedGC)
	m.RuntimeMetrics["NumGC"] = gauge(ms.NumGC)
	m.RuntimeMetrics["OtherSys"] = gauge(ms.OtherSys)
	m.RuntimeMetrics["PauseTotalNs"] = gauge(ms.PauseTotalNs)
	m.RuntimeMetrics["StackInuse"] = gauge(ms.StackInuse)
	m.RuntimeMetrics["StackSys"] = gauge(ms.StackSys)
	m.RuntimeMetrics["Sys"] = gauge(ms.Sys)
	m.RuntimeMetrics["TotalAlloc"] = gauge(ms.TotalAlloc)

	m.PollCount += 1
}
