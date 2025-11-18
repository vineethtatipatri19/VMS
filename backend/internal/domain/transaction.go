package domain

import "time"

// Transaction represents a sale or payment transaction
type Transaction struct {
	ID              string
	CustomerID      string
	Date            time.Time
	Type            string // sale, payment
	InvoiceNumber   string
	PaymentMethod   string
	PaymentAmount   float64
	PaymentRef      string
	DueDate         *time.Time
	IsOverdue       bool
	Status          string // pending, completed, cancelled
	TotalAmount     float64
	DiscountAmount  float64
	TaxAmount       float64
	BalanceAfter    float64
	SaleType        string
	ReceiptSent     bool
	DeliveryStatus  string
	DeliveryDate    *time.Time
	DeliveryAddress string
	Notes           string
	Details         map[string]interface{}
	CreatedAt       time.Time
	DeletedAt       *time.Time
	DeletedBy       string
	DeletionReason  string
}

// SaleItem represents a line item in a sale transaction
type SaleItem struct {
	ID              string
	TransactionID   string
	InventoryLotID  string
	ItemName        string
	BatchNumber     string
	ExpiryDate      string
	Quantity        float64
	Unit            string
	PricePerUnit    float64
	CostPerUnit     float64
	Profit          float64
	DiscountPercent float64
	TaxPercent      float64
	HSNCode         string
	Total           float64
	DeletedAt       *time.Time
	DeletedBy       string
	DeletionReason  string
}

// TransactionFilters represents query filters for transactions
type TransactionFilters struct {
	CustomerID string
	Type       string
	StartDate  string
	EndDate    string
	IsOverdue  *bool
	Limit      int
	Offset     int
}

// Validate validates transaction business rules
func (t *Transaction) Validate() error {
	if t.CustomerID == "" {
		return ErrInvalidInput("customer ID is required")
	}
	if t.Type == "" {
		return ErrInvalidInput("transaction type is required")
	}
	if t.Type != "sale" && t.Type != "payment" {
		return ErrInvalidInput("invalid transaction type")
	}
	if t.Type == "sale" && t.TotalAmount < 0 {
		return ErrInvalidInput("total amount cannot be negative")
	}
	if t.Type == "payment" && t.PaymentAmount <= 0 {
		return ErrInvalidInput("payment amount must be positive")
	}
	return nil
}

// CalculateProfit calculates total profit from sale items
func (t *Transaction) CalculateProfit(items []SaleItem) float64 {
	var totalProfit float64
	for _, item := range items {
		totalProfit += item.Profit
	}
	return totalProfit
}

// Validate validates sale item business rules
func (si *SaleItem) Validate() error {
	if si.ItemName == "" {
		return ErrInvalidInput("item name is required")
	}
	if si.Quantity <= 0 {
		return ErrInvalidInput("quantity must be positive")
	}
	if si.PricePerUnit < 0 {
		return ErrInvalidInput("price cannot be negative")
	}
	return nil
}

// CalculateProfit calculates profit for this sale item
func (si *SaleItem) CalculateProfit() {
	si.Profit = (si.PricePerUnit - si.CostPerUnit) * si.Quantity
}

// CalculateTotal calculates total amount for this line item
func (si *SaleItem) CalculateTotal() {
	subtotal := si.PricePerUnit * si.Quantity
	discount := subtotal * (si.DiscountPercent / 100)
	afterDiscount := subtotal - discount
	tax := afterDiscount * (si.TaxPercent / 100)
	si.Total = afterDiscount + tax
}
