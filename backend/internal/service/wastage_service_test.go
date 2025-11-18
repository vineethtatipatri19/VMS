package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestWastageService_RecordWastage(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockWastageRepo := &mocks.MockWastageRepository{
			CreateFunc: func(ctx context.Context, wastage *domain.WastageLog) error {
				return nil
			},
		}
		mockInventoryRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{
					ID:       "i1",
					Quantity: 100.0,
				}, nil
			},
			UpdateQuantityFunc: func(ctx context.Context, id string, delta float64) error {
				return nil
			},
		}

		service := NewWastageService(mockWastageRepo, mockInventoryRepo)
		wastage := &domain.WastageLog{
			InventoryID: "i1",
			Quantity:    10.0,
			Reason:      "Spoilage",
			CostValue:   50.0,
		}

		err := service.RecordWastage(ctx, wastage)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("insufficient stock", func(t *testing.T) {
		mockWastageRepo := &mocks.MockWastageRepository{}
		mockInventoryRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{
					ID:       "i1",
					Quantity: 5.0,
				}, nil
			},
		}

		service := NewWastageService(mockWastageRepo, mockInventoryRepo)
		wastage := &domain.WastageLog{
			InventoryID: "i1",
			Quantity:    10.0,
			Reason:      "Spoilage",
		}

		err := service.RecordWastage(ctx, wastage)
		if err == nil {
			t.Error("Expected error for insufficient stock, got nil")
		}
	})
}

func TestWastageService_GetWastage(t *testing.T) {
	ctx := context.Background()

	mockWastageRepo := &mocks.MockWastageRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.WastageLog, error) {
			return &domain.WastageLog{
				ID:          id,
				InventoryID: "i1",
				Quantity:    10.0,
				Reason:      "Spoilage",
			}, nil
		},
	}

	service := NewWastageService(mockWastageRepo, nil)
	wastage, err := service.GetWastage(ctx, "w1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if wastage.ID != "w1" {
		t.Errorf("Expected ID w1, got %s", wastage.ID)
	}
}

func TestWastageService_ListWastage(t *testing.T) {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	mockWastageRepo := &mocks.MockWastageRepository{
		ListFunc: func(ctx context.Context, sd, ed time.Time) ([]*domain.WastageLog, error) {
			return []*domain.WastageLog{
				{ID: "w1", Quantity: 10.0},
				{ID: "w2", Quantity: 5.0},
			}, nil
		},
	}

	service := NewWastageService(mockWastageRepo, nil)
	wastages, err := service.ListWastage(ctx, startDate, endDate)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(wastages) != 2 {
		t.Errorf("Expected 2 wastages, got %d", len(wastages))
	}
}

func TestWastageService_UpdateWastage(t *testing.T) {
	ctx := context.Background()

	mockWastageRepo := &mocks.MockWastageRepository{
		UpdateFunc: func(ctx context.Context, wastage *domain.WastageLog) error {
			return nil
		},
	}

	service := NewWastageService(mockWastageRepo, nil)
	wastage := &domain.WastageLog{
		ID:       "w1",
		Quantity: 15.0,
		Reason:   "Updated reason",
	}

	err := service.UpdateWastage(ctx, wastage)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWastageService_DeleteWastage(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockWastageRepo := &mocks.MockWastageRepository{
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				return nil
			},
		}

		service := NewWastageService(mockWastageRepo, nil)
		req := &domain.DeleteRequest{
			Reason:      "Duplicate entry",
			Attestation: "Confirmed by manager",
		}

		err := service.DeleteWastage(ctx, "w1", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("missing reason", func(t *testing.T) {
		mockWastageRepo := &mocks.MockWastageRepository{}
		service := NewWastageService(mockWastageRepo, nil)
		req := &domain.DeleteRequest{
			Attestation: "Confirmed",
		}

		err := service.DeleteWastage(ctx, "w1", req)
		if err == nil {
			t.Error("Expected validation error for missing reason")
		}
	})
}

func TestWastageService_CalculateTotalWastageCost(t *testing.T) {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	mockWastageRepo := &mocks.MockWastageRepository{
		ListFunc: func(ctx context.Context, sd, ed time.Time) ([]*domain.WastageLog, error) {
			return []*domain.WastageLog{
				{ID: "w1", CostValue: 100.0},
				{ID: "w2", CostValue: 50.0},
				{ID: "w3", CostValue: 25.0},
			}, nil
		},
	}

	service := NewWastageService(mockWastageRepo, nil)
	totalCost, err := service.CalculateTotalWastageCost(ctx, startDate, endDate)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if totalCost != 175.0 {
		t.Errorf("Expected total cost 175.0, got %.2f", totalCost)
	}
}
