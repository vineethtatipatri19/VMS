package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type WastageService struct {
	wastageRepo   repository.WastageRepository
	inventoryRepo repository.InventoryRepository
}

func NewWastageService(
	wastageRepo repository.WastageRepository,
	inventoryRepo repository.InventoryRepository,
) *WastageService {
	return &WastageService{
		wastageRepo:   wastageRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *WastageService) RecordWastage(ctx context.Context, wastage *domain.WastageLog) error {
	invItem, err := s.inventoryRepo.GetByID(ctx, wastage.InventoryID)
	if err != nil {
		return fmt.Errorf("failed to get inventory item: %w", err)
	}

	if invItem.Quantity < wastage.Quantity {
		return &domain.BusinessError{
			Code:    "INSUFFICIENT_STOCK",
			Message: fmt.Sprintf("insufficient stock: requested %.2f, available %.2f", wastage.Quantity, invItem.Quantity),
		}
	}

	// Generate UUID for new wastage record
	wastage.ID = uuid.New().String()

	now := time.Now()
	wastage.RecordedAt = now
	wastage.CreatedAt = now
	wastage.UpdatedAt = now

	if err := s.wastageRepo.Create(ctx, wastage); err != nil {
		return fmt.Errorf("failed to record wastage: %w", err)
	}

	delta := -wastage.Quantity
	if err := s.inventoryRepo.UpdateQuantity(ctx, wastage.InventoryID, delta); err != nil {
		return fmt.Errorf("failed to update inventory quantity: %w", err)
	}

	return nil
}

func (s *WastageService) GetWastage(ctx context.Context, id string) (*domain.WastageLog, error) {
	return s.wastageRepo.GetByID(ctx, id)
}

func (s *WastageService) ListWastage(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error) {
	wastages, err := s.wastageRepo.List(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list wastages: %w", err)
	}
	return wastages, nil
}

func (s *WastageService) UpdateWastage(ctx context.Context, wastage *domain.WastageLog) error {
	wastage.UpdatedAt = time.Now()

	if err := s.wastageRepo.Update(ctx, wastage); err != nil {
		return fmt.Errorf("failed to update wastage: %w", err)
	}

	return nil
}

func (s *WastageService) DeleteWastage(ctx context.Context, id string, req *domain.DeleteRequest) error {
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

	if err := s.wastageRepo.Delete(ctx, id, req); err != nil {
		return fmt.Errorf("failed to delete wastage: %w", err)
	}

	return nil
}

func (s *WastageService) GetWastageReport(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error) {
	return s.ListWastage(ctx, startDate, endDate)
}

func (s *WastageService) CalculateTotalWastageCost(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	wastages, err := s.GetWastageReport(ctx, startDate, endDate)
	if err != nil {
		return 0, err
	}

	var totalCost float64
	for _, w := range wastages {
		totalCost += w.CostValue
	}

	return totalCost, nil
}
