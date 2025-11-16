package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// InventoryItem represents an inventory item/lot
type InventoryItem struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Variant          string    `json:"variant,omitempty"`
	LotNumber        string    `json:"lotNumber"`
	Quantity         float64   `json:"quantity"`
	Unit             string    `json:"unit"`
	PurchaseDate     string    `json:"purchaseDate"`
	ExpiryDate       string    `json:"expiryDate"`
	Category         string    `json:"category,omitempty"`
	SubCategory      string    `json:"subCategory,omitempty"`
	CostPrice        float64   `json:"costPrice,omitempty"`
	SellingPrice     float64   `json:"sellingPrice,omitempty"`
	MarginPercentage float64   `json:"marginPercentage,omitempty"`
	SupplierID       string    `json:"supplierId,omitempty"`
	SupplierName     string    `json:"supplierName,omitempty"`
	PurchaseInvoice  string    `json:"purchaseInvoice,omitempty"`
	MinStockLevel    float64   `json:"minStockLevel,omitempty"`
	ReorderPoint     float64   `json:"reorderPoint,omitempty"`
	ShelfLifeDays    int       `json:"shelfLifeDays,omitempty"`
	StorageLocation  string    `json:"storageLocation,omitempty"`
	Barcode          string    `json:"barcode,omitempty"`
	SKU              string    `json:"sku,omitempty"`
	HSNCode          string    `json:"hsnCode,omitempty"`
	GSTRate          float64   `json:"gstRate,omitempty"`
	Status           string    `json:"status,omitempty"`
	WeightPerUnit    float64   `json:"weightPerUnit,omitempty"`
	PackagingType    string    `json:"packagingType,omitempty"`
	ImageURL         string    `json:"imageUrl,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	TotalSold        float64   `json:"totalSold,omitempty"`
	TotalWasted      float64   `json:"totalWasted,omitempty"`
	LastRestockDate  string    `json:"lastRestockDate,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// Handler for listing inventory items with FEFO sorting
