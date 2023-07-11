package servermetrics

import "strconv"

const (
	GaugeType   string = "gauge"
	CounterType string = "counter"
)

type Gauge float64
type Counter int64

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func IfHasCorrectType(s string) bool {
	if s == GaugeType || s == CounterType {
		return true
	}

	return false
}

func GetGaugeMetricsResponseNames() map[string]string {
	mapOfMetricsNames := make(map[string]string)
	mapOfMetricsNames["alloc"] = "Alloc"
	mapOfMetricsNames["buckhashsys"] = "BuckHashSys"
	mapOfMetricsNames["frees"] = "Frees"
	mapOfMetricsNames["gccpufraction"] = "GCCPUFraction"
	mapOfMetricsNames["gcsys"] = "GCSys"
	mapOfMetricsNames["heapalloc"] = "HeapAlloc"
	mapOfMetricsNames["heapidle"] = "HeapIdle"
	mapOfMetricsNames["heapinuse"] = "HeapInuse"
	mapOfMetricsNames["heapobjects"] = "HeapObjects"
	mapOfMetricsNames["heapreleased"] = "HeapReleased"
	mapOfMetricsNames["heapsys"] = "HeapSys"
	mapOfMetricsNames["lastgc"] = "LastGC"
	mapOfMetricsNames["lookups"] = "Lookups"
	mapOfMetricsNames["mcacheinuse"] = "MCacheInuse"
	mapOfMetricsNames["mcachesys"] = "MCacheSys"
	mapOfMetricsNames["mspaninuse"] = "MSpanInuse"
	mapOfMetricsNames["mspansys"] = "MSpanSys"
	mapOfMetricsNames["mallocs"] = "Mallocs"
	mapOfMetricsNames["nextgc"] = "NextGC"
	mapOfMetricsNames["numforcedgc"] = "NumForcedGC"
	mapOfMetricsNames["numgc"] = "NumGC"
	mapOfMetricsNames["othersys"] = "OtherSys"
	mapOfMetricsNames["pausetotalns"] = "PauseTotalNs"
	mapOfMetricsNames["stackinuse"] = "StackInuse"
	mapOfMetricsNames["stacksys"] = "StackSys"
	mapOfMetricsNames["sys"] = "Sys"
	mapOfMetricsNames["totalalloc"] = "TotalAlloc"
	mapOfMetricsNames["randomvalue"] = "RandomValue"
	return mapOfMetricsNames
}

func GetCounterMetricsResponseNames() map[string]string {
	mapOfMetricsNames := make(map[string]string)
	mapOfMetricsNames["PollCount"] = "PollCount"
	return mapOfMetricsNames
}

func GaugeToString(gaugeValue Gauge) string {
	value := strconv.FormatFloat(float64(gaugeValue), 'f', -1, 64)

	return value
}

func StringToGauge(s string, bitSize int) (Gauge, error) {
	v, e := strconv.ParseFloat(s, bitSize)
	if e != nil {
		return 0.0, e
	}
	return Gauge(v), nil
}

func CounterToString(cv Counter) string {
	value := strconv.Itoa(int(cv))
	return value
}

func StringToCounter(s string) (Counter, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return Counter(value), nil
}

type MetricsGauge struct {
	RuntimeMetrics map[string]Gauge
	RandomValue    Gauge
}

type MetricsCount struct {
	RuntimeMetrics map[string]Counter
}

func (mg *MetricsGauge) New() {
	mg.RuntimeMetrics = make(map[string]Gauge)
	metricsNames := GetGaugeMetricsResponseNames()

	for _, v := range metricsNames {
		mg.RuntimeMetrics[v] = 0.0
	}
}

func (mc *MetricsCount) New() {
	mc.RuntimeMetrics = make(map[string]Counter)
	metricsNames := GetCounterMetricsResponseNames()
	for _, v := range metricsNames {
		mc.RuntimeMetrics[v] = 0
	}
}
