package service

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestInventoryService_CreateItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success - valid item", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{
			CreateFunc: func(ctx context.Context, item *domain.InventoryItem) error {
				if item.Name == "" {
					t.Error("Expected name to be set")
				}
				return nil
			},
		}

		service := NewInventoryService(mockRepo)
		item := &domain.InventoryItem{
			Name:         "Test Product",
			LotNumber:    "LOT001",
			Quantity:     100,
			Unit:         "kg",
			CostPrice:    50,
			SellingPrice: 75,
		}

		err := service.CreateItem(ctx, item)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{}
		service := NewInventoryService(mockRepo)

		item := &domain.InventoryItem{
			LotNumber: "LOT001",
			Quantity:  100,
			Unit:      "kg",
		}

		err := service.CreateItem(ctx, item)
		if err == nil {
			t.Error("Expected validation error for missing name")
		}
	})
}

func TestInventoryService_GetItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success - item found", func(t *testing.T) {
		expected := &domain.InventoryItem{
			ID:   "inv123",
			Name: "Test Product",
		}

		mockRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				if id == "inv123" {
					return expected, nil
				}
				return nil, domain.ErrNotFound
			},
		}

		service := NewInventoryService(mockRepo)
		result, err := service.GetItem(ctx, "inv123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
		}
	})

	t.Run("error - item not found", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return nil, domain.ErrNotFound
			},
		}

		service := NewInventoryService(mockRepo)
		_, err := service.GetItem(ctx, "nonexistent")

		if err != domain.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}

func TestInventoryService_UpdateItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success - valid update", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{
					ID:       "inv123",
					Name:     "Old Name",
					Quantity: 100,
				}, nil
			},
			UpdateFunc: func(ctx context.Context, item *domain.InventoryItem) error {
				return nil
			},
		}

		service := NewInventoryService(mockRepo)
		item := &domain.InventoryItem{
			ID:        "inv123",
			Name:      "New Name",
			LotNumber: "LOT001",
			Quantity:  100,
			Unit:      "kg",
		}

		err := service.UpdateItem(ctx, item)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestInventoryService_DeleteItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success - item with zero stock", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{
					ID:       "inv123",
					Quantity: 0,
				}, nil
			},
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				return nil
			},
		}

		service := NewInventoryService(mockRepo)
		req := &domain.DeleteRequest{
			Reason:      "No longer needed",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteItem(ctx, "inv123", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error - item has stock", func(t *testing.T) {
		mockRepo := &mocks.MockInventoryRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{
					ID:       "inv123",
					Quantity: 50,
				}, nil
			},
		}

		service := NewInventoryService(mockRepo)
		req := &domain.DeleteRequest{
			Reason:      "Test",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteItem(ctx, "inv123", req)
		if err == nil {
			t.Error("Expected error for item with stock")
		}
	})
}


