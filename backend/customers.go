package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Customer represents a customer in the system with enhanced fields
type Customer struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email,omitempty"`
	Address             string     `json:"address,omitempty"`
	ContactNumber       string     `json:"contactNumber,omitempty"`
	AlternateContact    string     `json:"alternateContact,omitempty"`
	WhatsappNumber      string     `json:"whatsappNumber,omitempty"`
	PhotoURL            string     `json:"photoUrl,omitempty"`
	BusinessName        string     `json:"businessName,omitempty"`
	GSTIN               string     `json:"gstin,omitempty"`
	CustomerType        string     `json:"customerType"` // b2b, b2c, retail, wholesale
	AadhaarVerified     bool       `json:"aadhaarVerified"`
	KYCDocumentType     string     `json:"kycDocumentType,omitempty"`
	KYCDocumentNumber   string     `json:"kycDocumentNumber,omitempty"`
	CreditLimit         float64    `json:"creditLimit"`
	CurrentBalance      float64    `json:"currentBalance"`
	PaymentTermsDays    int        `json:"paymentTermsDays"`
	InterestRate        float64    `json:"interestRate"`
	Status              string     `json:"status"` // active, inactive, blocked
	Notes               string     `json:"notes,omitempty"`
	Tags                []string   `json:"tags,omitempty"`
	LastTransactionDate *time.Time `json:"lastTransactionDate,omitempty"`
	TotalPurchases      float64    `json:"totalPurchases"`
	LoyaltyPoints       int        `json:"loyaltyPoints"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// Handler for listing customers
func listCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT id, name, email, address, contact_number, alternate_contact, whatsapp_number,
		       photo_url, business_name, gstin, customer_type, aadhaar_verified, 
		       kyc_document_type, kyc_document_number, credit_limit, current_balance, 
		       payment_terms_days, interest_rate, status, notes, tags, 
		       last_transaction_date, total_purchases, loyalty_points, created_at, updated_at 
		FROM customers 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var c Customer
		var email, address, contactNumber, alternateContact, whatsappNumber, photoURL sql.NullString
		var businessName, gstin, kycDocType, kycDocNumber, notes sql.NullString
		var lastTransDate sql.NullTime

		if err := rows.Scan(&c.ID, &c.Name, &email, &address, &contactNumber, &alternateContact,
			&whatsappNumber, &photoURL, &businessName, &gstin, &c.CustomerType, &c.AadhaarVerified,
			&kycDocType, &kycDocNumber, &c.CreditLimit, &c.CurrentBalance,
			&c.PaymentTermsDays, &c.InterestRate, &c.Status, &notes, pq.Array(&c.Tags),
			&lastTransDate, &c.TotalPurchases, &c.LoyaltyPoints, &c.CreatedAt, &c.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Map nullable values
		if email.Valid {
			c.Email = email.String
		}
		if address.Valid {
			c.Address = address.String
		}
		if contactNumber.Valid {
			c.ContactNumber = contactNumber.String
		}
		if alternateContact.Valid {
			c.AlternateContact = alternateContact.String
		}
		if whatsappNumber.Valid {
			c.WhatsappNumber = whatsappNumber.String
		}
		if photoURL.Valid {
			c.PhotoURL = photoURL.String
		}
		if businessName.Valid {
			c.BusinessName = businessName.String
		}
		if gstin.Valid {
			c.GSTIN = gstin.String
		}
		if kycDocType.Valid {
			c.KYCDocumentType = kycDocType.String
		}
		if kycDocNumber.Valid {
			c.KYCDocumentNumber = kycDocNumber.String
		}
		if notes.Valid {
			c.Notes = notes.String
		}
		if lastTransDate.Valid {
			c.LastTransactionDate = &lastTransDate.Time
		}

		// Initialize empty array if nil
		if c.Tags == nil {
			c.Tags = []string{}
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
	var email, address, contactNumber, alternateContact, whatsappNumber, photoURL sql.NullString
	var businessName, gstin, kycDocType, kycDocNumber, notes sql.NullString
	var lastTransDate sql.NullTime

	err := db.QueryRowContext(r.Context(), `
		SELECT id, name, email, address, contact_number, alternate_contact, whatsapp_number,
		       photo_url, business_name, gstin, customer_type, aadhaar_verified, 
		       kyc_document_type, kyc_document_number, credit_limit, current_balance, 
		       payment_terms_days, interest_rate, status, notes, tags, 
		       last_transaction_date, total_purchases, loyalty_points, created_at, updated_at 
		FROM customers WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&c.ID, &c.Name, &email, &address, &contactNumber,
		&alternateContact, &whatsappNumber, &photoURL, &businessName, &gstin, &c.CustomerType,
		&c.AadhaarVerified, &kycDocType, &kycDocNumber, &c.CreditLimit, &c.CurrentBalance,
		&c.PaymentTermsDays, &c.InterestRate, &c.Status, &notes, pq.Array(&c.Tags),
		&lastTransDate, &c.TotalPurchases, &c.LoyaltyPoints, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "customer not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Map nullable values
	if email.Valid {
		c.Email = email.String
	}
	if address.Valid {
		c.Address = address.String
	}
	if contactNumber.Valid {
		c.ContactNumber = contactNumber.String
	}
	if alternateContact.Valid {
		c.AlternateContact = alternateContact.String
	}
	if whatsappNumber.Valid {
		c.WhatsappNumber = whatsappNumber.String
	}
	if photoURL.Valid {
		c.PhotoURL = photoURL.String
	}
	if businessName.Valid {
		c.BusinessName = businessName.String
	}
	if gstin.Valid {
		c.GSTIN = gstin.String
	}
	if kycDocType.Valid {
		c.KYCDocumentType = kycDocType.String
	}
	if kycDocNumber.Valid {
		c.KYCDocumentNumber = kycDocNumber.String
	}
	if notes.Valid {
		c.Notes = notes.String
	}
	if lastTransDate.Valid {
		c.LastTransactionDate = &lastTransDate.Time
	}

	// Initialize empty array if nil
	if c.Tags == nil {
		c.Tags = []string{}
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
