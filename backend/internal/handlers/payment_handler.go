package handlers

import (
"net/http"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/pkg/httputil"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

type PaymentHandler struct {
paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
return &PaymentHandler{
paymentService: paymentService,
}
}

func (h *PaymentHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
var schedule domain.PaymentSchedule
if err := httputil.DecodeJSON(r, &schedule); err != nil {
httputil.SendError(w, err)
return
}

if err := h.paymentService.CreateSchedule(r.Context(), &schedule); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, schedule)
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

schedule, err := h.paymentService.GetSchedule(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, schedule)
}

func (h *PaymentHandler) ListByCustomer(w http.ResponseWriter, r *http.Request) {
customerID := mux.Vars(r)["customerId"]

schedules, err := h.paymentService.ListSchedules(r.Context(), customerID)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, schedules)
}

func (h *PaymentHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var req struct {
PaidAmount  float64 `json:"paidAmount"`
PaymentDate string  `json:"paymentDate"`
}
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

paymentDate := time.Now()
if req.PaymentDate != "" {
paymentDate, _ = time.Parse("2006-01-02", req.PaymentDate)
}

if err := h.paymentService.RecordPayment(r.Context(), id, req.PaidAmount, paymentDate); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "payment recorded successfully"})
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var schedule domain.PaymentSchedule
if err := httputil.DecodeJSON(r, &schedule); err != nil {
httputil.SendError(w, err)
return
}

schedule.ID = id

if err := h.paymentService.UpdateSchedule(r.Context(), &schedule); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, schedule)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

if err := h.paymentService.DeleteSchedule(r.Context(), id); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "payment schedule deleted successfully"})
}

func (h *PaymentHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
customerID := mux.Vars(r)["customerId"]

schedules, err := h.paymentService.GetOverdueSchedules(r.Context(), customerID)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, schedules)
}
