package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS middleware
	handler := CORS(testHandler)

	t.Run("adds CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		// Verify CORS headers are set
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin to be '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
		}

		if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("Expected Access-Control-Allow-Methods header, got '%s'", w.Header().Get("Access-Control-Allow-Methods"))
		}

		if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("Expected Access-Control-Allow-Headers header, got '%s'", w.Header().Get("Access-Control-Allow-Headers"))
		}

		// Verify handler was called
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("handles preflight OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		// Verify CORS headers are set
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("Expected Access-Control-Allow-Origin header for preflight")
		}

		// Verify preflight returns 200 OK
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for preflight, got %d", w.Code)
		}

		// Verify handler body is empty (preflight doesn't call next handler)
		if w.Body.String() != "" {
			t.Errorf("Expected empty body for preflight, got '%s'", w.Body.String())
		}
	})

	t.Run("allows all HTTP methods", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

		for _, method := range methods {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Header().Get("Access-Control-Allow-Origin") != "*" {
				t.Errorf("Expected CORS headers for %s method", method)
			}

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}
		}
	})
}
