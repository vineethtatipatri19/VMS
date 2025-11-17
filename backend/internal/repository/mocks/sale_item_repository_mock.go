package mocks

import (
	"context"

	"github.com/example/pgvms/internal/domain"
)

// MockSaleItemRepository is a mock implementation of repository.SaleItemRepository
type MockSaleItemRepository struct {
	CreateFunc              func(ctx context.Context, item *domain.SaleItem) error
	GetByIDFunc             func(ctx context.Context, id string) (*domain.SaleItem, error)
	ListByTransactionFunc   func(ctx context.Context, transactionID string) ([]*domain.SaleItem, error)
	UpdateFunc              func(ctx context.Context, item *domain.SaleItem) error
	DeleteFunc              func(ctx context.Context, id string, req *domain.DeleteRequest) error
	DeleteByTransactionFunc func(ctx context.Context, transactionID string, req *domain.DeleteRequest) error
}

func (m *MockSaleItemRepository) Create(ctx context.Context, item *domain.SaleItem) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, item)
	}
	return nil
}

func (m *MockSaleItemRepository) GetByID(ctx context.Context, id string) (*domain.SaleItem, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockSaleItemRepository) ListByTransaction(ctx context.Context, transactionID string) ([]*domain.SaleItem, error) {
	if m.ListByTransactionFunc != nil {
		return m.ListByTransactionFunc(ctx, transactionID)
	}
	return []*domain.SaleItem{}, nil
}

func (m *MockSaleItemRepository) Update(ctx context.Context, item *domain.SaleItem) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, item)
	}
	return nil
}

func (m *MockSaleItemRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}

func (m *MockSaleItemRepository) DeleteByTransaction(ctx context.Context, transactionID string, req *domain.DeleteRequest) error {
	if m.DeleteByTransactionFunc != nil {
		return m.DeleteByTransactionFunc(ctx, transactionID, req)
	}
	return nil
}
