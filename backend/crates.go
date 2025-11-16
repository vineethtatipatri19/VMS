package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CrateLedgerEntry represents a crate transaction
type CrateLedgerEntry struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customerId"`
	Date           time.Time `json:"date"`
	CratesIssued   int       `json:"cratesIssued"`
	CratesReturned int       `json:"cratesReturned"`
	Balance        int       `json:"balance"`
	Notes          string    `json:"notes,omitempty"`
	CrateType      string    `json:"crateType,omitempty"`
	CrateValue     float64   `json:"crateValue,omitempty"`
	TransactionID  string    `json:"transactionId,omitempty"`
}

// Handler for listing crate ledger entries
func listCrates(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customerId")

	query := `
		SELECT id, customer_id, date, crates_issued, crates_returned, balance, notes, crate_type, crate_value, transaction_id 
		FROM crate_ledger 
		WHERE deleted_at IS NULL`

	args := []interface{}{}
	if customerID != "" {
		query += ` AND customer_id=$1`
		args = append(args, customerID)
	}

	query += ` ORDER BY date DESC`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	entries := []CrateLedgerEntry{}
	for rows.Next() {
		var entry CrateLedgerEntry
		var notes, crateType sql.NullString
		var crateValue sql.NullFloat64
		var transactionID sql.NullString

		if err := rows.Scan(&entry.ID, &entry.CustomerID, &entry.Date, &entry.CratesIssued, &entry.CratesReturned,
			&entry.Balance, &notes, &crateType, &crateValue, &transactionID); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if notes.Valid {
			entry.Notes = notes.String
		}
		if crateType.Valid {
			entry.CrateType = crateType.String
		}
		if crateValue.Valid {
			entry.CrateValue = crateValue.Float64
		}
		if transactionID.Valid {
			entry.TransactionID = transactionID.String
		}

		entries = append(entries, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// Handler for creating a crate entry
func createCrateEntry(w http.ResponseWriter, r *http.Request) {
	var entry CrateLedgerEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if entry.CustomerID == "" {
		http.Error(w, "customerId is required", 400)
		return
	}

	// Get current balance for customer
	var currentBalance int
	err := db.QueryRowContext(r.Context(), `
		SELECT COALESCE(balance, 0) 
		FROM crate_ledger 
		WHERE customer_id=$1 AND deleted_at IS NULL
		ORDER BY date DESC 
		LIMIT 1`, entry.CustomerID).Scan(&currentBalance)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), 500)
		return
	}

	// Calculate new balance
	entry.Balance = currentBalance + entry.CratesIssued - entry.CratesReturned

	if entry.Balance < 0 {
		http.Error(w, "cannot return more crates than issued", 400)
		return
	}

	entry.ID = uuid.New().String()
	entry.Date = time.Now()

	_, err = db.ExecContext(r.Context(), `
		INSERT INTO crate_ledger (id, customer_id, date, crates_issued, crates_returned, balance, notes, crate_type, crate_value, transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		entry.ID, entry.CustomerID, entry.Date, entry.CratesIssued, entry.CratesReturned, entry.Balance, entry.Notes,
		entry.CrateType, entry.CrateValue, entry.TransactionID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// Handler for getting crate balance for a customer
func getCrateBalance(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerId")

	var balance int
	err := db.QueryRowContext(r.Context(), `
		SELECT COALESCE(balance, 0) 
		FROM crate_ledger 
		WHERE customer_id=$1 AND deleted_at IS NULL
		ORDER BY date DESC 
		LIMIT 1`, customerID).Scan(&balance)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"customerId": customerID,
		"balance":    balance,
	})
}
