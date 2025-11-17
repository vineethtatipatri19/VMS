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

func TestSaleItemHandler_Create(t *testing.T) {
	mockRepo := &mockSaleItemRepo{}
	mockInvRepo := &mockInventoryRepo{
		item: &domain.InventoryItem{
			ID:        "i1",
			Name:      "Test Item",
			LotNumber: "LOT123",
			Unit:      "kg",
			Quantity:  100,
		},
	}
	saleItemService := service.NewSaleItemService(mockRepo, mockInvRepo)
	handler := NewSaleItemHandler(saleItemService)

	item := domain.SaleItem{
		TransactionID:  "t1",
		InventoryLotID: "i1",
		ItemName:       "Test Item",
		Quantity:       2,
		PricePerUnit:   75.00,
		Total:          150.00,
	}

	body, _ := json.Marshal(item)
req := httptest.NewRequest(http.MethodPost, "/api/sale-items", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.Create(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
}
}

func TestSaleItemHandler_GetByID(t *testing.T) {
mockRepo := &mockSaleItemRepo{
saleItem: &domain.SaleItem{
ID:            "si1",
TransactionID: "t1",
Quantity:      2,
},
}
mockInvRepo := &mockInventoryRepo{}
saleItemService := service.NewSaleItemService(mockRepo, mockInvRepo)
handler := NewSaleItemHandler(saleItemService)

req := httptest.NewRequest(http.MethodGet, "/api/sale-items/si1", nil)
req = mux.SetURLVars(req, map[string]string{"id": "si1"})
w := httptest.NewRecorder()

handler.GetByID(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected 200, got %d", w.Code)
}
}
