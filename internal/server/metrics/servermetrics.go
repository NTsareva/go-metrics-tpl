package servermetrics

import (
	"strconv"
)

const (
	Gauge   string = "gauge"
	Counter string = "counter"
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
	value := strconv.FormatFloat(float64(gv), 'f', 64, 64)
	return value, nil
}

func GaugeToString(gv gauge) string {
	value := strconv.FormatFloat(float64(gv), 'f', 64, 64)
	return value
}

type MetricsGauge struct {
	RuntimeMetrics map[string]gauge
	RandomValue    gauge
}

type MetricsCount struct {
	RuntimeMetrics map[string]counter
}

func (mg *MetricsGauge) New() {
	mg.RuntimeMetrics = make(map[string]gauge)
	mg.RuntimeMetrics["alloc"] = 0.0
	mg.RuntimeMetrics["buckhashsys"] = 0.0
	mg.RuntimeMetrics["frees"] = 0.0
	mg.RuntimeMetrics["gccpufraction"] = 0.0
	mg.RuntimeMetrics["gcsys"] = 0.0
	mg.RuntimeMetrics["heapalloc"] = 0.0
	mg.RuntimeMetrics["heapidle"] = 0.0
	mg.RuntimeMetrics["heapinuse"] = 0.0
	mg.RuntimeMetrics["heapobjects"] = 0.0
	mg.RuntimeMetrics["heapreleased"] = 0.0
	mg.RuntimeMetrics["heapsys"] = 0.0
	mg.RuntimeMetrics["lastgc"] = 0.0
	mg.RuntimeMetrics["lookups"] = 0.0
	mg.RuntimeMetrics["mccacheinuse"] = 0.0
	mg.RuntimeMetrics["mcachesys"] = 0.0
	mg.RuntimeMetrics["mspaninuse"] = 0.0
	mg.RuntimeMetrics["mspansys"] = 0.0
	mg.RuntimeMetrics["mallocs"] = 0.0
	mg.RuntimeMetrics["nextgc"] = 0.0
	mg.RuntimeMetrics["numcorcedgc"] = 0.0
	mg.RuntimeMetrics["numgc"] = 0.0
	mg.RuntimeMetrics["othersys"] = 0.0
	mg.RuntimeMetrics["pausetotalns"] = 0.0
	mg.RuntimeMetrics["stackinuse"] = 0.0
	mg.RuntimeMetrics["stacksys"] = 0.0
	mg.RuntimeMetrics["sys"] = 0.0
	mg.RuntimeMetrics["totalalloc"] = 0.0
}

func (mc *MetricsCount) New() {
	mc.RuntimeMetrics = make(map[string]counter)
	mc.RuntimeMetrics["pollcount"] = 0
}
