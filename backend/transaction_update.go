package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler for updating a transaction
func updateTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system" // Default
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		PaymentAmount    *float64 `json:"paymentAmount,omitempty"`
		TotalAmount      *float64 `json:"totalAmount,omitempty"`
		PaymentMethod    string   `json:"paymentMethod,omitempty"`
		PaymentReference string   `json:"paymentReference,omitempty"`
		DueDate          string   `json:"dueDate,omitempty"`
		DiscountAmount   *float64 `json:"discountAmount,omitempty"`
		TaxAmount        *float64 `json:"taxAmount,omitempty"`
		Notes            string   `json:"notes,omitempty"`
		InvoiceNumber    string   `json:"invoiceNumber,omitempty"`
		SaleType         string   `json:"saleType,omitempty"`
		DeliveryStatus   string   `json:"deliveryStatus,omitempty"`
		DeliveryDate     string   `json:"deliveryDate,omitempty"`
		DeliveryAddress  string   `json:"deliveryAddress,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Build dynamic UPDATE query
	query := "UPDATE transactions SET updated_at=NOW(), updated_by=$1"
	args := []interface{}{userID}
	argCount := 1

	if req.PaymentAmount != nil {
		argCount++
		query += ", payment_amount=$" + string(rune('0'+argCount))
		args = append(args, req.PaymentAmount)
	}
	if req.TotalAmount != nil {
		argCount++
		query += ", total_amount=$" + string(rune('0'+argCount))
		args = append(args, req.TotalAmount)
	}
	if req.PaymentMethod != "" {
		argCount++
		query += ", payment_method=$" + string(rune('0'+argCount))
		args = append(args, req.PaymentMethod)
	}
	if req.PaymentReference != "" {
		argCount++
		query += ", payment_reference=$" + string(rune('0'+argCount))
		args = append(args, req.PaymentReference)
	}
	if req.DueDate != "" {
		argCount++
		query += ", due_date=$" + string(rune('0'+argCount))
		args = append(args, req.DueDate)
	}
	if req.DiscountAmount != nil {
		argCount++
		query += ", discount_amount=$" + string(rune('0'+argCount))
		args = append(args, req.DiscountAmount)
	}
	if req.TaxAmount != nil {
		argCount++
		query += ", tax_amount=$" + string(rune('0'+argCount))
		args = append(args, req.TaxAmount)
	}
	if req.Notes != "" {
		argCount++
		query += ", notes=$" + string(rune('0'+argCount))
		args = append(args, req.Notes)
	}
	if req.InvoiceNumber != "" {
		argCount++
		query += ", invoice_number=$" + string(rune('0'+argCount))
		args = append(args, req.InvoiceNumber)
	}
	if req.SaleType != "" {
		argCount++
		query += ", sale_type=$" + string(rune('0'+argCount))
		args = append(args, req.SaleType)
	}
	if req.DeliveryStatus != "" {
		argCount++
		query += ", delivery_status=$" + string(rune('0'+argCount))
		args = append(args, req.DeliveryStatus)
	}
	if req.DeliveryDate != "" {
		argCount++
		query += ", delivery_date=$" + string(rune('0'+argCount))
		args = append(args, req.DeliveryDate)
	}
	if req.DeliveryAddress != "" {
		argCount++
		query += ", delivery_address=$" + string(rune('0'+argCount))
		args = append(args, req.DeliveryAddress)
	}

	argCount++
	query += " WHERE id=$" + string(rune('0'+argCount))
	args = append(args, id)

	result, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "transaction not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "transaction updated successfully", "id": id})
}
