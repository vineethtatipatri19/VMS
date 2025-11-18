package e2e

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestDirectDatabaseCreate_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	// Create repository directly
	repo := postgres.NewCustomerRepository(ctx.DB)

	customer := &domain.Customer{
		Name:          "Direct DB Test",
		ContactNumber: "+9876543210",
		Status:        "active",
		CustomerType:  "b2c",
	}

	// Try to create directly
	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Fatalf("Failed to create customer directly: %v", err)
	}

	t.Logf("Successfully created customer with ID: %s", customer.ID)
}
