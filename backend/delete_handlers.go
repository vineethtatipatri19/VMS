package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Soft delete handler for inventory
func deleteInventory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE inventory_items 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "inventory item not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "inventory item deleted successfully"})
}

// Soft delete handler for customers
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE customers 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "customer not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "customer deleted successfully"})
}

// Soft delete handler for transactions
func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE transactions 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "transaction not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "transaction deleted successfully"})
}

// Soft delete handler for crate entries
func deleteCrateEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE crate_ledger 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "crate entry not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "crate entry deleted successfully"})
}

// Soft delete handler for wastage logs
func deleteWastage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE wastage_log 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "wastage entry not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "wastage entry deleted successfully"})
}

// Soft delete handler for expiry alerts
func deleteExpiryAlert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := "system"
	if uid := r.Context().Value("userID"); uid != nil {
		userID = uid.(string)
	}

	var req struct {
		Reason      string `json:"reason"`
		Attestation string `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "deletion reason required", 400)
		return
	}

	if req.Reason == "" {
		http.Error(w, "deletion reason cannot be empty", 400)
		return
	}

	if req.Attestation == "" {
		http.Error(w, "attestation required", 400)
		return
	}

	result, err := db.ExecContext(r.Context(), `
		UPDATE expiry_alerts 
		SET deleted_at=NOW(), deleted_by=$1, deletion_reason=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		userID, req.Reason+" | Attestation: "+req.Attestation, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "expiry alert not found or already deleted", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "expiry alert deleted successfully"})
}
