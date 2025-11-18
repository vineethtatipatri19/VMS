package service

import (
	"context"
	"fmt"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type SaleItemService struct {
	saleItemRepo  repository.SaleItemRepository
	inventoryRepo repository.InventoryRepository
}

func NewSaleItemService(
	saleItemRepo repository.SaleItemRepository,
	inventoryRepo repository.InventoryRepository,
) *SaleItemService {
	return &SaleItemService{
		saleItemRepo:  saleItemRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *SaleItemService) CreateSaleItem(ctx context.Context, item *domain.SaleItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	// Generate UUID for new sale item
	item.ID = uuid.New().String()

	invItem, err := s.inventoryRepo.GetByID(ctx, item.InventoryLotID)
	if err != nil {
		return fmt.Errorf("failed to get inventory item: %w", err)
	}

	if invItem.Quantity < float64(item.Quantity) {
		return &domain.BusinessError{
			Code:    "INSUFFICIENT_STOCK",
			Message: fmt.Sprintf("insufficient stock: requested %.0f, available %.0f", item.Quantity, invItem.Quantity),
		}
	}

	if err := s.saleItemRepo.Create(ctx, item); err != nil {
		return fmt.Errorf("failed to create sale item: %w", err)
	}

	return nil
}

func (s *SaleItemService) GetSaleItem(ctx context.Context, id string) (*domain.SaleItem, error) {
	return s.saleItemRepo.GetByID(ctx, id)
}

func (s *SaleItemService) ListItemsForTransaction(ctx context.Context, transactionID string) ([]*domain.SaleItem, error) {
	items, err := s.saleItemRepo.ListByTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sale items: %w", err)
	}
	return items, nil
}

func (s *SaleItemService) UpdateSaleItem(ctx context.Context, item *domain.SaleItem) error {
	if err := item.Validate(); err != nil {
		return err
	}
	if err := s.saleItemRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to update sale item: %w", err)
	}

	return nil
}

func (s *SaleItemService) DeleteSaleItem(ctx context.Context, id string, req *domain.DeleteRequest) error {
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

	if err := s.saleItemRepo.Delete(ctx, id, req); err != nil {
		return fmt.Errorf("failed to delete sale item: %w", err)
	}

	return nil
}

func (s *SaleItemService) CalculateTotalProfit(ctx context.Context, transactionID string) (float64, error) {
	items, err := s.saleItemRepo.ListByTransaction(ctx, transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to list sale items: %w", err)
	}

	var totalProfit float64
	for _, item := range items {
		totalProfit += item.Profit
	}

	return totalProfit, nil
}
