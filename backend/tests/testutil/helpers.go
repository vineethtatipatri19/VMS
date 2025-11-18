package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string for testing
func GenerateUUID() string {
	return uuid.New().String()
}

// MakeRequest creates and executes an HTTP request for testing
func MakeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		bodyReader = strings.NewReader(string(jsonBody))
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	return w
}

// ParseJSONResponse parses JSON response into target struct
func ParseJSONResponse(t *testing.T, resp *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("Failed to decode response: %v\nBody: %s", err, resp.Body.String())
	}
}

// AssertStatusCode checks if response has expected status code
func AssertStatusCode(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("Expected status code %d, got %d\nBody: %s", expected, resp.Code, resp.Body.String())
	}
}

// AssertNoError fails the test if error is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError fails the test if error is nil
func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual checks if two values are not equal
func AssertNotEqual(t *testing.T, notExpected, actual interface{}) {
	t.Helper()

	if notExpected == actual {
		t.Errorf("Expected values to be different, both are %v", actual)
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, str, substr string) {
	t.Helper()

	if !strings.Contains(str, substr) {
		t.Errorf("Expected string to contain %q, got %q", substr, str)
	}
}

// AssertJSONEqual compares two JSON strings
func AssertJSONEqual(t *testing.T, expected, actual string) {
	t.Helper()

	var expectedObj, actualObj interface{}

	if err := json.Unmarshal([]byte(expected), &expectedObj); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualObj); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	expectedJSON, _ := json.Marshal(expectedObj)
	actualJSON, _ := json.Marshal(actualObj)

	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("JSON not equal:\nExpected: %s\nActual: %s", expectedJSON, actualJSON)
	}
}

// SetAuthHeader sets Authorization header with JWT token
func SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

// GetTestToken generates a test JWT token (implement based on your auth)
func GetTestToken(t *testing.T, userID string) string {
	t.Helper()
	// TODO: Implement based on your JWT generation logic
	return "test-token-" + userID
}
