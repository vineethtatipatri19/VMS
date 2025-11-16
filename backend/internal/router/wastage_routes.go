package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerWastageRoutes registers all wastage-related routes
func registerWastageRoutes(r *mux.Router, h *handlers.WastageHandler) {
r.HandleFunc("/wastage", h.RecordWastage).Methods(http.MethodPost)
r.HandleFunc("/wastage", h.List).Methods(http.MethodGet)
r.HandleFunc("/wastage/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/wastage/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/wastage/{id}", h.Delete).Methods(http.MethodDelete)
r.HandleFunc("/wastage/report", h.GetReport).Methods(http.MethodGet)
}
