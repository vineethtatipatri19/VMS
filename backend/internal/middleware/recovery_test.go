package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	t.Run("recovers from panic and returns 500", func(t *testing.T) {
		// Capture log output
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		// Create handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something went wrong!")
		})

		// Wrap with recovery middleware
		handler := Recovery(panicHandler)

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Should not panic
		handler.ServeHTTP(w, req)

		// Verify 500 status
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}

		// Verify error message
		if !strings.Contains(w.Body.String(), "Internal Server Error") {
			t.Errorf("Expected 'Internal Server Error', got '%s'", w.Body.String())
		}

		// Verify panic was logged
		logOutput := buf.String()
		if !strings.Contains(logOutput, "PANIC") {
			t.Errorf("Expected log to contain 'PANIC', got: %s", logOutput)
		}

		if !strings.Contains(logOutput, "something went wrong!") {
			t.Errorf("Expected log to contain panic message, got: %s", logOutput)
		}
	})

	t.Run("logs stack trace on panic", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		handler := Recovery(panicHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		logOutput := buf.String()

		// Stack trace should contain goroutine info
		if !strings.Contains(logOutput, "goroutine") {
			t.Errorf("Expected stack trace to contain 'goroutine', got: %s", logOutput)
		}
	})

	t.Run("does not interfere with normal requests", func(t *testing.T) {
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
		})

		handler := Recovery(normalHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		if w.Body.String() != "Success" {
			t.Errorf("Expected 'Success', got '%s'", w.Body.String())
		}
	})

	t.Run("recovers from panic with string value", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("unexpected error")
		})

		handler := Recovery(panicHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})

	t.Run("recovers from panic with error value", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(http.ErrAbortHandler)
		})

		handler := Recovery(panicHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}

		logOutput := buf.String()
		if !strings.Contains(logOutput, "PANIC") {
			t.Errorf("Expected panic to be logged")
		}
	})

	t.Run("recovers from panic with custom struct", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		type customError struct {
			message string
			code    int
		}

		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(customError{message: "custom error", code: 999})
		})

		handler := Recovery(panicHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}
