package main

import (
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(res http.ResponseWriter, req *http.Request) {
	methodTypeHandler(res, req)

	relPathString := req.URL.Path
	pathParamsArray := strings.Split(relPathString, "/")

	if len(pathParamsArray) == 4 {
		http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
		return
	}

	if len(pathParamsArray) == 3 {
		http.Error(res, "incorrect type of metrics", http.StatusNotFound)
		return
	}

	metricsType := pathParamsArray[2]
	checkType(res, metricsType)
	//metricsName := pathParamsArray[3]
	//обработчик <ЗНАЧЕНИЕ_МЕТРИКИ>
	metricsValue := pathParamsArray[4]
	checkMetricsValue(res, metricsValue)
}

func checkMetricsValue(res http.ResponseWriter, metricValue string) {
	_, e := strconv.Atoi(metricValue)
	if e != nil {
		http.Error(res, "incorrect value of metrics", http.StatusBadRequest)
		return
	}
}

func methodTypeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "only POST methods are allowed", http.StatusMethodNotAllowed)
		return
	}
}

func checkType(res http.ResponseWriter, metricsType string) {
	switch metricsType {
	case "gauge":
		//Работаем с типом по его логике
	case "counter":
		//Работаем с типом по его логике
	default:
		http.Error(res, "incorrect type of metrics "+metricsType, http.StatusBadRequest)
		return
	}
}
