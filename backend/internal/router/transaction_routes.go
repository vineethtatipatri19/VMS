package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerTransactionRoutes registers all transaction-related routes
func registerTransactionRoutes(r *mux.Router, h *handlers.TransactionHandler) {
r.HandleFunc("/transactions", h.List).Methods(http.MethodGet)
r.HandleFunc("/transactions", h.Create).Methods(http.MethodPost)
r.HandleFunc("/transactions/sale", h.CreateSale).Methods(http.MethodPost)
r.HandleFunc("/transactions/payment", h.CreatePayment).Methods(http.MethodPost)
r.HandleFunc("/transactions/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/transactions/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/transactions/{id}", h.Delete).Methods(http.MethodDelete)
r.HandleFunc("/transactions/customer/{customerId}", h.ListByCustomer).Methods(http.MethodGet)
}
