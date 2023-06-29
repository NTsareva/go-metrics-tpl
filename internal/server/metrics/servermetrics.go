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

func GaugeToString(gaugeValue Gauge) string {
	value := strconv.FormatFloat(float64(gaugeValue), 'f', 3, 64)
	if value[len(value)-1] == '0' && value[len(value)-2] == '0' {
		value = strconv.FormatFloat(float64(gaugeValue), 'f', 1, 64)
	}
	if value[len(value)-1] == '0' {
		value = strconv.FormatFloat(float64(gaugeValue), 'f', 2, 64)
	}

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
	mc.RuntimeMetrics = make(map[string]Counter)
	mc.RuntimeMetrics["pollcount"] = 0
}
