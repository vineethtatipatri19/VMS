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

func TestPaymentHandler_CreateSchedule(t *testing.T) {
	mockPmtRepo := &mockPaymentRepo{}
	mockTxRepo := &mockTransactionRepo{}
	mockCustRepo := &mockCustomerRepo{
		customer: &domain.Customer{
			ID:           "c1",
			Name:         "Test Customer",
			CustomerType: "b2b",
			Status:       "active",
		},
	}
	paymentService := service.NewPaymentService(mockPmtRepo, mockTxRepo, mockCustRepo)
	handler := NewPaymentHandler(paymentService)

	schedule := domain.PaymentSchedule{
		CustomerID: "c1",
		AmountDue:  1000.00,
		DueDate:    time.Now().AddDate(0, 0, 30),
	}

	body, _ := json.Marshal(schedule)
	req := httptest.NewRequest(http.MethodPost, "/api/payment-schedules", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateSchedule(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPaymentHandler_GetByID(t *testing.T) {
	mockPmtRepo := &mockPaymentRepo{
		paymentSchedule: &domain.PaymentSchedule{
			ID:         "ps1",
			CustomerID: "c1",
			AmountDue:  1000.00,
		},
	}
	mockTxRepo := &mockTransactionRepo{}
	mockCustRepo := &mockCustomerRepo{}
	paymentService := service.NewPaymentService(mockPmtRepo, mockTxRepo, mockCustRepo)
	handler := NewPaymentHandler(paymentService)

	req := httptest.NewRequest(http.MethodGet, "/api/payment-schedules/ps1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "ps1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
