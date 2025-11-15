
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// sale request structure expects items with inventoryLotId, quantity, pricePerUnit
type SaleItemRequest struct {
	InventoryLotId string  `json:"inventoryLotId"`
	ItemName string           `json:"itemName"`
	Quantity float64         `json:"quantity"`
	Unit string               `json:"unit"`
	PricePerUnit float64      `json:"pricePerUnit"`
}

type TransactionRequest struct {
	CustomerId string `json:"customerId"`
	Type string `json:"type"` // 'sale' or 'payment'
	PaymentAmount *float64 `json:"paymentAmount,omitempty"`
	Items []SaleItemRequest `json:"items,omitempty"`
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var tr TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&tr); err!=nil { http.Error(w, err.Error(), 400); return }
	if tr.Type != "sale" && tr.Type != "payment" { http.Error(w, "invalid type", 400); return }
	// For sales: validate and perform in DB transaction
	ctx := r.Context()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err!=nil { http.Error(w, err.Error(), 500); return }
	defer func() {
		if p := recover(); p!=nil {
			tx.Rollback()
			panic(p)
		}
	}()
	id := uuid.New().String()
	now := time.Now()
	var total float64 = 0
	if tr.Type == "sale" {
		// process each item: check inventory and decrement
		for _, it := range tr.Items {
			// fetch current qty
			var currentQty float64
			err := tx.QueryRowContext(ctx, `SELECT quantity FROM inventory_items WHERE id=$1 FOR UPDATE`, it.InventoryLotId).Scan(&currentQty)
			if err==sql.ErrNoRows { tx.Rollback(); http.Error(w, "inventory not found", 400); return }
			if err!=nil { tx.Rollback(); http.Error(w, err.Error(), 500); return }
			if currentQty < it.Quantity { tx.Rollback(); http.Error(w, "insufficient inventory", 400); return }
			// decrement
			newQty := currentQty - it.Quantity
			_, err = tx.ExecContext(ctx, `UPDATE inventory_items SET quantity=$1, updated_at=now() WHERE id=$2`, newQty, it.InventoryLotId)
			if err!=nil { tx.Rollback(); http.Error(w, err.Error(), 500); return }
			// insert sale_items later after transaction inserted
			total += it.Quantity * it.PricePerUnit
		}
	}
	// insert transaction row
	details, _ := json.Marshal(tr.Items)
	_, err = tx.ExecContext(ctx, `INSERT INTO transactions(id,customer_id,date,type,payment_amount,total_amount,details,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		id, tr.CustomerId, now, tr.Type, tr.PaymentAmount, total, string(details), now)
	if err!=nil { tx.Rollback(); http.Error(w, err.Error(), 500); return }
	// insert sale_items records if sale
	if tr.Type == "sale" {
		for _, it := range tr.Items {
			sid := uuid.New().String()
			totalLine := it.Quantity * it.PricePerUnit
			_, err := tx.ExecContext(ctx, `INSERT INTO sale_items(id,transaction_id,inventory_lot_id,item_name,quantity,unit,price_per_unit,total) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				sid, id, it.InventoryLotId, it.ItemName, it.Quantity, it.Unit, it.PricePerUnit, totalLine)
			if err!=nil { tx.Rollback(); http.Error(w, err.Error(), 500); return }
		}
	}
	if err := tx.Commit(); err!=nil { tx.Rollback(); http.Error(w, err.Error(), 500); return }
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
