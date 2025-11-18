package testutil

import (
	"time"

	"github.com/example/pgvms/internal/domain"
)

// FixtureCustomer returns a test customer
func FixtureCustomer(overrides ...func(*domain.Customer)) *domain.Customer {
	customer := &domain.Customer{
		ID:             "cust-test-001",
		Name:           "Test Customer",
		ContactNumber:  "1234567890",
		Address:        "123 Test Street",
		CustomerType:   "wholesale",
		CreditLimit:    10000.0,
		CurrentBalance: 0.0,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	for _, override := range overrides {
		override(customer)
	}

	return customer
}

// FixtureInventoryItem returns a test inventory item
func FixtureInventoryItem(overrides ...func(*domain.InventoryItem)) *domain.InventoryItem {
	item := &domain.InventoryItem{
		ID:            "inv-test-001",
		Name:          "Test Product",
		Category:      "test-category",
		Unit:          "kg",
		Quantity:      100.0,
		MinStockLevel: 10.0,
		CostPrice:     50.0,
		SellingPrice:  75.0,
		PurchaseDate:  time.Now().Format("2006-01-02"),
		ExpiryDate:    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		LotNumber:     "BATCH-001",
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	for _, override := range overrides {
		override(item)
	}

	return item
}

// FixtureTransaction returns a test transaction
func FixtureTransaction(overrides ...func(*domain.Transaction)) *domain.Transaction {
	txn := &domain.Transaction{
		ID:            "txn-test-001",
		CustomerID:    "cust-test-001",
		Type:          "sale",
		TotalAmount:   500.0,
		PaymentAmount: 500.0,
		PaymentMethod: "cash",
		Notes:         "Test transaction",
		Date:          time.Now(),
		CreatedAt:     time.Now(),
	}

	for _, override := range overrides {
		override(txn)
	}

	return txn
}

// FixtureSaleItem returns a test sale item
func FixtureSaleItem(overrides ...func(*domain.SaleItem)) *domain.SaleItem {
	item := &domain.SaleItem{
		ID:             "sale-item-001",
		TransactionID:  "txn-test-001",
		InventoryLotID: "inv-test-001",
		ItemName:       "Test Item",
		Quantity:       5.0,
		PricePerUnit:   75.0,
		CostPerUnit:    50.0,
		Profit:         125.0,
		Total:          375.0,
	}

	for _, override := range overrides {
		override(item)
	}

	return item
}

// FixtureCrateEntry returns a test crate entry
func FixtureCrateEntry(overrides ...func(*domain.CrateEntry)) *domain.CrateEntry {
	entry := &domain.CrateEntry{
		ID:              "crate-001",
		CustomerID:      "cust-test-001",
		TransactionType: "out",
		Quantity:        10,
		Notes:           "Test crate issue",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	for _, override := range overrides {
		override(entry)
	}

	return entry
}

// FixtureWastageLog returns a test wastage log
func FixtureWastageLog(overrides ...func(*domain.WastageLog)) *domain.WastageLog {
	log := &domain.WastageLog{
		ID:          "wastage-001",
		InventoryID: "inv-test-001",
		Quantity:    2.0,
		Reason:      "spoiled",
		CostValue:   100.0,
		RecordedAt:  time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, override := range overrides {
		override(log)
	}

	return log
}

// FixtureExpiryAlert returns a test expiry alert
func FixtureExpiryAlert(overrides ...func(*domain.ExpiryAlert)) *domain.ExpiryAlert {
	alert := &domain.ExpiryAlert{
		ID:              "alert-001",
		InventoryItemID: "inv-test-001",
		AlertDate:       time.Now(),
		ExpiryDate:      time.Now().AddDate(0, 0, 5),
		DaysUntilExpiry: 5,
		Acknowledged:    false,
		CreatedAt:       time.Now(),
	}

	for _, override := range overrides {
		override(alert)
	}

	return alert
}

// FixturePaymentSchedule returns a test payment schedule
func FixturePaymentSchedule(overrides ...func(*domain.PaymentSchedule)) *domain.PaymentSchedule {
	schedule := &domain.PaymentSchedule{
		ID:         "payment-001",
		CustomerID: "cust-test-001",
		AmountDue:  1000.0,
		DueDate:    time.Now().AddDate(0, 0, 30),
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	for _, override := range overrides {
		override(schedule)
	}

	return schedule
}
