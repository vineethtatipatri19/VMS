package repository_test

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestSaleItemRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewSaleItemRepository(db)
	txnRepo := postgres.NewTransactionRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	invRepo := postgres.NewInventoryRepository(db)
	ctx := context.Background()

	t.Run("Create and Get SaleItem", func(t *testing.T) {
		// Create prerequisites
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create sale item
		saleItem := testutil.FixtureSaleItem(func(s *domain.SaleItem) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.InventoryLotID = invItem.ID
			s.ItemName = "Test Item"
			s.Quantity = 10.0
			s.PricePerUnit = 100.0
			s.Total = 1000.0
		})

		err := repo.Create(ctx, saleItem)
		testutil.AssertNoError(t, err)

		// Get sale item
		retrieved, err := repo.GetByID(ctx, saleItem.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, saleItem.TransactionID, retrieved.TransactionID)
		testutil.AssertEqual(t, saleItem.ItemName, retrieved.ItemName)
		testutil.AssertEqual(t, saleItem.Quantity, retrieved.Quantity)
	})

	t.Run("ListByTransaction returns sale items for transaction", func(t *testing.T) {
		// Create prerequisites
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create multiple sale items
		for i := 1; i <= 3; i++ {
			saleItem := testutil.FixtureSaleItem(func(s *domain.SaleItem) {
				s.ID = testutil.GenerateUUID()
				s.TransactionID = txn.ID
				s.InventoryLotID = invItem.ID
			})
			testutil.AssertNoError(t, repo.Create(ctx, saleItem))
		}

		// List items
		items, err := repo.ListByTransaction(ctx, txn.ID)
		testutil.AssertNoError(t, err)

		if len(items) < 3 {
			t.Errorf("Expected at least 3 sale items, got %d", len(items))
		}
	})

	t.Run("Update SaleItem", func(t *testing.T) {
		// Create prerequisites
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		saleItem := testutil.FixtureSaleItem(func(s *domain.SaleItem) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.InventoryLotID = invItem.ID
			s.Quantity = 10.0
			s.PricePerUnit = 100.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, saleItem))

		// Update
		saleItem.Quantity = 15.0
		saleItem.PricePerUnit = 150.0
		saleItem.Total = 2250.0
		err := repo.Update(ctx, saleItem)
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, saleItem.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 15.0, retrieved.Quantity)
		testutil.AssertEqual(t, 150.0, retrieved.PricePerUnit)
	})

	t.Run("Delete SaleItem", func(t *testing.T) {
		// Create prerequisites
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		saleItem := testutil.FixtureSaleItem(func(s *domain.SaleItem) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.InventoryLotID = invItem.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, saleItem))

		// Delete
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, saleItem.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, saleItem.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted sale item")
		}
	})

	t.Run("DeleteByTransaction removes all items", func(t *testing.T) {
		// Create prerequisites
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create multiple sale items
		for i := 1; i <= 3; i++ {
			saleItem := testutil.FixtureSaleItem(func(s *domain.SaleItem) {
				s.ID = testutil.GenerateUUID()
				s.TransactionID = txn.ID
				s.InventoryLotID = invItem.ID
			})
			testutil.AssertNoError(t, repo.Create(ctx, saleItem))
		}

		// Delete by transaction
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test bulk deletion",
			Attestation: "Integration test",
		}
		err := repo.DeleteByTransaction(ctx, txn.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify all deleted
		items, err := repo.ListByTransaction(ctx, txn.ID)
		testutil.AssertNoError(t, err)
		if len(items) != 0 {
			t.Errorf("Expected 0 items after bulk delete, got %d", len(items))
		}
	})
}
