package repository

import (
	"context"
	"time"

	"github.com/example/pgvms/internal/domain"
)

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	GetByContactNumber(ctx context.Context, contactNumber string) (*domain.Customer, error)
	List(ctx context.Context) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	GetBalance(ctx context.Context, customerID string) (float64, error)
	UpdateBalance(ctx context.Context, customerID string, delta float64) error
	UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error
}

// InventoryRepository defines the interface for inventory data access
type InventoryRepository interface {
	Create(ctx context.Context, item *domain.InventoryItem) error
	GetByID(ctx context.Context, id string) (*domain.InventoryItem, error)
	List(ctx context.Context, status string, sortBy string) ([]*domain.InventoryItem, error)
	Update(ctx context.Context, item *domain.InventoryItem) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	UpdateQuantity(ctx context.Context, id string, delta float64) error
	GetExpiringSoon(ctx context.Context, days int) ([]*domain.InventoryItem, error)
	GetLowStock(ctx context.Context) ([]*domain.InventoryItem, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error)
	List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error)
	Update(ctx context.Context, tx *domain.Transaction) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
}

// SaleItemRepository defines the interface for sale item data access
type SaleItemRepository interface {
	Create(ctx context.Context, item *domain.SaleItem) error
	GetByID(ctx context.Context, id string) (*domain.SaleItem, error)
	ListByTransaction(ctx context.Context, transactionID string) ([]*domain.SaleItem, error)
	Update(ctx context.Context, item *domain.SaleItem) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	DeleteByTransaction(ctx context.Context, transactionID string, req *domain.DeleteRequest) error
}

// CrateRepository defines the interface for crate data access
type CrateRepository interface {
	Create(ctx context.Context, crate *domain.CrateEntry) error
	GetByID(ctx context.Context, id string) (*domain.CrateEntry, error)
	List(ctx context.Context) ([]*domain.CrateEntry, error)
	ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error)
	Update(ctx context.Context, crate *domain.CrateEntry) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	GetBalance(ctx context.Context, customerID string) (int, error)
}

// WastageRepository defines the interface for wastage data access
type WastageRepository interface {
	Create(ctx context.Context, wastage *domain.WastageLog) error
	GetByID(ctx context.Context, id string) (*domain.WastageLog, error)
	List(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error)
	Update(ctx context.Context, wastage *domain.WastageLog) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
}

// ExpiryAlertRepository defines the interface for expiry alert data access
type ExpiryAlertRepository interface {
	Create(ctx context.Context, alert *domain.ExpiryAlert) error
	GetByID(ctx context.Context, id string) (*domain.ExpiryAlert, error)
	List(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error)
	Update(ctx context.Context, alert *domain.ExpiryAlert) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	Acknowledge(ctx context.Context, id string, userID string) error
}

// PaymentScheduleRepository defines the interface for payment schedule data access
type PaymentScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.PaymentSchedule) error
	GetByID(ctx context.Context, id string) (*domain.PaymentSchedule, error)
	ListByCustomer(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error)
	Update(ctx context.Context, schedule *domain.PaymentSchedule) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
}
