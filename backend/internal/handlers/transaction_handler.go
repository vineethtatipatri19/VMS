package handlers

import (
"net/http"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/httputil"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

type TransactionHandler struct {
transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
return &TransactionHandler{
transactionService: transactionService,
}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
var txn domain.Transaction
if err := httputil.DecodeJSON(r, &txn); err != nil {
httputil.SendError(w, err)
return
}

if err := h.transactionService.CreateTransaction(r.Context(), &txn); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
var req struct {
Transaction domain.Transaction `json:"transaction"`
Items       []domain.SaleItem  `json:"items"`
}
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

var items []*domain.SaleItem
for i := range req.Items {
items = append(items, &req.Items[i])
}

if err := h.transactionService.CreateSale(r.Context(), &req.Transaction, items); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, req.Transaction)
}

func (h *TransactionHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
var txn domain.Transaction
if err := httputil.DecodeJSON(r, &txn); err != nil {
httputil.SendError(w, err)
return
}

if err := h.transactionService.CreatePayment(r.Context(), &txn); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, txn)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

txn, err := h.transactionService.GetTransaction(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, txn)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
txType := r.URL.Query().Get("type")
startDateStr := r.URL.Query().Get("startDate")
endDateStr := r.URL.Query().Get("endDate")

var startDate, endDate time.Time
if startDateStr != "" {
startDate, _ = time.Parse("2006-01-02", startDateStr)
}
if endDateStr != "" {
endDate, _ = time.Parse("2006-01-02", endDateStr)
}

transactions, err := h.transactionService.ListTransactions(r.Context(), txType, startDate, endDate)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) ListByCustomer(w http.ResponseWriter, r *http.Request) {
customerID := mux.Vars(r)["customerId"]

transactions, err := h.transactionService.ListCustomerTransactions(r.Context(), customerID)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var txn domain.Transaction
if err := httputil.DecodeJSON(r, &txn); err != nil {
httputil.SendError(w, err)
return
}

txn.ID = id

if err := h.transactionService.UpdateTransaction(r.Context(), &txn); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, txn)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var req domain.DeleteRequest
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

if err := h.transactionService.DeleteTransaction(r.Context(), id, &req); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "transaction deleted successfully"})
}
