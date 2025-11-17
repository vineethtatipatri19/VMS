package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

func TestTransactionHandler_Create(t *testing.T) {
	mockTxRepo := &mockTransactionRepo{}
	mockCustRepo := &mockCustomerRepo{}
	mockInvRepo := &mockInventoryRepo{}
	mockSaleRepo := &mockSaleItemRepo{}
	transactionService := service.NewTransactionService(mockTxRepo, mockCustRepo, mockInvRepo, mockSaleRepo)
	handler := NewTransactionHandler(transactionService)

	txn := domain.Transaction{
		CustomerID:  "c1",
		Type:        "sale",
		TotalAmount: 1000.00,
		Date:        time.Now(),
	}

	body, _ := json.Marshal(txn)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
	mockTxRepo := &mockTransactionRepo{
		transaction: &domain.Transaction{
			ID:          "t1",
			CustomerID:  "c1",
			Type:        "sale",
			TotalAmount: 1000.00,
		},
	}
	mockCustRepo := &mockCustomerRepo{}
	mockInvRepo := &mockInventoryRepo{}
	mockSaleRepo := &mockSaleItemRepo{}
	transactionService := service.NewTransactionService(mockTxRepo, mockCustRepo, mockInvRepo, mockSaleRepo)
	handler := NewTransactionHandler(transactionService)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions/t1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "t1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
