package service

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/mocks"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("success - valid transaction", func(t *testing.T) {
		mockTxRepo := &mocks.MockTransactionRepository{
			CreateFunc: func(ctx context.Context, tx *domain.Transaction) error {
				return nil
			},
		}

		service := NewTransactionService(mockTxRepo, nil, nil, nil)

		tx := &domain.Transaction{
			CustomerID:  "cust123",
			Type:        "sale",
			TotalAmount: 100,
		}

		err := service.CreateTransaction(ctx, tx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("validation error - missing customer ID", func(t *testing.T) {
		service := NewTransactionService(nil, nil, nil, nil)

		tx := &domain.Transaction{
			Type:        "sale",
			TotalAmount: 100,
		}

		err := service.CreateTransaction(ctx, tx)
		if err == nil {
			t.Error("Expected validation error for missing customer ID")
		}
	})
}

func TestTransactionService_GetTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("success - transaction found", func(t *testing.T) {
		expected := &domain.Transaction{
			ID:          "tx123",
			CustomerID:  "cust123",
			TotalAmount: 100,
		}

		mockRepo := &mocks.MockTransactionRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Transaction, error) {
				if id == "tx123" {
					return expected, nil
				}
				return nil, domain.ErrNotFound
			},
		}

		service := NewTransactionService(mockRepo, nil, nil, nil)
		result, err := service.GetTransaction(ctx, "tx123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID != expected.ID {
			t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
		}
	})

	t.Run("error - transaction not found", func(t *testing.T) {
		mockRepo := &mocks.MockTransactionRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Transaction, error) {
				return nil, domain.ErrNotFound
			},
		}

		service := NewTransactionService(mockRepo, nil, nil, nil)
		_, err := service.GetTransaction(ctx, "nonexistent")

		if err != domain.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}

func TestTransactionService_ListByCustomer(t *testing.T) {
	ctx := context.Background()

	t.Run("success - get customer transactions", func(t *testing.T) {
		expected := []*domain.Transaction{
			{ID: "tx1", CustomerID: "cust123"},
			{ID: "tx2", CustomerID: "cust123"},
		}

		mockRepo := &mocks.MockTransactionRepository{
			ListByCustomerFunc: func(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
				if customerID == "cust123" {
					return expected, nil
				}
				return []*domain.Transaction{}, nil
			},
		}

		service := NewTransactionService(mockRepo, nil, nil, nil)
		result, err := service.ListCustomerTransactions(ctx, "cust123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("Expected 2 transactions, got %d", len(result))
		}
	})
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("success - valid deletion", func(t *testing.T) {
		mockTxRepo := &mocks.MockTransactionRepository{
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				return nil
			},
		}
		mockSaleItemRepo := &mocks.MockSaleItemRepository{
			DeleteByTransactionFunc: func(ctx context.Context, transactionID string, req *domain.DeleteRequest) error {
				return nil
			},
		}

		service := NewTransactionService(mockTxRepo, nil, nil, mockSaleItemRepo)
		req := &domain.DeleteRequest{
			Reason:      "Error in transaction",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteTransaction(ctx, "tx123", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error - missing reason", func(t *testing.T) {
		mockTxRepo := &mocks.MockTransactionRepository{}
		mockSaleItemRepo := &mocks.MockSaleItemRepository{}

		service := NewTransactionService(mockTxRepo, nil, nil, mockSaleItemRepo)
		req := &domain.DeleteRequest{
			Reason:      "",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteTransaction(ctx, "tx123", req)
		if err == nil {
			t.Error("Expected validation error for missing reason")
		}
	})
}
