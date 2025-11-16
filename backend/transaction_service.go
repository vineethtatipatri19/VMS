package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// sale request structure expects items with inventoryLotId, quantity, pricePerUnit
type SaleItemRequest struct {
	InventoryLotId string  `json:"inventoryLotId"`
	ItemName       string  `json:"itemName"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	PricePerUnit   float64 `json:"pricePerUnit"`
}

type TransactionRequest struct {
	CustomerId    string            `json:"customerId"`
	Type          string            `json:"type"` // 'sale' or 'payment'
	PaymentAmount *float64          `json:"paymentAmount,omitempty"`
	Items         []SaleItemRequest `json:"items,omitempty"`
}

type Transaction struct {
	ID               string          `json:"id"`
	CustomerID       string          `json:"customerId"`
	Date             time.Time       `json:"date"`
	Type             string          `json:"type"`
	PaymentAmount    *float64        `json:"paymentAmount,omitempty"`
	TotalAmount      float64         `json:"totalAmount"`
	Details          json.RawMessage `json:"details,omitempty"`
	PaymentMethod    string          `json:"paymentMethod,omitempty"`
	PaymentReference string          `json:"paymentReference,omitempty"`
	DueDate          string          `json:"dueDate,omitempty"`
	IsOverdue        bool            `json:"isOverdue"`
	DiscountAmount   float64         `json:"discountAmount,omitempty"`
	TaxAmount        float64         `json:"taxAmount,omitempty"`
	Notes            string          `json:"notes,omitempty"`
	InvoiceNumber    string          `json:"invoiceNumber,omitempty"`
	ReceiptSent      bool            `json:"receiptSent"`
	BalanceAfter     float64         `json:"balanceAfter,omitempty"`
	SaleType         string          `json:"saleType,omitempty"`
	DeliveryStatus   string          `json:"deliveryStatus,omitempty"`
	DeliveryDate     string          `json:"deliveryDate,omitempty"`
	DeliveryAddress  string          `json:"deliveryAddress,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        *time.Time      `json:"updatedAt,omitempty"`
	UpdatedBy        string          `json:"updatedBy,omitempty"`
}

// Handler for listing transactions
func listTransactions(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customerId")
	txType := r.URL.Query().Get("type")

	query := `
		SELECT id, customer_id, date, type, payment_amount, total_amount, details, payment_method,
		       payment_reference, due_date, is_overdue, discount_amount, tax_amount, notes, invoice_number,
		       receipt_sent, balance_after, sale_type, delivery_status, delivery_date, delivery_address, 
		       created_at, updated_at, updated_by
		FROM transactions 
		WHERE deleted_at IS NULL`

	args := []interface{}{}
	argCount := 0

	if customerID != "" {
		argCount++
		query += ` AND customer_id=$` + string(rune('0'+argCount))
		args = append(args, customerID)
	}

	if txType != "" {
		argCount++
		query += ` AND type=$` + string(rune('0'+argCount))
		args = append(args, txType)
	}

	query += ` ORDER BY date DESC, created_at DESC`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		var tx Transaction
		var paymentMethod, paymentRef, notes, invoiceNum sql.NullString
		var saleType, deliveryStatus, deliveryAddr sql.NullString
		var paymentAmt, balanceAfter sql.NullFloat64
		var discountAmt, taxAmt sql.NullFloat64
		var dueDate sql.NullTime
		var deliveryDate sql.NullTime
		var updatedAt sql.NullTime
		var updatedBy sql.NullString
		var details []byte

		if err := rows.Scan(&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &paymentAmt, &tx.TotalAmount,
			&details, &paymentMethod, &paymentRef, &dueDate, &tx.IsOverdue,
			&discountAmt, &taxAmt, &notes, &invoiceNum, &tx.ReceiptSent,
			&balanceAfter, &saleType, &deliveryStatus, &deliveryDate, &deliveryAddr,
			&tx.CreatedAt, &updatedAt, &updatedBy); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if len(details) > 0 {
			tx.Details = json.RawMessage(details)
		}
		if paymentAmt.Valid {
			pa := paymentAmt.Float64
			tx.PaymentAmount = &pa
		}
		if paymentMethod.Valid {
			tx.PaymentMethod = paymentMethod.String
		}
		if paymentRef.Valid {
			tx.PaymentReference = paymentRef.String
		}
		if dueDate.Valid {
			tx.DueDate = dueDate.Time.Format("2006-01-02")
		}
		if notes.Valid {
			tx.Notes = notes.String
		}
		if invoiceNum.Valid {
			tx.InvoiceNumber = invoiceNum.String
		}
		if balanceAfter.Valid {
			tx.BalanceAfter = balanceAfter.Float64
		}
		if saleType.Valid {
			tx.SaleType = saleType.String
		}
		if deliveryStatus.Valid {
			tx.DeliveryStatus = deliveryStatus.String
		}
		if deliveryDate.Valid {
			tx.DeliveryDate = deliveryDate.Time.Format(time.RFC3339)
		}
		if deliveryAddr.Valid {
			tx.DeliveryAddress = deliveryAddr.String
		}
		if discountAmt.Valid {
			tx.DiscountAmount = discountAmt.Float64
		}
		if taxAmt.Valid {
			tx.TaxAmount = taxAmt.Float64
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			tx.UpdatedAt = &t
		}
		if updatedBy.Valid {
			tx.UpdatedBy = updatedBy.String
		}

		transactions = append(transactions, tx)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// Handler for getting a single transaction
func getTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var tx Transaction
	err := db.QueryRowContext(r.Context(), `
		SELECT id, customer_id, date, type, payment_amount, total_amount, details, created_at 
		FROM transactions WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount, &tx.Details, &tx.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "transaction not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var tr TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if tr.Type != "sale" && tr.Type != "payment" {
		http.Error(w, "invalid type", 400)
		return
	}
	// For sales: validate and perform in DB transaction
	ctx := r.Context()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer func() {
		if p := recover(); p != nil {
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
			if err == sql.ErrNoRows {
				tx.Rollback()
				http.Error(w, "inventory not found", 400)
				return
			}
			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 500)
				return
			}
			if currentQty < it.Quantity {
				tx.Rollback()
				http.Error(w, "insufficient inventory", 400)
				return
			}
			// decrement
			newQty := currentQty - it.Quantity
			_, err = tx.ExecContext(ctx, `UPDATE inventory_items SET quantity=$1, updated_at=now() WHERE id=$2`, newQty, it.InventoryLotId)
			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 500)
				return
			}
			// insert sale_items later after transaction inserted
			total += it.Quantity * it.PricePerUnit
		}
	} else if tr.Type == "payment" && tr.PaymentAmount != nil {
		total = *tr.PaymentAmount
	}
	// insert transaction row
	details, _ := json.Marshal(tr.Items)
	_, err = tx.ExecContext(ctx, `INSERT INTO transactions(id,customer_id,date,type,payment_amount,total_amount,details,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		id, tr.CustomerId, now, tr.Type, tr.PaymentAmount, total, string(details), now)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}
	// insert sale_items records if sale
	if tr.Type == "sale" {
		for _, it := range tr.Items {
			sid := uuid.New().String()
			totalLine := it.Quantity * it.PricePerUnit
			_, err := tx.ExecContext(ctx, `INSERT INTO sale_items(id,transaction_id,inventory_lot_id,item_name,quantity,unit,price_per_unit,total) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				sid, id, it.InventoryLotId, it.ItemName, it.Quantity, it.Unit, it.PricePerUnit, totalLine)
			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