func listInventory(w http.ResponseWriter, r *http.Request) {
	// Support filtering by status
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort")

	query := `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, 
		       category, sub_category, cost_price, selling_price, margin_percentage,
		       supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		       shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		       weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		       last_restock_date, created_at, updated_at 
		FROM inventory_items 
		WHERE deleted_at IS NULL`

	// Add status filter
	if status != "" {
		query += ` AND status = '` + status + `'`
	}

	// Default FEFO sorting (earliest expiry first)
	if sortBy == "expiry" || sortBy == "" {
		query += ` ORDER BY expiry_date ASC, created_at DESC`
	} else if sortBy == "name" {
		query += ` ORDER BY name ASC`
	} else if sortBy == "quantity" {
		query += ` ORDER BY quantity ASC`
	} else {
		query += ` ORDER BY created_at DESC`
	}

	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	items := []InventoryItem{}
	for rows.Next() {
		var item InventoryItem
		var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
		var storageLocation, barcode, sku, hsnCode, status, packagingType, imageUrl, notes sql.NullString
		var costPrice, sellingPrice, marginPct, minStock, reorderPoint, weightPerUnit sql.NullFloat64
		var totalSold, totalWasted sql.NullFloat64
		var shelfLifeDays sql.NullInt64
		var gstRate sql.NullFloat64
		var lastRestockDate sql.NullTime

		if err := rows.Scan(&item.ID, &item.Name, &variant, &item.LotNumber, &item.Quantity, &item.Unit,
			&item.PurchaseDate, &item.ExpiryDate, &category, &subCategory, &costPrice,
			&sellingPrice, &marginPct, &supplierId, &supplierName,
			&purchaseInvoice, &minStock, &reorderPoint, &shelfLifeDays,
			&storageLocation, &barcode, &sku, &hsnCode, &gstRate, &status,
			&weightPerUnit, &packagingType, &imageUrl, &notes, &totalSold,
			&totalWasted, &lastRestockDate, &item.CreatedAt, &item.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if variant.Valid {
			item.Variant = variant.String
		}
		if category.Valid {
			item.Category = category.String
		}
		if subCategory.Valid {
			item.SubCategory = subCategory.String
		}
		if costPrice.Valid {
			item.CostPrice = costPrice.Float64
		}
		if sellingPrice.Valid {
			item.SellingPrice = sellingPrice.Float64
		}
		if marginPct.Valid {
			item.MarginPercentage = marginPct.Float64
		}
		if supplierId.Valid {
			item.SupplierID = supplierId.String
		}
		if supplierName.Valid {
			item.SupplierName = supplierName.String
		}
		if purchaseInvoice.Valid {
			item.PurchaseInvoice = purchaseInvoice.String
		}
		if minStock.Valid {
			item.MinStockLevel = minStock.Float64
		}
		if reorderPoint.Valid {
			item.ReorderPoint = reorderPoint.Float64
		}
		if shelfLifeDays.Valid {
			item.ShelfLifeDays = int(shelfLifeDays.Int64)
		}
		if storageLocation.Valid {
			item.StorageLocation = storageLocation.String
		}
		if barcode.Valid {
			item.Barcode = barcode.String
		}
		if sku.Valid {
			item.SKU = sku.String
		}
		if hsnCode.Valid {
			item.HSNCode = hsnCode.String
		}
		if gstRate.Valid {
			item.GSTRate = gstRate.Float64
		}
		if status.Valid {
			item.Status = status.String
		}
		if weightPerUnit.Valid {
			item.WeightPerUnit = weightPerUnit.Float64
		}
		if packagingType.Valid {
			item.PackagingType = packagingType.String
		}
		if imageUrl.Valid {
			item.ImageURL = imageUrl.String
		}
		if notes.Valid {
			item.Notes = notes.String
		}
		if totalSold.Valid {
			item.TotalSold = totalSold.Float64
		}
		if totalWasted.Valid {
			item.TotalWasted = totalWasted.Float64
		}
		if lastRestockDate.Valid {
			item.LastRestockDate = lastRestockDate.Time.Format("2006-01-02")
		}

		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Handler for getting a single inventory item
func getInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var item InventoryItem
	var variant, category, subCategory, supplierId, supplierName, purchaseInvoice sql.NullString
	var storageLocation, barcode, sku, hsnCode, status, packagingType, imageUrl, notes sql.NullString
	var costPrice, sellingPrice, marginPct, minStock, reorderPoint, weightPerUnit sql.NullFloat64
	var totalSold, totalWasted sql.NullFloat64
	var shelfLifeDays sql.NullInt64
	var gstRate sql.NullFloat64
	var lastRestockDate sql.NullTime

	err := db.QueryRowContext(r.Context(), `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		       category, sub_category, cost_price, selling_price, margin_percentage,
		       supplier_id, supplier_name, purchase_invoice, min_stock_level, reorder_point,
		       shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate, status,
		       weight_per_unit, packaging_type, image_url, notes, total_sold, total_wasted,
		       last_restock_date, created_at, updated_at 
		FROM inventory_items WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&item.ID, &item.Name, &variant, &item.LotNumber,
		&item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpiryDate, &category, &subCategory,
		&costPrice, &sellingPrice, &marginPct, &supplierId, &supplierName,
		&purchaseInvoice, &minStock, &reorderPoint, &shelfLifeDays,
		&storageLocation, &barcode, &sku, &hsnCode, &gstRate, &status,
		&weightPerUnit, &packagingType, &imageUrl, &notes, &totalSold,
		&totalWasted, &lastRestockDate, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "inventory item not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Map nullable values
	if variant.Valid {
		item.Variant = variant.String
	}
	if category.Valid {
		item.Category = category.String
	}
	if subCategory.Valid {
		item.SubCategory = subCategory.String
	}
	if costPrice.Valid {
		item.CostPrice = costPrice.Float64
	}
	if sellingPrice.Valid {
		item.SellingPrice = sellingPrice.Float64
	}
	if marginPct.Valid {
		item.MarginPercentage = marginPct.Float64
	}
	if supplierId.Valid {
		item.SupplierID = supplierId.String
	}
	if supplierName.Valid {
		item.SupplierName = supplierName.String
	}
	if purchaseInvoice.Valid {
		item.PurchaseInvoice = purchaseInvoice.String
	}
	if minStock.Valid {
		item.MinStockLevel = minStock.Float64
	}
	if reorderPoint.Valid {
		item.ReorderPoint = reorderPoint.Float64
	}
	if shelfLifeDays.Valid {
		item.ShelfLifeDays = int(shelfLifeDays.Int64)
	}
	if storageLocation.Valid {
		item.StorageLocation = storageLocation.String
	}
	if barcode.Valid {
		item.Barcode = barcode.String
	}
	if sku.Valid {
		item.SKU = sku.String
	}
	if hsnCode.Valid {
		item.HSNCode = hsnCode.String
	}
	if gstRate.Valid {
		item.GSTRate = gstRate.Float64
	}
	if status.Valid {
		item.Status = status.String
	}
	if weightPerUnit.Valid {
		item.WeightPerUnit = weightPerUnit.Float64
	}
	if packagingType.Valid {
		item.PackagingType = packagingType.String
	}
	if imageUrl.Valid {
		item.ImageURL = imageUrl.String
	}
	if notes.Valid {
		item.Notes = notes.String
	}
	if totalSold.Valid {
		item.TotalSold = totalSold.Float64
	}
	if totalWasted.Valid {
		item.TotalWasted = totalWasted.Float64
	}
	if lastRestockDate.Valid {
		item.LastRestockDate = lastRestockDate.Time.Format("2006-01-02")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Handler for creating inventory
func createInventory(w http.ResponseWriter, r *http.Request) {
	var item InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Validation
	if item.Name == "" || item.Quantity <= 0 || item.Unit == "" {
		http.Error(w, "name, quantity and unit are required", 400)
		return
	}

	if item.Unit != "kg" && item.Unit != "lot" && item.Unit != "piece" && item.Unit != "box" {
		http.Error(w, "unit must be 'kg', 'lot', 'piece', or 'box'", 400)
		return
	}

	// Generate lot number if not provided
	if item.LotNumber == "" {
		item.LotNumber = generateLotNumber()
	}

	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO inventory_items (id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
		    category, sub_category, cost_price, selling_price, supplier_id, supplier_name, purchase_invoice,
		    min_stock_level, reorder_point, shelf_life_days, storage_location, barcode, sku, hsn_code, gst_rate,
		    weight_per_unit, packaging_type, image_url, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`,
		item.ID, item.Name, item.Variant, item.LotNumber, item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate,
		item.Category, item.SubCategory, item.CostPrice, item.SellingPrice, item.SupplierID, item.SupplierName,
		item.PurchaseInvoice, item.MinStockLevel, item.ReorderPoint, item.ShelfLifeDays, item.StorageLocation,
		item.Barcode, item.SKU, item.HSNCode, item.GSTRate, item.WeightPerUnit, item.PackagingType, item.ImageURL,
		item.Notes, item.CreatedAt, item.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// Handler for updating inventory
func updateInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var item InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	item.UpdatedAt = time.Now()

	result, err := db.ExecContext(r.Context(), `
		UPDATE inventory_items 
		SET name=$1, variant=$2, quantity=$3, unit=$4, purchase_date=$5, expiry_date=$6,
		    category=$7, sub_category=$8, cost_price=$9, selling_price=$10, supplier_id=$11, supplier_name=$12,
		    purchase_invoice=$13, min_stock_level=$14, reorder_point=$15, shelf_life_days=$16, storage_location=$17,
		    barcode=$18, sku=$19, hsn_code=$20, gst_rate=$21, weight_per_unit=$22, packaging_type=$23, image_url=$24,
		    notes=$25, updated_at=$26, updated_by=$27
		WHERE id=$28 AND deleted_at IS NULL`,
		item.Name, item.Variant, item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate,
		item.Category, item.SubCategory, item.CostPrice, item.SellingPrice, item.SupplierID, item.SupplierName,
		item.PurchaseInvoice, item.MinStockLevel, item.ReorderPoint, item.ShelfLifeDays, item.StorageLocation,
		item.Barcode, item.SKU, item.HSNCode, item.GSTRate, item.WeightPerUnit, item.PackagingType, item.ImageURL,
		item.Notes, item.UpdatedAt, userID, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "inventory item not found", 404)
		return
	}

	item.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Handler for soft deleting inventory
