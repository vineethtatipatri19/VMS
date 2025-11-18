package service

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestCrateService_IssueCrates(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockCrateRepo := &mocks.MockCrateRepository{
			CreateFunc: func(ctx context.Context, crate *domain.Crate) error {
				return nil
			},
		}
		mockCustomerRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:     "c1",
					Name:   "Test Customer",
					Status: "active",
				}, nil
			},
		}

		service := NewCrateService(mockCrateRepo, mockCustomerRepo)
		crate := &domain.Crate{
			CustomerID:      "c1",
			TransactionType: "out",
			Quantity:        10,
		}

		err := service.IssueCrates(ctx, crate)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("inactive customer", func(t *testing.T) {
		mockCrateRepo := &mocks.MockCrateRepository{}
		mockCustomerRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:     "c1",
					Status: "inactive",
				}, nil
			},
		}

		service := NewCrateService(mockCrateRepo, mockCustomerRepo)
		crate := &domain.Crate{
			CustomerID:      "c1",
			TransactionType: "out",
			Quantity:        10,
		}

		err := service.IssueCrates(ctx, crate)
		if err == nil {
			t.Error("Expected error for inactive customer")
		}
	})
}

func TestCrateService_ReturnCrates(t *testing.T) {
	ctx := context.Background()

	mockCrateRepo := &mocks.MockCrateRepository{
		GetBalanceFunc: func(ctx context.Context, customerID string) (int, error) {
			return 10, nil
		},
		CreateFunc: func(ctx context.Context, crate *domain.Crate) error {
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

	service := NewCrateService(mockCrateRepo, mockCustomerRepo)
	crate := &domain.Crate{
		CustomerID:      "c1",
		TransactionType: "in",
		Quantity:        5,
	}

	err := service.ReturnCrates(ctx, crate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCrateService_GetCrateBalance(t *testing.T) {
	ctx := context.Background()

	mockCrateRepo := &mocks.MockCrateRepository{
		GetBalanceFunc: func(ctx context.Context, customerID string) (int, error) {
			return 15, nil
		},
	}

	service := NewCrateService(mockCrateRepo, nil)
	balance, err := service.GetCrateBalance(ctx, "c1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if balance != 15 {
		t.Errorf("Expected balance 15, got %d", balance)
	}
}

func TestCrateService_GetCrateHistory(t *testing.T) {
	ctx := context.Background()

	crates := []*domain.Crate{
		{ID: "cr1", CustomerID: "c1", Quantity: 10},
		{ID: "cr2", CustomerID: "c1", Quantity: 5},
	}
	mockCrateRepo := &mocks.MockCrateRepository{
		ListByCustomerFunc: func(ctx context.Context, customerID string) ([]*domain.Crate, error) {
			return crates, nil
		},
	}

	service := NewCrateService(mockCrateRepo, nil)
	result, err := service.GetCrateHistory(ctx, "c1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 crates, got %d", len(result))
	}
}

func TestCrateService_DeleteCrate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := &mocks.MockCrateRepository{
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				return nil
			},
		}

		service := NewCrateService(mockRepo, nil)
		req := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteCrate(ctx, "cr1", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
