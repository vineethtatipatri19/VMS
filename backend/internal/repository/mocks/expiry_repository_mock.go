package mocks

import (
	"context"

	"github.com/example/pgvms/internal/domain"
)

// MockExpiryAlertRepository is a mock implementation of repository.ExpiryAlertRepository
type MockExpiryAlertRepository struct {
	CreateFunc      func(ctx context.Context, alert *domain.ExpiryAlert) error
	GetByIDFunc     func(ctx context.Context, id string) (*domain.ExpiryAlert, error)
	ListFunc        func(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error)
	UpdateFunc      func(ctx context.Context, alert *domain.ExpiryAlert) error
	DeleteFunc      func(ctx context.Context, id string, req *domain.DeleteRequest) error
	AcknowledgeFunc func(ctx context.Context, id string, userID string) error
}

func (m *MockExpiryAlertRepository) Create(ctx context.Context, alert *domain.ExpiryAlert) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, alert)
	}
	return nil
}

func (m *MockExpiryAlertRepository) GetByID(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockExpiryAlertRepository) List(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, acknowledged)
	}
	return []*domain.ExpiryAlert{}, nil
}

func (m *MockExpiryAlertRepository) Update(ctx context.Context, alert *domain.ExpiryAlert) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, alert)
	}
	return nil
}

func (m *MockExpiryAlertRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, req)
	}
	return nil
}

func (m *MockExpiryAlertRepository) Acknowledge(ctx context.Context, id string, userID string) error {
	if m.AcknowledgeFunc != nil {
		return m.AcknowledgeFunc(ctx, id, userID)
	}
	return nil
}
