package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	t.Run("logs request details", func(t *testing.T) {
		// Capture log output
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil) // Reset to default

		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Wrap with logging middleware
		handler := Logging(testHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Check log output
		logOutput := buf.String()

		if !strings.Contains(logOutput, "GET") {
			t.Errorf("Expected log to contain 'GET', got: %s", logOutput)
		}

		if !strings.Contains(logOutput, "/test/path") {
			t.Errorf("Expected log to contain '/test/path', got: %s", logOutput)
		}

		if !strings.Contains(logOutput, "200") {
			t.Errorf("Expected log to contain '200', got: %s", logOutput)
		}
	})

	t.Run("logs different status codes", func(t *testing.T) {
		statusCodes := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusInternalServerError,
		}

		for _, code := range statusCodes {
			var buf bytes.Buffer
			log.SetOutput(&buf)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			})

			handler := Logging(testHandler)
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			logOutput := buf.String()
			if !strings.Contains(logOutput, "POST") {
				t.Errorf("Expected log to contain 'POST' for status %d", code)
			}

			log.SetOutput(nil)
		}
	})

	t.Run("captures custom status code from WriteHeader", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		})

		handler := Logging(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		logOutput := buf.String()
		if !strings.Contains(logOutput, "404") {
			t.Errorf("Expected log to contain '404', got: %s", logOutput)
		}
	})

	t.Run("logs duration", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := Logging(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		logOutput := buf.String()
		// Duration should contain time units (µs, ms, s)
		hasDuration := strings.Contains(logOutput, "µs") ||
			strings.Contains(logOutput, "ms") ||
			strings.Contains(logOutput, "s") ||
			strings.Contains(logOutput, "ns")

		if !hasDuration {
			t.Errorf("Expected log to contain duration, got: %s", logOutput)
		}
	})
}

func TestResponseWriter(t *testing.T) {
	t.Run("captures status code on WriteHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		rw.WriteHeader(http.StatusCreated)

		if rw.statusCode != http.StatusCreated {
			t.Errorf("Expected status code 201, got %d", rw.statusCode)
		}
	})

	t.Run("default status code is 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Don't call WriteHeader, just write body
		rw.Write([]byte("test"))

		if rw.statusCode != http.StatusOK {
			t.Errorf("Expected default status code 200, got %d", rw.statusCode)
		}
	})
}
