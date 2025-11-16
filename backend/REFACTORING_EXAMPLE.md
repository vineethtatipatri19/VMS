# Complete Reference Example: Customer Feature

This document shows a complete, working implementation of the Customer feature using Clean Architecture. Use this as a reference when implementing other features.

## Layer 1: Domain Model

**File**: `internal/domain/customer.go` (already exists)

```go
package domain

import "time"

type Customer struct {
	ID                  string
	Name                string
	Email               string
	Address             string
	ContactNumber       string
	AlternateContact    string
	WhatsappNumber      string
	PhotoURL            string
	BusinessName        string
	GSTIN               string
	CustomerType        string // b2b, b2c, retail, wholesale
	AadhaarVerified     bool
	KYCDocumentType     string
	KYCDocumentNumber   string
	CreditLimit         float64
	CurrentBalance      float64
	PaymentTermsDays    int
	InterestRate        float64
	Status              string // active, inactive, blocked
	Notes               string
	Tags                []string
	LastTransactionDate *time.Time
	TotalPurchases      float64
	LoyaltyPoints       int
	CreatedAt           time.Time
	DeletedAt           *time.Time
	DeletedBy           string
	DeletionReason      string
}

func (c *Customer) Validate() error {
	if c.Name == "" {
		return ErrInvalidInput("name is required")
	}
	if c.CustomerType != "" && !isValidCustomerType(c.CustomerType) {
		return ErrInvalidInput("invalid customer type")
	}
	if c.Status != "" && !isValidStatus(c.Status) {
		return ErrInvalidInput("invalid status")
	}
	if c.CreditLimit < 0 {
		return ErrInvalidInput("credit limit cannot be negative")
	}
	return nil
}

func (c *Customer) CanPurchase(amount float64) bool {
	if c.Status != "active" {
		return false
	}
	return c.CurrentBalance+amount <= c.CreditLimit
}
```

## Layer 2: Repository Interface

**File**: `internal/repository/interfaces.go` (already exists)

```go
package repository

import (
	"context"
	"github.com/example/pgvms/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	List(ctx context.Context) ([]*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id string, req *domain.DeleteRequest) error
	GetBalance(ctx context.Context, customerID string) (float64, error)
	UpdateBalance(ctx context.Context, customerID string, amount float64) error
	UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error
}
```

## Layer 3: Repository Implementation

**File**: `internal/repository/postgres/customer.go` (to create)

```go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/lib/pq"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (
			id, name, email, address, contact_number, alternate_contact, whatsapp_number,
			photo_url, business_name, gstin, customer_type, aadhaar_verified, 
			kyc_document_type, kyc_document_number, credit_limit, current_balance, 
			payment_terms_days, interest_rate, status, notes, tags, 
			last_transaction_date, total_purchases, loyalty_points, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name,
		nullString(customer.Email), nullString(customer.Address),
		nullString(customer.ContactNumber), nullString(customer.AlternateContact),
		nullString(customer.WhatsappNumber), nullString(customer.PhotoURL),
		nullString(customer.BusinessName), nullString(customer.GSTIN),
		customer.CustomerType, customer.AadhaarVerified,
		nullString(customer.KYCDocumentType), nullString(customer.KYCDocumentNumber),
		customer.CreditLimit, customer.CurrentBalance,
		customer.PaymentTermsDays, customer.InterestRate, customer.Status,
		nullString(customer.Notes), pq.Array(customer.Tags),
		nullTime(customer.LastTransactionDate), customer.TotalPurchases,
		customer.LoyaltyPoints, customer.CreatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
		SELECT id, name, email, address, contact_number, alternate_contact, whatsapp_number,
		       photo_url, business_name, gstin, customer_type, aadhaar_verified, 
		       kyc_document_type, kyc_document_number, credit_limit, current_balance, 
		       payment_terms_days, interest_rate, status, notes, tags, 
		       last_transaction_date, total_purchases, loyalty_points, created_at,
		       deleted_at, deleted_by, deletion_reason
		FROM customers 
		WHERE id = $1 AND deleted_at IS NULL`

	var c domain.Customer
	var email, address, contactNumber, alternateContact, whatsappNumber, photoURL sql.NullString
	var businessName, gstin, kycDocType, kycDocNumber, notes sql.NullString
	var lastTransDate, deletedAt sql.NullTime
	var deletedBy, deletionReason sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &email, &address, &contactNumber, &alternateContact,
		&whatsappNumber, &photoURL, &businessName, &gstin, &c.CustomerType, &c.AadhaarVerified,
		&kycDocType, &kycDocNumber, &c.CreditLimit, &c.CurrentBalance,
		&c.PaymentTermsDays, &c.InterestRate, &c.Status, &notes, pq.Array(&c.Tags),
		&lastTransDate, &c.TotalPurchases, &c.LoyaltyPoints, &c.CreatedAt,
		&deletedAt, &deletedBy, &deletionReason,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Map nullable fields
	c.Email = fromNullString(email)
	c.Address = fromNullString(address)
	c.ContactNumber = fromNullString(contactNumber)
	c.AlternateContact = fromNullString(alternateContact)
	c.WhatsappNumber = fromNullString(whatsappNumber)
	c.PhotoURL = fromNullString(photoURL)
	c.BusinessName = fromNullString(businessName)
	c.GSTIN = fromNullString(gstin)
	c.KYCDocumentType = fromNullString(kycDocType)
	c.KYCDocumentNumber = fromNullString(kycDocNumber)
	c.Notes = fromNullString(notes)
	c.LastTransactionDate = fromNullTime(lastTransDate)
	c.DeletedAt = fromNullTime(deletedAt)
	c.DeletedBy = fromNullString(deletedBy)
	c.DeletionReason = fromNullString(deletionReason)

	return &c, nil
}

