package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Customer represents a customer in the system
type Customer struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Address         string    `json:"address,omitempty"`
	ContactNumber   string    `json:"contactNumber,omitempty"`
	PhotoURL        string    `json:"photoUrl,omitempty"`
	AadhaarVerified bool      `json:"aadhaarVerified"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Handler for listing customers
func listCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT id, name, address, contact_number, photo_url, aadhaar_verified, created_at, updated_at 
		FROM customers 
		ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	
	customers := []Customer{}
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.ContactNumber, &c.PhotoURL, &c.AadhaarVerified, &c.CreatedAt, &c.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		customers = append(customers, c)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// Handler for getting a single customer
func getCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c Customer
	err := db.QueryRowContext(r.Context(), `
		SELECT id, name, address, contact_number, photo_url, aadhaar_verified, created_at, updated_at 
		FROM customers WHERE id=$1`, id).Scan(&c.ID, &c.Name, &c.Address, &c.ContactNumber, &c.PhotoURL, &c.AadhaarVerified, &c.CreatedAt, &c.UpdatedAt)
	
	if err == sql.ErrNoRows {
		http.Error(w, "customer not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// Handler for creating a customer
func createCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	if c.Name == "" {
		http.Error(w, "name is required", 400)
		return
	}
	
	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO customers (id, name, address, contact_number, photo_url, aadhaar_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.ID, c.Name, c.Address, c.ContactNumber, c.PhotoURL, c.AadhaarVerified, c.CreatedAt, c.UpdatedAt)
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// Handler for updating a customer
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	c.UpdatedAt = time.Now()
	
	result, err := db.ExecContext(r.Context(), `
		UPDATE customers 
		SET name=$1, address=$2, contact_number=$3, photo_url=$4, aadhaar_verified=$5, updated_at=$6
		WHERE id=$7`,
		c.Name, c.Address, c.ContactNumber, c.PhotoURL, c.AadhaarVerified, c.UpdatedAt, id)
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "customer not found", 404)
		return
	}
	
	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// Handler for deleting a customer
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	result, err := db.ExecContext(r.Context(), `DELETE FROM customers WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "customer not found", 404)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}