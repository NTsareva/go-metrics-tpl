package handlers

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"log"
)

func WriteMemstorageToFile() {
	for k, v := range memStorage.GaugeStorage {
		floatMetricValue := float64(v)
		int64MetricDelta := int64(0)
		metric := memstorage.Metrics{
			ID:    k,
			MType: memstorage.GaugeType,
			Delta: &int64MetricDelta,
			Value: &floatMetricValue,
		}

		Producer.WriteMetric(&metric)
	}
	for k, v := range memStorage.CounterStorage {
		floatMetricValue := float64(0.0)
		int64MetricDelta := int64(v)
		metric := memstorage.Metrics{
			ID:    k,
			MType: memstorage.CounterType,
			Delta: &int64MetricDelta,
			Value: &floatMetricValue,
		}

		Producer.WriteMetric(&metric)
	}

	log.Println("Saving done")
}
