package domain

import "time"

// Crate represents a crate transaction (issue/return)
type Crate struct {
	ID             string     `json:"id"`
	CustomerID     string     `json:"customer_id"`
	CustomerName   string     `json:"customer_name,omitempty"`
	TransactionID  string     `json:"transaction_id,omitempty"`
	Date           time.Time  `json:"date"`
	CratesIssued   int        `json:"crates_issued"`
	CratesReturned int        `json:"crates_returned"`
	Balance        int        `json:"balance"`
	Notes          string     `json:"notes,omitempty"`
	CrateType      string     `json:"crate_type,omitempty"`
	CrateValue     float64    `json:"crate_value,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      string     `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      string     `json:"deleted_by,omitempty"`
	DeletionReason string     `json:"deletion_reason,omitempty"`
}

// CrateEntry represents a crate issue/return entry (legacy alias)
type CrateEntry = Crate

// Wastage represents a wastage entry
type Wastage struct {
	ID          string
	InventoryID string
	ItemName    string
	LotNumber   string
	Quantity    float64
	Unit        string
	Reason      string
	CostValue   float64
	ExpiryDate  time.Time
	Notes       string
	RecordedBy  string
	RecordedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WastageLog represents a wastage entry (legacy alias)
type WastageLog = Wastage

// ExpiryAlert represents an expiry alert
type ExpiryAlert struct {
	ID              string
	InventoryItemID string
	AlertDate       time.Time
	ExpiryDate      time.Time
	DaysUntilExpiry int
	Acknowledged    bool
	AcknowledgedAt  *time.Time
	AcknowledgedBy  string
	CreatedAt       time.Time
	DeletedAt       *time.Time
	DeletedBy       string
	DeletionReason  string
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
