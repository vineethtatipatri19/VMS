package mocks

import (
	"context"
	"time"

	"github.com/example/pgvms/internal/domain"
)

// MockWastageRepository is a mock implementation of repository.WastageRepository
type MockWastageRepository struct {
	CreateFunc  func(ctx context.Context, wastage *domain.WastageLog) error
	GetByIDFunc func(ctx context.Context, id string) (*domain.WastageLog, error)
	ListFunc    func(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error)
	UpdateFunc  func(ctx context.Context, wastage *domain.WastageLog) error
	DeleteFunc  func(ctx context.Context, id string, req *domain.DeleteRequest) error
}

func (m *MockWastageRepository) Create(ctx context.Context, wastage *domain.WastageLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, wastage)
	}
	return nil
}

func (m *MockWastageRepository) GetByID(ctx context.Context, id string) (*domain.WastageLog, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockWastageRepository) List(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, startDate, endDate)
	}
	return []*domain.WastageLog{}, nil
}

func (m *MockWastageRepository) Update(ctx context.Context, wastage *domain.WastageLog) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, wastage)
	}
	return nil
}

func (m *MockWastageRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}
