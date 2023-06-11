package main

import (
	"runtime"
	"time"
)

type gauge float64
type counter int64

const pollInterval = 2
const reportInterval = 10

type MetricsGauge struct {
	runtimeMetrics map[string]float64
	PollCount      counter
	RandomValue    gauge
}

func (m *MetricsGauge) New() {
	m.runtimeMetrics = make(map[string]float64)
	m.runtimeMetrics["Alloc"] = 0
	m.runtimeMetrics["BuckHashSys"] = 0
	m.runtimeMetrics["Frees"] = 0
	m.runtimeMetrics["GCCPUFraction"] = 0
	m.runtimeMetrics["GCSys"] = 0
	m.runtimeMetrics["HeapAlloc"] = 0
	m.runtimeMetrics["HeapIdle"] = 0
	m.runtimeMetrics["HeapInuse"] = 0
	m.runtimeMetrics["HeapObjects"] = 0
	m.runtimeMetrics["HeapReleased"] = 0
	m.runtimeMetrics["HeapSys"] = 0
	m.runtimeMetrics["LastGC"] = 0
	m.runtimeMetrics["Lookups"] = 0
	m.runtimeMetrics["MCacheInuse"] = 0
	m.runtimeMetrics["MCacheSys"] = 0
	m.runtimeMetrics["MSpanInuse"] = 0
	m.runtimeMetrics["MSpanSys"] = 0
	m.runtimeMetrics["Mallocs"] = 0
	m.runtimeMetrics["NextGC"] = 0
	m.runtimeMetrics["NumForcedGC"] = 0
	m.runtimeMetrics["NumGC"] = 0
	m.runtimeMetrics["OtherSys"] = 0
	m.runtimeMetrics["PauseTotalNs"] = 0
	m.runtimeMetrics["StackInuse"] = 0
	m.runtimeMetrics["StackSys"] = 0
	m.runtimeMetrics["Sys"] = 0
	m.runtimeMetrics["TotalAlloc"] = 0
}

func (m *MetricsGauge) Renew() {
	time.Sleep(pollInterval * time.Second)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.runtimeMetrics["Alloc"] = float64(ms.Alloc)
	m.runtimeMetrics["BuckHashSys"] = float64(ms.BuckHashSys)
	m.runtimeMetrics["Frees"] = float64(ms.Frees)
	m.runtimeMetrics["GCCPUFraction"] = float64(ms.GCCPUFraction)
	m.runtimeMetrics["GCSys"] = float64(ms.GCSys)
	m.runtimeMetrics["HeapAlloc"] = float64(ms.HeapAlloc)
	m.runtimeMetrics["HeapIdle"] = float64(ms.HeapIdle)
	m.runtimeMetrics["HeapInuse"] = float64(ms.HeapInuse)
	m.runtimeMetrics["HeapObjects"] = float64(ms.HeapObjects)
	m.runtimeMetrics["HeapReleased"] = float64(ms.HeapReleased)
	m.runtimeMetrics["HeapSys"] = float64(ms.HeapSys)
	m.runtimeMetrics["LastGC"] = float64(ms.LastGC)
	m.runtimeMetrics["Lookups"] = float64(ms.Lookups)
	m.runtimeMetrics["MCacheInuse"] = float64(ms.MCacheInuse)
	m.runtimeMetrics["MCacheSys"] = float64(ms.MCacheSys)
	m.runtimeMetrics["MSpanInuse"] = float64(ms.MSpanInuse)
	m.runtimeMetrics["MSpanSys"] = float64(ms.MSpanSys)
	m.runtimeMetrics["Mallocs"] = float64(ms.Mallocs)
	m.runtimeMetrics["NextGC"] = float64(ms.NextGC)
	m.runtimeMetrics["NumForcedGC"] = float64(ms.NumForcedGC)
	m.runtimeMetrics["NumGC"] = float64(ms.NumGC)
	m.runtimeMetrics["OtherSys"] = float64(ms.OtherSys)
	m.runtimeMetrics["PauseTotalNs"] = float64(ms.PauseTotalNs)
	m.runtimeMetrics["StackInuse"] = float64(ms.StackInuse)
	m.runtimeMetrics["StackSys"] = float64(ms.StackSys)
	m.runtimeMetrics["Sys"] = float64(ms.Sys)
	m.runtimeMetrics["TotalAlloc"] = float64(ms.TotalAlloc)

	m.PollCount += 1
}
