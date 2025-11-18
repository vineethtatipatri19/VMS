package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type crateRepository struct {
	db *sql.DB
}

func NewCrateRepository(db *sql.DB) repository.CrateRepository {
	return &crateRepository{db: db}
}

func (r *crateRepository) Create(ctx context.Context, crate *domain.CrateEntry) error {
	query := `INSERT INTO crate_ledger (
		id, customer_id, transaction_id, transaction_type,
		quantity, unit_price, total_price, balance, notes, created_by, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	// For transaction_type "in", balance increases; for "out", it decreases
	balance := crate.Quantity
	if crate.TransactionType == "out" {
		balance = -crate.Quantity
	}

	_, err := r.db.ExecContext(ctx, query,
		crate.ID, crate.CustomerID, toNullString(crate.TransactionID),
		crate.TransactionType, crate.Quantity, crate.UnitPrice, crate.TotalPrice,
		balance, toNullString(crate.Notes), toNullString(crate.CreatedBy),
		crate.CreatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create crate: %w", err)
	}
	return nil
}

func (r *crateRepository) GetByID(ctx context.Context, id string) (*domain.CrateEntry, error) {
	query := `SELECT id, customer_id, transaction_id, transaction_type,
		quantity, unit_price, total_price, notes, created_by, created_at, updated_at
	FROM crate_ledger 
	WHERE id = $1`

	var crate domain.CrateEntry
	var notes, createdBy, transactionID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&crate.ID, &crate.CustomerID, &transactionID,
		&crate.TransactionType, &crate.Quantity, &crate.UnitPrice, &crate.TotalPrice,
		&notes, &createdBy, &crate.CreatedAt, &crate.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get crate: %w", err)
	}

	crate.TransactionID = fromNullString(transactionID)
	crate.Notes = fromNullString(notes)
	crate.CreatedBy = fromNullString(createdBy)

	return &crate, nil
}

func (r *crateRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	query := `SELECT id, customer_id, transaction_id, transaction_type,
		quantity, unit_price, total_price, notes, created_by, created_at, updated_at
	FROM crate_ledger 
	WHERE customer_id = $1
	ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list crates: %w", err)
	}
	defer rows.Close()

	crates := []*domain.CrateEntry{}
	for rows.Next() {
		var crate domain.CrateEntry
		var notes, createdBy, transactionID sql.NullString

		if err := rows.Scan(
			&crate.ID, &crate.CustomerID, &transactionID,
			&crate.TransactionType, &crate.Quantity, &crate.UnitPrice, &crate.TotalPrice,
			&notes, &createdBy, &crate.CreatedAt, &crate.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan crate: %w", err)
		}

		crate.TransactionID = fromNullString(transactionID)
		crate.Notes = fromNullString(notes)
		crate.CreatedBy = fromNullString(createdBy)

		crates = append(crates, &crate)
	}

	return crates, nil
}

func (r *crateRepository) Update(ctx context.Context, crate *domain.CrateEntry) error {
	query := `UPDATE crate_ledger SET
		transaction_type = $2, quantity = $3, unit_price = $4,
		total_price = $5, balance = $6, notes = $7, updated_at = $8
	WHERE id = $1 AND deleted_at IS NULL`

	// Recalculate balance based on transaction type
	balance := crate.Quantity
	if crate.TransactionType == "out" {
		balance = -crate.Quantity
	}

	result, err := r.db.ExecContext(ctx, query,
		crate.ID, crate.TransactionType, crate.Quantity, crate.UnitPrice,
		crate.TotalPrice, balance, toNullString(crate.Notes), time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update crate: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *crateRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `DELETE FROM crate_ledger WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete crate: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *crateRepository) GetBalance(ctx context.Context, customerID string) (int, error) {
	query := `SELECT COALESCE(SUM(
		CASE 
			WHEN transaction_type = 'out' THEN quantity
			WHEN transaction_type = 'in' THEN -quantity
			ELSE 0
		END
	), 0) as balance
	FROM crate_ledger
	WHERE customer_id = $1 AND deleted_at IS NULL`

	var balance int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get crate balance: %w", err)
	}

	return balance, nil
}
