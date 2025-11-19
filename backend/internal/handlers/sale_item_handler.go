package handlers

import (
"net/http"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/pkg/httputil"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

type SaleItemHandler struct {
saleItemService *service.SaleItemService
}

func NewSaleItemHandler(saleItemService *service.SaleItemService) *SaleItemHandler {
return &SaleItemHandler{
saleItemService: saleItemService,
}
}

func (h *SaleItemHandler) Create(w http.ResponseWriter, r *http.Request) {
var item domain.SaleItem
if err := httputil.DecodeJSON(r, &item); err != nil {
httputil.SendError(w, err)
return
}

if err := h.saleItemService.CreateSaleItem(r.Context(), &item); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, item)
}

func (h *SaleItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

item, err := h.saleItemService.GetSaleItem(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, item)
}

func (h *SaleItemHandler) ListByTransaction(w http.ResponseWriter, r *http.Request) {
transactionID := mux.Vars(r)["transactionId"]

items, err := h.saleItemService.ListItemsForTransaction(r.Context(), transactionID)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, items)
}

func (h *SaleItemHandler) Update(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var item domain.SaleItem
if err := httputil.DecodeJSON(r, &item); err != nil {
httputil.SendError(w, err)
return
}

item.ID = id

if err := h.saleItemService.UpdateSaleItem(r.Context(), &item); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, item)
}

func (h *SaleItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var req domain.DeleteRequest
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

if err := h.saleItemService.DeleteSaleItem(r.Context(), id, &req); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "sale item deleted successfully"})
}
