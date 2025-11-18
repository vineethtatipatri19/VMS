package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/example/pgvms/tests/testutil"
)

// TestCompleteCustomerFlow_E2E tests the complete customer lifecycle
func TestCompleteCustomerFlow_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var customerID string

	// Step 1: Create Customer
	t.Run("Create customer", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "E2E Customer",
			"Email":         "e2e@example.com",
			"ContactNumber": "+1234567890",
			"Address":       "123 Test St",
			"CustomerType":  "wholesale",
			"CreditLimit":   50000.0,
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		customerID = data["ID"].(string)
		t.Logf("Created customer: %s", customerID)

		testutil.AssertEqual(t, "E2E Customer", data["Name"])
		testutil.AssertEqual(t, "wholesale", data["CustomerType"])
	})

	// Step 2: Get Customer
	t.Run("Get customer by ID", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		testutil.AssertEqual(t, customerID, data["ID"])
		testutil.AssertEqual(t, "E2E Customer", data["Name"])
	})

	// Step 3: Update Customer
	t.Run("Update customer", func(t *testing.T) {
		// First get the existing customer
		getResp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, getResp, http.StatusOK)
		customer := ctx.ExtractData(t, getResp)

		// Modify the fields
		customer["Name"] = "Updated E2E Customer"
		customer["CreditLimit"] = 75000.0

		// Send the update
		resp := ctx.MakeRequest(t, "PUT", fmt.Sprintf("/api/v1/customers/%s", customerID), customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		testutil.AssertEqual(t, "Updated E2E Customer", data["Name"])
		testutil.AssertEqual(t, 75000.0, data["CreditLimit"])
	})

	// Step 4: Delete Customer (soft delete)
	t.Run("Delete customer", func(t *testing.T) {
		deleteReq := map[string]interface{}{
			"reason":      "E2E test cleanup",
			"attestation": "test",
		}

		resp := ctx.MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/customers/%s", customerID), deleteReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)
	})

	// Step 5: Verify Deleted
	t.Run("Verify customer deleted", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusNotFound)
	})
}

// TestCompleteInventoryFlow_E2E tests the complete inventory lifecycle
func TestCompleteInventoryFlow_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var inventoryID string

	// Step 1: Add Inventory Item
	t.Run("Add inventory item", func(t *testing.T) {
		item := map[string]interface{}{
			"Name": "Test Product",
			// "Category":      "beverages",
			"Unit":          "bottles",
			"Quantity":      100.0,
			"CostPrice":     50.0,
			"SellingPrice":  75.0,
			"SupplierName":  "Test Supplier",
			"LotNumber":     testutil.GenerateUUID(),
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"MinStockLevel": 20.0,
			// "ReorderPoint":  30.0,
			"Status": "available",
		}

		t.Logf("Request payload: %+v", item)
		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		t.Logf("Response status: %d", resp.StatusCode)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		inventoryID = data["ID"].(string)
		t.Logf("Created inventory item: %s", inventoryID)

		testutil.AssertEqual(t, "Test Product", data["Name"])
		testutil.AssertEqual(t, 100.0, data["Quantity"])
	})

	// Step 2: Get Inventory Item
	t.Run("Get inventory by ID", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/inventory/%s", inventoryID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		testutil.AssertEqual(t, inventoryID, data["ID"])
		testutil.AssertEqual(t, "Test Product", data["Name"])
	})

	// Step 3: Deduct Stock
	t.Run("Deduct stock", func(t *testing.T) {
		deductReq := map[string]interface{}{
			"quantity": 25.0,
			"reason":   "Sale",
		}

		resp := ctx.MakeRequest(t, "POST", fmt.Sprintf("/api/v1/inventory/%s/deduct", inventoryID), deductReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		testutil.AssertEqual(t, 75.0, data["Quantity"])
	})

	// Step 4: Check Low Stock
	t.Run("Check low stock items", func(t *testing.T) {
		// Deduct more to go below reorder point (30)
		deductReq := map[string]interface{}{
			"quantity": 50.0,
			"reason":   "Sale",
		}
		resp := ctx.MakeRequest(t, "POST", fmt.Sprintf("/api/v1/inventory/%s/deduct", inventoryID), deductReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		// Now check low stock
		resp = ctx.MakeRequest(t, "GET", "/api/v1/inventory/low-stock", nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		items, ok := data["items"].([]interface{})
		if !ok || len(items) == 0 {
			t.Fatal("Expected low stock items")
		}
		t.Logf("Found %d low stock items", len(items))
	})

	// Step 5: Delete Inventory Item
	t.Run("Delete inventory item", func(t *testing.T) {
		deleteReq := map[string]interface{}{
			"reason":      "E2E test cleanup",
			"attestation": "test",
		}

		resp := ctx.MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/inventory/%s", inventoryID), deleteReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)
	})
}
