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
		id, customer_id, transaction_id, date,
		crates_issued, crates_returned, balance, notes, crate_type, crate_value, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		crate.ID, crate.CustomerID, toNullString(crate.TransactionID),
		crate.Date, crate.CratesIssued, crate.CratesReturned,
		crate.Balance, toNullString(crate.Notes), toNullString(crate.CrateType),
		crate.CrateValue, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create crate: %w", err)
	}
	return nil
}

func (r *crateRepository) GetByID(ctx context.Context, id string) (*domain.CrateEntry, error) {
	query := `SELECT 
		cl.id, cl.customer_id, c.name as customer_name, cl.transaction_id, cl.date,
		cl.crates_issued, cl.crates_returned, cl.balance, cl.notes, 
		cl.crate_type, cl.crate_value, cl.updated_at, cl.updated_by
	FROM crate_ledger cl
	LEFT JOIN customers c ON cl.customer_id = c.id
	WHERE cl.id = $1 AND cl.deleted_at IS NULL`

	var crate domain.CrateEntry
	var notes, crateType, updatedBy, transactionID sql.NullString
	var customerName sql.NullString
	var crateValue sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&crate.ID, &crate.CustomerID, &customerName, &transactionID, &crate.Date,
		&crate.CratesIssued, &crate.CratesReturned, &crate.Balance, &notes,
		&crateType, &crateValue, &crate.UpdatedAt, &updatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get crate: %w", err)
	}

	crate.CustomerName = fromNullString(customerName)
	crate.TransactionID = fromNullString(transactionID)
	crate.Notes = fromNullString(notes)
	crate.CrateType = fromNullString(crateType)
	crate.CrateValue = crateValue.Float64
	crate.UpdatedBy = fromNullString(updatedBy)

	return &crate, nil
}

func (r *crateRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.CrateEntry, error) {
	query := `SELECT 
		cl.id, cl.customer_id, c.name as customer_name, cl.transaction_id, cl.date,
		cl.crates_issued, cl.crates_returned, cl.balance, cl.notes, 
		cl.crate_type, cl.crate_value, cl.updated_at, cl.updated_by
	FROM crate_ledger cl
	LEFT JOIN customers c ON cl.customer_id = c.id
	WHERE cl.customer_id = $1 AND cl.deleted_at IS NULL
	ORDER BY cl.date DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list crates: %w", err)
	}
	defer rows.Close()

	crates := []*domain.CrateEntry{}
	for rows.Next() {
		var crate domain.CrateEntry
		var notes, crateType, updatedBy, transactionID sql.NullString
		var customerName sql.NullString
		var crateValue sql.NullFloat64

		if err := rows.Scan(
			&crate.ID, &crate.CustomerID, &customerName, &transactionID, &crate.Date,
			&crate.CratesIssued, &crate.CratesReturned, &crate.Balance, &notes,
			&crateType, &crateValue, &crate.UpdatedAt, &updatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan crate: %w", err)
		}

		crate.CustomerName = fromNullString(customerName)
		crate.TransactionID = fromNullString(transactionID)
		crate.Notes = fromNullString(notes)
		crate.CrateType = fromNullString(crateType)
		crate.CrateValue = crateValue.Float64
		crate.UpdatedBy = fromNullString(updatedBy)

		crates = append(crates, &crate)
	}

	return crates, nil
}

func (r *crateRepository) List(ctx context.Context) ([]*domain.CrateEntry, error) {
	query := `SELECT 
		cl.id, cl.customer_id, c.name as customer_name, cl.transaction_id, cl.date,
		cl.crates_issued, cl.crates_returned, cl.balance, cl.notes, 
		cl.crate_type, cl.crate_value, cl.updated_at, cl.updated_by
	FROM crate_ledger cl
	LEFT JOIN customers c ON cl.customer_id = c.id
	WHERE cl.deleted_at IS NULL
	ORDER BY cl.date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all crates: %w", err)
	}
	defer rows.Close()

	crates := []*domain.CrateEntry{}
	for rows.Next() {
		var crate domain.CrateEntry
		var notes, crateType, updatedBy, transactionID sql.NullString
		var customerName sql.NullString
		var crateValue sql.NullFloat64

		if err := rows.Scan(
			&crate.ID, &crate.CustomerID, &customerName, &transactionID, &crate.Date,
			&crate.CratesIssued, &crate.CratesReturned, &crate.Balance, &notes,
			&crateType, &crateValue, &crate.UpdatedAt, &updatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan crate: %w", err)
		}

		crate.CustomerName = fromNullString(customerName)
		crate.TransactionID = fromNullString(transactionID)
		crate.Notes = fromNullString(notes)
		crate.CrateType = fromNullString(crateType)
		crate.CrateValue = crateValue.Float64
		crate.UpdatedBy = fromNullString(updatedBy)

		crates = append(crates, &crate)
	}

	return crates, nil
}

func (r *crateRepository) Update(ctx context.Context, crate *domain.CrateEntry) error {
	query := `UPDATE crate_ledger SET
		crates_issued = $2, crates_returned = $3, balance = $4, 
		notes = $5, crate_type = $6, crate_value = $7, updated_at = $8
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		crate.ID, crate.CratesIssued, crate.CratesReturned, crate.Balance,
		toNullString(crate.Notes), toNullString(crate.CrateType), crate.CrateValue, time.Now(),
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
	query := `UPDATE crate_ledger SET 
		deleted_at = $2, deletion_reason = $3
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), req.Reason)
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
	query := `SELECT COALESCE(SUM(balance), 0) as total_balance
	FROM crate_ledger
	WHERE customer_id = $1 AND deleted_at IS NULL`

	var balance int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get crate balance: %w", err)
	}

	return balance, nil
}
