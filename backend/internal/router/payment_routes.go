package router

import (
"net/http"

"github.com/example/pgvms/internal/handlers"
"github.com/gorilla/mux"
)

// registerPaymentRoutes registers all payment-schedule-related routes
func registerPaymentRoutes(r *mux.Router, h *handlers.PaymentHandler) {
r.HandleFunc("/payment-schedules", h.CreateSchedule).Methods(http.MethodPost)
r.HandleFunc("/payment-schedules/{id}", h.GetByID).Methods(http.MethodGet)
r.HandleFunc("/payment-schedules/customer/{customerId}", h.ListByCustomer).Methods(http.MethodGet)
r.HandleFunc("/payment-schedules/{id}/pay", h.RecordPayment).Methods(http.MethodPost)
r.HandleFunc("/payment-schedules/{id}", h.Update).Methods(http.MethodPut)
r.HandleFunc("/payment-schedules/{id}", h.Delete).Methods(http.MethodDelete)
r.HandleFunc("/payment-schedules/overdue/{customerId}", h.GetOverdue).Methods(http.MethodGet)
}
