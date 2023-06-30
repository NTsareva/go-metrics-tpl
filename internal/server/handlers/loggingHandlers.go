package handlers

import (
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var sugar zap.SugaredLogger

type responseData struct {
	status int
	size   int
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

var loggingResponse = LoggingResponseWriter{responseData: nil}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Print(err, "#7")
		}

		defer logger.Sync()

		sugar := *logger.Sugar()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		loggingResponse = LoggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&loggingResponse, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}
