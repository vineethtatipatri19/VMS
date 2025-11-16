package service

import (
"context"
"fmt"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/repository"
)

type CustomerService struct {
customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) *CustomerService {
return &CustomerService{
customerRepo: customerRepo,
}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
if err := customer.Validate(); err != nil {
return err
}

existing, err := s.customerRepo.List(ctx)
if err != nil {
return fmt.Errorf("failed to check existing customer: %w", err)
}
if len(existing) > 0 {
return domain.ErrAlreadyExists
}

now := time.Now()
customer.CreatedAt = now
customer.UpdatedAt = now
customer.Status = "active"
customer.CurrentBalance = 0
if err := s.customerRepo.Create(ctx, customer); err != nil {
return fmt.Errorf("failed to create customer: %w", err)
}

return nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
return s.customerRepo.GetByID(ctx, id)
}

func (s *CustomerService) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
customers, err := s.customerRepo.List(ctx)
if err != nil {
return nil, fmt.Errorf("failed to list customers: %w", err)
}
return customers, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
if err := customer.Validate(); err != nil {
return err
}

existing, err := s.customerRepo.GetByID(ctx, customer.ID)
if err != nil {
return err
}

customer.UpdatedAt = time.Now()
customer.CurrentBalance = existing.CurrentBalance
customer.LastTransactionDate = existing.LastTransactionDate

if err := s.customerRepo.Update(ctx, customer); err != nil {
return fmt.Errorf("failed to update customer: %w", err)
}

return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string, req *domain.DeleteRequest) error {
if req.Reason == "" {
return &domain.ValidationError{
Field:   "reason",
Message: "reason is required",
}
}
if req.Attestation == "" {
return &domain.ValidationError{
Field:   "attestation",
Message: "attestation is required",
}
}

customer, err := s.customerRepo.GetByID(ctx, id)
if err != nil {
return err
}

if customer.CurrentBalance != 0 {
return &domain.BusinessError{
Code:    "OUTSTANDING_BALANCE",
Message: fmt.Sprintf("cannot delete customer with outstanding balance: %.2f", customer.CurrentBalance),
}
}

if err := s.customerRepo.Delete(ctx, id, req); err != nil {
return fmt.Errorf("failed to delete customer: %w", err)
}

return nil
}

func (s *CustomerService) GetCustomerBalance(ctx context.Context, customerID string) (float64, error) {
return s.customerRepo.GetBalance(ctx, customerID)
}

func (s *CustomerService) CheckCreditLimit(ctx context.Context, customerID string, amount float64) error {
customer, err := s.customerRepo.GetByID(ctx, customerID)
if err != nil {
return err
}

if customer.Status != "active" {
return &domain.BusinessError{
Code:    "INACTIVE_CUSTOMER",
Message: "customer is not active",
}
}

if !customer.CanPurchase(amount) {
return &domain.BusinessError{
Code:    "CREDIT_LIMIT_EXCEEDED",
Message: fmt.Sprintf("purchase amount %.2f exceeds available credit", amount),
}
}

if customer.IsOverdue() {
return &domain.BusinessError{
Code:    "CUSTOMER_OVERDUE",
Message: "customer has overdue payments",
}
}

return nil
}

func (s *CustomerService) UpdateCustomerBalance(ctx context.Context, customerID string, amount float64, transactionType string) error {	delta := amount
	if transactionType == "payment" {
		delta = -amount
	}

	if err := s.customerRepo.UpdateBalance(ctx, customerID, delta); err != nil {
return fmt.Errorf("failed to update customer balance: %w", err)
}

now := time.Now()
if err := s.customerRepo.UpdateLastTransaction(ctx, customerID, now); err != nil {
return fmt.Errorf("failed to update last transaction date: %w", err)
}

return nil
}

func (s *CustomerService) SearchCustomers(ctx context.Context, query string) ([]*domain.Customer, error) {
	// Note: Search would need custom implementation
	return s.ListCustomers(ctx)
}

func (s *CustomerService) GetActiveCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return s.ListCustomers(ctx)
}

func (s *CustomerService) GetOverdueCustomers(ctx context.Context) ([]*domain.Customer, error) {
customers, err := s.GetActiveCustomers(ctx)
if err != nil {
return nil, err
}

var overdueCustomers []*domain.Customer
for _, customer := range customers {
if customer.IsOverdue() {
overdueCustomers = append(overdueCustomers, customer)
}
}

return overdueCustomers, nil
}
