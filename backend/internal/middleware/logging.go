package middleware

import (
"log"
"net/http"
"time"
)

// Logging logs HTTP requests with method, path, and duration
func Logging(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
start := time.Now()

// Create a response writer wrapper to capture status code
wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

next.ServeHTTP(wrapped, r)

duration := time.Since(start)
log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, duration)
})
}

type responseWriter struct {
http.ResponseWriter
statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
rw.statusCode = code
rw.ResponseWriter.WriteHeader(code)
}
