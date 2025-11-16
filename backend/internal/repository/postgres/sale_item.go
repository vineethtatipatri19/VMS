package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type saleItemRepository struct {
	db *sql.DB
}

func NewSaleItemRepository(db *sql.DB) repository.SaleItemRepository {
	return &saleItemRepository{db: db}
}

func (r *saleItemRepository) Create(ctx context.Context, item *domain.SaleItem) error {
	query := `INSERT INTO sale_items (
		id, transaction_id, inventory_lot_id, item_name, batch_number, quantity, unit,
		price_per_unit, total, cost_per_unit, profit, expiry_date
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.TransactionID, item.InventoryLotID, item.ItemName,
		item.BatchNumber, item.Quantity, item.Unit, item.PricePerUnit,
		item.Total, item.CostPerUnit, item.Profit, item.ExpiryDate,
	)

	if err != nil {
		return fmt.Errorf("failed to create sale item: %w", err)
	}
	return nil
}

func (r *saleItemRepository) GetByID(ctx context.Context, id string) (*domain.SaleItem, error) {
	query := `SELECT id, transaction_id, inventory_lot_id, item_name, batch_number, quantity,
		unit, price_per_unit, total, cost_per_unit, profit, expiry_date
	FROM sale_items 
	WHERE id = $1 AND deleted_at IS NULL`

	var item domain.SaleItem

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.TransactionID, &item.InventoryLotID, &item.ItemName,
		&item.BatchNumber, &item.Quantity, &item.Unit, &item.PricePerUnit,
		&item.Total, &item.CostPerUnit, &item.Profit, &item.ExpiryDate,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sale item: %w", err)
	}

	return &item, nil
}

func (r *saleItemRepository) ListByTransaction(ctx context.Context, transactionID string) ([]*domain.SaleItem, error) {
	query := `SELECT id, transaction_id, inventory_lot_id, item_name, batch_number, quantity,
		unit, price_per_unit, total, cost_per_unit, profit, expiry_date
	FROM sale_items 
	WHERE transaction_id = $1 AND deleted_at IS NULL
	ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sale items: %w", err)
	}
	defer rows.Close()

	items := []*domain.SaleItem{}
	for rows.Next() {
		var item domain.SaleItem

		if err := rows.Scan(
			&item.ID, &item.TransactionID, &item.InventoryLotID, &item.ItemName,
			&item.BatchNumber, &item.Quantity, &item.Unit, &item.PricePerUnit,
			&item.Total, &item.CostPerUnit, &item.Profit, &item.ExpiryDate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sale item: %w", err)
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *saleItemRepository) Update(ctx context.Context, item *domain.SaleItem) error {
	query := `UPDATE sale_items SET
		quantity = $2, price_per_unit = $3, total = $4, 
		cost_per_unit = $5, profit = $6
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		item.ID, item.Quantity, item.PricePerUnit, item.Total,
		item.CostPerUnit, item.Profit,
	)

	if err != nil {
		return fmt.Errorf("failed to update sale item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *saleItemRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `UPDATE sale_items 
	SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete sale item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *saleItemRepository) DeleteByTransaction(ctx context.Context, transactionID string, req *domain.DeleteRequest) error {
	query := `UPDATE sale_items 
	SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
	WHERE transaction_id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, transactionID, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete sale items by transaction: %w", err)
	}

	return nil
}
