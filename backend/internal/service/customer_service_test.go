package service

import (
"context"
"testing"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/repository/mocks"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
ctx := context.Background()

t.Run("success - valid customer", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{
ListFunc: func(ctx context.Context) ([]*domain.Customer, error) {
return []*domain.Customer{}, nil // No existing customers
},
CreateFunc: func(ctx context.Context, c *domain.Customer) error {
if c.Name == "" {
t.Error("Expected name to be set")
}
return nil
},
}

service := NewCustomerService(mockRepo)
customer := &domain.Customer{
Name:         "Test Customer",
CustomerType: "b2b",
Status:       "active",
}

err := service.CreateCustomer(ctx, customer)
if err != nil {
t.Errorf("Expected no error, got %v", err)
}
})

t.Run("validation error - missing name", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{}
service := NewCustomerService(mockRepo)

customer := &domain.Customer{
CustomerType: "b2b",
}

err := service.CreateCustomer(ctx, customer)
if err == nil {
t.Error("Expected validation error for missing name")
}
})

t.Run("validation error - invalid customer type", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{}
service := NewCustomerService(mockRepo)

customer := &domain.Customer{
Name:         "Test",
CustomerType: "invalid",
}

err := service.CreateCustomer(ctx, customer)
if err == nil {
t.Error("Expected validation error for invalid customer type")
}
})
}

func TestCustomerService_GetCustomer(t *testing.T) {
ctx := context.Background()

t.Run("success - customer found", func(t *testing.T) {
expected := &domain.Customer{
ID:   "cust123",
Name: "Test Customer",
}

mockRepo := &mocks.MockCustomerRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
if id == "cust123" {
return expected, nil
}
return nil, domain.ErrNotFound
},
}

service := NewCustomerService(mockRepo)
result, err := service.GetCustomer(ctx, "cust123")

if err != nil {
t.Errorf("Expected no error, got %v", err)
}
if result.ID != expected.ID {
t.Errorf("Expected ID %s, got %s", expected.ID, result.ID)
}
})

t.Run("error - customer not found", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
return nil, domain.ErrNotFound
},
}

service := NewCustomerService(mockRepo)
_, err := service.GetCustomer(ctx, "nonexistent")

if err != domain.ErrNotFound {
t.Errorf("Expected ErrNotFound, got %v", err)
}
})
}

func TestCustomerService_UpdateCustomer(t *testing.T) {
ctx := context.Background()

t.Run("success - valid update", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
return &domain.Customer{
ID:     "cust123",
Name:   "Old Name",
Status: "active",
}, nil
},
UpdateFunc: func(ctx context.Context, c *domain.Customer) error {
if c.Name == "" {
t.Error("Expected name to be set")
}
return nil
},
}

service := NewCustomerService(mockRepo)
customer := &domain.Customer{
ID:   "cust123",
Name: "New Name",
}

err := service.UpdateCustomer(ctx, customer)
if err != nil {
t.Errorf("Expected no error, got %v", err)
}
})

t.Run("error - customer not found", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{
GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
return nil, domain.ErrNotFound
},
}

service := NewCustomerService(mockRepo)
customer := &domain.Customer{
ID:   "nonexistent",
Name: "Test",
}

err := service.UpdateCustomer(ctx, customer)
if err != domain.ErrNotFound {
t.Errorf("Expected ErrNotFound, got %v", err)
}
})
}

func TestCustomerService_DeleteCustomer(t *testing.T) {
ctx := context.Background()

	t.Run("success - valid deletion", func(t *testing.T) {
		mockRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:             "cust123",
					Name:           "Test",
					CurrentBalance: 0, // Zero balance required for deletion
				}, nil
			},
			DeleteFunc: func(ctx context.Context, id string, req *domain.DeleteRequest) error {
				if req.Reason == "" {
					t.Error("Expected reason to be set")
				}
				return nil
			},
		}

		service := NewCustomerService(mockRepo)
		req := &domain.DeleteRequest{
			Reason:      "No longer needed",
			Attestation: "I CONFIRM DELETE",
		}

		err := service.DeleteCustomer(ctx, "cust123", req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error - invalid attestation", func(t *testing.T) {
mockRepo := &mocks.MockCustomerRepository{}
service := NewCustomerService(mockRepo)

req := &domain.DeleteRequest{
Reason:      "Test",
Attestation: "invalid",
}

err := service.DeleteCustomer(ctx, "cust123", req)
if err == nil {
t.Error("Expected validation error for invalid attestation")
}
})
}

