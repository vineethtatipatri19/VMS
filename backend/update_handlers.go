package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

// Update customer handler
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	query := "UPDATE customers SET updated_at=NOW(), updated_by=$1"
	args := []interface{}{userID}
	argCount := 1

	// Build dynamic UPDATE query based on provided fields
	fieldMap := map[string]string{
		"name": "name", "phone": "phone", "email": "email", "address": "address",
		"city": "city", "state": "state", "pincode": "pincode", "gstNumber": "gst_number",
		"panNumber": "pan_number", "creditLimit": "credit_limit", "paymentTerms": "payment_terms",
		"status": "status", "crateLimit": "crate_limit", "businessType": "business_type",
		"contactPerson": "contact_person", "alternatePhone": "alternate_phone",
		"deliveryAddress": "delivery_address", "notes": "notes",
	}

	for jsonField, dbField := range fieldMap {
		if val, exists := req[jsonField]; exists {
			argCount++
			query += ", " + dbField + "=$" + string(rune('0'+argCount))
			args = append(args, val)
		}
	}

	// Handle tags array separately
	if tags, exists := req["tags"]; exists {
		if tagArray, ok := tags.([]interface{}); ok {
			strTags := make([]string, len(tagArray))
			for i, tag := range tagArray {
				strTags[i] = tag.(string)
			}
			argCount++
			query += ", tags=$" + string(rune('0'+argCount))
			args = append(args, pq.Array(strTags))
		}
	}

	argCount++
	query += " WHERE id=$" + string(rune('0'+argCount)) + " AND deleted_at IS NULL"
	args = append(args, id)

	result, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "customer not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "customer updated successfully", "id": id})
}

// Update crate entry handler
func updateCrateEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	query := "UPDATE crate_ledger SET updated_at=NOW(), updated_by=$1"
	args := []interface{}{userID}
	argCount := 1

	fieldMap := map[string]string{
		"quantity": "quantity", "type": "type", "notes": "notes",
	}

	for jsonField, dbField := range fieldMap {
		if val, exists := req[jsonField]; exists {
			argCount++
			query += ", " + dbField + "=$" + string(rune('0'+argCount))
			args = append(args, val)
		}
	}

	argCount++
	query += " WHERE id=$" + string(rune('0'+argCount)) + " AND deleted_at IS NULL"
	args = append(args, id)

	result, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "crate entry not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "crate entry updated successfully", "id": id})
}

// Update wastage log handler
func updateWastage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	query := "UPDATE wastage_log SET updated_at=NOW(), updated_by=$1"
	args := []interface{}{userID}
	argCount := 1

	fieldMap := map[string]string{
		"quantity": "quantity", "reason": "reason", "reasonDetails": "reason_details",
		"costValue": "cost_value", "notes": "notes",
	}

	for jsonField, dbField := range fieldMap {
		if val, exists := req[jsonField]; exists {
			argCount++
			query += ", " + dbField + "=$" + string(rune('0'+argCount))
			args = append(args, val)
		}
	}

	argCount++
	query += " WHERE id=$" + string(rune('0'+argCount)) + " AND deleted_at IS NULL"
	args = append(args, id)

	result, err := db.ExecContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "wastage entry not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "wastage entry updated successfully", "id": id})
}
