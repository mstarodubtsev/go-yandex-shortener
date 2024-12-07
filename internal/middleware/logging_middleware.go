package middleware

import (
	"github.com/mstarodubtsev/go-yandex-shortener/internal/log"
	"net/http"
	"time"
)

type (
	// struct for storing response data
	responseData struct {
		status int
		size   int
	}

	// add implementation http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // inject original http.ResponseWriter
		responseData        *responseData
	}
)

// Write
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// write response using the original http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// write status code using the original http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// WithLogging add new code to log request/response data and return new http.Handler
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // inject implementation http.ResponseWriter

		duration := time.Since(start)

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
