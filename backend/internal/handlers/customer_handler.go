package handlers

import (
"net/http"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/httputil"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

type CustomerHandler struct {
customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
return &CustomerHandler{
customerService: customerService,
}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
var customer domain.Customer
if err := httputil.DecodeJSON(r, &customer); err != nil {
httputil.SendError(w, err)
return
}

if err := h.customerService.CreateCustomer(r.Context(), &customer); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

customer, err := h.customerService.GetCustomer(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
customers, err := h.customerService.ListCustomers(r.Context())
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var customer domain.Customer
if err := httputil.DecodeJSON(r, &customer); err != nil {
httputil.SendError(w, err)
return
}

customer.ID = id

if err := h.customerService.UpdateCustomer(r.Context(), &customer); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var req domain.DeleteRequest
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

if err := h.customerService.DeleteCustomer(r.Context(), id, &req); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "customer deleted successfully"})
}

func (h *CustomerHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

balance, err := h.customerService.GetCustomerBalance(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]float64{"balance": balance})
}

func (h *CustomerHandler) Search(w http.ResponseWriter, r *http.Request) {
query := r.URL.Query().Get("q")

customers, err := h.customerService.SearchCustomers(r.Context(), query)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) GetOverdue(w http.ResponseWriter, r *http.Request) {
customers, err := h.customerService.GetOverdueCustomers(r.Context())
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, customers)
}
