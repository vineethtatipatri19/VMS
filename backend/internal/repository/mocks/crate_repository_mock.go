package mocks

import (
	"context"

	"github.com/example/pgvms/internal/domain"
)

// MockCrateRepository is a mock implementation of repository.CrateRepository
type MockCrateRepository struct {
	CreateFunc         func(ctx context.Context, crate *domain.CrateEntry) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.CrateEntry, error)
	ListFunc           func(ctx context.Context) ([]*domain.CrateEntry, error)
	ListByCustomerFunc func(ctx context.Context, customerID string) ([]*domain.CrateEntry, error)
	UpdateFunc         func(ctx context.Context, crate *domain.CrateEntry) error
	DeleteFunc         func(ctx context.Context, id string, req *domain.DeleteRequest) error
	GetBalanceFunc     func(ctx context.Context, customerID string) (int, error)
}

func (m *MockCrateRepository) Create(ctx context.Context, crate *domain.CrateEntry) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, crate)
	}
	return nil
}

func (m *MockCrateRepository) GetByID(ctx context.Context, id string) (*domain.CrateEntry, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockCrateRepository) List(ctx context.Context) ([]*domain.CrateEntry, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []*domain.CrateEntry{}, nil
}

func (m *MockCrateRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	if m.ListByCustomerFunc != nil {
		return m.ListByCustomerFunc(ctx, customerID)
	}
	return []*domain.CrateEntry{}, nil
}

func (m *MockCrateRepository) Update(ctx context.Context, crate *domain.CrateEntry) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, crate)
	}
	return nil
}

func (m *MockCrateRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}

func (m *MockCrateRepository) GetBalance(ctx context.Context, customerID string) (int, error) {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc(ctx, customerID)
	}
	return 0, nil
}
