package mocks

import (
"context"
"time"

"github.com/example/pgvms/internal/domain"
)

// MockCustomerRepository is a mock implementation of repository.CustomerRepository
type MockCustomerRepository struct {
CreateFunc                func(ctx context.Context, customer *domain.Customer) error
GetByIDFunc               func(ctx context.Context, id string) (*domain.Customer, error)
ListFunc                  func(ctx context.Context) ([]*domain.Customer, error)
UpdateFunc                func(ctx context.Context, customer *domain.Customer) error
DeleteFunc                func(ctx context.Context, id string, req *domain.DeleteRequest) error
GetBalanceFunc            func(ctx context.Context, customerID string) (float64, error)
UpdateBalanceFunc         func(ctx context.Context, customerID string, delta float64) error
UpdateLastTransactionFunc func(ctx context.Context, customerID string, date time.Time) error
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
if m.CreateFunc != nil {
return m.CreateFunc(ctx, customer)
}
return nil
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
if m.GetByIDFunc != nil {
return m.GetByIDFunc(ctx, id)
}
return nil, domain.ErrNotFound
}

func (m *MockCustomerRepository) List(ctx context.Context) ([]*domain.Customer, error) {
if m.ListFunc != nil {
return m.ListFunc(ctx)
}
return []*domain.Customer{}, nil
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
if m.UpdateFunc != nil {
return m.UpdateFunc(ctx, customer)
}
return nil
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
if m.DeleteFunc != nil {
return m.DeleteFunc(ctx, id, req)
}
return nil
}

func (m *MockCustomerRepository) GetBalance(ctx context.Context, customerID string) (float64, error) {
if m.GetBalanceFunc != nil {
return m.GetBalanceFunc(ctx, customerID)
}
return 0, nil
}

func (m *MockCustomerRepository) UpdateBalance(ctx context.Context, customerID string, delta float64) error {
if m.UpdateBalanceFunc != nil {
return m.UpdateBalanceFunc(ctx, customerID, delta)
}
return nil
}

func (m *MockCustomerRepository) UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error {
if m.UpdateLastTransactionFunc != nil {
return m.UpdateLastTransactionFunc(ctx, customerID, date)
}
return nil
}
