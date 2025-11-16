package service

import (
"context"
"fmt"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/repository"
)

type InventoryService struct {
inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) *InventoryService {
return &InventoryService{
inventoryRepo: inventoryRepo,
}
}

func (s *InventoryService) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
if err := item.Validate(); err != nil {
return err
}

now := time.Now()
item.CreatedAt = now
item.UpdatedAt = now
item.Status = "active"

if err := s.inventoryRepo.Create(ctx, item); err != nil {
return fmt.Errorf("failed to create inventory item: %w", err)
}

return nil
}

func (s *InventoryService) GetItem(ctx context.Context, id string) (*domain.InventoryItem, error) {
return s.inventoryRepo.GetByID(ctx, id)
}

func (s *InventoryService) ListItems(ctx context.Context, status, sortBy string) ([]*domain.InventoryItem, error) {
items, err := s.inventoryRepo.List(ctx, status, sortBy)
if err != nil {
return nil, fmt.Errorf("failed to list inventory items: %w", err)
}
return items, nil
}

func (s *InventoryService) UpdateItem(ctx context.Context, item *domain.InventoryItem) error {
if err := item.Validate(); err != nil {
return err
}

existing, err := s.inventoryRepo.GetByID(ctx, item.ID)
if err != nil {
return err
}

item.UpdatedAt = time.Now()
item.Quantity = existing.Quantity
if err := s.inventoryRepo.Update(ctx, item); err != nil {
return fmt.Errorf("failed to update inventory item: %w", err)
}

return nil
}

func (s *InventoryService) DeleteItem(ctx context.Context, id string, req *domain.DeleteRequest) error {
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

item, err := s.inventoryRepo.GetByID(ctx, id)
if err != nil {
return err
}

if item.Quantity > 0 {
return &domain.BusinessError{
Code:    "STOCK_EXISTS",
Message: fmt.Sprintf("cannot delete item with stock: %d units remaining", item.Quantity),
}
}

if err := s.inventoryRepo.Delete(ctx, id, req); err != nil {
return fmt.Errorf("failed to delete inventory item: %w", err)
}

return nil
}

func (s *InventoryService) DeductStock(ctx context.Context, itemID string, quantity int) error {
item, err := s.inventoryRepo.GetByID(ctx, itemID)
if err != nil {
return err
}

if item.Status != "active" {
return &domain.BusinessError{
Code:    "INACTIVE_ITEM",
Message: "inventory item is not active",
}
}

if item.Quantity < float64(quantity) {
return &domain.BusinessError{
Code:    "INSUFFICIENT_STOCK",
Message: fmt.Sprintf("insufficient stock: requested %d, available %d", quantity, item.Quantity),
}
}

if item.IsExpired() {
return &domain.BusinessError{
Code:    "ITEM_EXPIRED",
Message: fmt.Sprintf("item expired on %s", item.ExpiryDate),
}
}

delta := float64(-quantity)
	if err := s.inventoryRepo.UpdateQuantity(ctx, itemID, delta); err != nil {
return fmt.Errorf("failed to deduct stock: %w", err)
}

return nil
}

func (s *InventoryService) AddStock(ctx context.Context, itemID string, quantity int) error {	delta := float64(quantity)
	if err := s.inventoryRepo.UpdateQuantity(ctx, itemID, delta); err != nil {
return fmt.Errorf("failed to add stock: %w", err)
}

return nil
}

func (s *InventoryService) CheckStock(ctx context.Context, itemID string, requiredQty int) error {
item, err := s.inventoryRepo.GetByID(ctx, itemID)
if err != nil {
return err
}

if item.Status != "active" {
return &domain.BusinessError{
Code:    "INACTIVE_ITEM",
Message: "inventory item is not active",
}
}

if item.Quantity < float64(requiredQty) {
return &domain.BusinessError{
Code:    "INSUFFICIENT_STOCK",
Message: fmt.Sprintf("insufficient stock: requested %d, available %d", requiredQty, item.Quantity),
}
}

if item.IsExpired() {
return &domain.BusinessError{
Code:    "ITEM_EXPIRED",
Message: "item has expired",
}
}

return nil
}

func (s *InventoryService) GetExpiringItems(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
items, err := s.inventoryRepo.GetExpiringSoon(ctx, days)
if err != nil {
return nil, fmt.Errorf("failed to get expiring items: %w", err)
}
return items, nil
}

func (s *InventoryService) GetLowStockItems(ctx context.Context) ([]*domain.InventoryItem, error) {
items, err := s.inventoryRepo.GetLowStock(ctx)
if err != nil {
return nil, fmt.Errorf("failed to get low stock items: %w", err)
}
return items, nil
}

func (s *InventoryService) SearchItems(ctx context.Context, query string) ([]*domain.InventoryItem, error) {
	// Note: Search would need custom implementation or use status="" sortBy=""
	return s.ListItems(ctx, "", "")
}
