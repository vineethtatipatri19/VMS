package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerCustomerRoutes registers all customer-related routes
func registerCustomerRoutes(r *mux.Router, h *handlers.CustomerHandler) {
r.HandleFunc("/customers", h.List).Methods(http.MethodGet)
r.HandleFunc("/customers", h.Create).Methods(http.MethodPost)
r.HandleFunc("/customers/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/customers/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/customers/{id}", h.Delete).Methods(http.MethodDelete)
r.HandleFunc("/customers/{id}/balance", h.GetBalance).Methods(http.MethodGet)
r.HandleFunc("/customers/search", h.Search).Methods(http.MethodGet)
r.HandleFunc("/customers/overdue", h.GetOverdue).Methods(http.MethodGet)
}
