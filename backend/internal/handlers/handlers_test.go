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
customer *domain.Customer
customers []*domain.Customer
balance  float64
err error
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
item *domain.InventoryItem
items []*domain.InventoryItem
err error
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
txn *domain.Transaction
txns []*domain.Transaction
err error
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
method: http.MethodPost,
path:   "/api/customers",
body:   domain.Customer{Name: "Test", CustomerType: "b2b"},
wantStatus: http.StatusCreated,
},
{
name: "Get returns 200",
setupRepo: func() *mockCustomerRepo {
return &mockCustomerRepo{
customer: &domain.Customer{ID: "c1", Name: "Test"},
}
},
method: http.MethodGet,
path:   "/api/customers/c1",
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
