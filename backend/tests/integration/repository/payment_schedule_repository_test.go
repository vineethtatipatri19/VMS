package repository_test

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/tests/testutil"
)

func TestPaymentScheduleRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	repo := postgres.NewPaymentScheduleRepository(db)
	txnRepo := postgres.NewTransactionRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)
	ctx := context.Background()

	t.Run("Create and Get PaymentSchedule", func(t *testing.T) {
		// Create customer first
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.ContactNumber = testutil.GenerateUUID() // Unique phone
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction
		txn := testutil.FixtureTransaction(func(t *domain.Transaction) {
			t.ID = testutil.GenerateUUID()
			t.CustomerID = customer.ID
			t.Type = "payment"
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		// Create payment schedule
		schedule := testutil.FixturePaymentSchedule(func(s *domain.PaymentSchedule) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.CustomerID = customer.ID
			s.InstallmentNumber = 1
			s.AmountDue = 100.0
			s.AmountPaid = 0.0
			s.Status = "pending"
		})

		err := repo.Create(ctx, schedule)
		testutil.AssertNoError(t, err)

		// Get schedule
		retrieved, err := repo.GetByID(ctx, schedule.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, schedule.TransactionID, retrieved.TransactionID)
		testutil.AssertEqual(t, schedule.CustomerID, retrieved.CustomerID)
		testutil.AssertEqual(t, schedule.InstallmentNumber, retrieved.InstallmentNumber)
		testutil.AssertEqual(t, "pending", retrieved.Status)
	})

	t.Run("ListByCustomer returns customer payment schedules", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.ContactNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction
		txn := testutil.FixtureTransaction(func(t *domain.Transaction) {
			t.ID = testutil.GenerateUUID()
			t.CustomerID = customer.ID
			t.Type = "payment"
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		// Create multiple payment schedules
		for i := 1; i <= 3; i++ {
			schedule := testutil.FixturePaymentSchedule(func(s *domain.PaymentSchedule) {
				s.ID = testutil.GenerateUUID()
				s.TransactionID = txn.ID
				s.CustomerID = customer.ID
				s.InstallmentNumber = i
			})
			testutil.AssertNoError(t, repo.Create(ctx, schedule))
		}

		// List by customer
		schedules, err := repo.ListByCustomer(ctx, customer.ID)
		testutil.AssertNoError(t, err)

		if len(schedules) != 3 {
			t.Errorf("Expected 3 payment schedules, got %d", len(schedules))
		}
	})

	t.Run("Update PaymentSchedule", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.ContactNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction
		txn := testutil.FixtureTransaction(func(t *domain.Transaction) {
			t.ID = testutil.GenerateUUID()
			t.CustomerID = customer.ID
			t.Type = "payment"
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		// Create schedule
		schedule := testutil.FixturePaymentSchedule(func(s *domain.PaymentSchedule) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.CustomerID = customer.ID
			s.AmountDue = 100.0
			s.AmountPaid = 0.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, schedule))

		// Update
		schedule.AmountPaid = 50.0
		err := repo.Update(ctx, schedule)
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, schedule.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, 50.0, retrieved.AmountPaid)
	})

	t.Run("Delete PaymentSchedule", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.ContactNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction
		txn := testutil.FixtureTransaction(func(t *domain.Transaction) {
			t.ID = testutil.GenerateUUID()
			t.CustomerID = customer.ID
			t.Type = "payment"
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		// Create schedule
		schedule := testutil.FixturePaymentSchedule(func(s *domain.PaymentSchedule) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.CustomerID = customer.ID
		})
		testutil.AssertNoError(t, repo.Create(ctx, schedule))

		// Delete
		err := repo.Delete(ctx, schedule.ID)
		testutil.AssertNoError(t, err)

		// Verify soft delete
		_, err = repo.GetByID(ctx, schedule.ID)
		if err == nil {
			t.Error("Expected error when getting soft-deleted schedule")
		}
	})

	t.Run("UpdateStatus changes payment schedule status", func(t *testing.T) {
		// Create customer
		customer := testutil.FixtureCustomer(func(c *domain.Customer) {
			c.ID = testutil.GenerateUUID()
			c.ContactNumber = testutil.GenerateUUID()
		})
		testutil.AssertNoError(t, customerRepo.Create(ctx, customer))

		// Create transaction
		txn := testutil.FixtureTransaction(func(t *domain.Transaction) {
			t.ID = testutil.GenerateUUID()
			t.CustomerID = customer.ID
			t.Type = "payment"
		})
		testutil.AssertNoError(t, txnRepo.Create(ctx, txn))

		// Create pending schedule
		schedule := testutil.FixturePaymentSchedule(func(s *domain.PaymentSchedule) {
			s.ID = testutil.GenerateUUID()
			s.TransactionID = txn.ID
			s.CustomerID = customer.ID
			s.Status = "pending"
			s.AmountDue = 100.0
			s.AmountPaid = 50.0
		})
		testutil.AssertNoError(t, repo.Create(ctx, schedule))

		// Update status to partial
		err := repo.UpdateStatus(ctx, schedule.ID, "partial")
		testutil.AssertNoError(t, err)

		// Verify
		retrieved, err := repo.GetByID(ctx, schedule.ID)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, "partial", retrieved.Status)
	})
}
