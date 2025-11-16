package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// WastageLog represents wastage/spoilage tracking
type WastageLog struct {
	ID              string    `json:"id"`
	InventoryItemID string    `json:"inventoryItemId,omitempty"`
	ItemName        string    `json:"itemName"`
	Quantity        float64   `json:"quantity"`
	Unit            string    `json:"unit"`
	Reason          string    `json:"reason"`
	ReasonDetails   string    `json:"reasonDetails,omitempty"`
	CostValue       float64   `json:"costValue,omitempty"`
	LoggedBy        string    `json:"loggedBy,omitempty"`
	LoggedAt        time.Time `json:"loggedAt"`
	PhotoUrl        string    `json:"photoUrl,omitempty"`
}

// ExpiryAlert represents expiry warning for items
type ExpiryAlert struct {
	ID              string     `json:"id"`
	InventoryItemID string     `json:"inventoryItemId,omitempty"`
	ItemName        string     `json:"itemName,omitempty"`
	AlertDate       string     `json:"alertDate"`
	ExpiryDate      string     `json:"expiryDate"`
	DaysUntilExpiry int        `json:"daysUntilExpiry"`
	Acknowledged    bool       `json:"acknowledged"`
	AcknowledgedAt  *time.Time `json:"acknowledgedAt,omitempty"`
	AcknowledgedBy  string     `json:"acknowledgedBy,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// PaymentSchedule represents installment plans
type PaymentSchedule struct {
	ID             string    `json:"id"`
	TransactionID  string    `json:"transactionId"`
	CustomerID     string    `json:"customerId"`
	InstallmentNum int       `json:"installmentNum"`
	DueDate        string    `json:"dueDate"`
	Amount         float64   `json:"amount"`
	PaidAmount     float64   `json:"paidAmount"`
	Status         string    `json:"status"`
	PaymentDate    string    `json:"paymentDate,omitempty"`
	PaymentMethod  string    `json:"paymentMethod,omitempty"`
	PaymentRef     string    `json:"paymentRef,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// PricingTier represents bulk/wholesale pricing
type PricingTier struct {
	ID            string    `json:"id"`
	ItemName      string    `json:"itemName"`
	MinQuantity   float64   `json:"minQuantity"`
	MaxQuantity   float64   `json:"maxQuantity,omitempty"`
	PricePerUnit  float64   `json:"pricePerUnit"`
	EffectiveDate string    `json:"effectiveDate"`
	CreatedAt     time.Time `json:"createdAt"`
}

// PriceHistory tracks price changes
type PriceHistory struct {
	ID          string    `json:"id"`
	InventoryID string    `json:"inventoryId"`
	ItemName    string    `json:"itemName"`
	OldPrice    float64   `json:"oldPrice"`
	NewPrice    float64   `json:"newPrice"`
	ChangeDate  time.Time `json:"changeDate"`
	ChangedBy   string    `json:"changedBy,omitempty"`
}

// Handler for listing wastage logs
func listWastage(w http.ResponseWriter, r *http.Request) {
	reason := r.URL.Query().Get("reason")

	query := `
		SELECT id, inventory_item_id, item_name, quantity, unit, 
		       reason, reason_details, cost_value, logged_by, logged_at, photo_url 
		FROM wastage_log 
		WHERE deleted_at IS NULL`

	args := []interface{}{}
	if reason != "" {
		query += ` AND reason=$1`
		args = append(args, reason)
	}

	query += ` ORDER BY logged_at DESC`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	logs := []WastageLog{}
	for rows.Next() {
		var log WastageLog
		var inventoryItemID, reasonDetails, loggedBy, photoUrl sql.NullString
		var costValue sql.NullFloat64

		if err := rows.Scan(&log.ID, &inventoryItemID, &log.ItemName, &log.Quantity,
			&log.Unit, &log.Reason, &reasonDetails, &costValue, &loggedBy,
			&log.LoggedAt, &photoUrl); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if inventoryItemID.Valid {
			log.InventoryItemID = inventoryItemID.String
		}
		if reasonDetails.Valid {
			log.ReasonDetails = reasonDetails.String
		}
		if costValue.Valid {
			log.CostValue = costValue.Float64
		}
		if loggedBy.Valid {
			log.LoggedBy = loggedBy.String
		}
		if photoUrl.Valid {
			log.PhotoUrl = photoUrl.String
		}

		logs = append(logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// Handler for creating wastage log entry
func createWastage(w http.ResponseWriter, r *http.Request) {
	var log WastageLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if log.ItemName == "" || log.Quantity <= 0 || log.Reason == "" {
		http.Error(w, "itemName, quantity, and reason are required", 400)
		return
	}

	log.ID = uuid.New().String()
	log.LoggedAt = time.Now()

	var inventoryItemID interface{}
	if log.InventoryItemID != "" {
		inventoryItemID = log.InventoryItemID
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO wastage_log (id, inventory_item_id, item_name, quantity, unit, 
		                         reason, reason_details, cost_value, logged_by, logged_at, photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		log.ID, inventoryItemID, log.ItemName, log.Quantity, log.Unit,
		log.Reason, log.ReasonDetails, log.CostValue, log.LoggedBy, log.LoggedAt, log.PhotoUrl)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(log)
}

// Handler for listing expiry alerts
func listExpiryAlerts(w http.ResponseWriter, r *http.Request) {
	acknowledged := r.URL.Query().Get("acknowledged")

	query := `
		SELECT ea.id, ea.inventory_item_id, ii.name, ea.alert_date, ea.expiry_date, 
		       ea.days_until_expiry, ea.acknowledged, ea.acknowledged_at, ea.acknowledged_by, ea.created_at 
		FROM expiry_alerts ea
		LEFT JOIN inventory_items ii ON ea.inventory_item_id = ii.id
		WHERE ea.deleted_at IS NULL`

	args := []interface{}{}
	if acknowledged != "" {
		query += ` AND ea.acknowledged=$1`
		args = append(args, acknowledged == "true")
	}

	query += ` ORDER BY ea.days_until_expiry ASC, ea.created_at DESC`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	alerts := []ExpiryAlert{}
	for rows.Next() {
		var alert ExpiryAlert
		var inventoryItemID, acknowledgedBy, itemName sql.NullString
		var acknowledgedAt sql.NullTime
		var alertDate, expiryDate sql.NullTime

		if err := rows.Scan(&alert.ID, &inventoryItemID, &itemName, &alertDate, &expiryDate,
			&alert.DaysUntilExpiry, &alert.Acknowledged, &acknowledgedAt, &acknowledgedBy,
			&alert.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if inventoryItemID.Valid {
			alert.InventoryItemID = inventoryItemID.String
		}
		if itemName.Valid {
			alert.ItemName = itemName.String
		}
		if alertDate.Valid {
			alert.AlertDate = alertDate.Time.Format("2006-01-02")
		}
		if expiryDate.Valid {
			alert.ExpiryDate = expiryDate.Time.Format("2006-01-02")
		}
		if acknowledgedAt.Valid {
			t := acknowledgedAt.Time
			alert.AcknowledgedAt = &t
		}
		if acknowledgedBy.Valid {
			alert.AcknowledgedBy = acknowledgedBy.String
		}

		alerts = append(alerts, alert)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// Handler for updating expiry alert status
func updateExpiryAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(string) // From JWT middleware

	var payload struct {
		Acknowledged bool `json:"acknowledged"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE expiry_alerts 
		SET acknowledged=$1, acknowledged_at=NOW(), acknowledged_by=$2 
		WHERE id=$3`,
		payload.Acknowledged, userID, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "alert not found", 404)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handler for listing payment schedules
func listPaymentSchedules(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customerId")
	status := r.URL.Query().Get("status")

	query := `
		SELECT id, transaction_id, customer_id, installment_num, due_date, amount, paid_amount, 
		       status, payment_date, payment_method, payment_ref, notes, created_at, updated_at 
		FROM payment_schedules 
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if customerID != "" {
		argCount++
		query += ` AND customer_id=$` + string(rune('0'+argCount))
		args = append(args, customerID)
	}

	if status != "" {
		argCount++
		query += ` AND status=$` + string(rune('0'+argCount))
		args = append(args, status)
	}

	query += ` ORDER BY due_date ASC`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	schedules := []PaymentSchedule{}
	for rows.Next() {
		var sched PaymentSchedule
		if err := rows.Scan(&sched.ID, &sched.TransactionID, &sched.CustomerID, &sched.InstallmentNum,
			&sched.DueDate, &sched.Amount, &sched.PaidAmount, &sched.Status, &sched.PaymentDate,
			&sched.PaymentMethod, &sched.PaymentRef, &sched.Notes, &sched.CreatedAt, &sched.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		schedules = append(schedules, sched)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// Handler for getting overdue payments view
func getOverduePayments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT customer_id, customer_name, transaction_id, invoice_number, 
		       total_amount, due_date, days_overdue, balance_after 
		FROM overdue_payments 
		ORDER BY days_overdue DESC`)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type OverduePayment struct {
		CustomerID    string  `json:"customerId"`
		CustomerName  string  `json:"customerName"`
		TransactionID string  `json:"transactionId"`
		InvoiceNumber string  `json:"invoiceNumber"`
		TotalAmount   float64 `json:"totalAmount"`
		DueDate       string  `json:"dueDate"`
		DaysOverdue   int     `json:"daysOverdue"`
		BalanceAfter  float64 `json:"balanceAfter"`
	}

	overdueList := []OverduePayment{}
	for rows.Next() {
		var od OverduePayment
		if err := rows.Scan(&od.CustomerID, &od.CustomerName, &od.TransactionID, &od.InvoiceNumber,
			&od.TotalAmount, &od.DueDate, &od.DaysOverdue, &od.BalanceAfter); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		overdueList = append(overdueList, od)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overdueList)
}

// Handler for getting wastage summary view
func getWastageSummary(w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT item_name, total_quantity_wasted, total_value_lost, 
		       most_common_reason, waste_count 
		FROM wastage_summary 
		ORDER BY total_value_lost DESC`)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type WastageSummary struct {
		ItemName            string  `json:"itemName"`
		TotalQuantityWasted float64 `json:"totalQuantityWasted"`
		TotalValueLost      float64 `json:"totalValueLost"`
		MostCommonReason    string  `json:"mostCommonReason"`
		WasteCount          int     `json:"wasteCount"`
	}

	summaryList := []WastageSummary{}
	for rows.Next() {
		var ws WastageSummary
		if err := rows.Scan(&ws.ItemName, &ws.TotalQuantityWasted, &ws.TotalValueLost,
			&ws.MostCommonReason, &ws.WasteCount); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		summaryList = append(summaryList, ws)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaryList)
}
