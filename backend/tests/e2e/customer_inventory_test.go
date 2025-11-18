package e2e_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/example/pgvms/tests/testutil"
)

func TestCustomerRegistrationFlow_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var customerID string

	t.Run("Create new customer", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "E2E Test Customer",
			"Email":         "e2e@test.com",
			"ContactNumber": "+1234567890",
			"Address":       "123 Test Street",
			"Status":        "active",
			"CustomerType":  "b2c",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		var wrapper map[string]interface{}
		ctx.DecodeResponse(t, resp, &wrapper)

		// Extract data from wrapper
		data, ok := wrapper["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Response does not contain data field. Got: %v", wrapper)
		}

		if id, ok := data["ID"].(string); ok {
			customerID = id
			t.Logf("Created customer with ID: %s", customerID)
		} else {
			t.Fatal("Response data does not contain customer ID")
		}

		testutil.AssertEqual(t, "E2E Test Customer", data["Name"])
		testutil.AssertEqual(t, "e2e@test.com", data["Email"])
	})

	t.Run("Get customer by ID", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var customer map[string]interface{}
		ctx.DecodeResponse(t, resp, &customer)

		// Verify customer data
		testutil.AssertEqual(t, customerID, customer["ID"])
		testutil.AssertEqual(t, "E2E Test Customer", customer["Name"])
		testutil.AssertEqual(t, "e2e@test.com", customer["Email"])
		testutil.AssertEqual(t, 50000.0, customer["CreditLimit"])
		testutil.AssertEqual(t, 0.0, customer["CurrentBalance"])
	})

	t.Run("List customers includes new customer", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", "/api/v1/customers", nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var customers []map[string]interface{}
		ctx.DecodeResponse(t, resp, &customers)

		// Verify our customer is in the list
		found := false
		for _, c := range customers {
			if c["ID"] == customerID {
				found = true
				testutil.AssertEqual(t, "E2E Test Customer", c["Name"])
				break
			}
		}

		if !found {
			t.Error("Created customer not found in customer list")
		}
	})

	t.Run("Update customer information", func(t *testing.T) {
		updates := map[string]interface{}{
			"Name":        "Updated E2E Customer",
			"CreditLimit": 75000.0,
			"Status":      "active",
		}

		resp := ctx.MakeRequest(t, "PUT", fmt.Sprintf("/api/v1/customers/%s", customerID), updates)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var updated map[string]interface{}
		ctx.DecodeResponse(t, resp, &updated)

		testutil.AssertEqual(t, "Updated E2E Customer", updated["Name"])
		testutil.AssertEqual(t, 75000.0, updated["CreditLimit"])
	})

	t.Run("Verify update persisted", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var customer map[string]interface{}
		ctx.DecodeResponse(t, resp, &customer)

		testutil.AssertEqual(t, "Updated E2E Customer", customer["Name"])
		testutil.AssertEqual(t, 75000.0, customer["CreditLimit"])
	})

	t.Run("Delete customer (soft delete)", func(t *testing.T) {
		deleteReq := map[string]interface{}{
			"reason":      "E2E test cleanup",
			"attestation": "Test completed",
		}

		resp := ctx.MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/customers/%s", customerID), deleteReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)
	})

	t.Run("Deleted customer not accessible", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusNotFound)
	})
}

func TestInventoryManagementFlow_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var inventoryID string

	t.Run("Add new inventory item", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "E2E Test Product",
			"category":      "Beverages",
			"Quantity":      100.0,
			"Unit":          "pieces",
			"CostPrice":     50.0,
			"SellingPrice":  75.0,
			"SupplierName":  "Test Supplier",
			"LotNumber":     testutil.GenerateUUID(),
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"), // 6 months from now
			"MinStockLevel": 20.0,
			"ReorderPoint":  30.0,
			"Status":        "available",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		var result map[string]interface{}
		ctx.DecodeResponse(t, resp, &result)

		if id, ok := result["ID"].(string); ok {
			inventoryID = id
			t.Logf("Created inventory item with ID: %s", inventoryID)
		} else {
			t.Fatal("Response does not contain inventory ID")
		}

		testutil.AssertEqual(t, "E2E Test Product", result["Name"])
		testutil.AssertEqual(t, 100.0, result["Quantity"])
	})

	t.Run("Get inventory item by ID", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/inventory/%s", inventoryID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var item map[string]interface{}
		ctx.DecodeResponse(t, resp, &item)

		testutil.AssertEqual(t, inventoryID, item["ID"])
		testutil.AssertEqual(t, "E2E Test Product", item["Name"])
		testutil.AssertEqual(t, 100.0, item["Quantity"])
		testutil.AssertEqual(t, "available", item["Status"])
	})

	t.Run("Update inventory quantity", func(t *testing.T) {
		quantityUpdate := map[string]interface{}{
			"quantity_delta": -10.0, // Reduce by 10
			"reason":         "E2E test sale",
		}

		resp := ctx.MakeRequest(t, "POST", fmt.Sprintf("/api/v1/inventory/%s/quantity", inventoryID), quantityUpdate)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var updated map[string]interface{}
		ctx.DecodeResponse(t, resp, &updated)

		testutil.AssertEqual(t, 90.0, updated["Quantity"])
	})

	t.Run("Check low stock items", func(t *testing.T) {
		// First reduce stock to below reorder point
		quantityUpdate := map[string]interface{}{
			"quantity_delta": -65.0, // Bring to 25, below reorder_point of 30
			"reason":         "E2E low stock test",
		}

		resp := ctx.MakeRequest(t, "POST", fmt.Sprintf("/api/v1/inventory/%s/quantity", inventoryID), quantityUpdate)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		// Now check low stock
		resp = ctx.MakeRequest(t, "GET", "/api/v1/inventory/low-stock", nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var items []map[string]interface{}
		ctx.DecodeResponse(t, resp, &items)

		// Verify our item is in low stock list
		found := false
		for _, item := range items {
			if item["ID"] == inventoryID {
				found = true
				testutil.AssertEqual(t, 25.0, item["Quantity"])
				t.Log("Item correctly identified as low stock")
				break
			}
		}

		if !found {
			t.Error("Low stock item not found in low stock list")
		}
	})

	t.Run("Check expiring soon items", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", "/api/v1/inventory/expiring?days=180", nil) // 6 months
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		var items []map[string]interface{}
		ctx.DecodeResponse(t, resp, &items)

		// Our item should be in this list (expires in 6 months)
		found := false
		for _, item := range items {
			if item["ID"] == inventoryID {
				found = true
				t.Log("Item correctly identified as expiring soon")
				break
			}
		}

		if !found {
			t.Error("Expiring item not found in expiring soon list")
		}
	})

	t.Run("Delete inventory item", func(t *testing.T) {
		deleteReq := map[string]interface{}{
			"reason":      "E2E test cleanup",
			"attestation": "Test completed",
		}

		resp := ctx.MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/inventory/%s", inventoryID), deleteReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)
	})
}
