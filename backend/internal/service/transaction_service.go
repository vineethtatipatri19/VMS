package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
	inventoryRepo   repository.InventoryRepository
	saleItemRepo    repository.SaleItemRepository
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
	inventoryRepo repository.InventoryRepository,
	saleItemRepo repository.SaleItemRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
		inventoryRepo:   inventoryRepo,
		saleItemRepo:    saleItemRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, txn *domain.Transaction) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	// Generate UUID for new transaction
	txn.ID = uuid.New().String()

	now := time.Now()
	txn.CreatedAt = now
	txn.Date = now

	if err := s.transactionRepo.Create(ctx, txn); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (s *TransactionService) CreateSale(ctx context.Context, txn *domain.Transaction, items []*domain.SaleItem) error {
	if txn.Type != "sale" {
		return &domain.ValidationError{
			Field:   "type",
			Message: "transaction type must be 'sale'",
		}
	}

	customer, err := s.customerRepo.GetByID(ctx, txn.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer.Status != "active" {
		return &domain.BusinessError{
			Code:    "INACTIVE_CUSTOMER",
			Message: "customer is not active",
		}
	}

	// Pre-calculate item totals for validation
	var totalAmount float64
	for _, item := range items {
		// Calculate total: Quantity * PricePerUnit (simplified - discounts/taxes applied later)
		if item.Total == 0 {
			item.Total = item.Quantity * item.PricePerUnit
		}
		totalAmount += item.Total
	}

	if !customer.CanPurchase(totalAmount) {
		return &domain.BusinessError{
			Code:    "CREDIT_LIMIT_EXCEEDED",
			Message: fmt.Sprintf("purchase amount %.2f exceeds available credit", totalAmount),
		}
	}

	for _, item := range items {
		invItem, err := s.inventoryRepo.GetByID(ctx, item.InventoryLotID)
		if err != nil {
			return fmt.Errorf("failed to get inventory item %s: %w", item.InventoryLotID, err)
		}

		if invItem.Quantity < item.Quantity {
			return &domain.BusinessError{
				Code: "INSUFFICIENT_STOCK",
				Message: fmt.Sprintf("insufficient stock for item %s: requested %.2f, available %.2f",
					invItem.Name, item.Quantity, invItem.Quantity),
			}
		}

		if invItem.IsExpired() {
			return &domain.BusinessError{
				Code:    "ITEM_EXPIRED",
				Message: fmt.Sprintf("item %s has expired", invItem.Name),
			}
		}
	}

	now := time.Now()
	txn.ID = uuid.New().String()
	txn.Date = now
	txn.CreatedAt = now
	txn.TotalAmount = totalAmount

	if err := s.transactionRepo.Create(ctx, txn); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	for _, item := range items {
		// Get inventory item to populate sale item details
		invItem, err := s.inventoryRepo.GetByID(ctx, item.InventoryLotID)
		if err != nil {
			return fmt.Errorf("failed to get inventory item %s: %w", item.InventoryLotID, err)
		}

		// Populate sale item fields from inventory
		item.ID = uuid.New().String()
		item.TransactionID = txn.ID
		item.ItemName = invItem.Name
		item.BatchNumber = invItem.LotNumber
		item.ExpiryDate = invItem.ExpiryDate
		item.Unit = invItem.Unit
		item.CostPerUnit = invItem.CostPrice
		item.HSNCode = invItem.HSNCode
		item.TaxPercent = invItem.GSTRate

		// Calculate total and profit
		item.CalculateTotal()
		item.CalculateProfit()

		if err := item.Validate(); err != nil {
			return err
		}

		if err := s.saleItemRepo.Create(ctx, item); err != nil {
			return fmt.Errorf("failed to create sale item: %w", err)
		}

		delta := float64(-item.Quantity)
		if err := s.inventoryRepo.UpdateQuantity(ctx, item.InventoryLotID, delta); err != nil {
			return fmt.Errorf("failed to update inventory quantity: %w", err)
		}
	}

	if err := s.customerRepo.UpdateBalance(ctx, customer.ID, totalAmount); err != nil {
		return fmt.Errorf("failed to update customer balance: %w", err)
	}

	if err := s.customerRepo.UpdateLastTransaction(ctx, customer.ID, now); err != nil {
		return fmt.Errorf("failed to update last transaction date: %w", err)
	}

	return nil
}

func (s *TransactionService) CreatePayment(ctx context.Context, txn *domain.Transaction) error {
	if txn.Type != "payment" {
		return &domain.ValidationError{
			Field:   "type",
			Message: "transaction type must be 'payment'",
		}
	}

	customer, err := s.customerRepo.GetByID(ctx, txn.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if txn.TotalAmount > customer.CurrentBalance {
		return &domain.BusinessError{
			Code:    "PAYMENT_EXCEEDS_BALANCE",
			Message: fmt.Sprintf("payment amount %.2f exceeds customer balance %.2f", txn.TotalAmount, customer.CurrentBalance),
		}
	}

	now := time.Now()
	txn.Date = now
	txn.CreatedAt = now
	if err := s.transactionRepo.Create(ctx, txn); err != nil {
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}

	if err := s.customerRepo.UpdateBalance(ctx, customer.ID, -txn.TotalAmount); err != nil {
		return fmt.Errorf("failed to update customer balance: %w", err)
	}

	if err := s.customerRepo.UpdateLastTransaction(ctx, customer.ID, now); err != nil {
		return fmt.Errorf("failed to update last transaction date: %w", err)
	}

	return nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

func (s *TransactionService) ListTransactions(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.List(ctx, txType, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

func (s *TransactionService) ListCustomerTransactions(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list customer transactions: %w", err)
	}
	return transactions, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, txn *domain.Transaction) error {
	if err := txn.Validate(); err != nil {
		return err
	}
	if err := s.transactionRepo.Update(ctx, txn); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id string, req *domain.DeleteRequest) error {
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

	if err := s.transactionRepo.Delete(ctx, id, req); err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}
