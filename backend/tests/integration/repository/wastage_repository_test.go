package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestWastageRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewWastageRepository(db)
	invRepo := postgres.NewInventoryRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Wastage", func(t *testing.T) {
		// Create inventory item first
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create wastage
		wastage := testutil.FixtureWastageLog(func(w *domain.WastageLog) {
			w.ID = testutil.GenerateUUID()
			w.InventoryID = invItem.ID
			w.ItemName = "Test Wasted Item"
			w.Quantity = 5.0
			w.Reason = "expired"
		})

		err := repo.Create(ctx, wastage)
		testutil.AssertNoError(t, err)

		// Get wastage
		retrieved, err := repo.GetByID(ctx, wastage.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, wastage.InventoryID, retrieved.InventoryID)
		testutil.AssertEqual(t, wastage.ItemName, retrieved.ItemName)
		testutil.AssertEqual(t, wastage.Quantity, retrieved.Quantity)
	})

	t.Run("List wastage entries", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create multiple wastage entries
		for i := 1; i <= 3; i++ {
			wastage := testutil.FixtureWastageLog(func(w *domain.WastageLog) {
				w.ID = testutil.GenerateUUID()
				w.InventoryID = invItem.ID
			})
			testutil.AssertNoError(t, repo.Create(ctx, wastage))
		}

		// List with date range
		startDate := time.Now().AddDate(0, 0, -7)
		endDate := time.Now().AddDate(0, 0, 1)

		wastages, err := repo.List(ctx, startDate, endDate)
		testutil.AssertNoError(t, err)

		if len(wastages) < 3 {
			t.Errorf("Expected at least 3 wastages, got %d", len(wastages))
		}
	})

	t.Run("Update Wastage", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create wastage
		wastage := testutil.FixtureWastageLog(func(w *domain.WastageLog) {
			w.ID = testutil.GenerateUUID()
			w.InventoryID = invItem.ID
			w.Quantity = 5.0
			w.Reason = "damaged"
		})
		testutil.AssertNoError(t, repo.Create(ctx, wastage))

		// Update
		wastage.Quantity = 7.0
		wastage.Reason = "other"
		err := repo.Update(ctx, wastage)
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, wastage.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 7.0, retrieved.Quantity)
		testutil.AssertEqual(t, "other", retrieved.Reason)
	})

	t.Run("Delete Wastage", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create wastage
		wastage := testutil.FixtureWastageLog(func(w *domain.WastageLog) {
			w.ID = testutil.GenerateUUID()
			w.InventoryID = invItem.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, wastage))

		// Delete
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, wastage.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, wastage.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted wastage")
		}
	})
}
