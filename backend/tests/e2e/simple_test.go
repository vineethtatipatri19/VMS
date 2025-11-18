package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/example/pgvms/tests/testutil"
)

func TestSimpleCustomerCreate_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	customer := map[string]interface{}{
		"Name":          "Simple Test",
		"ContactNumber": "+1234567890",
		"Status":        "active",
		"CustomerType":  "b2c",
	}

	body, _ := json.Marshal(customer)
	t.Logf("Request body: %s", string(body))

	resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)

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
