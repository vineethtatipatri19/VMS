package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/example/pgvms/tests/testutil"
)

// TestCrateTrackingFlow_E2E tests crate issue, return, and balance tracking
func TestCrateTrackingFlow_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var customerID string

	// Step 1: Create customer
	t.Run("Create customer", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "Crate Test Customer",
			"ContactNumber": "+5555555555",
			"CustomerType":  "wholesale",
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		customerID = data["ID"].(string)
		t.Logf("Created customer: %s", customerID)
	})

	// Step 2: Issue crates to customer
	t.Run("Issue crates", func(t *testing.T) {
		issueReq := map[string]interface{}{
			"CustomerID":      customerID,
			"TransactionType": "out", // out = issued to customer
			"Quantity":        20,
			"UnitPrice":       10.0,
			"TotalPrice":      200.0,
			"Notes":           "Initial crate issue",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/crates/issue", issueReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Issued crates, transaction: %v", data["ID"])
	})

	// Step 3: Check crate balance (should be 20)
	t.Run("Check initial balance", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/crates/balance/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		balance := data["balance"].(float64) // Note: lowercase 'balance'

		testutil.AssertEqual(t, 20.0, balance)
		t.Logf("Crate balance after issue: %.0f", balance)
	})

	// Step 4: Customer returns some crates
	t.Run("Return crates", func(t *testing.T) {
		returnReq := map[string]interface{}{
			"CustomerID":      customerID,
			"TransactionType": "in", // in = returned from customer
			"Quantity":        8,
			"UnitPrice":       10.0,
			"TotalPrice":      80.0,
			"Notes":           "Partial crate return",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/crates/return", returnReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Returned crates, transaction: %v", data["ID"])
	})

	// Step 5: Check updated balance (should be 12)
	t.Run("Check balance after return", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/crates/balance/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		balance := data["balance"].(float64) // Note: lowercase 'balance'

		testutil.AssertEqual(t, 12.0, balance) // 20 - 8 = 12
		t.Logf("Crate balance after return: %.0f", balance)
	})

	// Step 6: Get crate transaction history
	t.Run("Check crate history", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/crates/history/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		// The data field IS the array directly
		var wrapper map[string]interface{}
		ctx.DecodeResponse(t, resp, &wrapper)

		history, ok := wrapper["data"].([]interface{})
		if !ok {
			t.Fatalf("Expected data to be an array, got: %T", wrapper["data"])
		}

		if len(history) < 2 {
			t.Errorf("Expected at least 2 transactions (issue + return), got %d", len(history))
		}
		t.Logf("Found %d crate transactions in history", len(history))
	})

	// Step 7: Issue more crates
	t.Run("Issue additional crates", func(t *testing.T) {
		issueReq := map[string]interface{}{
			"CustomerID":      customerID,
			"TransactionType": "out",
			"Quantity":        15,
			"UnitPrice":       10.0,
			"TotalPrice":      150.0,
			"Notes":           "Additional crate issue",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/crates/issue", issueReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)
	})

	// Step 8: Verify final balance (should be 27)
	t.Run("Check final balance", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/crates/balance/%s", customerID), nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		balance := data["balance"].(float64) // Note: lowercase 'balance'

		testutil.AssertEqual(t, 27.0, balance) // 12 + 15 = 27
		t.Logf("Final crate balance: %.0f", balance)
	})
}

// TestExpiryAlertFlow_E2E tests expiry alert generation and acknowledgment
func TestExpiryAlertFlow_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	var inventoryID, alertID string

	// Step 1: Add inventory item expiring soon
	t.Run("Add inventory expiring soon", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":         "Expiring Product",
			"LotNumber":    testutil.GenerateUUID(),
			"Unit":         "units",
			"Quantity":     50.0,
			"CostPrice":    10.0,
			"SellingPrice": 15.0,
			"PurchaseDate": "2025-11-01",
			"ExpiryDate":   "2025-12-15", // Expiring in ~1 month
			"Status":       "available",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		inventoryID = data["ID"].(string)
		t.Logf("Created inventory item: %s", inventoryID)
	})

	// Step 2: Generate expiry alerts
	t.Run("Generate expiry alerts", func(t *testing.T) {
		generateReq := map[string]interface{}{
			"DaysThreshold": 60, // Alert for items expiring within 60 days
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/expiry-alerts/generate", generateReq)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Generated alerts: %v", data)
	})

	// Step 3: List pending expiry alerts
	t.Run("List pending alerts", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", "/api/v1/expiry-alerts/pending", nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		data := ctx.ExtractData(t, resp)
		alerts, ok := data["alerts"].([]interface{})
		if !ok {
			t.Fatal("Expected alerts array in response")
		}

		if len(alerts) == 0 {
			t.Error("Expected at least one pending alert")
		} else {
			firstAlert := alerts[0].(map[string]interface{})
			alertID = firstAlert["ID"].(string)
			t.Logf("Found %d pending alert(s), first alert ID: %s", len(alerts), alertID)
		}
	})

	// Step 4: Acknowledge an alert
	if alertID != "" {
		t.Run("Acknowledge alert", func(t *testing.T) {
			ackReq := map[string]interface{}{
				"AcknowledgedBy": "test-user",
				"ActionTaken":    "Marked for discount sale",
			}

			resp := ctx.MakeRequest(t, "POST", fmt.Sprintf("/api/v1/expiry-alerts/%s/acknowledge", alertID), ackReq)
			testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

			data := ctx.ExtractData(t, resp)
			t.Logf("Acknowledged alert: %v", data["Status"])
		})

		// Step 5: Verify alert is acknowledged
		t.Run("Verify alert acknowledged", func(t *testing.T) {
			resp := ctx.MakeRequest(t, "GET", fmt.Sprintf("/api/v1/expiry-alerts/%s", alertID), nil)
			testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

			data := ctx.ExtractData(t, resp)
			status := data["Status"].(string)

			if status != "acknowledged" {
				t.Errorf("Expected status 'acknowledged', got '%s'", status)
			}
			t.Logf("Alert status: %s", status)
		})
	}
}
