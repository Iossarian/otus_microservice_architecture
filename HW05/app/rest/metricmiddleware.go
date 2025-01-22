package rest

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var (
	// Гистограмма для времени отклика
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Счетчик для запросов
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path"},
	)

	// Счетчик для ошибок 500
	httpErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP 500 errors.",
		},
		[]string{"method", "path"},
	)
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Записываем время начала
		start := time.Now()

		// Получаем путь и метод
		path, _ := mux.CurrentRoute(r).GetPathTemplate()
		method := r.Method

		if strings.Contains(path, "/metrics") {
			next.ServeHTTP(w, r)
			return
		}

		// Создаем кастомный ResponseWriter, чтобы захватывать статус код
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		// Считаем длительность
		duration := time.Since(start).Seconds()

		// Увеличиваем счетчики
		httpRequests.WithLabelValues(method, path).Inc()
		httpDuration.WithLabelValues(method, path).Observe(duration)

		if rw.statusCode == http.StatusInternalServerError || rw.statusCode == http.StatusNotFound {
			httpErrors.WithLabelValues(method, path).Inc()
		}
	})
}

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpErrors)
	prometheus.MustRegister(httpDuration)
}