func (r *customerRepository) List(ctx context.Context) ([]*domain.Customer, error) {
	query := `
		SELECT id, name, email, contact_number, customer_type, status, 
		       current_balance, credit_limit, created_at
		FROM customers 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	customers := []*domain.Customer{}
	for rows.Next() {
		var c domain.Customer
		var email, contactNumber sql.NullString

		if err := rows.Scan(
			&c.ID, &c.Name, &email, &contactNumber, &c.CustomerType, &c.Status,
			&c.CurrentBalance, &c.CreditLimit, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		c.Email = fromNullString(email)
		c.ContactNumber = fromNullString(contactNumber)
		customers = append(customers, &c)
	}

	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers SET
			name = $2, email = $3, address = $4, contact_number = $5,
			customer_type = $6, status = $7, credit_limit = $8,
			notes = $9, updated_at = $10
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, nullString(customer.Email),
		nullString(customer.Address), nullString(customer.ContactNumber),
		customer.CustomerType, customer.Status, customer.CreditLimit,
		nullString(customer.Notes), time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `
		UPDATE customers 
		SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) GetBalance(ctx context.Context, customerID string) (float64, error) {
	var balance float64
	query := `SELECT current_balance FROM customers WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, domain.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (r *customerRepository) UpdateBalance(ctx context.Context, customerID string, amount float64) error {
	query := `
		UPDATE customers 
		SET current_balance = current_balance + $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, customerID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error {
	query := `
		UPDATE customers 
		SET last_transaction_date = $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, customerID, date, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Helper functions
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
```

## Layer 4: Service

**File**: `internal/service/customer_service.go` (to create)

```go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type CustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	// Validate
	if err := customer.Validate(); err != nil {
		return err
	}

	// Set defaults
	if customer.ID == "" {
		customer.ID = uuid.New().String()
	}
	if customer.Status == "" {
		customer.Status = "active"
	}
	if customer.CustomerType == "" {
		customer.CustomerType = "retail"
	}
	customer.CreatedAt = time.Now()

	// Create
	return s.repo.Create(ctx, customer)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput("customer ID is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *CustomerService) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return s.repo.List(ctx)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
	// Validate
	if err := customer.Validate(); err != nil {
		return err
	}

	// Check exists
	existing, err := s.repo.GetByID(ctx, customer.ID)
	if err != nil {
		return err
	}

	// Preserve created timestamp
	customer.CreatedAt = existing.CreatedAt

	// Update
	return s.repo.Update(ctx, customer)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string, req *domain.DeleteRequest) error {
	// Validate delete request
	if err := req.Validate(); err != nil {
		return err
	}

	// Check if customer exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check outstanding balance
	balance, err := s.repo.GetBalance(ctx, id)
	if err != nil {
		return err
	}

	if balance > 0 {
		return domain.NewBusinessError("OUTSTANDING_BALANCE",
			fmt.Sprintf("Cannot delete customer with outstanding balance: ₹%.2f", balance))
	}

	// Perform soft delete
	return s.repo.Delete(ctx, id, req)
}

func (s *CustomerService) RecordPurchase(ctx context.Context, id string, amount float64) error {
	if amount < 0 {
		return domain.ErrInvalidInput("amount cannot be negative")
	}

	// Get customer to check credit limit
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if purchase allowed
	if !customer.CanPurchase(amount) {
		if customer.Status != "active" {
			return domain.NewBusinessError("CUSTOMER_INACTIVE",
				fmt.Sprintf("Customer is %s", customer.Status))
		}
		return domain.NewBusinessError("CREDIT_LIMIT_EXCEEDED",
			fmt.Sprintf("Purchase would exceed credit limit (%.2f/%.2f)",
				customer.CurrentBalance+amount, customer.CreditLimit))
	}

	// Update balance
	if err := s.repo.UpdateBalance(ctx, id, amount); err != nil {
		return err
	}

	// Update last transaction date
	return s.repo.UpdateLastTransaction(ctx, id, time.Now())
}

func (s *CustomerService) RecordPayment(ctx context.Context, id string, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidInput("payment amount must be positive")
	}

	// Update balance (negative to reduce debt)
	if err := s.repo.UpdateBalance(ctx, id, -amount); err != nil {
		return err
	}

	// Update last transaction date
	return s.repo.UpdateLastTransaction(ctx, id, time.Now())
}
```

## Layer 5: HTTP Handler

**File**: `internal/handlers/customer_handler.go` (to create)

```go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/httputil"
	"github.com/example/pgvms/internal/service"
	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.ListCustomers(r.Context())
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	customer, err := h.service.GetCustomer(r.Context(), id)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		httputil.SendError(w, domain.ErrInvalidInput("invalid request body"))
		return
	}

	if err := h.service.CreateCustomer(r.Context(), &customer); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var customer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		httputil.SendError(w, domain.ErrInvalidInput("invalid request body"))
		return
	}

	customer.ID = id

	if err := h.service.UpdateCustomer(r.Context(), &customer); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, domain.ErrInvalidInput("invalid request body"))
		return
	}

	if err := h.service.DeleteCustomer(r.Context(), id, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}
```

## Layer 6: Router Setup

**File**: `internal/router/customer_routes.go` (to create)

```go
package router

import (
	"github.com/example/pgvms/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterCustomerRoutes(r chi.Router, h *handlers.CustomerHandler) {
	r.Route("/api/customers", func(r chi.Router) {
		r.Get("/", h.ListCustomers)
		r.Post("/", h.CreateCustomer)
		r.Get("/{id}", h.GetCustomer)
		r.Put("/{id}", h.UpdateCustomer)
		r.Delete("/{id}", h.DeleteCustomer)
	})
}
```

## Layer 7: Dependency Injection in main.go

**File**: `main.go` (modification)

```go
package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/example/pgvms/internal/config"
	"github.com/example/pgvms/internal/handlers"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/internal/router"
	"github.com/example/pgvms/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize repositories
	customerRepo := postgres.NewCustomerRepository(db)

	// Initialize services
	customerService := service.NewCustomerService(customerRepo)

	// Initialize handlers
	customerHandler := handlers.NewCustomerHandler(customerService)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register routes
	router.RegisterCustomerRoutes(r, customerHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("Server failed:", err)
	}
}
```

## Testing Examples

### Unit Test - Domain

**File**: `internal/domain/customer_test.go`

```go
package domain

import (
	"testing"
)

func TestCustomer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		customer Customer
		wantErr bool
	}{
		{
			name: "valid customer",
			customer: Customer{
				Name: "Test Customer",
				CustomerType: "retail",
				Status: "active",
				CreditLimit: 10000,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			customer: Customer{
				CustomerType: "retail",
			},
			wantErr: true,
		},
		{
			name: "negative credit limit",
			customer: Customer{
				Name: "Test",
				CreditLimit: -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.customer.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomer_CanPurchase(t *testing.T) {
	customer := &Customer{
		Status: "active",
		CreditLimit: 10000,
		CurrentBalance: 5000,
	}

	// Should allow purchase within limit
	if !customer.CanPurchase(3000) {
		t.Error("Should allow purchase within credit limit")
	}

	// Should not allow purchase exceeding limit
	if customer.CanPurchase(6000) {
		t.Error("Should not allow purchase exceeding credit limit")
	}

	// Inactive customer cannot purchase
	customer.Status = "inactive"
	if customer.CanPurchase(1000) {
		t.Error("Inactive customer should not be able to purchase")
	}
}
```

### Service Test - With Mocks

**File**: `internal/service/customer_service_test.go`

```go
package service

import (
	"context"
	"testing"

	"github.com/example/pgvms/internal/domain"
)

// Mock repository
type mockCustomerRepository struct {
	createCalled bool
	getByCalled  bool
	customer     *domain.Customer
	err          error
}

func (m *mockCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	m.createCalled = true
	return m.err
}

func (m *mockCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	m.getByCalled = true
	return m.customer, m.err
}

func (m *mockCustomerRepository) List(ctx context.Context) ([]*domain.Customer, error) {
	return nil, m.err
}

func (m *mockCustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	return m.err
}

func (m *mockCustomerRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	return m.err
}

func (m *mockCustomerRepository) GetBalance(ctx context.Context, id string) (float64, error) {
	return m.customer.CurrentBalance, m.err
}

func (m *mockCustomerRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	return m.err
}

func (m *mockCustomerRepository) UpdateLastTransaction(ctx context.Context, id string, date time.Time) error {
	return m.err
}

// Tests
func TestCustomerService_CreateCustomer(t *testing.T) {
	mockRepo := &mockCustomerRepository{}
	service := NewCustomerService(mockRepo)

	customer := &domain.Customer{
		Name: "Test Customer",
	}

	err := service.CreateCustomer(context.Background(), customer)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !mockRepo.createCalled {
		t.Error("Expected Create to be called on repository")
	}

	if customer.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if customer.Status != "active" {
		t.Error("Expected default status to be 'active'")
	}
}

func TestCustomerService_CreateCustomer_ValidationError(t *testing.T) {
	mockRepo := &mockCustomerRepository{}
	service := NewCustomerService(mockRepo)

	customer := &domain.Customer{
		Name: "", // Invalid: missing name
	}

	err := service.CreateCustomer(context.Background(), customer)

	if err == nil {
		t.Error("Expected validation error")
	}

	if mockRepo.createCalled {
		t.Error("Create should not be called with invalid customer")
	}
}
```

### Integration Test - Repository

**File**: `internal/repository/postgres/customer_test.go`

```go
package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/example/pgvms/internal/domain"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/test_db?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			customer_type VARCHAR(50),
			status VARCHAR(50),
			credit_limit DECIMAL(10,2),
			current_balance DECIMAL(10,2),
			created_at TIMESTAMP,
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255),
			deletion_reason TEXT
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func cleanupTestDB(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS customers")
	db.Close()
}

func TestCustomerRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	repo := NewCustomerRepository(db)

	customer := &domain.Customer{
		ID:          "test-123",
		Name:        "Test Customer",
		CustomerType: "retail",
		Status:      "active",
		CreditLimit: 10000,
		CreatedAt:   time.Now(),
	}

	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	// Verify customer was created
	found, err := repo.GetByID(context.Background(), "test-123")
	if err != nil {
		t.Fatalf("Failed to get customer: %v", err)
	}

	if found.Name != "Test Customer" {
		t.Errorf("Expected name 'Test Customer', got '%s'", found.Name)
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	repo := NewCustomerRepository(db)

	// Create customer
	customer := &domain.Customer{
		ID:        "test-delete",
		Name:      "Delete Me",
		CreatedAt: time.Now(),
	}
	repo.Create(context.Background(), customer)

	// Delete customer
	req := &domain.DeleteRequest{
		Reason:      "Test deletion",
		Attestation: "I CONFIRM DELETE",
	}
	err := repo.Delete(context.Background(), "test-delete", req)
	if err != nil {
		t.Fatalf("Failed to delete customer: %v", err)
	}

	// Verify customer is deleted (not found)
	_, err = repo.GetByID(context.Background(), "test-delete")
	if err != domain.ErrNotFound {
		t.Error("Expected ErrNotFound for deleted customer")
	}
}
```

## Summary

This complete example shows:

1. ✅ Domain model with business logic
2. ✅ Repository interface (contract)
3. ✅ Repository implementation (PostgreSQL)
4. ✅ Service with business rules
5. ✅ Thin HTTP handlers
6. ✅ Router setup
7. ✅ Dependency injection in main.go
8. ✅ Unit tests (domain)
9. ✅ Service tests (with mocks)
10. ✅ Integration tests (real database)

## Key Takeaways

- **Each layer has a single responsibility**
- **Dependencies point inward** (handler → service → repository)
- **Errors are domain-specific** (ErrNotFound, ErrInvalidInput)
- **Business logic lives in services** (credit limit checks, validation)
- **Handlers are thin** (parse, call service, respond)
- **Tests are isolated** (unit, service, integration)

Use this pattern for all other features:
- Inventory
- Transactions
- Crates
- Wastage
- Expiry Alerts
- Payment Schedules
- Reports/Dashboard
