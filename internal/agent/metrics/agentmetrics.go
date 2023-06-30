package agentmetrics

import (
	"math/rand"
	"runtime"
	"strconv"
)

type gauge float64
type counter int64

type MetricsGauge struct {
	RuntimeMetrics map[string]gauge
	RandomValue    gauge
}

type MetricsCount struct {
	RuntimeMetrics map[string]counter
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

func CounterToString(counterValue counter) string {
	return strconv.Itoa(int(counterValue))
}

func (m *MetricsGauge) New() {
	m.RuntimeMetrics = make(map[string]gauge)
	m.RuntimeMetrics["Alloc"] = 0.0
	m.RuntimeMetrics["BuckHashSys"] = 0.0
	m.RuntimeMetrics["Frees"] = 0.0
	m.RuntimeMetrics["GCCPUFraction"] = 0.0
	m.RuntimeMetrics["GCSys"] = 0.0
	m.RuntimeMetrics["HeapAlloc"] = 0.0
	m.RuntimeMetrics["HeapIdle"] = 0.0
	m.RuntimeMetrics["HeapInuse"] = 0.0
	m.RuntimeMetrics["HeapObjects"] = 0.0
	m.RuntimeMetrics["HeapReleased"] = 0.0
	m.RuntimeMetrics["HeapSys"] = 0.0
	m.RuntimeMetrics["LastGC"] = 0.0
	m.RuntimeMetrics["Lookups"] = 0.0
	m.RuntimeMetrics["MCacheInuse"] = 0.0
	m.RuntimeMetrics["MCacheSys"] = 0.0
	m.RuntimeMetrics["MSpanInuse"] = 0.0
	m.RuntimeMetrics["MSpanSys"] = 0.0
	m.RuntimeMetrics["Mallocs"] = 0.0
	m.RuntimeMetrics["NextGC"] = 0.0
	m.RuntimeMetrics["NumForcedGC"] = 0.0
	m.RuntimeMetrics["NumGC"] = 0.0
	m.RuntimeMetrics["OtherSys"] = 0.0
	m.RuntimeMetrics["PauseTotalNs"] = 0.0
	m.RuntimeMetrics["RandomValue"] = 0.0
	m.RuntimeMetrics["StackInuse"] = 0.0
	m.RuntimeMetrics["StackSys"] = 0.0
	m.RuntimeMetrics["Sys"] = 0.0
	m.RuntimeMetrics["TotalAlloc"] = 0.0
}

func (m *MetricsGauge) Renew() {
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
	m.RuntimeMetrics["RandomValue"] = gauge(rand.Float64())
	m.RuntimeMetrics["StackInuse"] = gauge(ms.StackInuse)
	m.RuntimeMetrics["StackSys"] = gauge(ms.StackSys)
	m.RuntimeMetrics["Sys"] = gauge(ms.Sys)
	m.RuntimeMetrics["TotalAlloc"] = gauge(ms.TotalAlloc)
}

func (mc *MetricsCount) New() {
	mc.RuntimeMetrics = make(map[string]counter)
	mc.RuntimeMetrics["PollCount"] = 0
}

func (mc *MetricsCount) Renew() {
	mc.RuntimeMetrics["PollCount"] += 1
}
