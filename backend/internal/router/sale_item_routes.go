package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerSaleItemRoutes registers all sale-item-related routes
func registerSaleItemRoutes(r *mux.Router, h *handlers.SaleItemHandler) {
r.HandleFunc("/sale-items", h.Create).Methods(http.MethodPost)
r.HandleFunc("/sale-items/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/sale-items/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/sale-items/{id}", h.Delete).Methods(http.MethodDelete)
r.HandleFunc("/sale-items/transaction/{transactionId}", h.ListByTransaction).Methods(http.MethodGet)
}
