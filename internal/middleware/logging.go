package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs incoming requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request details
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log response details
		duration := time.Since(start)
		log.Printf("Response: %s %s - %d (%s)", r.Method, r.URL.Path, rw.status, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}