package router

import (
	"net/http"

	"github.com/example/pgvms/internal/handlers"
	"github.com/gorilla/mux"
)

// registerInventoryRoutes registers all inventory-related routes
func registerInventoryRoutes(r *mux.Router, h *handlers.InventoryHandler) {
	r.HandleFunc("/inventory", h.List).Methods(http.MethodGet)
	r.HandleFunc("/inventory", h.Create).Methods(http.MethodPost)
	// Register specific routes BEFORE parameterized routes to avoid conflicts
	r.HandleFunc("/inventory/expiring", h.GetExpiring).Methods(http.MethodGet)
	r.HandleFunc("/inventory/low-stock", h.GetLowStock).Methods(http.MethodGet)
	// Parameterized routes come after specific routes
	r.HandleFunc("/inventory/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/inventory/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/inventory/{id}/deduct", h.DeductStock).Methods(http.MethodPost)
}
