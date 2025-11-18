package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/google/uuid"
)

type CrateService struct {
	crateRepo    repository.CrateRepository
	customerRepo repository.CustomerRepository
}

func NewCrateService(
	crateRepo repository.CrateRepository,
	customerRepo repository.CustomerRepository,
) *CrateService {
	return &CrateService{
		crateRepo:    crateRepo,
		customerRepo: customerRepo,
	}
}

func (s *CrateService) IssueCrates(ctx context.Context, crate *domain.CrateEntry) error {
	if crate.CratesIssued <= 0 {
		return &domain.ValidationError{
			Field:   "crates_issued",
			Message: "crates issued must be greater than 0",
		}
	}

	customer, err := s.customerRepo.GetByID(ctx, crate.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer.Status != "active" {
		return &domain.BusinessError{
			Code:    "INACTIVE_CUSTOMER",
			Message: "customer is not active",
		}
	}

	// Get current balance
	currentBalance, err := s.crateRepo.GetBalance(ctx, crate.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get crate balance: %w", err)
	}

	// Generate UUID for new crate entry
	crate.ID = uuid.New().String()
	crate.CratesReturned = 0
	crate.Balance = currentBalance + crate.CratesIssued
	crate.Date = time.Now()
	crate.UpdatedAt = time.Now()

	if err := s.crateRepo.Create(ctx, crate); err != nil {
		return fmt.Errorf("failed to issue crates: %w", err)
	}

	return nil
}

func (s *CrateService) ReturnCrates(ctx context.Context, crate *domain.CrateEntry) error {
	if crate.CratesReturned <= 0 {
		return &domain.ValidationError{
			Field:   "crates_returned",
			Message: "crates returned must be greater than 0",
		}
	}

	balance, err := s.crateRepo.GetBalance(ctx, crate.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get crate balance: %w", err)
	}

	if crate.CratesReturned > balance {
		return &domain.BusinessError{
			Code:    "INSUFFICIENT_CRATES",
			Message: fmt.Sprintf("cannot return %d crates, customer only has %d", crate.CratesReturned, balance),
		}
	}

	// Generate UUID for new crate entry
	crate.ID = uuid.New().String()
	crate.CratesIssued = 0
	crate.Balance = balance - crate.CratesReturned
	crate.Date = time.Now()
	crate.UpdatedAt = time.Now()

	if err := s.crateRepo.Create(ctx, crate); err != nil {
		return fmt.Errorf("failed to return crates: %w", err)
	}

	return nil
}

func (s *CrateService) GetCrateBalance(ctx context.Context, customerID string) (int, error) {
	balance, err := s.crateRepo.GetBalance(ctx, customerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get crate balance: %w", err)
	}
	return balance, nil
}

func (s *CrateService) GetCrateHistory(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	crates, err := s.crateRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get crate history: %w", err)
	}
	return crates, nil
}

func (s *CrateService) GetAllCrates(ctx context.Context) ([]*domain.CrateEntry, error) {
	crates, err := s.crateRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all crates: %w", err)
	}
	return crates, nil
}

func (s *CrateService) GetCrate(ctx context.Context, id string) (*domain.CrateEntry, error) {
	return s.crateRepo.GetByID(ctx, id)
}

func (s *CrateService) UpdateCrate(ctx context.Context, crate *domain.CrateEntry) error {
	crate.UpdatedAt = time.Now()

	if err := s.crateRepo.Update(ctx, crate); err != nil {
		return fmt.Errorf("failed to update crate: %w", err)
	}

	return nil
}

func (s *CrateService) DeleteCrate(ctx context.Context, id string, req *domain.DeleteRequest) error {
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

	if err := s.crateRepo.Delete(ctx, id, req); err != nil {
		return fmt.Errorf("failed to delete crate: %w", err)
	}

	return nil
}
