package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PostHandler(res http.ResponseWriter, req *http.Request) {
	//Обработчик, что пост
	methodTypeHandler(res, req)
	//Обработчик, что не правильный path
	//Обработчик, что не правильный тип контента
	//Обработчик, что правильный URL
	//Обработчик<ТИП_МЕТРИКИ>/
	relPathString := req.URL.Path
	pathParamsArray := strings.Split(relPathString, "/")
	if len(pathParamsArray) != 5 {
		http.Error(res, "incorrect request", http.StatusNotFound)
	}
	metricsType := pathParamsArray[2]
	checkType(res, metricsType)
	//Обработчик<ИМЯ_МЕТРИКИ>
	//metricsName := pathParamsArray[3]
	//обработчик <ЗНАЧЕНИЕ_МЕТРИКИ>

	body := req.Method

	for k, v := range req.Header {
		body += fmt.Sprintf("%s: % v\r\n", k, v)
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
