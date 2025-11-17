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

func TestCrateHandler_IssueCrates(t *testing.T) {
	mockRepo := &mockCrateRepo{}
	mockCustRepo := &mockCustomerRepo{
		customer: &domain.Customer{
			ID:           "c1",
			Name:         "Test Customer",
			CustomerType: "b2b",
			Status:       "active",
		},
	}
	crateService := service.NewCrateService(mockRepo, mockCustRepo)
	handler := NewCrateHandler(crateService)

	crate := domain.Crate{
CustomerID:      "c1",
TransactionType: "out",
Quantity:        10,
}

body, _ := json.Marshal(crate)
req := httptest.NewRequest(http.MethodPost, "/api/crates/issue", bytes.NewReader(body))
w := httptest.NewRecorder()

handler.IssueCrates(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
}
}

func TestCrateHandler_GetByID(t *testing.T) {
mockRepo := &mockCrateRepo{
crate: &domain.Crate{
ID:              "cr1",
CustomerID:      "c1",
TransactionType: "out",
Quantity:        10,
},
}
mockCustRepo := &mockCustomerRepo{}
crateService := service.NewCrateService(mockRepo, mockCustRepo)
handler := NewCrateHandler(crateService)

req := httptest.NewRequest(http.MethodGet, "/api/crates/cr1", nil)
req = mux.SetURLVars(req, map[string]string{"id": "cr1"})
w := httptest.NewRecorder()

handler.GetByID(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected 200, got %d", w.Code)
}
}
