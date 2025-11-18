package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type ExpiryService struct {
	expiryRepo    repository.ExpiryAlertRepository
	inventoryRepo repository.InventoryRepository
}

func NewExpiryService(
	expiryRepo repository.ExpiryAlertRepository,
	inventoryRepo repository.InventoryRepository,
) *ExpiryService {
	return &ExpiryService{
		expiryRepo:    expiryRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *ExpiryService) GenerateAlerts(ctx context.Context, daysThreshold int) error {
	items, err := s.inventoryRepo.GetExpiringSoon(ctx, daysThreshold)
	if err != nil {
		return fmt.Errorf("failed to get expiring items: %w", err)
	}

	now := time.Now()
	for _, item := range items {
		// Parse expiry date string
		expiryDate, err := time.Parse("2006-01-02", item.ExpiryDate)
		if err != nil {
			continue
		}

		daysUntil := int(expiryDate.Sub(now).Hours() / 24)

		alert := &domain.ExpiryAlert{
			ID:              uuid.New().String(),
			InventoryItemID: item.ID,
			AlertDate:       now,
			ExpiryDate:      expiryDate,
			DaysUntilExpiry: daysUntil,
			Acknowledged:    false,
			CreatedAt:       now,
		}

		if err := s.expiryRepo.Create(ctx, alert); err != nil {
			return fmt.Errorf("failed to create alert for item %s: %w", item.Name, err)
		}
	}

	return nil
}

func (s *ExpiryService) GetAlert(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
	return s.expiryRepo.GetByID(ctx, id)
}

func (s *ExpiryService) ListAlerts(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
	alerts, err := s.expiryRepo.List(ctx, acknowledged)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}
	return alerts, nil
}

func (s *ExpiryService) GetPendingAlerts(ctx context.Context) ([]*domain.ExpiryAlert, error) {
	return s.ListAlerts(ctx, false)
}

func (s *ExpiryService) GetAcknowledgedAlerts(ctx context.Context) ([]*domain.ExpiryAlert, error) {
	return s.ListAlerts(ctx, true)
}

func (s *ExpiryService) AcknowledgeAlert(ctx context.Context, id string, acknowledgedBy string) error {
	if acknowledgedBy == "" {
		return &domain.ValidationError{
			Field:   "acknowledged_by",
			Message: "acknowledged_by is required",
		}
	}

	if err := s.expiryRepo.Acknowledge(ctx, id, acknowledgedBy); err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

func (s *ExpiryService) UpdateAlert(ctx context.Context, alert *domain.ExpiryAlert) error {
	if err := s.expiryRepo.Update(ctx, alert); err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	return nil
}

func (s *ExpiryService) DeleteAlert(ctx context.Context, id string, req *domain.DeleteRequest) error {
	if req.Reason == "" {
		return &domain.ValidationError{
			Field:   "reason",
			Message: "reason is required",
		}
	}
	if req.Attestation == "" {
		return &domain.ValidationError{
			Field:   "attestation",
			Message: "attestation is required",
		}
	}

	if err := s.expiryRepo.Delete(ctx, id, req); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}
