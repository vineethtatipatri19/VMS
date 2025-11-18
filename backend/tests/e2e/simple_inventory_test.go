package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/example/pgvms/tests/testutil"
)

func TestSimpleInventoryCreate_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	item := map[string]interface{}{
		"Name":         "Simple Product",
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
	t.Logf("Request body: %s", string(body))

	resp := ctx.MakeRequest(t, "POST", "/api/v1/inventory", item)

	// Read response body
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	t.Logf("Response status: %d", resp.StatusCode)
	t.Logf("Response body: %s", responseBody)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201, got %d: %s", resp.StatusCode, responseBody)
	}
}
