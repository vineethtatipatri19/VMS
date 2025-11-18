package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

// TestCustomerRepository_Integration demonstrates integration testing with real database
func TestCustomerRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	} // Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	// Create repository
	repo := postgres.NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Customer", func(t *testing.T) {
		// Create test customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.Name = "Integration Test Customer"
		})

		// Test Create
		err := repo.Create(ctx, customer)
		testutil.AssertNoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, customer.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, customer.Name, retrieved.Name)
		testutil.AssertEqual(t, customer.CustomerType, retrieved.CustomerType)
	})

	t.Run("List Customers", func(t *testing.T) {
		// Create multiple customers
		for i := 1; i <= 3; i++ {
			customer := testutil.FixtureCustomer(func(c *domain.Customer) {
				c.ID = testutil.GenerateUUID()
				c.Name = fmt.Sprintf("Customer %d", i)
			})
			testutil.AssertNoError(t, repo.Create(ctx, customer))
		}

		// List all customers
		customers, err := repo.List(ctx)
		testutil.AssertNoError(t, err)

		if len(customers) < 3 {
			t.Errorf("Expected at least 3 customers, got %d", len(customers))
		}
	})

	t.Run("Update Customer", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.Name = "Original Name"
			c.CreditLimit = 5000.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, customer))

		// Update customer
		customer.Name = "Updated Name"
		customer.CreditLimit = 20000.0
		err := repo.Update(ctx, customer)
		testutil.AssertNoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(ctx, customer.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, "Updated Name", retrieved.Name)
		testutil.AssertEqual(t, 20000.0, retrieved.CreditLimit)
	})

	t.Run("Delete Customer", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, repo.Create(ctx, customer))

		// Delete customer
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, customer.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify deletion (should return ErrNotFound or soft delete)
		_, err = repo.GetByID(ctx, customer.ID)
		if err == nil {
			t.Error("Expected error when getting deleted customer")
		}
	})
}
