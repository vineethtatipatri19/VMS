package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/httputil"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

func TestCustomerHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		customer := domain.Customer{Name: "John Doe", CustomerType: "b2b"}
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
		}

		var resp httputil.Response
		json.NewDecoder(w.Body).Decode(&resp)
		if !resp.Success {
			t.Error("Expected success response")
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected error on invalid JSON")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		customer := domain.Customer{CustomerType: "b2b"} // Missing name
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected error for missing name")
		}
	})
}

func TestCustomerHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			customer: &domain.Customer{ID: "test123", Name: "Test Customer"},
		}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/test123", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "test123"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		var resp httputil.Response
		json.NewDecoder(w.Body).Decode(&resp)
		if !resp.Success {
			t.Error("Expected success response")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := &mockCustomerRepo{err: domain.ErrNotFound}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})
}

func TestCustomerHandler_List(t *testing.T) {
	repo := &mockCustomerRepo{
		customers: []*domain.Customer{
			{ID: "c1", Name: "Customer 1"},
			{ID: "c2", Name: "Customer 2"},
		},
	}
	svc := service.NewCustomerService(repo)
	handler := NewCustomerHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/customers", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp httputil.Response
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Error("Expected success response")
	}
}

func TestCustomerHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		customer := domain.Customer{Name: "Updated Name", CustomerType: "b2b"}
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPut, "/api/customers/c1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "c1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestCustomerHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		delReq := domain.DeleteRequest{Reason: "test reason", Attestation: "I CONFIRM DELETE"}
		body, _ := json.Marshal(delReq)
		req := httptest.NewRequest(http.MethodDelete, "/api/customers/c1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "c1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("MissingAttestation", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		delReq := domain.DeleteRequest{Reason: "test reason"}
		body, _ := json.Marshal(delReq)
		req := httptest.NewRequest(http.MethodDelete, "/api/customers/c1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "c1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})
}

func TestCustomerHandler_GetBalance(t *testing.T) {
	repo := &mockCustomerRepo{balance: 150.50}
	svc := service.NewCustomerService(repo)
	handler := NewCustomerHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/customers/c1/balance", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	w := httptest.NewRecorder()

	handler.GetBalance(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp httputil.Response
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Error("Expected success response")
	}
}
