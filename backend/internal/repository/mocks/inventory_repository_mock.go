package mocks

import (
	"context"

	"github.com/example/pgvms/internal/domain"
)

// MockInventoryRepository is a mock implementation of repository.InventoryRepository
type MockInventoryRepository struct {
	CreateFunc          func(ctx context.Context, item *domain.InventoryItem) error
	GetByIDFunc         func(ctx context.Context, id string) (*domain.InventoryItem, error)
	ListFunc            func(ctx context.Context, status string, sortBy string) ([]*domain.InventoryItem, error)
	UpdateFunc          func(ctx context.Context, item *domain.InventoryItem) error
	DeleteFunc          func(ctx context.Context, id string, req *domain.DeleteRequest) error
	UpdateQuantityFunc  func(ctx context.Context, id string, delta float64) error
	GetExpiringSoonFunc func(ctx context.Context, days int) ([]*domain.InventoryItem, error)
	GetLowStockFunc     func(ctx context.Context) ([]*domain.InventoryItem, error)
}

func (m *MockInventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, item)
	}
	return nil
}

func (m *MockInventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockInventoryRepository) List(ctx context.Context, status string, sortBy string) ([]*domain.InventoryItem, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, status, sortBy)
	}
	return []*domain.InventoryItem{}, nil
}

func (m *MockInventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, item)
	}
	return nil
}

func (m *MockInventoryRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}

func (m *MockInventoryRepository) UpdateQuantity(ctx context.Context, id string, delta float64) error {
	if m.UpdateQuantityFunc != nil {
		return m.UpdateQuantityFunc(ctx, id, delta)
	}
	return nil
}

func (m *MockInventoryRepository) GetExpiringSoon(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
	if m.GetExpiringSoonFunc != nil {
		return m.GetExpiringSoonFunc(ctx, days)
	}
	return []*domain.InventoryItem{}, nil
}

func (m *MockInventoryRepository) GetLowStock(ctx context.Context) ([]*domain.InventoryItem, error) {
	if m.GetLowStockFunc != nil {
		return m.GetLowStockFunc(ctx)
	}
	return []*domain.InventoryItem{}, nil
}
