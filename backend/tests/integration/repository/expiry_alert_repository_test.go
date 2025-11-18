package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestExpiryAlertRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewExpiryAlertRepository(db)
	invRepo := postgres.NewInventoryRepository(db)
	ctx := context.Background()

	t.Run("Create and Get ExpiryAlert", func(t *testing.T) {
		// Create inventory item first
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create expiry alert
		alert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
			a.AlertDate = time.Now()
			a.ExpiryDate = time.Now().AddDate(0, 0, 7)
			a.DaysUntilExpiry = 7
			a.Acknowledged = false
		})

		err := repo.Create(ctx, alert)
		testutil.AssertNoError(t, err)

		// Get alert
		retrieved, err := repo.GetByID(ctx, alert.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, alert.InventoryItemID, retrieved.InventoryItemID)
		testutil.AssertEqual(t, alert.DaysUntilExpiry, retrieved.DaysUntilExpiry)
		testutil.AssertEqual(t, false, retrieved.Acknowledged)
	})

	t.Run("List expiry alerts with filter", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create acknowledged and unacknowledged alerts
		acknowledgedAlert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
			a.Acknowledged = true
			now := time.Now()
			a.AcknowledgedAt = &now
			a.AcknowledgedBy = "test-user"
		})
		testutil.AssertNoError(t, repo.Create(ctx, acknowledgedAlert))

		unacknowledgedAlert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
			a.Acknowledged = false
		})
		testutil.AssertNoError(t, repo.Create(ctx, unacknowledgedAlert))

		// List only unacknowledged
		alerts, err := repo.List(ctx, false)
		testutil.AssertNoError(t, err)

		// Should have at least our unacknowledged alert
		found := false
		for _, alert := range alerts {
			if alert.ID == unacknowledgedAlert.ID {
				found = true
				testutil.AssertEqual(t, false, alert.Acknowledged)
			}
		}
		if !found {
			t.Error("Expected to find unacknowledged alert in list")
		}
	})

	t.Run("Update ExpiryAlert", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create alert
		alert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
			a.DaysUntilExpiry = 7
		})
		testutil.AssertNoError(t, repo.Create(ctx, alert))

		// Update
		alert.DaysUntilExpiry = 3
		err := repo.Update(ctx, alert)
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, alert.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 3, retrieved.DaysUntilExpiry)
	})

	t.Run("Delete ExpiryAlert", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create alert
		alert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, alert))

		// Delete
		deleteReq := &domain.DeleteRequest{
			Reason:      "Test deletion",
			Attestation: "Integration test",
		}
		err := repo.Delete(ctx, alert.ID, deleteReq)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, alert.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted alert")
		}
	})

	t.Run("Acknowledge ExpiryAlert", func(t *testing.T) {
		// Create inventory item
		invItem := testutil.FixtureInventoryItem(func(i *domain.InventoryItem) {
			i.ID = testutil.GenerateUUID()
			i.LotNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, invRepo.Create(ctx, invItem))

		// Create unacknowledged alert
		alert := testutil.FixtureExpiryAlert(func(a *domain.ExpiryAlert) {
			a.ID = testutil.GenerateUUID()
			a.InventoryItemID = invItem.ID
			a.Acknowledged = false
		})
		testutil.AssertNoError(t, repo.Create(ctx, alert))

		// Acknowledge
		err := repo.Acknowledge(ctx, alert.ID, "test-user")
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, alert.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, true, retrieved.Acknowledged)
		testutil.AssertEqual(t, "test-user", retrieved.AcknowledgedBy)
		if retrieved.AcknowledgedAt == nil {
			t.Error("Expected AcknowledgedAt to be set")
		}
	})
}
