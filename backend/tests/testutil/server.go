package testutil

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// SetupTestServer creates a test HTTP server with gorilla/mux
func SetupTestServer(t *testing.T) (*mux.Router, *httptest.Server) {
	t.Helper()

	// Create router
	router := mux.NewRouter()

	// Create test server
	server := httptest.NewServer(router)

	return router, server
}

// TeardownTestServer closes the test server
func TeardownTestServer(t *testing.T, server *httptest.Server) {
	t.Helper()

	if server != nil {
		server.Close()
	}
}

// NewTestRouter creates a mux router for testing
func NewTestRouter() *mux.Router {
	return mux.NewRouter()
}
