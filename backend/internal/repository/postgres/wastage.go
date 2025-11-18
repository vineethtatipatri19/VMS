package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type wastageRepository struct {
	db *sql.DB
}

func NewWastageRepository(db *sql.DB) repository.WastageRepository {
	return &wastageRepository{db: db}
}

func (r *wastageRepository) Create(ctx context.Context, wastage *domain.WastageLog) error {
	query := `
		INSERT INTO wastage_log (
			id, inventory_item_id, item_name, quantity, unit, reason,
			reason_details, cost_value, logged_by, logged_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		wastage.ID, wastage.InventoryID, wastage.ItemName,
		wastage.Quantity, wastage.Unit, wastage.Reason, wastage.CostValue,
		toNullString(wastage.Notes), toNullString(wastage.RecordedBy),
		wastage.RecordedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create wastage: %w", err)
	}
	return nil
}

func (r *wastageRepository) GetByID(ctx context.Context, id string) (*domain.WastageLog, error) {
	query := `SELECT id, inventory_item_id, item_name, quantity, unit, reason,
		reason_details, cost_value, logged_by, logged_at
	FROM wastage_log 
	WHERE id = $1`

	var wastage domain.WastageLog
	var notes, loggedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&wastage.ID, &wastage.InventoryID, &wastage.ItemName,
		&wastage.Quantity, &wastage.Unit, &wastage.Reason,
		&notes, &wastage.CostValue, &loggedBy, &wastage.RecordedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get wastage: %w", err)
	}

	wastage.Notes = fromNullString(notes)
	wastage.RecordedBy = fromNullString(loggedBy)

	return &wastage, nil
}

func (r *wastageRepository) List(ctx context.Context, startDate, endDate time.Time) ([]*domain.WastageLog, error) {
	query := `SELECT id, inventory_item_id, item_name, quantity, unit, reason,
		reason_details, cost_value, logged_by, logged_at
	FROM wastage_log 
	WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if !startDate.IsZero() {
		query += fmt.Sprintf(` AND logged_at >= $%d`, argCount)
		args = append(args, startDate)
		argCount++
	}

	if !endDate.IsZero() {
		query += fmt.Sprintf(` AND logged_at <= $%d`, argCount)
		args = append(args, endDate)
		argCount++
	}

	query += ` ORDER BY logged_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list wastage: %w", err)
	}
	defer rows.Close()

	wastages := []*domain.WastageLog{}
	for rows.Next() {
		var wastage domain.WastageLog
		var notes, loggedBy sql.NullString

		if err := rows.Scan(
			&wastage.ID, &wastage.InventoryID, &wastage.ItemName,
			&wastage.Quantity, &wastage.Unit, &wastage.Reason,
			&notes, &wastage.CostValue, &loggedBy, &wastage.RecordedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan wastage: %w", err)
		}

		wastage.Notes = fromNullString(notes)
		wastage.RecordedBy = fromNullString(loggedBy)

		wastages = append(wastages, &wastage)
	}

	return wastages, nil
}

func (r *wastageRepository) Update(ctx context.Context, wastage *domain.WastageLog) error {
	query := `UPDATE wastage_log SET
		quantity = $2, reason = $3, cost_value = $4, reason_details = $5
	WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		wastage.ID, wastage.Quantity, wastage.Reason, wastage.CostValue,
		toNullString(wastage.Notes),
	)

	if err != nil {
		return fmt.Errorf("failed to update wastage: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *wastageRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `DELETE FROM wastage_log WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete wastage: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
