package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerExpiryRoutes registers all expiry-alert-related routes
func registerExpiryRoutes(r *mux.Router, h *handlers.ExpiryHandler) {
r.HandleFunc("/expiry-alerts/generate", h.GenerateAlerts).Methods(http.MethodPost)
r.HandleFunc("/expiry-alerts", h.List).Methods(http.MethodGet)
r.HandleFunc("/expiry-alerts/pending", h.GetPending).Methods(http.MethodGet)
r.HandleFunc("/expiry-alerts/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/expiry-alerts/{id}/acknowledge", h.Acknowledge).Methods(http.MethodPost)
r.HandleFunc("/expiry-alerts/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/expiry-alerts/{id}", h.Delete).Methods(http.MethodDelete)
}
