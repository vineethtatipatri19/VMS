package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

func TestInventoryHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		item := domain.InventoryItem{
			Name:          "Test Item",
			LotNumber:     "LOT123",
			Unit:          "kg",
			Quantity:      100,
			MinStockLevel: 10,
		}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected error on invalid JSON")
		}
	})
}

func TestInventoryHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockInventoryRepo{
			item: &domain.InventoryItem{ID: "inv999", Name: "Test Item"},
		}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/inv999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "inv999"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := &mockInventoryRepo{err: domain.ErrNotFound}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})
}

func TestInventoryHandler_List(t *testing.T) {
	t.Run("WithFilters", func(t *testing.T) {
		repo := &mockInventoryRepo{
			items: []*domain.InventoryItem{
				{ID: "i1", Name: "Item 1"},
			},
		}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory?status=active&sortBy=name", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestInventoryHandler_Update(t *testing.T) {
	repo := &mockInventoryRepo{
		item: &domain.InventoryItem{ID: "i1", Name: "Old Item", LotNumber: "LOT123", Unit: "kg"},
	}
	svc := service.NewInventoryService(repo)
	handler := NewInventoryHandler(svc)

	item := domain.InventoryItem{
		Name:      "Updated Item",
		LotNumber: "LOT123",
		Unit:      "kg",
	}
	body, _ := json.Marshal(item)
	req := httptest.NewRequest(http.MethodPut, "/api/inventory/i1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "i1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInventoryHandler_Delete(t *testing.T) {
	repo := &mockInventoryRepo{
		item: &domain.InventoryItem{ID: "i1", Name: "Test Item", LotNumber: "LOT123", Unit: "kg"},
	}
	svc := service.NewInventoryService(repo)
	handler := NewInventoryHandler(svc)

	delReq := domain.DeleteRequest{Reason: "discontinued", Attestation: "I CONFIRM DELETE"}
	body, _ := json.Marshal(delReq)
	req := httptest.NewRequest(http.MethodDelete, "/api/inventory/i1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "i1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInventoryHandler_GetExpiring(t *testing.T) {
	repo := &mockInventoryRepo{items: []*domain.InventoryItem{}}
	svc := service.NewInventoryService(repo)
	handler := NewInventoryHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/inventory/expiring?days=7", nil)
	w := httptest.NewRecorder()

	handler.GetExpiring(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestInventoryHandler_GetLowStock(t *testing.T) {
	repo := &mockInventoryRepo{}
	svc := service.NewInventoryService(repo)
	handler := NewInventoryHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/inventory/low-stock", nil)
	w := httptest.NewRecorder()

	handler.GetLowStock(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestInventoryHandler_DeductStock(t *testing.T) {
	repo := &mockInventoryRepo{
		item: &domain.InventoryItem{ID: "i1", Name: "Test", Quantity: 100, Status: "active"},
	}
	svc := service.NewInventoryService(repo)
	handler := NewInventoryHandler(svc)

	reqBody := map[string]int{"quantity": 10}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/inventory/i1/deduct", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "i1"})
	w := httptest.NewRecorder()

	handler.DeductStock(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
