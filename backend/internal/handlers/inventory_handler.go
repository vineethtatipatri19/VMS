package handlers

import (
	"net/http"
	"strconv"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/httputil"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item domain.InventoryItem
	if err := httputil.DecodeJSON(r, &item); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.inventoryService.CreateItem(r.Context(), &item); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, item)
}

func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	item, err := h.inventoryService.GetItem(r.Context(), id)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sortBy")

	items, err := h.inventoryService.ListItems(r.Context(), status, sortBy)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, items)
}

func (h *InventoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var item domain.InventoryItem
	if err := httputil.DecodeJSON(r, &item); err != nil {
		httputil.SendError(w, err)
		return
	}

	item.ID = id

	if err := h.inventoryService.UpdateItem(r.Context(), &item); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req domain.DeleteRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.inventoryService.DeleteItem(r.Context(), id, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "inventory item deleted successfully"})
}

func (h *InventoryHandler) GetExpiring(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	items, err := h.inventoryService.GetExpiringItems(r.Context(), days)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, items)
}

func (h *InventoryHandler) GetLowStock(w http.ResponseWriter, r *http.Request) {
	items, err := h.inventoryService.GetLowStockItems(r.Context())
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	// Wrap in object for consistent API response
	response := map[string]interface{}{
		"items": items,
	}
	httputil.SendJSON(w, http.StatusOK, response)
}

func (h *InventoryHandler) DeductStock(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.inventoryService.DeductStock(r.Context(), id, req.Quantity); err != nil {
		httputil.SendError(w, err)
		return
	}

	// Return updated inventory item
	item, err := h.inventoryService.GetItem(r.Context(), id)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, item)
}
