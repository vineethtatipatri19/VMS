package domain

import "time"

// Customer represents a customer entity with complete business fields
type Customer struct {
	ID                  string
	Name                string
	Email               string
	Address             string
	ContactNumber       string
	AlternateContact    string
	WhatsappNumber      string
	PhotoURL            string
	BusinessName        string
	GSTIN               string
	CustomerType        string // b2b, b2c, retail, wholesale
	AadhaarVerified     bool
	KYCDocumentType     string
	KYCDocumentNumber   string
	CreditLimit         float64
	CurrentBalance      float64
	PaymentTermsDays    int
	InterestRate        float64
	Status              string // active, inactive, blocked
	Notes               string
	Tags                []string
	LastTransactionDate *time.Time
	TotalPurchases      float64
	LoyaltyPoints       int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
	DeletedBy           string
	DeletionReason      string
}

// CustomerFilters represents query filters for customers
type CustomerFilters struct {
	Status       string
	CustomerType string
	SearchTerm   string
	Limit        int
	Offset       int
}

// Validate validates customer business rules
func (c *Customer) Validate() error {
	if c.Name == "" {
		return ErrInvalidInput("name is required")
	}
	if c.CustomerType != "" && !isValidCustomerType(c.CustomerType) {
		return ErrInvalidInput("invalid customer type")
	}
	if c.Status != "" && !isValidStatus(c.Status) {
		return ErrInvalidInput("invalid status")
	}
	if c.CreditLimit < 0 {
		return ErrInvalidInput("credit limit cannot be negative")
	}
	return nil
}

// CanPurchase checks if customer can make a purchase
func (c *Customer) CanPurchase(amount float64) bool {
	if c.Status != "active" {
		return false
	}
	return c.CurrentBalance+amount <= c.CreditLimit
}

// IsOverdue checks if customer has overdue payments
// This would need transaction data - placeholder logic
func (c *Customer) IsOverdue() bool {
	return c.CurrentBalance > 0
}

func isValidCustomerType(t string) bool {
	validTypes := map[string]bool{
		"b2b":       true,
		"b2c":       true,
		"retail":    true,
		"wholesale": true,
	}
	return validTypes[t]
}

func isValidStatus(s string) bool {
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"blocked":  true,
	}
	return validStatuses[s]
}
