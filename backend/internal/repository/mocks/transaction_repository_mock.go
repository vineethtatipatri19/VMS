package mocks

import (
	"context"
	"time"

	"github.com/example/pgvms/internal/domain"
)

// MockTransactionRepository is a mock implementation of repository.TransactionRepository
type MockTransactionRepository struct {
	CreateFunc         func(ctx context.Context, tx *domain.Transaction) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Transaction, error)
	ListByCustomerFunc func(ctx context.Context, customerID string) ([]*domain.Transaction, error)
	ListFunc           func(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error)
	UpdateFunc         func(ctx context.Context, tx *domain.Transaction) error
	DeleteFunc         func(ctx context.Context, id string, req *domain.DeleteRequest) error
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, tx)
	}
	return nil
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockTransactionRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	if m.ListByCustomerFunc != nil {
		return m.ListByCustomerFunc(ctx, customerID)
	}
	return []*domain.Transaction{}, nil
}

func (m *MockTransactionRepository) List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, txType, startDate, endDate)
	}
	return []*domain.Transaction{}, nil
}

func (m *MockTransactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, tx)
	}
	return nil
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}
