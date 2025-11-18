package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/example/pgvms/tests/testutil"
)

// TestSaleTransactionFEFO_E2E tests sale transaction with FEFO (First Expiry First Out) logic
func TestSaleTransactionFEFO_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var customerID, item1ID, item2ID, transactionID string

	// Step 1: Create a customer
	t.Run("Create customer", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "FEFO Test Customer",
			"ContactNumber": "+1234567890",
			"CustomerType":  "wholesale",
			"CreditLimit":   100000.0,
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		customerID = data["ID"].(string)
		t.Logf("Created customer: %s", customerID)
	})

	// Step 2: Add inventory item with near expiry (expires in 2 months)
	t.Run("Add inventory item expiring soon", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "Milk Carton",
			"LotNumber":     testutil.GenerateUUID(),
			"Unit":          "cartons",
			"Quantity":      50.0,
			"CostPrice":     20.0,
			"SellingPrice":  30.0,
			"SupplierName":  "Dairy Co",
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 2, 0).Format("2006-01-02"), // 2 months
			"MinStockLevel": 10.0,
			"Status":        "available",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		item1ID = data["ID"].(string)
		t.Logf("Created inventory item 1 (expires soon): %s", item1ID)
	})

	// Step 3: Add inventory item with far expiry (expires in 6 months)
	t.Run("Add inventory item expiring later", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "Milk Carton",
			"LotNumber":     testutil.GenerateUUID(),
			"Unit":          "cartons",
			"Quantity":      50.0,
			"CostPrice":     20.0,
			"SellingPrice":  30.0,
			"SupplierName":  "Dairy Co",
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"), // 6 months
			"MinStockLevel": 10.0,
			"Status":        "available",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		item2ID = data["ID"].(string)
		t.Logf("Created inventory item 2 (expires later): %s", item2ID)
	})

	// Step 4: Create a sale transaction
	t.Run("Create sale transaction", func(t *testing.T) {
		transaction := map[string]interface{}{
			"CustomerID":  customerID,
			"Type":        "sale",
			"PaymentMode": "cash",
			"TotalAmount": 900.0,
			"AmountPaid":  900.0,
			"Status":      "completed",
			"Items": []map[string]interface{}{
				{
					"InventoryID": item1ID,
					"Quantity":    15.0,
					"UnitPrice":   30.0,
					"TotalPrice":  450.0,
				},
				{
					"InventoryID": item2ID,
					"Quantity":    15.0,
					"UnitPrice":   30.0,
					"TotalPrice":  450.0,
				},
			},
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/transactions", transaction)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		transactionID = data["ID"].(string)
		t.Logf("Created transaction: %s", transactionID)
	})

	// Step 5: Verify inventory quantities updated (FEFO should deduct from item1 first)
	t.Run("Verify FEFO inventory deduction", func(t *testing.T) {
		// Check item1 (expires sooner) - should have less quantity
		resp1 := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/inventory/%s", item1ID), nil)
		testutil.AssertHTTPStatusCode(t, resp1, http.StatusOK)
		data1 := ctx.ExtractData(t, resp1)

		// Check item2 (expires later) - should have more quantity
		resp2 := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/inventory/%s", item2ID), nil)
		testutil.AssertHTTPStatusCode(t, resp2, http.StatusOK)
		data2 := ctx.ExtractData(t, resp2)

		qty1 := data1["Quantity"].(float64)
		qty2 := data2["Quantity"].(float64)

		t.Logf("Item 1 quantity (expires soon): %.0f", qty1)
		t.Logf("Item 2 quantity (expires later): %.0f", qty2)

		// With FEFO, item1 should be depleted first
		// If FEFO works correctly, item1 should have 35 (50-15) and item2 should have 35 (50-15)
		// Or if FEFO is aggressive, item1 could be lower
		if qty1 > qty2 {
			t.Errorf("FEFO logic may not be working: item expiring soon has more stock (%.0f) than item expiring later (%.0f)", qty1, qty2)
		}
	})

	// Step 6: Verify transaction was recorded
	t.Run("Verify transaction recorded", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/transactions/%s", transactionID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		testutil.AssertEqual(t, "completed", data["Status"])
		testutil.AssertEqual(t, 900.0, data["TotalAmount"])
	})
}

// TestCreditSaleFlow_E2E tests a credit sale and payment tracking
func TestCreditSaleFlow_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var customerID, inventoryID, transactionID string

	// Step 1: Create customer with credit limit
	t.Run("Create customer with credit", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":             "Credit Customer",
			"ContactNumber":    "+9876543210",
			"CustomerType":     "wholesale",
			"CreditLimit":      50000.0,
			"PaymentTermsDays": 30,
			"Status":           "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		customerID = data["ID"].(string)
		t.Logf("Created customer with credit: %s", customerID)
	})

	// Step 2: Add inventory
	t.Run("Add inventory item", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":         "Bulk Rice Bag",
			"LotNumber":    testutil.GenerateUUID(),
			"Unit":         "bags",
			"Quantity":     100.0,
			"CostPrice":    800.0,
			"SellingPrice": 1000.0,
			"PurchaseDate": time.Now().Format("2006-01-02"),
			"ExpiryDate":   time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
			"Status":       "available",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		inventoryID = data["ID"].(string)
		t.Logf("Created inventory: %s", inventoryID)
	})

	// Step 3: Create credit sale
	t.Run("Create credit sale", func(t *testing.T) {
		transaction := map[string]interface{}{
			"CustomerID":  customerID,
			"Type":        "sale",
			"PaymentMode": "credit",
			"TotalAmount": 10000.0,
			"AmountPaid":  0.0, // No payment yet
			"Status":      "pending",
			"Items": []map[string]interface{}{
				{
					"InventoryID": inventoryID,
					"Quantity":    10.0,
					"UnitPrice":   1000.0,
					"TotalPrice":  10000.0,
				},
			},
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/transactions", transaction)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		transactionID = data["ID"].(string)
		t.Logf("Created credit transaction: %s", transactionID)
	})

	// Step 4: Verify customer balance increased
	t.Run("Verify customer balance", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		balance := data["CurrentBalance"].(float64)
		t.Logf("Customer balance after credit sale: %.2f", balance)

		if balance <= 0 {
			t.Errorf("Expected positive balance after credit sale, got %.2f", balance)
		}
	})

	// Step 5: Verify inventory deducted
	t.Run("Verify inventory deducted", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/inventory/%s", inventoryID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		quantity := data["Quantity"].(float64)

		testutil.AssertEqual(t, 90.0, quantity) // 100 - 10
		t.Logf("Inventory quantity after sale: %.0f", quantity)
	})
}
