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
	transaction  *domain.Transaction
	transactions []*domain.Transaction
	err          error
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
	if m.transaction != nil {
		return m.transaction, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockTransactionRepo) List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	return m.transactions, m.err
}

func (m *mockTransactionRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	return m.transactions, m.err
}

func (m *mockTransactionRepo) Update(ctx context.Context, txn *domain.Transaction) error {
	return m.err
}

func (m *mockTransactionRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

// mockSaleItemRepo for testing
type mockSaleItemRepo struct {
	saleItem  *domain.SaleItem
	saleItems []*domain.SaleItem
	err       error
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
	if m.saleItem != nil {
		return m.saleItem, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockSaleItemRepo) ListByTransaction(ctx context.Context, txnID string) ([]*domain.SaleItem, error) {
	return m.saleItems, m.err
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

// mockWastageRepo for testing
type mockWastageRepo struct {
	wastage   *domain.Wastage
	wastages  []*domain.Wastage
	totalCost float64
	err       error
}

func (m *mockWastageRepo) Create(ctx context.Context, wastage *domain.Wastage) error {
	if m.err != nil {
		return m.err
	}
	wastage.ID = "w123"
	return nil
}

func (m *mockWastageRepo) GetByID(ctx context.Context, id string) (*domain.Wastage, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.wastage != nil {
		return m.wastage, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockWastageRepo) List(ctx context.Context, startDate, endDate time.Time) ([]*domain.Wastage, error) {
	return m.wastages, m.err
}

func (m *mockWastageRepo) Update(ctx context.Context, wastage *domain.Wastage) error {
	return m.err
}

func (m *mockWastageRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockWastageRepo) CalculateTotalCost(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	return m.totalCost, m.err
}

// mockExpiryRepo for testing
type mockExpiryRepo struct {
	alert  *domain.ExpiryAlert
	alerts []*domain.ExpiryAlert
	err    error
}

func (m *mockExpiryRepo) Create(ctx context.Context, alert *domain.ExpiryAlert) error {
	if m.err != nil {
		return m.err
	}
	alert.ID = "e123"
	return nil
}

func (m *mockExpiryRepo) GetByID(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.alert != nil {
		return m.alert, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockExpiryRepo) List(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
	return m.alerts, m.err
}

func (m *mockExpiryRepo) GetPendingAlerts(ctx context.Context) ([]*domain.ExpiryAlert, error) {
	return m.alerts, m.err
}

func (m *mockExpiryRepo) Update(ctx context.Context, alert *domain.ExpiryAlert) error {
	return m.err
}

func (m *mockExpiryRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockExpiryRepo) Acknowledge(ctx context.Context, id string, acknowledgedBy string) error {
	return m.err
}

func (m *mockExpiryRepo) CreateBatch(ctx context.Context, alerts []*domain.ExpiryAlert) error {
	return m.err
}

// mockPaymentRepo for testing
type mockPaymentRepo struct {
	paymentSchedule *domain.PaymentSchedule
	schedules       []*domain.PaymentSchedule
	err             error
}

func (m *mockPaymentRepo) Create(ctx context.Context, schedule *domain.PaymentSchedule) error {
	if m.err != nil {
		return m.err
	}
	schedule.ID = "ps123"
	return nil
}

func (m *mockPaymentRepo) GetByID(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.paymentSchedule != nil {
		return m.paymentSchedule, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPaymentRepo) ListByCustomer(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	return m.schedules, m.err
}

func (m *mockPaymentRepo) GetOverdueSchedules(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	return m.schedules, m.err
}

func (m *mockPaymentRepo) Update(ctx context.Context, schedule *domain.PaymentSchedule) error {
	return m.err
}

func (m *mockPaymentRepo) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *mockPaymentRepo) RecordPayment(ctx context.Context, scheduleID string, paidAmount float64, paymentDate time.Time) error {
	return m.err
}

func (m *mockPaymentRepo) UpdateStatus(ctx context.Context, scheduleID string, status string) error {
	return m.err
}