func TestCustomerService_GetBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("success - get balance", func(t *testing.T) {
		mockRepo := &mocks.MockCustomerRepository{
			GetBalanceFunc: func(ctx context.Context, customerID string) (float64, error) {
				return 1500.50, nil
			},
		}

		service := NewCustomerService(mockRepo)
		balance, err := service.GetCustomerBalance(ctx, "cust123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if balance != 1500.50 {
			t.Errorf("Expected balance 1500.50, got %.2f", balance)
		}
	})
}

func TestCustomerService_UpdateBalance(t *testing.T) {
	ctx := context.Background()

	t.Run("success - sale transaction", func(t *testing.T) {
		called := false
		mockRepo := &mocks.MockCustomerRepository{
			UpdateBalanceFunc: func(ctx context.Context, customerID string, delta float64) error {
				called = true
				if delta != 500.0 {
					t.Errorf("Expected delta 500.0, got %.2f", delta)
				}
				return nil
			},
			UpdateLastTransactionFunc: func(ctx context.Context, customerID string, date time.Time) error {
				return nil
			},
		}

		service := NewCustomerService(mockRepo)
		err := service.UpdateCustomerBalance(ctx, "cust123", 500.0, "sale")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !called {
			t.Error("Expected UpdateBalance to be called")
		}
	})

	t.Run("success - payment transaction", func(t *testing.T) {
		var capturedDelta float64
		mockRepo := &mocks.MockCustomerRepository{
			UpdateBalanceFunc: func(ctx context.Context, customerID string, delta float64) error {
				capturedDelta = delta
				return nil
			},
			UpdateLastTransactionFunc: func(ctx context.Context, customerID string, date time.Time) error {
				return nil
			},
		}

		service := NewCustomerService(mockRepo)
		err := service.UpdateCustomerBalance(ctx, "cust123", 500.0, "payment")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if capturedDelta != -500.0 {
			t.Errorf("Expected delta -500.0 for payment, got %.2f", capturedDelta)
		}
	})
}

func TestCustomerService_CheckCreditLimit(t *testing.T) {
	ctx := context.Background()

	t.Run("success - within credit limit", func(t *testing.T) {
		mockRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:             "cust123",
					Status:         "active",
					CreditLimit:    10000,
					CurrentBalance: 0, // Zero balance = not overdue
				}, nil
			},
		}

		service := NewCustomerService(mockRepo)
		err := service.CheckCreditLimit(ctx, "cust123", 5000)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error - exceeds credit limit", func(t *testing.T) {
		mockRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:             "cust123",
					Status:         "active",
					CreditLimit:    10000,
					CurrentBalance: 8000,
				}, nil
			},
		}

		service := NewCustomerService(mockRepo)
		err := service.CheckCreditLimit(ctx, "cust123", 5000)

		if err == nil {
			t.Error("Expected credit limit error")
		}
	})

	t.Run("error - inactive customer", func(t *testing.T) {
		mockRepo := &mocks.MockCustomerRepository{
			GetByIDFunc: func(ctx context.Context, id string) (*domain.Customer, error) {
				return &domain.Customer{
					ID:      "cust123",
					Status:  "inactive",
				}, nil
			},
		}

		service := NewCustomerService(mockRepo)
		err := service.CheckCreditLimit(ctx, "cust123", 5000)

		if err == nil {
			t.Error("Expected inactive customer error")
		}
	})
}