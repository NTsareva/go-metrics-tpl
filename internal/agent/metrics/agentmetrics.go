package agentmetrics

import (
	"runtime"
	"strconv"
)

type gauge float64
type counter int64

const ReportInterval = 10

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

func (m *MetricsGauge) New() {
	m.RuntimeMetrics = make(map[string]gauge)
	m.RuntimeMetrics["alloc"] = 0.0
	m.RuntimeMetrics["buckhashsys"] = 0.0
	m.RuntimeMetrics["frees"] = 0.0
	m.RuntimeMetrics["gccpufraction"] = 0.0
	m.RuntimeMetrics["gcsys"] = 0.0
	m.RuntimeMetrics["heapalloc"] = 0.0
	m.RuntimeMetrics["heapidle"] = 0.0
	m.RuntimeMetrics["heapinuse"] = 0.0
	m.RuntimeMetrics["heapobjects"] = 0.0
	m.RuntimeMetrics["heapreleased"] = 0.0
	m.RuntimeMetrics["heapsys"] = 0.0
	m.RuntimeMetrics["lastgc"] = 0.0
	m.RuntimeMetrics["lookups"] = 0.0
	m.RuntimeMetrics["mccacheinuse"] = 0.0
	m.RuntimeMetrics["mcachesys"] = 0.0
	m.RuntimeMetrics["mspaninuse"] = 0.0
	m.RuntimeMetrics["mspansys"] = 0.0
	m.RuntimeMetrics["mallocs"] = 0.0
	m.RuntimeMetrics["nextgc"] = 0.0
	m.RuntimeMetrics["numcorcedgc"] = 0.0
	m.RuntimeMetrics["numc"] = 0.0
	m.RuntimeMetrics["othersys"] = 0.0
	m.RuntimeMetrics["pausetotalns"] = 0.0
	m.RuntimeMetrics["stackinuse"] = 0.0
	m.RuntimeMetrics["stacksys"] = 0.0
	m.RuntimeMetrics["sys"] = 0.0
	m.RuntimeMetrics["totalalloc"] = 0.0
}

func (m *MetricsGauge) Renew() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	m.RuntimeMetrics["alloc"] = gauge(ms.Alloc)
	m.RuntimeMetrics["buckhashsys"] = gauge(ms.BuckHashSys)
	m.RuntimeMetrics["frees"] = gauge(ms.Frees)
	m.RuntimeMetrics["gccpufraction"] = gauge(ms.GCCPUFraction)
	m.RuntimeMetrics["gcsys"] = gauge(ms.GCSys)
	m.RuntimeMetrics["heapalloc"] = gauge(ms.HeapAlloc)
	m.RuntimeMetrics["heapidle"] = gauge(ms.HeapIdle)
	m.RuntimeMetrics["heapinuse"] = gauge(ms.HeapInuse)
	m.RuntimeMetrics["heapobjects"] = gauge(ms.HeapObjects)
	m.RuntimeMetrics["heapreleased"] = gauge(ms.HeapReleased)
	m.RuntimeMetrics["heapsys"] = gauge(ms.HeapSys)
	m.RuntimeMetrics["lastgc"] = gauge(ms.LastGC)
	m.RuntimeMetrics["lookups"] = gauge(ms.Lookups)
	m.RuntimeMetrics["mcacheinuse"] = gauge(ms.MCacheInuse)
	m.RuntimeMetrics["mcachesys"] = gauge(ms.MCacheSys)
	m.RuntimeMetrics["mspaninuse"] = gauge(ms.MSpanInuse)
	m.RuntimeMetrics["mspansys"] = gauge(ms.MSpanSys)
	m.RuntimeMetrics["mallocs"] = gauge(ms.Mallocs)
	m.RuntimeMetrics["nextgc"] = gauge(ms.NextGC)
	m.RuntimeMetrics["numforcedgc"] = gauge(ms.NumForcedGC)
	m.RuntimeMetrics["numgc"] = gauge(ms.NumGC)
	m.RuntimeMetrics["othersys"] = gauge(ms.OtherSys)
	m.RuntimeMetrics["pausetotalns"] = gauge(ms.PauseTotalNs)
	m.RuntimeMetrics["stackinuse"] = gauge(ms.StackInuse)
	m.RuntimeMetrics["stacksys"] = gauge(ms.StackSys)
	m.RuntimeMetrics["sys"] = gauge(ms.Sys)
	m.RuntimeMetrics["totalslloc"] = gauge(ms.TotalAlloc)
}

func (mc *MetricsCount) New() {
	mc.RuntimeMetrics = make(map[string]counter)
	mc.RuntimeMetrics["pollcount"] = 0
}

func (mc *MetricsCount) Renew() {
	mc.RuntimeMetrics["pollcount"] += 1
}
