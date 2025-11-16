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
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Variant      string    `json:"variant,omitempty"`
	LotNumber    string    `json:"lotNumber"`
	Quantity     float64   `json:"quantity"`
	Unit         string    `json:"unit"`
	PurchaseDate string    `json:"purchaseDate"`
	ExpiryDate   string    `json:"expiryDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Handler for listing inventory items with FEFO sorting
func listInventory(w http.ResponseWriter, r *http.Request) {
	// Support filtering by status
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort")
	
	query := `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, created_at, updated_at 
		FROM inventory_items 
		WHERE 1=1`
	
	// Add status filter (expiring_soon, expired, fresh)
	if status == "expiring_soon" {
		query += ` AND expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days'`
	} else if status == "expired" {
		query += ` AND expiry_date < CURRENT_DATE`
	} else if status == "fresh" {
		query += ` AND expiry_date > CURRENT_DATE + INTERVAL '3 days'`
	}
	
	// Default FEFO sorting (earliest expiry first)
	if sortBy == "expiry" || sortBy == "" {
		query += ` ORDER BY expiry_date ASC, created_at DESC`
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
		if err := rows.Scan(&item.ID, &item.Name, &item.Variant, &item.LotNumber, &item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpiryDate, &item.CreatedAt, &item.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
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
	err := db.QueryRowContext(r.Context(), `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, created_at, updated_at 
		FROM inventory_items WHERE id=$1`, id).Scan(&item.ID, &item.Name, &item.Variant, &item.LotNumber, &item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpiryDate, &item.CreatedAt, &item.UpdatedAt)
	
	if err == sql.ErrNoRows {
		http.Error(w, "inventory item not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
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
	
	if item.Unit != "kg" && item.Unit != "lot" {
		http.Error(w, "unit must be 'kg' or 'lot'", 400)
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
		INSERT INTO inventory_items (id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		item.ID, item.Name, item.Variant, item.LotNumber, item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate, item.CreatedAt, item.UpdatedAt)
	
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
	var item InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	item.UpdatedAt = time.Now()
	
	result, err := db.ExecContext(r.Context(), `
		UPDATE inventory_items 
		SET name=$1, variant=$2, quantity=$3, unit=$4, purchase_date=$5, expiry_date=$6, updated_at=$7
		WHERE id=$8`,
		item.Name, item.Variant, item.Quantity, item.Unit, item.PurchaseDate, item.ExpiryDate, item.UpdatedAt, id)
	
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

// Handler for deleting inventory
func deleteInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	result, err := db.ExecContext(r.Context(), `DELETE FROM inventory_items WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "inventory item not found", 404)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}
