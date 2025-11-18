package e2e

import (
	"net/http"
	"testing"

	"github.com/example/pgvms/tests/testutil"
)

// TestMultipleCustomers_E2E verifies that multiple customers can be created (bug fix test)
func TestMultipleCustomers_E2E(t *testing.T) {
	ctx := testutil.SetupE2ETest(t)
	defer testutil.TeardownE2ETest(t, ctx)

	// Create first customer
	t.Run("Create first customer", func(t *testing.T) {
		customer1 := map[string]interface{}{
			"Name":          "Customer One",
			"ContactNumber": "+1111111111",
			"Email":         "customer1@test.com",
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer1)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Created first customer: %s", data["ID"].(string))
	})

	// Create second customer - this should NOT fail with ErrAlreadyExists
	t.Run("Create second customer", func(t *testing.T) {
		customer2 := map[string]interface{}{
			"Name":          "Customer Two",
			"ContactNumber": "+2222222222",
			"Email":         "customer2@test.com",
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer2)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Created second customer: %s", data["ID"].(string))
	})

	// Create third customer
	t.Run("Create third customer", func(t *testing.T) {
		customer3 := map[string]interface{}{
			"Name":          "Customer Three",
			"ContactNumber": "+3333333333",
			"Email":         "customer3@test.com",
			"Status":        "active",
		}

		resp := ctx.MakeRequest(t, "POST", "/api/v1/customers", customer3)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusCreated)

		data := ctx.ExtractData(t, resp)
		t.Logf("Created third customer: %s", data["ID"].(string))
	})

	// Verify all customers exist
	t.Run("List all customers", func(t *testing.T) {
		resp := ctx.MakeRequest(t, "GET", "/api/v1/customers", nil)
		testutil.AssertHTTPStatusCode(t, resp, http.StatusOK)

		// Response is {success:true, data:[...array...]}
		var wrapper map[string]interface{}
		ctx.DecodeResponse(t, resp, &wrapper)

		customers, ok := wrapper["data"].([]interface{})
		if !ok {
			t.Fatalf("Expected data to be an array")
		}

		if len(customers) < 3 {
			t.Errorf("Expected at least 3 customers, got %d", len(customers))
		} else {
			t.Logf("Successfully created and listed %d customers", len(customers))
		}
	})
}
