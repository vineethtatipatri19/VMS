package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestExpiryService_GenerateAlerts(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tomorrow := time.Now().AddDate(0, 0, 1)
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{
			CreateFunc: func(ctx context.Context, alert *domain.ExpiryAlert) error {
				return nil
			},
		}
		mockInventoryRepo := &mocks.MockInventoryRepository{
			GetExpiringSoonFunc: func(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
				return []*domain.InventoryItem{
					{
						ID:         "i1",
						Name:       "Test Item",
						ExpiryDate: tomorrow.Format("2006-01-02"),
					},
				}, nil
			},
		}

		service := NewExpiryService(mockExpiryRepo, mockInventoryRepo)
		err := service.GenerateAlerts(ctx, 7)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("no expiring items", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{}
		mockInventoryRepo := &mocks.MockInventoryRepository{
			GetExpiringSoonFunc: func(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
				return []*domain.InventoryItem{}, nil
			},
		}

		service := NewExpiryService(mockExpiryRepo, mockInventoryRepo)
		err := service.GenerateAlerts(ctx, 7)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestExpiryService_GetAlert(t *testing.T) {
	ctx := context.Background()

	mockExpiryRepo := &mocks.MockExpiryAlertRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
			return &domain.ExpiryAlert{
				ID:              id,
				InventoryItemID: "i1",
				DaysUntilExpiry: 3,
				Acknowledged:    false,
			}, nil
		},
	}

	service := NewExpiryService(mockExpiryRepo, nil)
	alert, err := service.GetAlert(ctx, "a1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if alert.ID != "a1" {
		t.Errorf("Expected ID a1, got %s", alert.ID)
	}
}

func TestExpiryService_ListAlerts(t *testing.T) {
	ctx := context.Background()

	mockExpiryRepo := &mocks.MockExpiryAlertRepository{
		ListFunc: func(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
			return []*domain.ExpiryAlert{
				{ID: "a1", Acknowledged: acknowledged},
				{ID: "a2", Acknowledged: acknowledged},
			}, nil
		},
	}

	service := NewExpiryService(mockExpiryRepo, nil)
	alerts, err := service.ListAlerts(ctx, false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(alerts) != 2 {
		t.Errorf("Expected 2 alerts, got %d", len(alerts))
	}
}

func TestExpiryService_GetPendingAlerts(t *testing.T) {
	ctx := context.Background()

	mockExpiryRepo := &mocks.MockExpiryAlertRepository{
		ListFunc: func(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
			if !acknowledged {
				return []*domain.ExpiryAlert{
					{ID: "a1", Acknowledged: false},
					{ID: "a2", Acknowledged: false},
				}, nil
			}
			return []*domain.ExpiryAlert{}, nil
		},
	}

	service := NewExpiryService(mockExpiryRepo, nil)
	alerts, err := service.GetPendingAlerts(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(alerts) != 2 {
		t.Errorf("Expected 2 pending alerts, got %d", len(alerts))
	}
}

func TestExpiryService_AcknowledgeAlert(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{
			AcknowledgeFunc: func(ctx context.Context, id string, userID string) error {
				return nil
			},
		}

		service := NewExpiryService(mockExpiryRepo, nil)
		err := service.AcknowledgeAlert(ctx, "a1", "user123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("missing acknowledged_by", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{}
		service := NewExpiryService(mockExpiryRepo, nil)

		err := service.AcknowledgeAlert(ctx, "a1", "")
		if err == nil {
			t.Error("Expected validation error for missing acknowledged_by")
		}
	})
}

func TestExpiryService_UpdateAlert(t *testing.T) {
	ctx := context.Background()

	mockExpiryRepo := &mocks.MockExpiryAlertRepository{
		UpdateFunc: func(ctx context.Context, alert *domain.ExpiryAlert) error {
			return nil
		},
	}

	service := NewExpiryService(mockExpiryRepo, nil)
	alert := &domain.ExpiryAlert{
		ID:              "a1",
		DaysUntilExpiry: 5,
		Acknowledged:    true,
	}

	err := service.UpdateAlert(ctx, alert)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExpiryService_DeleteAlert(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				return nil
			},
		}

		service := NewExpiryService(mockExpiryRepo, nil)
		req := &domain.DeleteRequest{
			Reason:      "Duplicate alert",
			Attestation: "Confirmed by manager",
		}

		err := service.DeleteAlert(ctx, "a1", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("missing reason", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{}
		service := NewExpiryService(mockExpiryRepo, nil)
		req := &domain.DeleteRequest{
			Attestation: "Confirmed",
		}

		err := service.DeleteAlert(ctx, "a1", req)
		if err == nil {
			t.Error("Expected validation error for missing reason")
		}
	})

	t.Run("missing attestation", func(t *testing.T) {
		mockExpiryRepo := &mocks.MockExpiryAlertRepository{}
		service := NewExpiryService(mockExpiryRepo, nil)
		req := &domain.DeleteRequest{
			Reason: "Duplicate",
		}

		err := service.DeleteAlert(ctx, "a1", req)
		if err == nil {
			t.Error("Expected validation error for missing attestation")
		}
	})
}
