package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type PaymentService struct {
	paymentRepo     repository.PaymentScheduleRepository
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
}

func NewPaymentService(
	paymentRepo repository.PaymentScheduleRepository,
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
	}
}

func (s *PaymentService) CreateSchedule(ctx context.Context, schedule *domain.PaymentSchedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	customer, err := s.customerRepo.GetByID(ctx, schedule.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer.Status != "active" {
		return &domain.BusinessError{
			Code:    "INACTIVE_CUSTOMER",
			Message: "customer is not active",
		}
	}

	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now
	schedule.Status = "pending"

	if err := s.paymentRepo.Create(ctx, schedule); err != nil {
		return fmt.Errorf("failed to create payment schedule: %w", err)
	}

	return nil
}

func (s *PaymentService) GetSchedule(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

func (s *PaymentService) ListSchedules(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	schedules, err := s.paymentRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment schedules: %w", err)
	}
	return schedules, nil
}

func (s *PaymentService) RecordPayment(ctx context.Context, scheduleID string, paidAmount float64, paymentDate time.Time) error {
	schedule, err := s.paymentRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.Status == "paid" {
		return &domain.BusinessError{
			Code:    "ALREADY_PAID",
			Message: "payment schedule already marked as paid",
		}
	}

	if paidAmount != schedule.AmountDue {
		return &domain.BusinessError{
			Code:    "AMOUNT_MISMATCH",
			Message: fmt.Sprintf("paid amount %.2f does not match scheduled amount %.2f", paidAmount, schedule.AmountDue),
		}
	}

	if err := s.paymentRepo.UpdateStatus(ctx, scheduleID, "paid"); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

func (s *PaymentService) UpdateSchedule(ctx context.Context, schedule *domain.PaymentSchedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	schedule.UpdatedAt = time.Now()

	if err := s.paymentRepo.Update(ctx, schedule); err != nil {
		return fmt.Errorf("failed to update payment schedule: %w", err)
	}

	return nil
}

func (s *PaymentService) DeleteSchedule(ctx context.Context, id string) error {
	schedule, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if schedule.Status == "paid" {
		return &domain.BusinessError{
			Code:    "CANNOT_DELETE_PAID",
			Message: "cannot delete a paid payment schedule",
		}
	}

	if err := s.paymentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete payment schedule: %w", err)
	}

	return nil
}

func (s *PaymentService) MarkOverdue(ctx context.Context, scheduleID string) error {
	schedule, err := s.paymentRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.Status == "paid" {
		return &domain.BusinessError{
			Code:    "ALREADY_PAID",
			Message: "cannot mark paid schedule as overdue",
		}
	}

	if err := s.paymentRepo.UpdateStatus(ctx, scheduleID, "overdue"); err != nil {
		return fmt.Errorf("failed to mark schedule as overdue: %w", err)
	}

	return nil
}

func (s *PaymentService) GetOverdueSchedules(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	schedules, err := s.paymentRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment schedules: %w", err)
	}

	var overdueSchedules []*domain.PaymentSchedule
	now := time.Now()
	for _, schedule := range schedules {
		if schedule.Status == "pending" && schedule.DueDate.Before(now) {
			overdueSchedules = append(overdueSchedules, schedule)
		}
	}

	return overdueSchedules, nil
}
