package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) repository.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		INSERT INTO inventory_items (
			id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
			category, sub_category, cost_price, selling_price,
			supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
			shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
			weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
			last_restock_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33
		)`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.Name, toNullString(item.Variant), item.LotNumber,
		item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate,
		toNullString(item.Category), toNullString(item.SubCategory),
		item.CostPrice, item.SellingPrice,
		toNullString(item.SupplierID), toNullString(item.SupplierName),
		toNullString(item.PurchaseInvoice), item.MinStockLevel, item.ReorderPoint,
		item.ShelfLifeDays, toNullString(item.StorageLocation),
		toNullString(item.Barcode), toNullString(item.SKU),
		toNullString(item.HSNCode), item.GSTRate, item.Status,
		item.WeightPerUnit, toNullString(item.PackagingType),
		toNullString(item.ImageURL), toNullString(item.Notes),
		item.TotalSold, item.TotalWasted, toNullString(item.LastRestockDate),
		item.CreatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create inventory item: %w", err)
	}
	return nil
}

func (r *inventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	query := `SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		category, sub_category, cost_price, selling_price, margin_percentage,
		supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		last_restock_date, created_at, updated_at
	FROM inventory_items 
	WHERE id = $1 AND deleted_at IS NULL`

	var item domain.InventoryItem
	var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
	var storageLocation, barcode, sku, hsnCode, packagingType, imageUrl, notes sql.NullString
	var lastRestockDate sql.NullString
	var weightPerUnit sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.Name, &variant, &item.LotNumber, &item.Quantity, &item.Unit,
		&item.PurchaseDate, &item.ExpiryDate,
		&category, &subCategory, &item.CostPrice, &item.SellingPrice, &item.MarginPercentage,
		&supplierId, &supplierName, &purchaseInvoice, &item.MinStockLevel, &item.ReorderPoint,
		&item.ShelfLifeDays, &storageLocation, &barcode, &sku, &hsnCode, &item.GSTRate, &item.Status,
		&weightPerUnit, &packagingType, &imageUrl, &notes, &item.TotalSold, &item.TotalWasted,
		&lastRestockDate, &item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	item.Variant = fromNullString(variant)
	item.Category = fromNullString(category)
	item.SubCategory = fromNullString(subCategory)
	item.SupplierID = fromNullString(supplierId)
	item.SupplierName = fromNullString(supplierName)
	item.PurchaseInvoice = fromNullString(purchaseInvoice)
	item.StorageLocation = fromNullString(storageLocation)
	item.Barcode = fromNullString(barcode)
	item.SKU = fromNullString(sku)
	item.HSNCode = fromNullString(hsnCode)
	item.PackagingType = fromNullString(packagingType)
	item.ImageURL = fromNullString(imageUrl)
	item.Notes = fromNullString(notes)
	item.LastRestockDate = fromNullString(lastRestockDate)
	if weightPerUnit.Valid {
		item.WeightPerUnit = weightPerUnit.Float64
	}

	return &item, nil
}

func (r *inventoryRepository) List(ctx context.Context, status string, sortBy string) ([]*domain.InventoryItem, error) {
	query := `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		       category, sub_category, cost_price, selling_price, margin_percentage,
		       supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		       shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		       weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		       last_restock_date, created_at, updated_at
		FROM inventory_items 
		WHERE deleted_at IS NULL`

	args := []interface{}{}
	if status != "" {
		query += ` AND status = $1`
		args = append(args, status)
	}

	switch sortBy {
	case "name":
		query += ` ORDER BY name ASC`
	case "quantity":
		query += ` ORDER BY quantity ASC`
	case "expiry", "":
		query += ` ORDER BY expiry_date ASC, created_at DESC`
	default:
		query += ` ORDER BY created_at DESC`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	defer rows.Close()

	items := []*domain.InventoryItem{}
	for rows.Next() {
		var item domain.InventoryItem
		var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
		var storageLocation, barcode, sku, hsnCode, packagingType, imageUrl, notes sql.NullString
		var lastRestockDate sql.NullString
		var weightPerUnit sql.NullFloat64

		if err := rows.Scan(
			&item.ID, &item.Name, &variant, &item.LotNumber, &item.Quantity, &item.Unit,
			&item.PurchaseDate, &item.ExpiryDate,
			&category, &subCategory, &item.CostPrice, &item.SellingPrice, &item.MarginPercentage,
			&supplierId, &supplierName, &purchaseInvoice, &item.MinStockLevel, &item.ReorderPoint,
			&item.ShelfLifeDays, &storageLocation, &barcode, &sku, &hsnCode, &item.GSTRate, &item.Status,
			&weightPerUnit, &packagingType, &imageUrl, &notes, &item.TotalSold, &item.TotalWasted,
			&lastRestockDate, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory item: %w", err)
		}

		item.Variant = fromNullString(variant)
		item.Category = fromNullString(category)
		item.SubCategory = fromNullString(subCategory)
		item.SupplierID = fromNullString(supplierId)
		item.SupplierName = fromNullString(supplierName)
		item.PurchaseInvoice = fromNullString(purchaseInvoice)
		item.StorageLocation = fromNullString(storageLocation)
		item.Barcode = fromNullString(barcode)
		item.SKU = fromNullString(sku)
		item.HSNCode = fromNullString(hsnCode)
		item.PackagingType = fromNullString(packagingType)
		item.ImageURL = fromNullString(imageUrl)
		item.Notes = fromNullString(notes)
		item.LastRestockDate = fromNullString(lastRestockDate)
		if weightPerUnit.Valid {
			item.WeightPerUnit = weightPerUnit.Float64
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *inventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		UPDATE inventory_items SET
			name = $2, variant = $3, lot_number = $4, quantity = $5, unit = $6,
			purchase_date = $7, expiry_date = $8, category = $9, sub_category = $10,
			cost_price = $11, selling_price = $12,
			supplier_id = $13, supplier_name = $14, purchase_invoice = $15,
			min_stock_level = $16, reorder_point = $17, shelf_life_days = $18,
			storage_location = $19, barcode = $20, sku = $21, hsn_code = $22,
			gst_rate = $23, status = $24, weight_per_unit = $25, packaging_type = $26,
			image_url = $27, notes = $28, total_sold = $29, total_wasted = $30,
			last_restock_date = $31, updated_at = $32
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		item.ID, item.Name, toNullString(item.Variant), item.LotNumber,
		item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate,
		toNullString(item.Category), toNullString(item.SubCategory),
		item.CostPrice, item.SellingPrice,
		toNullString(item.SupplierID), toNullString(item.SupplierName),
		toNullString(item.PurchaseInvoice), item.MinStockLevel, item.ReorderPoint,
		item.ShelfLifeDays, toNullString(item.StorageLocation),
		toNullString(item.Barcode), toNullString(item.SKU),
		toNullString(item.HSNCode), item.GSTRate, item.Status,
		item.WeightPerUnit, toNullString(item.PackagingType),
		toNullString(item.ImageURL), toNullString(item.Notes),
		item.TotalSold, item.TotalWasted, toNullString(item.LastRestockDate),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update inventory item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *inventoryRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `
		UPDATE inventory_items 
		SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete inventory item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *inventoryRepository) UpdateQuantity(ctx context.Context, id string, delta float64) error {
	query := `
		UPDATE inventory_items 
		SET quantity = quantity + $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, delta, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *inventoryRepository) GetExpiringSoon(ctx context.Context, days int) ([]*domain.InventoryItem, error) {
	query := `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		       category, sub_category, cost_price, selling_price, margin_percentage,
		       supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		       shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		       weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		       last_restock_date, created_at, updated_at
		FROM inventory_items 
		WHERE deleted_at IS NULL 
		  AND expiry_date IS NOT NULL
		  AND expiry_date <= CURRENT_DATE + $1
		  AND expiry_date >= CURRENT_DATE
		  AND quantity > 0
		ORDER BY expiry_date ASC`

	rows, err := r.db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring items: %w", err)
	}
	defer rows.Close()

	items := []*domain.InventoryItem{}
	for rows.Next() {
		var item domain.InventoryItem
		var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
		var storageLocation, barcode, sku, hsnCode, packagingType, imageUrl, notes sql.NullString
		var lastRestockDate sql.NullString
		var weightPerUnit sql.NullFloat64

		if err := rows.Scan(
			&item.ID, &item.Name, &variant, &item.LotNumber, &item.Quantity, &item.Unit,
			&item.PurchaseDate, &item.ExpiryDate,
			&category, &subCategory, &item.CostPrice, &item.SellingPrice, &item.MarginPercentage,
			&supplierId, &supplierName, &purchaseInvoice, &item.MinStockLevel, &item.ReorderPoint,
			&item.ShelfLifeDays, &storageLocation, &barcode, &sku, &hsnCode, &item.GSTRate, &item.Status,
			&weightPerUnit, &packagingType, &imageUrl, &notes, &item.TotalSold, &item.TotalWasted,
			&lastRestockDate, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expiring item: %w", err)
		}

		item.Variant = fromNullString(variant)
		item.Category = fromNullString(category)
		item.SubCategory = fromNullString(subCategory)
		item.SupplierID = fromNullString(supplierId)
		item.SupplierName = fromNullString(supplierName)
		item.PurchaseInvoice = fromNullString(purchaseInvoice)
		item.StorageLocation = fromNullString(storageLocation)
		item.Barcode = fromNullString(barcode)
		item.SKU = fromNullString(sku)
		item.HSNCode = fromNullString(hsnCode)
		item.PackagingType = fromNullString(packagingType)
		item.ImageURL = fromNullString(imageUrl)
		item.Notes = fromNullString(notes)
		item.LastRestockDate = fromNullString(lastRestockDate)
		if weightPerUnit.Valid {
			item.WeightPerUnit = weightPerUnit.Float64
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *inventoryRepository) GetLowStock(ctx context.Context) ([]*domain.InventoryItem, error) {
	query := `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		       category, sub_category, cost_price, selling_price, margin_percentage,
		       supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		       shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		       weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		       last_restock_date, created_at, updated_at
		FROM inventory_items 
		WHERE deleted_at IS NULL 
		  AND quantity <= COALESCE(reorder_point, min_stock_level, 0)
		ORDER BY quantity ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}
	defer rows.Close()

	items := []*domain.InventoryItem{}
	for rows.Next() {
		var item domain.InventoryItem
		var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
		var storageLocation, barcode, sku, hsnCode, packagingType, imageUrl, notes sql.NullString
		var lastRestockDate sql.NullString
		var weightPerUnit sql.NullFloat64

		if err := rows.Scan(
			&item.ID, &item.Name, &variant, &item.LotNumber, &item.Quantity, &item.Unit,
			&item.PurchaseDate, &item.ExpiryDate,
			&category, &subCategory, &item.CostPrice, &item.SellingPrice, &item.MarginPercentage,
			&supplierId, &supplierName, &purchaseInvoice, &item.MinStockLevel, &item.ReorderPoint,
			&item.ShelfLifeDays, &storageLocation, &barcode, &sku, &hsnCode, &item.GSTRate, &item.Status,
			&weightPerUnit, &packagingType, &imageUrl, &notes, &item.TotalSold, &item.TotalWasted,
			&lastRestockDate, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan low stock item: %w", err)
		}

		item.Variant = fromNullString(variant)
		item.Category = fromNullString(category)
		item.SubCategory = fromNullString(subCategory)
		item.SupplierID = fromNullString(supplierId)
		item.SupplierName = fromNullString(supplierName)
		item.PurchaseInvoice = fromNullString(purchaseInvoice)
		item.StorageLocation = fromNullString(storageLocation)
		item.Barcode = fromNullString(barcode)
		item.SKU = fromNullString(sku)
		item.HSNCode = fromNullString(hsnCode)
		item.PackagingType = fromNullString(packagingType)
		item.ImageURL = fromNullString(imageUrl)
		item.Notes = fromNullString(notes)
		item.LastRestockDate = fromNullString(lastRestockDate)
		if weightPerUnit.Valid {
			item.WeightPerUnit = weightPerUnit.Float64
		}

		items = append(items, &item)
	}

	return items, nil
}
