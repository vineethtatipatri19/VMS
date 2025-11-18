package service

import (
"context"
"testing"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/repository/mocks"
)

func TestSaleItemService_CreateSaleItem(t *testing.T) {
ctx := context.Background()

t.Run("success", func(t *testing.T) {
mockSaleItemRepo := &mocks.MockSaleItemRepository{
CreateFunc: func(ctx context.Context, item *domain.SaleItem) error {
return nil
},
}
mockInventoryRepo := &mocks.MockInventoryRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
return &domain.InventoryItem{
ID:       "i1",
Quantity: 100,
}, nil
},
}

service := NewSaleItemService(mockSaleItemRepo, mockInventoryRepo)
item := &domain.SaleItem{
TransactionID:  "t1",
InventoryLotID: "i1",
ItemName:       "Test Item",
Quantity:       5,
PricePerUnit:   10.0,
}

err := service.CreateSaleItem(ctx, item)
if err != nil {
t.Errorf("Expected no error, got %v", err)
}
})

t.Run("insufficient stock", func(t *testing.T) {
mockSaleItemRepo := &mocks.MockSaleItemRepository{}
mockInventoryRepo := &mocks.MockInventoryRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.InventoryItem, error) {
return &domain.InventoryItem{
ID:       "i1",
Quantity: 3,
}, nil
},
}

service := NewSaleItemService(mockSaleItemRepo, mockInventoryRepo)
item := &domain.SaleItem{
TransactionID:  "t1",
InventoryLotID: "i1",
ItemName:       "Test Item",
Quantity:       5,
PricePerUnit:   10.0,
}

err := service.CreateSaleItem(ctx, item)
if err == nil {
t.Error("Expected insufficient stock error")
}
})
}

func TestSaleItemService_GetSaleItem(t *testing.T) {
ctx := context.Background()

t.Run("success", func(t *testing.T) {
expected := &domain.SaleItem{
ID:       "si1",
ItemName: "Test Item",
}
mockRepo := &mocks.MockSaleItemRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.SaleItem, error) {
return expected, nil
},
}

service := NewSaleItemService(mockRepo, nil)
result, err := service.GetSaleItem(ctx, "si1")

if err != nil {
t.Errorf("Expected no error, got %v", err)
}
if result.ID != expected.ID {
t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
}
})
}

func TestSaleItemService_ListItemsForTransaction(t *testing.T) {
ctx := context.Background()

items := []*domain.SaleItem{
{ID: "si1", TransactionID: "t1"},
{ID: "si2", TransactionID: "t1"},
}
mockRepo := &mocks.MockSaleItemRepository{
ListByTransactionFunc: func(ctx context.Context, txnID string) ([]*domain.SaleItem, error) {
return items, nil
},
}

service := NewSaleItemService(mockRepo, nil)
result, err := service.ListItemsForTransaction(ctx, "t1")

if err != nil {
t.Errorf("Expected no error, got %v", err)
}
if len(result) != 2 {
t.Errorf("Expected 2 items, got %d", len(result))
}
}

func TestSaleItemService_DeleteSaleItem(t *testing.T) {
ctx := context.Background()

t.Run("success", func(t *testing.T) {
mockRepo := &mocks.MockSaleItemRepository{
DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
return nil
},
}

service := NewSaleItemService(mockRepo, nil)
req := &domain.DeleteRequest{
Reason:      "Test deletion",
Attestation: "I CONFIRM DELETE",
}

err := service.DeleteSaleItem(ctx, "si1", req)
if err != nil {
t.Errorf("Expected no error, got %v", err)
}
})

t.Run("missing reason", func(t *testing.T) {
mockRepo := &mocks.MockSaleItemRepository{}
service := NewSaleItemService(mockRepo, nil)
req := &domain.DeleteRequest{
Attestation: "I CONFIRM DELETE",
}

err := service.DeleteSaleItem(ctx, "si1", req)
if err == nil {
t.Error("Expected validation error for missing reason")
}
})
}
