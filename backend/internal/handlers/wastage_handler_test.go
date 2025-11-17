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

func TestWastageHandler_RecordWastage(t *testing.T) {
	mockRepo := &mockWastageRepo{}
	mockInvRepo := &mockInventoryRepo{
		item: &domain.InventoryItem{
			ID:        "i1",
			Name:      "Test Item",
			LotNumber: "LOT123",
			Unit:      "kg",
			Quantity:  100,
		},
	}
	wastageService := service.NewWastageService(mockRepo, mockInvRepo)
	handler := NewWastageHandler(wastageService)

	wastage := domain.Wastage{
InventoryID: "i1",
Quantity:    5,
Reason:      "Damaged",
}

body, _ := json.Marshal(wastage)
req := httptest.NewRequest(http.MethodPost, "/api/wastage", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.RecordWastage(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
}
}

func TestWastageHandler_GetByID(t *testing.T) {
mockRepo := &mockWastageRepo{
wastage: &domain.Wastage{
ID:          "w1",
InventoryID: "i1",
Quantity:    5,
},
}
mockInvRepo := &mockInventoryRepo{}
wastageService := service.NewWastageService(mockRepo, mockInvRepo)
handler := NewWastageHandler(wastageService)

req := httptest.NewRequest(http.MethodGet, "/api/wastage/w1", nil)
req = mux.SetURLVars(req, map[string]string{"id": "w1"})
w := httptest.NewRecorder()

handler.GetByID(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected 200, got %d", w.Code)
}
}
