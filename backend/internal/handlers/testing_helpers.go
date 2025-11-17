package handlers

import (
	"context"
	"time"

	"github.com/example/pgvms/internal/domain"
)

// Test helper mocks for all handlers

// mockCustomerRepo for testing
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
	c.CreatedAt = time.Now()
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

// mockInventoryRepo for testing
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

// mockTransactionRepo for testing
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

// mockSaleItemRepo for testing
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

// mockCrateRepo for testing
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

func (m *mockCrateRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	return m.crates, m.err
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
