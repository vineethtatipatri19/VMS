package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestPaymentService_CreateSchedule(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			CreateFunc: func(ctx context.Context, schedule *domain.PaymentSchedule) error {
				return nil
			},
		}
		mockCustomerRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:     "c1",
					Status: "active",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, mockCustomerRepo)
		schedule := &domain.PaymentSchedule{
			CustomerID: "c1",
			AmountDue:  100.0,
			DueDate:    time.Now().AddDate(0, 0, 7),
		}

		err := service.CreateSchedule(ctx, schedule)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("inactive customer", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{}
		mockCustomerRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:     "c1",
					Status: "inactive",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, mockCustomerRepo)
		schedule := &domain.PaymentSchedule{
			CustomerID: "c1",
			AmountDue:  100.0,
			DueDate:    time.Now().AddDate(0, 0, 7),
		}

		err := service.CreateSchedule(ctx, schedule)
		if err == nil {
			t.Error("Expected error for inactive customer")
		}
	})
}

func TestPaymentService_GetSchedule(t *testing.T) {
	ctx := context.Background()

	mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
			return &domain.PaymentSchedule{
				ID:         id,
				CustomerID: "c1",
				AmountDue:  100.0,
				Status:     "pending",
			}, nil
		},
	}

	service := NewPaymentService(mockPaymentRepo, nil, nil)
	schedule, err := service.GetSchedule(ctx, "ps1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if schedule.ID != "ps1" {
		t.Errorf("Expected ID ps1, got %s", schedule.ID)
	}
}

func TestPaymentService_ListSchedules(t *testing.T) {
	ctx := context.Background()

	mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
		ListByCustomerFunc: func(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
			return []*domain.PaymentSchedule{
				{ID: "ps1", CustomerID: customerID},
				{ID: "ps2", CustomerID: customerID},
			}, nil
		},
	}

	service := NewPaymentService(mockPaymentRepo, nil, nil)
	schedules, err := service.ListSchedules(ctx, "c1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(schedules) != 2 {
		t.Errorf("Expected 2 schedules, got %d", len(schedules))
	}
}

func TestPaymentService_RecordPayment(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:        id,
					AmountDue: 100.0,
					Status:    "pending",
				}, nil
			},
			UpdateStatusFunc: func(ctx context.Context, id string, status string) error {
				return nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.RecordPayment(ctx, "ps1", 100.0, time.Now())

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("already paid", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:        id,
					AmountDue: 100.0,
					Status:    "paid",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.RecordPayment(ctx, "ps1", 100.0, time.Now())

		if err == nil {
			t.Error("Expected error for already paid schedule")
		}
	})

	t.Run("amount mismatch", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:        id,
					AmountDue: 100.0,
					Status:    "pending",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.RecordPayment(ctx, "ps1", 50.0, time.Now())

		if err == nil {
			t.Error("Expected error for amount mismatch")
		}
	})
}

func TestPaymentService_UpdateSchedule(t *testing.T) {
	ctx := context.Background()

	mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
		UpdateFunc: func(ctx context.Context, schedule *domain.PaymentSchedule) error {
			return nil
		},
	}

	service := NewPaymentService(mockPaymentRepo, nil, nil)
	schedule := &domain.PaymentSchedule{
		ID:         "ps1",
		CustomerID: "c1",
		AmountDue:  150.0,
		DueDate:    time.Now().AddDate(0, 0, 10),
	}

	err := service.UpdateSchedule(ctx, schedule)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPaymentService_DeleteSchedule(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:     id,
					Status: "pending",
				}, nil
			},
			DeleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.DeleteSchedule(ctx, "ps1")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("cannot delete paid", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:     id,
					Status: "paid",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.DeleteSchedule(ctx, "ps1")

		if err == nil {
			t.Error("Expected error when deleting paid schedule")
		}
	})
}

func TestPaymentService_MarkOverdue(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:     id,
					Status: "pending",
				}, nil
			},
			UpdateStatusFunc: func(ctx context.Context, id string, status string) error {
				return nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.MarkOverdue(ctx, "ps1")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("cannot mark paid as overdue", func(t *testing.T) {
		mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
				return &domain.PaymentSchedule{
					ID:     id,
					Status: "paid",
				}, nil
			},
		}

		service := NewPaymentService(mockPaymentRepo, nil, nil)
		err := service.MarkOverdue(ctx, "ps1")

		if err == nil {
			t.Error("Expected error when marking paid schedule as overdue")
		}
	})
}

func TestPaymentService_GetOverdueSchedules(t *testing.T) {
	ctx := context.Background()

	pastDate := time.Now().AddDate(0, 0, -5)
	futureDate := time.Now().AddDate(0, 0, 5)

	mockPaymentRepo := &mocks.MockPaymentScheduleRepository{
		ListByCustomerFunc: func(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
			return []*domain.PaymentSchedule{
				{ID: "ps1", Status: "pending", DueDate: pastDate},
				{ID: "ps2", Status: "pending", DueDate: futureDate},
				{ID: "ps3", Status: "paid", DueDate: pastDate},
			}, nil
		},
	}

	service := NewPaymentService(mockPaymentRepo, nil, nil)
	overdueSchedules, err := service.GetOverdueSchedules(ctx, "c1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(overdueSchedules) != 1 {
		t.Errorf("Expected 1 overdue schedule, got %d", len(overdueSchedules))
	}
	if overdueSchedules[0].ID != "ps1" {
		t.Errorf("Expected overdue schedule ID ps1, got %s", overdueSchedules[0].ID)
	}
}
