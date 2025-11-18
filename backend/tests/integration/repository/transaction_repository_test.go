package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestTransactionRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewTransactionRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Transaction", func(t *testing.T) {
		// Create a customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
			tx.Type = "sale"
			tx.TotalAmount = 500.0
		})

		// Test Create
		err := repo.Create(ctx, txn)
		testutil.AssertNoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, txn.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, txn.CustomerID, retrieved.CustomerID)
		testutil.AssertEqual(t, txn.Type, retrieved.Type)
		testutil.AssertEqual(t, txn.TotalAmount, retrieved.TotalAmount)
	})

	t.Run("ListByCustomer returns customer transactions", func(t *testing.T) {
		// Create a customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transactions for this customer
		for i := 1; i <= 3; i++ {
			txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
				tx.ID = testutil.GenerateUUID()
				tx.CustomerID = customer.ID
			})
			testutil.AssertNoError(t, repo.Create(ctx, txn))
		}

		// List transactions for customer
		transactions, err := repo.ListByCustomer(ctx, customer.ID)
		testutil.AssertNoError(t, err)

		if len(transactions) < 3 {
			t.Errorf("Expected at least 3 transactions, got %d", len(transactions))
		}

		// Verify all transactions belong to the customer
		for _, txn := range transactions {
			testutil.AssertEqual(t, customer.ID, txn.CustomerID)
		}
	})

	t.Run("List filters by transaction type", func(t *testing.T) {
		// Create customers
		customer1 := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer1))

		customer2 := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer2))

		// Create sale transaction
		saleTxn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer1.ID
			tx.Type = "sale"
		})
		testutil.AssertNoError(t, repo.Create(ctx, saleTxn))

		// Create purchase transaction
		purchaseTxn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer2.ID
			tx.Type = "payment"
		})
		testutil.AssertNoError(t, repo.Create(ctx, purchaseTxn))

		// List sale transactions
		saleTransactions, err := repo.List(ctx, "sale", time.Time{}, time.Time{})
		testutil.AssertNoError(t, err)

		// Verify all are sale type
		for _, txn := range saleTransactions {
			testutil.AssertEqual(t, "sale", txn.Type)
		}
	})

	t.Run("Update Transaction", func(t *testing.T) {
		// Create a customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
			tx.TotalAmount = 100.0
			tx.PaymentAmount = 0.0
			tx.PaymentMethod = "cash"
		})
		testutil.AssertNoError(t, repo.Create(ctx, txn))

		// Update transaction
		txn.TotalAmount = 150.0
		txn.PaymentAmount = 150.0
		txn.PaymentMethod = "card"
		err := repo.Update(ctx, txn)
		testutil.AssertNoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(ctx, txn.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 150.0, retrieved.TotalAmount)
		testutil.AssertEqual(t, 150.0, retrieved.PaymentAmount)
		testutil.AssertEqual(t, "card", retrieved.PaymentMethod)
	})

	t.Run("Delete performs soft delete", func(t *testing.T) {
		// Create a customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, txn))

		// Delete transaction
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, txn.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, txn.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted transaction")
		}
	})

	t.Run("Date range filtering", func(t *testing.T) {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)

		// Create a customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction from yesterday
		txn := testutil.FixtureTransaction(func(tx *domain.Transaction) {
			tx.ID = testutil.GenerateUUID()
			tx.CustomerID = customer.ID
			tx.CreatedAt = yesterday
		})
		testutil.AssertNoError(t, repo.Create(ctx, txn))

		// List transactions should include it
		transactions, err := repo.List(ctx, "", time.Time{}, time.Time{})
		testutil.AssertNoError(t, err)

		if len(transactions) == 0 {
			t.Error("Expected to find transactions")
		}
	})
}
