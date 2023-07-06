package handlers

import (
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/filestorage"
	"github.com/NTsareva/go-metrics-tpl.git/cmd/storage/memstorage"
	"log"
)

func WriteMemstorageToFile(filePath string) {
	if filePath != "" {
		var err error
		Producer, err = filestorage.NewProducer(filePath)
		//log.Println(filePath)
		if err != nil {
			log.Println(err)
			log.Println("prod")
		}
	}
	if Producer != nil {
		if memStorage.GaugeStorage != nil {
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
		}
		if memStorage.CounterStorage != nil {
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
		}

		Producer.Close()
		log.Println("Saving done")
	}
}
