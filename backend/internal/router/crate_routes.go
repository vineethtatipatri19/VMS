package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerCrateRoutes registers all crate-related routes
func registerCrateRoutes(r *mux.Router, h *handlers.CrateHandler) {
r.HandleFunc("/crates/issue", h.IssueCrates).Methods(http.MethodPost)
r.HandleFunc("/crates/return", h.ReturnCrates).Methods(http.MethodPost)
r.HandleFunc("/crates/balance/{customerId}", h.GetBalance).Methods(http.MethodGet)
r.HandleFunc("/crates/history/{customerId}", h.GetHistory).Methods(http.MethodGet)
r.HandleFunc("/crates/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/crates/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/crates/{id}", h.Delete).Methods(http.MethodDelete)
}
