package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

// mockRepo provides minimal repository implementation for handler tests
type mockRepo struct{}

func (m *mockRepo) Create(ctx context.Context, item interface{}) error {
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (interface{}, error) {
	return nil, domain.ErrNotFound
}

func (m *mockRepo) List(ctx context.Context, args ...interface{}) (interface{}, error) {
	return []interface{}{}, nil
}

func (m *mockRepo) Update(ctx context.Context, item interface{}) error {
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return nil
}

// mockCustomerRepo minimal implementation
type mockCustomerRepo struct {
	customer  *domain.Customer
	customers []*domain.Customer
	balance   float64
	err       error
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *domain.Customer) error {
	if m.err != nil {
		return m.err
	}
	c.ID = "cust123"
	return nil
}

func (m *mockCustomerRepo) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.customer != nil {
		return m.customer, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockCustomerRepo) List(ctx context.Context) ([]*domain.Customer, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.customers, nil
}

func (m *mockCustomerRepo) Update(ctx context.Context, c *domain.Customer) error {
	return m.err
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockCustomerRepo) GetBalance(ctx context.Context, customerID string) (float64, error) {
	return m.balance, m.err
}

func (m *mockCustomerRepo) UpdateBalance(ctx context.Context, customerID string, delta float64) error {
	return m.err
}

func (m *mockCustomerRepo) UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error {
	return nil
}

// mockInventoryRepo minimal implementation
type mockInventoryRepo struct {
	item  *domain.InventoryItem
	items []*domain.InventoryItem
	err   error
}

func (m *mockInventoryRepo) Create(ctx context.Context, item *domain.InventoryItem) error {
	if m.err != nil {
		return m.err
	}
	item.ID = "inv123"
	return nil
}

func (m *mockInventoryRepo) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.item != nil {
		return m.item, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockInventoryRepo) List(ctx context.Context, status, sortBy string) ([]*domain.InventoryItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

func (m *mockInventoryRepo) Update(ctx context.Context, item *domain.InventoryItem) error {
	return m.err
}

func (m *mockInventoryRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockInventoryRepo) UpdateQuantity(ctx context.Context, id string, delta float64) error {
	return m.err
}

func (m *mockInventoryRepo) GetExpiringSoon(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
	return m.items, m.err
}

func (m *mockInventoryRepo) GetLowStock(ctx context.Context) ([]*domain.InventoryItem, error) {
	return m.items, m.err
}

// mockTransactionRepo minimal implementation
type mockTransactionRepo struct {
	txn  *domain.Transaction
	txns []*domain.Transaction
	err  error
}

func (m *mockTransactionRepo) Create(ctx context.Context, txn *domain.Transaction) error {
	if m.err != nil {
		return m.err
	}
	txn.ID = "tx123"
	return nil
}

func (m *mockTransactionRepo) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.txn != nil {
		return m.txn, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockTransactionRepo) List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	return m.txns, m.err
}

func (m *mockTransactionRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	return m.txns, m.err
}

func (m *mockTransactionRepo) Update(ctx context.Context, txn *domain.Transaction) error {
	return m.err
}

func (m *mockTransactionRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func TestHandlers_HTTPMethods(t *testing.T) {
	// Test that handlers properly decode JSON, call services, and encode responses

	t.Run("CustomerHandler_InvalidJSON", func(t *testing.T) {
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

	t.Run("CustomerHandler_GetByID_NotFound", func(t *testing.T) {
		repo := &mockCustomerRepo{err: domain.ErrNotFound}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code == http.StatusOK {
			t.Error("Expected 404 for not found")
		}
	})

	t.Run("CustomerHandler_List", func(t *testing.T) {
		repo := &mockCustomerRepo{
			customers: []*domain.Customer{
				{ID: "c1", Name: "Customer 1"},
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
	})

	t.Run("InventoryHandler_InvalidJSON", func(t *testing.T) {
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

	t.Run("InventoryHandler_List_WithQueryParams", func(t *testing.T) {
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

	t.Run("InventoryHandler_GetExpiring", func(t *testing.T) {
		repo := &mockInventoryRepo{
			items: []*domain.InventoryItem{},
		}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/expiring?days=7", nil)
		w := httptest.NewRecorder()

		handler.GetExpiring(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("InventoryHandler_GetLowStock", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/low-stock", nil)
		w := httptest.NewRecorder()

		handler.GetLowStock(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("TransactionHandler_InvalidJSON", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		var saleItemRepo repository.SaleItemRepository = nil

		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleItemRepo)
		handler := NewTransactionHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected error on invalid JSON")
		}
	})

	t.Run("TransactionHandler_List", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{
			txns: []*domain.Transaction{},
		}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		var saleItemRepo repository.SaleItemRepository = nil

		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleItemRepo)
		handler := NewTransactionHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestHandlers_URLVariables(t *testing.T) {
	// Test that handlers correctly extract URL variables

	t.Run("CustomerHandler_GetByID_ExtractsID", func(t *testing.T) {
		repo := &mockCustomerRepo{
			customer: &domain.Customer{ID: "test123", Name: "Test"},
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

		var result domain.Customer
		json.NewDecoder(w.Body).Decode(&result)
		if result.ID != "test123" {
			t.Errorf("Expected ID test123, got %s", result.ID)
		}
	})

	t.Run("InventoryHandler_GetByID_ExtractsID", func(t *testing.T) {
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
}

func TestHandlers_StatusCodes(t *testing.T) {
	// Test that handlers return correct HTTP status codes

	tests := []struct {
		name       string
		setupRepo  func() *mockCustomerRepo
		method     string
		path       string
		body       interface{}
		wantStatus int
	}{
		{
			name: "Create returns 201",
			setupRepo: func() *mockCustomerRepo {
				return &mockCustomerRepo{}
			},
			method:     http.MethodPost,
			path:       "/api/customers",
			body:       domain.Customer{Name: "Test", CustomerType: "b2b"},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Get returns 200",
			setupRepo: func() *mockCustomerRepo {
				return &mockCustomerRepo{
					customer: &domain.Customer{ID: "c1", Name: "Test"},
				}
			},
			method:     http.MethodGet,
			path:       "/api/customers/c1",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			svc := service.NewCustomerService(repo)
			handler := NewCustomerHandler(svc)

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
			if tt.method == http.MethodGet {
				req = mux.SetURLVars(req, map[string]string{"id": "c1"})
			}
			w := httptest.NewRecorder()

			if tt.method == http.MethodPost {
				handler.Create(w, req)
			} else if tt.method == http.MethodGet {
				handler.GetByID(w, req)
			}

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

// mockSaleItemRepo minimal implementation
type mockSaleItemRepo struct {
	item  *domain.SaleItem
	items []*domain.SaleItem
	err   error
}

func (m *mockSaleItemRepo) Create(ctx context.Context, item *domain.SaleItem) error {
	if m.err != nil {
		return m.err
	}
	item.ID = "si123"
	return nil
}

func (m *mockSaleItemRepo) GetByID(ctx context.Context, id string) (*domain.SaleItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.item != nil {
		return m.item, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockSaleItemRepo) ListByTransaction(ctx context.Context, txnID string) ([]*domain.SaleItem, error) {
	return m.items, m.err
}

func (m *mockSaleItemRepo) Update(ctx context.Context, item *domain.SaleItem) error {
	return m.err
}

func (m *mockSaleItemRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockSaleItemRepo) DeleteByTransaction(ctx context.Context, transactionID string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockSaleItemRepo) CreateBatch(ctx context.Context, items []*domain.SaleItem) error {
	return m.err
}

// mockCrateRepo minimal implementation
type mockCrateRepo struct {
	crate   *domain.CrateEntry
	crates  []*domain.CrateEntry
	balance int
	err     error
}

func (m *mockCrateRepo) Create(ctx context.Context, crate *domain.CrateEntry) error {
	if m.err != nil {
		return m.err
	}
	crate.ID = "cr123"
	return nil
}

func (m *mockCrateRepo) GetByID(ctx context.Context, id string) (*domain.CrateEntry, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.crate != nil {
		return m.crate, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockCrateRepo) Update(ctx context.Context, crate *domain.CrateEntry) error {
	return m.err
}

func (m *mockCrateRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockCrateRepo) GetBalance(ctx context.Context, customerID string) (int, error) {
	return m.balance, m.err
}

func (m *mockCrateRepo) GetHistory(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	return m.crates, m.err
}

func (m *mockCrateRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	return m.crates, m.err
}

// mockWastageRepo minimal implementation
type mockWastageRepo struct {
	wastage  *domain.WastageLog
	wastages []*domain.WastageLog
	err      error
}

func (m *mockWastageRepo) Create(ctx context.Context, w *domain.WastageLog) error {
	if m.err != nil {
		return m.err
	}
	w.ID = "w123"
	return nil
}

func (m *mockWastageRepo) GetByID(ctx context.Context, id string) (*domain.WastageLog, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.wastage != nil {
		return m.wastage, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockWastageRepo) List(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error) {
	return m.wastages, m.err
}

func (m *mockWastageRepo) Update(ctx context.Context, w *domain.WastageLog) error {
	return m.err
}

func (m *mockWastageRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

// mockExpiryRepo minimal implementation
type mockExpiryRepo struct {
	expiry   *domain.ExpiryAlert
	expiries []*domain.ExpiryAlert
	err      error
}

func (m *mockExpiryRepo) Create(ctx context.Context, e *domain.ExpiryAlert) error {
	if m.err != nil {
		return m.err
	}
	e.ID = "e123"
	return nil
}

func (m *mockExpiryRepo) GetByID(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.expiry != nil {
		return m.expiry, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockExpiryRepo) ListByInventory(ctx context.Context, inventoryID string) ([]*domain.ExpiryAlert, error) {
	return m.expiries, m.err
}

func (m *mockExpiryRepo) ListPending(ctx context.Context, days int) ([]*domain.ExpiryAlert, error) {
	return m.expiries, m.err
}

func (m *mockExpiryRepo) Update(ctx context.Context, e *domain.ExpiryAlert) error {
	return m.err
}

func (m *mockExpiryRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

// mockPaymentRepo minimal implementation
type mockPaymentRepo struct {
	payment  *domain.PaymentSchedule
	payments []*domain.PaymentSchedule
	err      error
}

func (m *mockPaymentRepo) Create(ctx context.Context, p *domain.PaymentSchedule) error {
	if m.err != nil {
		return m.err
	}
	p.ID = "p123"
	return nil
}

func (m *mockPaymentRepo) GetByID(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.payment != nil {
		return m.payment, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPaymentRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	return m.payments, m.err
}

func (m *mockPaymentRepo) ListByTransaction(ctx context.Context, transactionID string) ([]*domain.PaymentSchedule, error) {
	return m.payments, m.err
}

func (m *mockPaymentRepo) Update(ctx context.Context, p *domain.PaymentSchedule) error {
	return m.err
}

func (m *mockPaymentRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

// Comprehensive Handler Tests

func TestCustomerHandler_CRUD(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		repo := &mockCustomerRepo{}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		customer := domain.Customer{Name: "John Doe", CustomerType: "b2b"}
		body, _ := json.Marshal(customer)
		req := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("Create_ValidationError", func(t *testing.T) {
		repo := &mockCustomerRepo{err: domain.ErrInvalidField("name", "name is required")}
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

	t.Run("Update_Success", func(t *testing.T) {
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
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Delete_Success", func(t *testing.T) {
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
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("GetBalance_Success", func(t *testing.T) {
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
	})
}

func TestInventoryHandler_CRUD(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		item := domain.InventoryItem{
			Name:          "Test Item",
			Unit:          "kg",
			Quantity:      100,
			MinStockLevel: 10,
		}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/api/inventory", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("Update_Success", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		item := domain.InventoryItem{Name: "Updated Item", Unit: "kg"}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPut, "/api/inventory/i1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "i1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Delete_WithReason", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		delReq := domain.DeleteRequest{Reason: "discontinued", Attestation: "I CONFIRM DELETE"}
		body, _ := json.Marshal(delReq)
		req := httptest.NewRequest(http.MethodDelete, "/api/inventory/i1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "i1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("DeductStock_Success", func(t *testing.T) {
		repo := &mockInventoryRepo{}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		reqBody := map[string]int{"quantity": 10}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/inventory/i1/deduct", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "i1"})
		w := httptest.NewRecorder()

		handler.DeductStock(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestTransactionHandler_CRUD(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		txn := domain.Transaction{
			Type:        "sale",
			CustomerID:  "c1",
			TotalAmount: 100.0,
		}
		body, _ := json.Marshal(txn)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("CreateSale_WithItems", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		saleReq := map[string]interface{}{
			"transaction": domain.Transaction{Type: "sale", CustomerID: "c1"},
			"items": []domain.SaleItem{
				{InventoryLotID: "i1", Quantity: 10, PricePerUnit: 5.0},
			},
		}
		body, _ := json.Marshal(saleReq)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreateSale(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("CreatePayment_Success", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		txn := domain.Transaction{Type: "payment", CustomerID: "c1", TotalAmount: 50.0}
		body, _ := json.Marshal(txn)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions/payment", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreatePayment(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("ListByCustomer_Success", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{txns: []*domain.Transaction{}}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/c1/transactions", nil)
		req = mux.SetURLVars(req, map[string]string{"customerId": "c1"})
		w := httptest.NewRecorder()

		handler.ListByCustomer(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Update_Success", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		txn := domain.Transaction{Type: "sale", TotalAmount: 200.0}
		body, _ := json.Marshal(txn)
		req := httptest.NewRequest(http.MethodPut, "/api/transactions/tx1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "tx1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestSaleItemHandler_CRUD(t *testing.T) {
	t.Run("Create_Success", func(t *testing.T) {
		repo := &mockSaleItemRepo{}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		item := domain.SaleItem{
			TransactionID:  "tx1",
			InventoryLotID: "i1",
			Quantity:       5,
			PricePerUnit:   10.0,
		}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/api/sale-items", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		repo := &mockSaleItemRepo{
			item: &domain.SaleItem{ID: "si1", Quantity: 5},
		}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/sale-items/si1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "si1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("ListByTransaction_Success", func(t *testing.T) {
		repo := &mockSaleItemRepo{items: []*domain.SaleItem{}}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/transactions/tx1/items", nil)
		req = mux.SetURLVars(req, map[string]string{"transactionId": "tx1"})
		w := httptest.NewRecorder()

		handler.ListByTransaction(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Update_Success", func(t *testing.T) {
		repo := &mockSaleItemRepo{}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		item := domain.SaleItem{Quantity: 10, PricePerUnit: 15.0}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPut, "/api/sale-items/si1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "si1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Delete_Success", func(t *testing.T) {
		repo := &mockSaleItemRepo{}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		delReq := domain.DeleteRequest{Reason: "error", Attestation: "I CONFIRM DELETE"}
		body, _ := json.Marshal(delReq)
		req := httptest.NewRequest(http.MethodDelete, "/api/sale-items/si1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "si1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestCrateHandler_Operations(t *testing.T) {
	t.Run("IssueCrates_Success", func(t *testing.T) {
		repo := &mockCrateRepo{}
		custRepo := &mockCustomerRepo{}
		svc := service.NewCrateService(repo, custRepo)
		handler := NewCrateHandler(svc)

		crate := domain.Crate{
			CustomerID:      "c1",
			Quantity:        10,
			TransactionType: "out",
		}
		body, _ := json.Marshal(crate)
		req := httptest.NewRequest(http.MethodPost, "/api/crates/issue", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.IssueCrates(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("ReturnCrates_Success", func(t *testing.T) {
		repo := &mockCrateRepo{}
		custRepo := &mockCustomerRepo{}
		svc := service.NewCrateService(repo, custRepo)
		handler := NewCrateHandler(svc)

		crate := domain.Crate{
			CustomerID:      "c1",
			Quantity:        5,
			TransactionType: "in",
		}
		body, _ := json.Marshal(crate)
		req := httptest.NewRequest(http.MethodPost, "/api/crates/return", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.ReturnCrates(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", w.Code)
		}
	})

	t.Run("GetBalance_Success", func(t *testing.T) {
		repo := &mockCrateRepo{balance: 15}
		custRepo := &mockCustomerRepo{}
		svc := service.NewCrateService(repo, custRepo)
		handler := NewCrateHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/c1/crates/balance", nil)
		req = mux.SetURLVars(req, map[string]string{"customerId": "c1"})
		w := httptest.NewRecorder()

		handler.GetBalance(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("GetHistory_Success", func(t *testing.T) {
		repo := &mockCrateRepo{crates: []*domain.CrateEntry{}}
		custRepo := &mockCustomerRepo{}
		svc := service.NewCrateService(repo, custRepo)
		handler := NewCrateHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/c1/crates/history", nil)
		req = mux.SetURLVars(req, map[string]string{"customerId": "c1"})
		w := httptest.NewRecorder()

		handler.GetHistory(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		repo := &mockCrateRepo{err: domain.ErrNotFound}
		custRepo := &mockCustomerRepo{}
		svc := service.NewCrateService(repo, custRepo)
		handler := NewCrateHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/crates/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code == http.StatusOK {
			t.Error("Expected error for not found")
		}
	})
}

func TestHandlers_ErrorScenarios(t *testing.T) {
	t.Run("Customer_NotFound", func(t *testing.T) {
		repo := &mockCustomerRepo{err: domain.ErrNotFound}
		svc := service.NewCustomerService(repo)
		handler := NewCustomerHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/customers/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code == http.StatusOK {
			t.Error("Expected 404 for not found customer")
		}
	})

	t.Run("Inventory_NotFound", func(t *testing.T) {
		repo := &mockInventoryRepo{err: domain.ErrNotFound}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code == http.StatusOK {
			t.Error("Expected 404 for not found inventory")
		}
	})

	t.Run("Transaction_NotFound", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{err: domain.ErrNotFound}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/transactions/nonexistent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		if w.Code == http.StatusOK {
			t.Error("Expected 404 for not found transaction")
		}
	})

	t.Run("SaleItem_InvalidJSON", func(t *testing.T) {
		repo := &mockSaleItemRepo{}
		invRepo := &mockInventoryRepo{}
		svc := service.NewSaleItemService(repo, invRepo)
		handler := NewSaleItemHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/api/sale-items", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected error for invalid JSON")
		}
	})
}

func TestHandlers_QueryParameters(t *testing.T) {
	t.Run("Inventory_List_WithFilters", func(t *testing.T) {
		repo := &mockInventoryRepo{items: []*domain.InventoryItem{}}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory?status=active&sortBy=expiry_date", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Transaction_List_WithDateRange", func(t *testing.T) {
		txnRepo := &mockTransactionRepo{txns: []*domain.Transaction{}}
		custRepo := &mockCustomerRepo{}
		invRepo := &mockInventoryRepo{}
		saleRepo := &mockSaleItemRepo{}
		svc := service.NewTransactionService(txnRepo, custRepo, invRepo, saleRepo)
		handler := NewTransactionHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/transactions?startDate=2024-01-01&endDate=2024-12-31", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Inventory_GetExpiring_WithDays", func(t *testing.T) {
		repo := &mockInventoryRepo{items: []*domain.InventoryItem{}}
		svc := service.NewInventoryService(repo)
		handler := NewInventoryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/api/inventory/expiring?days=30", nil)
		w := httptest.NewRecorder()

		handler.GetExpiring(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}
