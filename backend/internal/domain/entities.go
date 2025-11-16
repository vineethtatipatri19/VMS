package domain

import "time"

// CrateEntry represents a crate issue/return entry
type CrateEntry struct {
	ID             string
	CustomerID     string
	Date           time.Time
	TransactionID  string
	CratesIssued   int
	CratesReturned int
	Balance        int
	CrateType      string
	CrateValue     float64
	Notes          string
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	DeletedBy      string
	DeletionReason string
}

// Validate validates crate entry business rules
func (c *CrateEntry) Validate() error {
	if c.CustomerID == "" {
		return ErrInvalidInput("customer ID is required")
	}
	if c.CratesIssued < 0 || c.CratesReturned < 0 {
		return ErrInvalidInput("crate counts cannot be negative")
	}
	return nil
}

// WastageLog represents a wastage entry
type WastageLog struct {
	ID              string
	InventoryItemID string
	ItemName        string
	Quantity        float64
	Unit            string
	Reason          string // expired, damaged, spoiled, pest, other
	ReasonDetails   string
	CostValue       float64
	LoggedBy        string
	LoggedAt        time.Time
	PhotoURL        string
	DeletedAt       *time.Time
	DeletedBy       string
	DeletionReason  string
}

// Validate validates wastage log business rules
func (w *WastageLog) Validate() error {
	if w.ItemName == "" {
		return ErrInvalidInput("item name is required")
	}
	if w.Quantity <= 0 {
		return ErrInvalidInput("quantity must be positive")
	}
	if w.Reason == "" {
		return ErrInvalidInput("reason is required")
	}
	if !isValidWastageReason(w.Reason) {
		return ErrInvalidInput("invalid wastage reason")
	}
	return nil
}

func isValidWastageReason(reason string) bool {
	validReasons := map[string]bool{
		"expired": true,
		"damaged": true,
		"spoiled": true,
		"pest":    true,
		"other":   true,
	}
	return validReasons[reason]
}

// ExpiryAlert represents an expiry alert
type ExpiryAlert struct {
	ID               string
	InventoryItemID  string
	AlertDate        time.Time
	ExpiryDate       time.Time
	DaysUntilExpiry  int
	Acknowledged     bool
	AcknowledgedAt   *time.Time
	AcknowledgedBy   string
	CreatedAt        time.Time
	DeletedAt        *time.Time
	DeletedBy        string
	DeletionReason   string
}

// Acknowledge marks alert as acknowledged
func (e *ExpiryAlert) Acknowledge(userID string) {
	e.Acknowledged = true
	now := time.Now()
	e.AcknowledgedAt = &now
	e.AcknowledgedBy = userID
}

// PaymentSchedule represents an installment payment
type PaymentSchedule struct {
	ID                string
	TransactionID     string
	CustomerID        string
	InstallmentNumber int
	DueDate           time.Time
	AmountDue         float64
	AmountPaid        float64
	Status            string // pending, partial, paid, overdue
	PaymentDate       *time.Time
	Notes             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Validate validates payment schedule business rules
func (p *PaymentSchedule) Validate() error {
	if p.CustomerID == "" {
		return ErrInvalidInput("customer ID is required")
	}
	if p.AmountDue <= 0 {
		return ErrInvalidInput("amount due must be positive")
	}
	if p.AmountPaid < 0 {
		return ErrInvalidInput("amount paid cannot be negative")
	}
	return nil
}

// UpdateStatus updates payment status based on amounts
func (p *PaymentSchedule) UpdateStatus() {
	if p.AmountPaid >= p.AmountDue {
		p.Status = "paid"
	} else if p.AmountPaid > 0 {
		p.Status = "partial"
	} else if time.Now().After(p.DueDate) {
		p.Status = "overdue"
	} else {
		p.Status = "pending"
	}
}
