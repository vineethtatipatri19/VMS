package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/example/pgvms/tests/testutil"
)

func TestDebugInventory_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	// Test 1: Base case (known to work)
	t.Run("Base inventory (works)", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":         "Base Product",
			"LotNumber":    testutil.GenerateUUID(),
			"Unit":         "kg",
			"Quantity":     100.0,
			"CostPrice":    50.0,
			"SellingPrice": 75.0,
			"PurchaseDate": time.Now().Format("2006-01-02"),
			"ExpiryDate":   time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"Status":       "available",
		}

		body, _ := json.Marshal(item)
		t.Logf("Request: %s", string(body))

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		t.Logf("Status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Base test failed with status %d", resp.StatusCode)
		}
	})

	// Test 2: Add SupplierName
	t.Run("With SupplierName", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":         "Product With Supplier",
			"LotNumber":    testutil.GenerateUUID(),
			"Unit":         "kg",
			"Quantity":     100.0,
			"CostPrice":    50.0,
			"SellingPrice": 75.0,
			"SupplierName": "Test Supplier",
			"PurchaseDate": time.Now().Format("2006-01-02"),
			"ExpiryDate":   time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"Status":       "available",
		}

		body, _ := json.Marshal(item)
		t.Logf("Request: %s", string(body))

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		t.Logf("Status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("SupplierName test failed with status %d", resp.StatusCode)
		}
	})

	// Test 3: Add MinStockLevel
	t.Run("With MinStockLevel", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "Product With MinStock",
			"LotNumber":     testutil.GenerateUUID(),
			"Unit":          "kg",
			"Quantity":      100.0,
			"CostPrice":     50.0,
			"SellingPrice":  75.0,
			"MinStockLevel": 20.0,
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"Status":        "available",
		}

		body, _ := json.Marshal(item)
		t.Logf("Request: %s", string(body))

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		t.Logf("Status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("MinStockLevel test failed with status %d", resp.StatusCode)
		}
	})

	// Test 5: Test unit="bottles"
	t.Run("With unit=bottles", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "Product With Bottles",
			"LotNumber":     testutil.GenerateUUID(),
			"Unit":          "bottles",
			"Quantity":      100.0,
			"CostPrice":     50.0,
			"SellingPrice":  75.0,
			"SupplierName":  "Test Supplier",
			"MinStockLevel": 20.0,
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"Status":        "available",
		}

		body, _ := json.Marshal(item)
		t.Logf("Request: %s", string(body))

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)

		// Read response body
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		responseBody := buf.String()

		t.Logf("Status: %d", resp.StatusCode)
		t.Logf("Response: %s", responseBody)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("bottles unit test failed with status %d: %s", resp.StatusCode, responseBody)
		}
	})

	// Test 6: Exact copy from complete_flows_test.go
	t.Run("Exact copy from complete_flows_test", func(t *testing.T) {
		item := map[string]interface{}{
			"Name":          "Test Product",
			"Unit":          "bottles",
			"Quantity":      100.0,
			"CostPrice":     50.0,
			"SellingPrice":  75.0,
			"SupplierName":  "Test Supplier",
			"LotNumber":     testutil.GenerateUUID(),
			"PurchaseDate":  time.Now().Format("2006-01-02"),
			"ExpiryDate":    time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
			"MinStockLevel": 20.0,
			"Status":        "available",
		}

		body, _ := json.Marshal(item)
		t.Logf("Request: %s", string(body))

		resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)
		t.Logf("Status: %d", resp.StatusCode)

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Exact copy test failed with status %d", resp.StatusCode)
		}
	})
}
