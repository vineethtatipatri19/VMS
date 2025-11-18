package mocks

import (
	"context"

	"github.com/example/pgvms/internal/domain"
)

// MockPaymentScheduleRepository is a mock implementation of repository.PaymentScheduleRepository
type MockPaymentScheduleRepository struct {
	CreateFunc         func(ctx context.Context, schedule *domain.PaymentSchedule) error
	GetByIDFunc        func(ctx context.Context, id string) (*domain.PaymentSchedule, error)
	ListByCustomerFunc func(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error)
	UpdateFunc         func(ctx context.Context, schedule *domain.PaymentSchedule) error
	DeleteFunc         func(ctx context.Context, id string) error
	UpdateStatusFunc   func(ctx context.Context, id string, status string) error
}

func (m *MockPaymentScheduleRepository) Create(ctx context.Context, schedule *domain.PaymentSchedule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, schedule)
	}
	return nil
}

func (m *MockPaymentScheduleRepository) GetByID(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, domain.ErrNotFound
}

func (m *MockPaymentScheduleRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	if m.ListByCustomerFunc != nil {
		return m.ListByCustomerFunc(ctx, customerID)
	}
	return []*domain.PaymentSchedule{}, nil
}

func (m *MockPaymentScheduleRepository) Update(ctx context.Context, schedule *domain.PaymentSchedule) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, schedule)
	}
	return nil
}

func (m *MockPaymentScheduleRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockPaymentScheduleRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}
