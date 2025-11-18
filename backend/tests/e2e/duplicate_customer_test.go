package e2e

import (
	"net/http"
	"testing"

	"github.com/example/pgvms/tests/testutil"
)

func TestDuplicateCustomerCheck_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	contactNumber := "+919876543210"

	// Step 1: Create first customer with contact number
	t.Run("Create first customer", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "John Doe",
			"ContactNumber": contactNumber,
			"Email":         "john@example.com",
			"Address":       "123 Main St",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Created first customer: %s", data["ID"].(string))
	})

	// Step 2: Try to create duplicate customer with same contact number
	t.Run("Attempt duplicate customer creation", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "Jane Doe",
			"ContactNumber": contactNumber, // Same contact number
			"Email":         "jane@example.com",
			"Address":       "456 Oak Ave",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)

		// Should get 422 Unprocessable Entity (business error)
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422 for duplicate, got %d", resp.StatusCode)
		}

		t.Logf("Duplicate correctly rejected with status %d", resp.StatusCode)
	})

	// Step 3: Create customer with different contact number (should succeed)
	t.Run("Create customer with different contact", func(t *testing.T) {
		customer := map[string]interface{}{
			"Name":          "Bob Smith",
			"ContactNumber": "+919876543211", // Different number
			"Email":         "bob@example.com",
			"Address":       "789 Pine Rd",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Created third customer: %s", data["ID"].(string))
	})
}
