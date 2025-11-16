package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type paymentScheduleRepository struct {
	db *sql.DB
}

func NewPaymentScheduleRepository(db *sql.DB) repository.PaymentScheduleRepository {
	return &paymentScheduleRepository{db: db}
}

func (r *paymentScheduleRepository) Create(ctx context.Context, schedule *domain.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (
		id, transaction_id, customer_id, installment_number, due_date,
		amount_due, amount_paid, status, payment_date, notes, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		schedule.ID, schedule.TransactionID, schedule.CustomerID,
		schedule.InstallmentNumber, schedule.DueDate, schedule.AmountDue,
		schedule.AmountPaid, schedule.Status,
		toNullTime(schedule.PaymentDate), toNullString(schedule.Notes),
		schedule.CreatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create payment schedule: %w", err)
	}
	return nil
}

func (r *paymentScheduleRepository) GetByID(ctx context.Context, id string) (*domain.PaymentSchedule, error) {
	query := `SELECT id, transaction_id, customer_id, installment_number, due_date,
		amount_due, amount_paid, status, payment_date, notes, created_at, updated_at
	FROM payment_schedules 
	WHERE id = $1`

	var schedule domain.PaymentSchedule
	var paymentDate sql.NullTime
	var notes sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID, &schedule.TransactionID, &schedule.CustomerID,
		&schedule.InstallmentNumber, &schedule.DueDate, &schedule.AmountDue,
		&schedule.AmountPaid, &schedule.Status,
		&paymentDate, &notes, &schedule.CreatedAt, &schedule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment schedule: %w", err)
	}

	schedule.PaymentDate = fromNullTime(paymentDate)
	schedule.Notes = fromNullString(notes)

	return &schedule, nil
}

func (r *paymentScheduleRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.PaymentSchedule, error) {
	query := `SELECT id, transaction_id, customer_id, installment_number, due_date,
		amount_due, amount_paid, status, payment_date, notes, created_at, updated_at
	FROM payment_schedules 
	WHERE customer_id = $1
	ORDER BY due_date ASC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment schedules: %w", err)
	}
	defer rows.Close()

	schedules := []*domain.PaymentSchedule{}
	for rows.Next() {
		var schedule domain.PaymentSchedule
		var paymentDate sql.NullTime
		var notes sql.NullString

		if err := rows.Scan(
			&schedule.ID, &schedule.TransactionID, &schedule.CustomerID,
			&schedule.InstallmentNumber, &schedule.DueDate, &schedule.AmountDue,
			&schedule.AmountPaid, &schedule.Status,
			&paymentDate, &notes, &schedule.CreatedAt, &schedule.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment schedule: %w", err)
		}

		schedule.PaymentDate = fromNullTime(paymentDate)
		schedule.Notes = fromNullString(notes)

		schedules = append(schedules, &schedule)
	}

	return schedules, nil
}

func (r *paymentScheduleRepository) Update(ctx context.Context, schedule *domain.PaymentSchedule) error {
	query := `UPDATE payment_schedules SET
		amount_paid = $2, status = $3, payment_date = $4,
		notes = $5, updated_at = $6
	WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		schedule.ID, schedule.AmountPaid, schedule.Status,
		toNullTime(schedule.PaymentDate), toNullString(schedule.Notes),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update payment schedule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *paymentScheduleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM payment_schedules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment schedule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *paymentScheduleRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE payment_schedules SET status = $2, updated_at = $3 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update payment schedule status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
