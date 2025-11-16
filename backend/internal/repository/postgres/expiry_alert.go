package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type expiryAlertRepository struct {
	db *sql.DB
}

func NewExpiryAlertRepository(db *sql.DB) repository.ExpiryAlertRepository {
	return &expiryAlertRepository{db: db}
}

func (r *expiryAlertRepository) Create(ctx context.Context, alert *domain.ExpiryAlert) error {
	query := `INSERT INTO expiry_alerts (
		id, inventory_item_id, alert_date, expiry_date, days_until_expiry,
		acknowledged, acknowledged_at, acknowledged_by, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.InventoryItemID, alert.AlertDate, alert.ExpiryDate,
		alert.DaysUntilExpiry, alert.Acknowledged,
		toNullTime(alert.AcknowledgedAt), toNullString(alert.AcknowledgedBy),
		alert.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create expiry alert: %w", err)
	}
	return nil
}

func (r *expiryAlertRepository) GetByID(ctx context.Context, id string) (*domain.ExpiryAlert, error) {
	query := `SELECT id, inventory_item_id, alert_date, expiry_date, days_until_expiry,
		acknowledged, acknowledged_at, acknowledged_by, created_at
	FROM expiry_alerts 
	WHERE id = $1 AND deleted_at IS NULL`

	var alert domain.ExpiryAlert
	var acknowledgedAt sql.NullTime
	var acknowledgedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID, &alert.InventoryItemID, &alert.AlertDate, &alert.ExpiryDate,
		&alert.DaysUntilExpiry, &alert.Acknowledged,
		&acknowledgedAt, &acknowledgedBy, &alert.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get expiry alert: %w", err)
	}

	alert.AcknowledgedAt = fromNullTime(acknowledgedAt)
	alert.AcknowledgedBy = fromNullString(acknowledgedBy)

	return &alert, nil
}

func (r *expiryAlertRepository) List(ctx context.Context, acknowledged bool) ([]*domain.ExpiryAlert, error) {
	query := `SELECT id, inventory_item_id, alert_date, expiry_date, days_until_expiry,
		acknowledged, acknowledged_at, acknowledged_by, created_at
	FROM expiry_alerts 
	WHERE deleted_at IS NULL AND acknowledged = $1
	ORDER BY expiry_date ASC`

	rows, err := r.db.QueryContext(ctx, query, acknowledged)
	if err != nil {
		return nil, fmt.Errorf("failed to list expiry alerts: %w", err)
	}
	defer rows.Close()

	alerts := []*domain.ExpiryAlert{}
	for rows.Next() {
		var alert domain.ExpiryAlert
		var acknowledgedAt sql.NullTime
		var acknowledgedBy sql.NullString

		if err := rows.Scan(
			&alert.ID, &alert.InventoryItemID, &alert.AlertDate, &alert.ExpiryDate,
			&alert.DaysUntilExpiry, &alert.Acknowledged,
			&acknowledgedAt, &acknowledgedBy, &alert.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expiry alert: %w", err)
		}

		alert.AcknowledgedAt = fromNullTime(acknowledgedAt)
		alert.AcknowledgedBy = fromNullString(acknowledgedBy)

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func (r *expiryAlertRepository) Update(ctx context.Context, alert *domain.ExpiryAlert) error {
	query := `UPDATE expiry_alerts SET
		days_until_expiry = $2, acknowledged = $3,
		acknowledged_at = $4, acknowledged_by = $5
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		alert.ID, alert.DaysUntilExpiry, alert.Acknowledged,
		toNullTime(alert.AcknowledgedAt), toNullString(alert.AcknowledgedBy),
	)

	if err != nil {
		return fmt.Errorf("failed to update expiry alert: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *expiryAlertRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `UPDATE expiry_alerts 
	SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete expiry alert: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *expiryAlertRepository) Acknowledge(ctx context.Context, id string, userID string) error {
	query := `UPDATE expiry_alerts SET
		acknowledged = true, acknowledged_at = $2, acknowledged_by = $3
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to acknowledge expiry alert: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
