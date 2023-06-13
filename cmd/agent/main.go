package main

//
//import (
//	"github.com/NTsareva/go-metrics-tpl.git/internal/metrics"
//	"log"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//const URL = "http://localhost:8080"
//
//func main() {
//	var gm metrics.MetricsGauge
//	gm.New()
//
//	for {
//		go gm.Renew()
//		SendRuntimeMetrics(&gm)
//	}
//}
//
//func SendRuntimeMetrics(m *metrics.MetricsGauge) {
//	time.Sleep(metrics.ReportInterval * time.Second)
//
//	url := URL
//	for k, v := range m.RuntimeMetrics {
//		url = url + "/update/gauge/" + k + "/" + strconv.FormatFloat(float64(v), 'f', 1, 64)
//
//		response, err := http.Post(url, "text/plain", nil)
//		if err != nil {
//			panic(err)
//		}
//
//		defer response.Body.Close()
//		log.Println(response)
//	}
//}
