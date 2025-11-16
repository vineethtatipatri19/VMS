package domain

import "time"

// InventoryItem represents an inventory item/lot with complete tracking
type InventoryItem struct {
	ID               string
	Name             string
	Variant          string
	LotNumber        string
	Quantity         float64
	Unit             string
	PurchaseDate     string
	ExpiryDate       string
	Category         string
	SubCategory      string
	CostPrice        float64
	SellingPrice     float64
	MarginPercentage float64
	SupplierID       string
	SupplierName     string
	PurchaseInvoice  string
	MinStockLevel    float64
	ReorderPoint     float64
	ShelfLifeDays    int
	StorageLocation  string
	Barcode          string
	SKU              string
	HSNCode          string
	GSTRate          float64
	Status           string
	WeightPerUnit    float64
	PackagingType    string
	ImageURL         string
	Notes            string
	TotalSold        float64
	TotalWasted      float64
	LastRestockDate  string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	DeletedBy        string
	DeletionReason   string
}

// InventoryFilters represents query filters for inventory
type InventoryFilters struct {
	Status     string
	Category   string
	SearchTerm string
	SortBy     string
	Limit      int
	Offset     int
}

// Validate validates inventory business rules
func (i *InventoryItem) Validate() error {
	if i.Name == "" {
		return ErrInvalidInput("name is required")
	}
	if i.LotNumber == "" {
		return ErrInvalidInput("lot number is required")
	}
	if i.Quantity < 0 {
		return ErrInvalidInput("quantity cannot be negative")
	}
	if i.Unit == "" {
		return ErrInvalidInput("unit is required")
	}
	if i.CostPrice < 0 || i.SellingPrice < 0 {
		return ErrInvalidInput("prices cannot be negative")
	}
	return nil
}

// CalculateMargin calculates profit margin percentage
func (i *InventoryItem) CalculateMargin() {
	if i.SellingPrice > 0 {
		i.MarginPercentage = ((i.SellingPrice - i.CostPrice) / i.SellingPrice) * 100
	}
}

// IsLowStock checks if item is below minimum stock level
func (i *InventoryItem) IsLowStock() bool {
	return i.MinStockLevel > 0 && i.Quantity <= i.MinStockLevel
}

// IsExpired checks if item has expired
func (i *InventoryItem) IsExpired() bool {
	if i.ExpiryDate == "" {
		return false
	}
	// Parse and check expiry date
	expiry, err := time.Parse("2006-01-02", i.ExpiryDate)
	if err != nil {
		return false
	}
	return time.Now().After(expiry)
}

// DaysUntilExpiry calculates days until expiry
func (i *InventoryItem) DaysUntilExpiry() int {
	if i.ExpiryDate == "" {
		return -1
	}
	expiry, err := time.Parse("2006-01-02", i.ExpiryDate)
	if err != nil {
		return -1
	}
	duration := time.Until(expiry)
	return int(duration.Hours() / 24)
}

// UpdateStatus updates status based on inventory conditions
func (i *InventoryItem) UpdateStatus() {
	if i.IsExpired() {
		i.Status = "expired"
	} else if i.Quantity == 0 {
		i.Status = "out_of_stock"
	} else if i.IsLowStock() {
		i.Status = "low_stock"
	} else {
		i.Status = "available"
	}
}
