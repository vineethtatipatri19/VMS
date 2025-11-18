package repository_test

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestCrateRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewCrateRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("Create and Get CrateEntry", func(t *testing.T) {
		// Create customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create crate entry
		crate := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
			c.ID = testutil.GenerateUUID()
			c.CustomerID = customer.ID
			c.TransactionType = "in"
			c.Quantity = 10
		})

		err := repo.Create(ctx, crate)
		testutil.AssertNoError(t, err)

		// Get crate entry
		retrieved, err := repo.GetByID(ctx, crate.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, crate.CustomerID, retrieved.CustomerID)
		testutil.AssertEqual(t, crate.TransactionType, retrieved.TransactionType)
		testutil.AssertEqual(t, crate.Quantity, retrieved.Quantity)
	})

	t.Run("ListByCustomer returns customer crates", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create multiple crate entries
		for i := 1; i <= 3; i++ {
			crate := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
				c.ID = testutil.GenerateUUID()
				c.CustomerID = customer.ID
				c.TransactionType = "in"
				c.Quantity = i * 5
			})
			testutil.AssertNoError(t, repo.Create(ctx, crate))
		}

		// List crates
		crates, err := repo.ListByCustomer(ctx, customer.ID)
		testutil.AssertNoError(t, err)

		if len(crates) < 3 {
			t.Errorf("Expected at least 3 crates, got %d", len(crates))
		}
	})

	t.Run("Update CrateEntry", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create crate
		crate := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
			c.ID = testutil.GenerateUUID()
			c.CustomerID = customer.ID
			c.TransactionType = "in"
			c.Quantity = 10
		})
		testutil.AssertNoError(t, repo.Create(ctx, crate))

		// Update
		crate.Quantity = 15
		crate.TransactionType = "out"
		err := repo.Update(ctx, crate)
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, crate.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 15, retrieved.Quantity)
		testutil.AssertEqual(t, "out", retrieved.TransactionType)
	})

	t.Run("Delete CrateEntry", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create crate
		crate := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
			c.ID = testutil.GenerateUUID()
			c.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, crate))

		// Delete
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, crate.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, crate.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted crate entry")
		}
	})

	t.Run("GetBalance calculates customer crate balance", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Issue 20 crates to customer (out from our perspective)
		issued := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
			c.ID = testutil.GenerateUUID()
			c.CustomerID = customer.ID
			c.TransactionType = "out" // Crates going OUT to customer
			c.Quantity = 20
		})
		testutil.AssertNoError(t, repo.Create(ctx, issued))

		// Customer returns 8 crates (in from our perspective)
		returned := testutil.FixtureCrateEntry(func(c *domain.CrateEntry) {
			c.ID = testutil.GenerateUUID()
			c.CustomerID = customer.ID
			c.TransactionType = "in" // Crates coming IN from customer
			c.Quantity = 8
		})
		testutil.AssertNoError(t, repo.Create(ctx, returned))

		// Check balance (should be 12 - customer owes us 12 crates)
		balance, err := repo.GetBalance(ctx, customer.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 12, balance)
	})
}
