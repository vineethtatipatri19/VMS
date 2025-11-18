package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (
		id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, status, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.CustomerID, tx.Date, tx.Type, tx.PaymentAmount, tx.TotalAmount,
		toNullString(tx.PaymentMethod), toNullString(tx.PaymentRef),
		toNullString(tx.Notes), toNullString(tx.Status), tx.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, status, created_at
	FROM transactions 
	WHERE id = $1 AND deleted_at IS NULL`

	var tx domain.Transaction
	var paymentMethod, paymentRef, notes, status sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
		&paymentMethod, &paymentRef, &notes, &status, &tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	tx.PaymentMethod = fromNullString(paymentMethod)
	tx.PaymentRef = fromNullString(paymentRef)
	tx.Notes = fromNullString(notes)
	tx.Status = fromNullString(status)

	return &tx, nil
}

func (r *transactionRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, created_at
	FROM transactions 
	WHERE customer_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		var tx domain.Transaction
		var paymentMethod, paymentRef, notes sql.NullString

		if err := rows.Scan(
			&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
			&paymentMethod, &paymentRef, &notes, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		tx.PaymentMethod = fromNullString(paymentMethod)
		tx.PaymentRef = fromNullString(paymentRef)
		tx.Notes = fromNullString(notes)

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (r *transactionRepository) List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, created_at
	FROM transactions 
	WHERE deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	if txType != "" {
		query += fmt.Sprintf(` AND type = $%d`, argCount)
		args = append(args, txType)
		argCount++
	}

	if !startDate.IsZero() {
		query += fmt.Sprintf(` AND date >= $%d`, argCount)
		args = append(args, startDate)
		argCount++
	}

	if !endDate.IsZero() {
		query += fmt.Sprintf(` AND date <= $%d`, argCount)
		args = append(args, endDate)
		argCount++
	}

	query += ` ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		var tx domain.Transaction
		var paymentMethod, paymentRef, notes sql.NullString

		if err := rows.Scan(
			&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
			&paymentMethod, &paymentRef, &notes, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		tx.PaymentMethod = fromNullString(paymentMethod)
		tx.PaymentRef = fromNullString(paymentRef)
		tx.Notes = fromNullString(notes)

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `UPDATE transactions SET
		payment_amount = $2, total_amount = $3, payment_method = $4,
		payment_reference = $5, notes = $6
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.PaymentAmount, tx.TotalAmount,
		toNullString(tx.PaymentMethod), toNullString(tx.PaymentRef),
		toNullString(tx.Notes),
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `UPDATE transactions 
	SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
