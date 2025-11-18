package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestInventoryRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	// Create repository
	repo := postgres.NewInventoryRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Inventory Item", func(t *testing.T) {
		item := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.Name = "Integration Test Item"
			i.LotNumber = testutil.GenerateUUID()
			i.Quantity = 100.0
		})

		// Test Create
		err := repo.Create(ctx, item)
		testutil.AssertNoError(t, err)

		// Test GetByID
		retrieved, err := repo.GetByID(ctx, item.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, item.Name, retrieved.Name)
		testutil.AssertEqual(t, item.LotNumber, retrieved.LotNumber)
		testutil.AssertEqual(t, item.Quantity, retrieved.Quantity)
	})

	t.Run("List Inventory Items", func(t *testing.T) {
		// Create multiple items
		for i := 1; i <= 3; i++ {
			item := testutil.FixtureInventoryItem(func(inv *domain.InventoryItem) {
				inv.ID = testutil.GenerateUUID()
				inv.LotNumber = testutil.GenerateUUID()
				inv.Name = fmt.Sprintf("Item %d", i)
				inv.Status = "available" // Use "available" which is set by DB trigger
			})
			testutil.AssertNoError(t, repo.Create(ctx, item))
		}

		// List all items (no status filter to see what we get)
		items, err := repo.List(ctx, "", "")
		testutil.AssertNoError(t, err)

		if len(items) < 3 {
			t.Errorf("Expected at least 3 items, got %d", len(items))
		}
	})

	t.Run("Update Inventory Item", func(t *testing.T) {
		item := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
			i.Name = "Original Name"
			i.Quantity = 50.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, item))

		// Update item
		item.Name = "Updated Name"
		item.Quantity = 75.0
		item.UpdatedAt = time.Now()
		err := repo.Update(ctx, item)
		testutil.AssertNoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(ctx, item.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, "Updated Name", retrieved.Name)
		testutil.AssertEqual(t, 75.0, retrieved.Quantity)
	})

	t.Run("UpdateQuantity adjusts stock", func(t *testing.T) {
		item := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
			i.Quantity = 100.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, item))

		// Decrease quantity by 25
		err := repo.UpdateQuantity(ctx, item.ID, -25.0)
		testutil.AssertNoError(t, err)

		// Verify quantity changed
		retrieved, err := repo.GetByID(ctx, item.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 75.0, retrieved.Quantity)

		// Increase quantity by 10
		err = repo.UpdateQuantity(ctx, item.ID, 10.0)
		testutil.AssertNoError(t, err)

		retrieved, err = repo.GetByID(ctx, item.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 85.0, retrieved.Quantity)
	})

	t.Run("GetExpiringSoon returns items expiring within days threshold", func(t *testing.T) {
		// Create item expiring in 3 days
		tomorrow := time.Now().AddDate(0, 0, 3)
		item1ID := testutil.GenerateUUID()
		item1 := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = item1ID
			i.LotNumber = testutil.GenerateUUID()
			i.Name = "Expiring Soon Item"
			i.ExpiryDate = tomorrow.Format("2006-01-02")
		})
		testutil.AssertNoError(t, repo.Create(ctx, item1))

		// Create item expiring in 30 days (should not be returned)
		laterDate := time.Now().AddDate(0, 0, 30)
		item2 := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
			i.Name = "Not Expiring Soon"
			i.ExpiryDate = laterDate.Format("2006-01-02")
		})
		testutil.AssertNoError(t, repo.Create(ctx, item2))

		// Get items expiring within 7 days
		expiringItems, err := repo.GetExpiringSoon(ctx, 7)
		testutil.AssertNoError(t, err)

		// Should find the first item
		found := false
		for _, item := range expiringItems {
			if item.ID == item1ID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find item expiring within 7 days")
		}
	})

	t.Run("GetLowStock returns items below min stock level", func(t *testing.T) {
		// Create low stock item
		lowStockID := testutil.GenerateUUID()
		item := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = lowStockID
			i.LotNumber = testutil.GenerateUUID()
			i.Name = "Low Stock Item"
			i.Quantity = 5.0
			i.MinStockLevel = 10.0
			i.ReorderPoint = 10.0 // GetLowStock checks reorder_point, not min_stock_level
		})
		testutil.AssertNoError(t, repo.Create(ctx, item))

		// Get low stock items
		lowStockItems, err := repo.GetLowStock(ctx)
		testutil.AssertNoError(t, err)

		// Should find our low stock item
		found := false
		for _, item := range lowStockItems {
			if item.ID == lowStockID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find low stock item")
		}
	})

	t.Run("Delete performs soft delete", func(t *testing.T) {
		item := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, repo.Create(ctx, item))

		// Delete item
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, item.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete (should return error on GetByID)
		_, err = repo.GetByID(ctx, item.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted item")
		}
	})
}
