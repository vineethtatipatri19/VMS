package service

import (
	"context"
	"errors"
	"testing"

	"github.com/example/pgvms/internal/domain"
)

// Mock repository for testing
type mockCustomerRepo struct {
	customers    map[string]*domain.Customer
	createErr    error
	getErr       error
	updateErr    error
	deleteErr    error
	listErr      error
	balanceErr   error
	updateBalErr error
	searchErr    error
}

func newMockCustomerRepo() *mockCustomerRepo {
	return &mockCustomerRepo{
		customers: make(map[string]*domain.Customer),
	}
}

func (m *mockCustomerRepo) Create(ctx context.Context, c *domain.Customer) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.customers[c.ID] = c
	return nil
}

func (m *mockCustomerRepo) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	c, ok := m.customers[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (m *mockCustomerRepo) List(ctx context.Context, status string) ([]*domain.Customer, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*domain.Customer
	for _, c := range m.customers {
		if status == "" || c.Status == status {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCustomerRepo) Update(ctx context.Context, c *domain.Customer) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.customers[c.ID]; !ok {
		return domain.ErrNotFound
	}
	m.customers[c.ID] = c
	return nil
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.customers[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.customers, id)
	return nil
}

func (m *mockCustomerRepo) GetBalance(ctx context.Context, id string) (float64, error) {
	if m.balanceErr != nil {
		return 0, m.balanceErr
	}
	c, ok := m.customers[id]
	if !ok {
		return 0, domain.ErrNotFound
	}
	return c.CurrentBalance, nil
}

func (m *mockCustomerRepo) UpdateBalance(ctx context.Context, id string, delta float64) error {
	if m.updateBalErr != nil {
		return m.updateBalErr
	}
	c, ok := m.customers[id]
	if !ok {
		return domain.ErrNotFound
	}
	c.CurrentBalance += delta
	return nil
}

func (m *mockCustomerRepo) Search(ctx context.Context, query string) ([]*domain.Customer, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	var result []*domain.Customer
	for _, c := range m.customers {
		result = append(result, c)
	}
	return result, nil
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	tests := []struct {
		name      string
		customer  *domain.Customer
		repoErr   error
		expectErr bool
	}{
		{
			name: "successful creation",
			customer: &domain.Customer{
				ID:             "cust-1",
				Name:           "John Doe",
				ContactNumber:  "1234567890",
				Status:         "active",
				CurrentBalance: 0,
			},
			repoErr:   nil,
			expectErr: false,
		},
		{
			name: "repository error",
			customer: &domain.Customer{
				ID:   "cust-2",
				Name: "Jane Doe",
			},
			repoErr:   errors.New("db error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockCustomerRepo()
			repo.createErr = tt.repoErr
			svc := NewCustomerService(repo)

			err := svc.CreateCustomer(context.Background(), tt.customer)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateCustomer() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestCustomerService_GetCustomer(t *testing.T) {
	repo := newMockCustomerRepo()
	svc := NewCustomerService(repo)

	// Setup test data
	customer := &domain.Customer{
		ID:             "cust-1",
		Name:           "John Doe",
		ContactNumber:  "1234567890",
		Status:         "active",
		CurrentBalance: 100.0,
	}
	repo.customers["cust-1"] = customer

	tests := []struct {
		name      string
		id        string
		expectErr bool
	}{
		{
			name:      "customer found",
			id:        "cust-1",
			expectErr: false,
		},
		{
			name:      "customer not found",
			id:        "cust-999",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetCustomer(context.Background(), tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetCustomer() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && result == nil {
				t.Error("GetCustomer() expected customer, got nil")
			}
		})
	}
}

func TestCustomerService_CheckCreditLimit(t *testing.T) {
	repo := newMockCustomerRepo()
	svc := NewCustomerService(repo)

	// Customer with credit limit
	customer := &domain.Customer{
		ID:             "cust-1",
		Name:           "John Doe",
		CreditLimit:    1000.0,
		CurrentBalance: 500.0, // owes 500
	}
	repo.customers["cust-1"] = customer

	tests := []struct {
		name        string
		id          string
		amount      float64
		canPurchase bool
	}{
		{
			name:        "within credit limit",
			id:          "cust-1",
			amount:      400.0,
			canPurchase: true,
		},
		{
			name:        "exceeds credit limit",
			id:          "cust-1",
			amount:      600.0,
			canPurchase: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CheckCreditLimit(context.Background(), tt.id, tt.amount)
			canPurchase := (err == nil)
			if canPurchase != tt.canPurchase {
				t.Errorf("CheckCreditLimit() canPurchase = %v, want %v", canPurchase, tt.canPurchase)
			}
		})
	}
}

func TestCustomerService_GetOverdueCustomers(t *testing.T) {
	repo := newMockCustomerRepo()
	svc := NewCustomerService(repo)

	// Setup test data
	overdue := &domain.Customer{
		ID:             "cust-1",
		Name:           "Overdue Customer",
		CurrentBalance: 100.0, // positive means owes money
		Status:         "active",
	}
	current := &domain.Customer{
		ID:             "cust-2",
		Name:           "Current Customer",
		CurrentBalance: -50.0,
		Status:         "active",
	}
	repo.customers["cust-1"] = overdue
	repo.customers["cust-2"] = current

	result, err := svc.GetOverdueCustomers(context.Background())
	if err != nil {
		t.Fatalf("GetOverdueCustomers() error = %v", err)
	}

	// Should return at least one overdue customer
	if len(result) == 0 {
		t.Error("GetOverdueCustomers() expected overdue customers, got none")
	}
}
