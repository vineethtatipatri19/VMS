package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
	"github.com/lib/pq"
)

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (
			id, name, email, address, contact_number, alternate_contact, whatsapp_number,
			photo_url, business_name, gstin, customer_type, aadhaar_verified,
			kyc_document_type, kyc_document_number, credit_limit, current_balance,
			payment_terms_days, interest_rate, status, notes, tags,
			last_transaction_date, total_purchases, loyalty_points,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, toNullString(customer.Email),
		toNullString(customer.Address), toNullString(customer.ContactNumber),
		toNullString(customer.AlternateContact), toNullString(customer.WhatsappNumber),
		toNullString(customer.PhotoURL), toNullString(customer.BusinessName),
		toNullString(customer.GSTIN), toNullString(customer.CustomerType),
		customer.AadhaarVerified, toNullString(customer.KYCDocumentType),
		toNullString(customer.KYCDocumentNumber), customer.CreditLimit,
		customer.CurrentBalance, customer.PaymentTermsDays, customer.InterestRate,
		toNullString(customer.Status), toNullString(customer.Notes), pq.Array(customer.Tags),
		toNullTime(customer.LastTransactionDate), customer.TotalPurchases,
		customer.LoyaltyPoints, customer.CreatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := `
		SELECT id, name, email, address, contact_number, alternate_contact, whatsapp_number,
		       photo_url, business_name, gstin, customer_type, aadhaar_verified, 
		       kyc_document_type, kyc_document_number, credit_limit, current_balance, 
		       payment_terms_days, interest_rate, status, notes, tags, 
		       last_transaction_date, total_purchases, loyalty_points, created_at, updated_at
		FROM customers 
		WHERE id = $1 AND deleted_at IS NULL`

	var c domain.Customer
	var email, address, contactNumber, alternateContact, whatsappNumber, photoURL sql.NullString
	var businessName, gstin, customerType, kycDocType, kycDocNumber, status, notes sql.NullString
	var lastTransDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &email, &address, &contactNumber, &alternateContact,
		&whatsappNumber, &photoURL, &businessName, &gstin, &customerType,
		&c.AadhaarVerified, &kycDocType, &kycDocNumber, &c.CreditLimit,
		&c.CurrentBalance, &c.PaymentTermsDays, &c.InterestRate, &status,
		&notes, pq.Array(&c.Tags), &lastTransDate, &c.TotalPurchases,
		&c.LoyaltyPoints, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Map nullable values
	c.Email = fromNullString(email)
	c.Address = fromNullString(address)
	c.ContactNumber = fromNullString(contactNumber)
	c.AlternateContact = fromNullString(alternateContact)
	c.WhatsappNumber = fromNullString(whatsappNumber)
	c.PhotoURL = fromNullString(photoURL)
	c.BusinessName = fromNullString(businessName)
	c.GSTIN = fromNullString(gstin)
	c.CustomerType = fromNullString(customerType)
	c.KYCDocumentType = fromNullString(kycDocType)
	c.KYCDocumentNumber = fromNullString(kycDocNumber)
	c.Status = fromNullString(status)
	c.Notes = fromNullString(notes)
	c.LastTransactionDate = fromNullTime(lastTransDate)

	return &c, nil
}

func (r *customerRepository) List(ctx context.Context) ([]*domain.Customer, error) {
	query := `
		SELECT id, name, email, address, contact_number, alternate_contact, whatsapp_number,
		       photo_url, business_name, gstin, customer_type, aadhaar_verified, 
		       kyc_document_type, kyc_document_number, credit_limit, current_balance, 
		       payment_terms_days, interest_rate, status, notes, tags, 
		       last_transaction_date, total_purchases, loyalty_points, created_at, updated_at
		FROM customers 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	customers := []*domain.Customer{}
	for rows.Next() {
		var c domain.Customer
		var email, address, contactNumber, alternateContact, whatsappNumber, photoURL sql.NullString
		var businessName, gstin, customerType, kycDocType, kycDocNumber, status, notes sql.NullString
		var lastTransDate sql.NullTime

		if err := rows.Scan(
			&c.ID, &c.Name, &email, &address, &contactNumber, &alternateContact,
			&whatsappNumber, &photoURL, &businessName, &gstin, &customerType,
			&c.AadhaarVerified, &kycDocType, &kycDocNumber, &c.CreditLimit,
			&c.CurrentBalance, &c.PaymentTermsDays, &c.InterestRate, &status,
			&notes, pq.Array(&c.Tags), &lastTransDate, &c.TotalPurchases,
			&c.LoyaltyPoints, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		// Map nullable values
		c.Email = fromNullString(email)
		c.Address = fromNullString(address)
		c.ContactNumber = fromNullString(contactNumber)
		c.AlternateContact = fromNullString(alternateContact)
		c.WhatsappNumber = fromNullString(whatsappNumber)
		c.PhotoURL = fromNullString(photoURL)
		c.BusinessName = fromNullString(businessName)
		c.GSTIN = fromNullString(gstin)
		c.CustomerType = fromNullString(customerType)
		c.KYCDocumentType = fromNullString(kycDocType)
		c.KYCDocumentNumber = fromNullString(kycDocNumber)
		c.Status = fromNullString(status)
		c.Notes = fromNullString(notes)
		c.LastTransactionDate = fromNullTime(lastTransDate)

		customers = append(customers, &c)
	}

	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers SET
			name = $2, email = $3, address = $4, contact_number = $5, alternate_contact = $6,
			whatsapp_number = $7, photo_url = $8, business_name = $9, gstin = $10,
			customer_type = $11, aadhaar_verified = $12, kyc_document_type = $13,
			kyc_document_number = $14, credit_limit = $15, current_balance = $16,
			payment_terms_days = $17, interest_rate = $18, status = $19, notes = $20,
			tags = $21, last_transaction_date = $22, total_purchases = $23,
			loyalty_points = $24, updated_at = $25, updated_by = $26
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, toNullString(customer.Email), toNullString(customer.Address),
		toNullString(customer.ContactNumber), toNullString(customer.AlternateContact),
		toNullString(customer.WhatsappNumber), toNullString(customer.PhotoURL),
		toNullString(customer.BusinessName), toNullString(customer.GSTIN),
		customer.CustomerType, customer.AadhaarVerified,
		toNullString(customer.KYCDocumentType), toNullString(customer.KYCDocumentNumber),
		customer.CreditLimit, customer.CurrentBalance,
		customer.PaymentTermsDays, customer.InterestRate, customer.Status,
		toNullString(customer.Notes), pq.Array(customer.Tags),
		toNullTime(customer.LastTransactionDate), customer.TotalPurchases, customer.LoyaltyPoints,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `
		UPDATE customers 
		SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) GetBalance(ctx context.Context, customerID string) (float64, error) {
	var balance float64
	query := `SELECT current_balance FROM customers WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, domain.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get customer balance: %w", err)
	}

	return balance, nil
}

func (r *customerRepository) UpdateBalance(ctx context.Context, customerID string, amount float64) error {
	query := `
		UPDATE customers 
		SET current_balance = current_balance + $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, customerID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update customer balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *customerRepository) UpdateLastTransaction(ctx context.Context, customerID string, date time.Time) error {
	query := `
		UPDATE customers 
		SET last_transaction_date = $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, customerID, date, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
