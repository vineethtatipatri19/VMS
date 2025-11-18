package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuth(t *testing.T) {
	jwtSecret := []byte("test-secret-key")

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("UserID: " + userID))
	})

	// Wrap with Auth middleware
	authMiddleware := Auth(jwtSecret)
	handler := authMiddleware(testHandler)

	// Helper to create a valid token
	createToken := func(userID, email string) string {
		claims := &Claims{
			UserID: userID,
			Email:  email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)
		return tokenString
	}

	t.Run("valid token allows access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		token := createToken("user123", "test@example.com")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		if w.Body.String() != "UserID: user123" {
			t.Errorf("Expected 'UserID: user123', got '%s'", w.Body.String())
		}
	})

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}

		if w.Body.String() != "missing authorization header\n" {
			t.Errorf("Unexpected error message: %s", w.Body.String())
		}
	})

	t.Run("invalid authorization header format returns 401", func(t *testing.T) {
		testCases := []struct {
			name   string
			header string
		}{
			{"no bearer prefix", "sometoken"},
			{"wrong prefix", "Basic sometoken"},
			{"empty token", "Bearer "},
			{"no token", "Bearer"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.Header.Set("Authorization", tc.header)
				w := httptest.NewRecorder()

				handler.ServeHTTP(w, req)

				if w.Code != http.StatusUnauthorized {
					t.Errorf("Expected status 401, got %d", w.Code)
				}
			})
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		claims := &Claims{
			UserID: "user123",
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for expired token, got %d", w.Code)
		}
	})

	t.Run("token signed with wrong secret returns 401", func(t *testing.T) {
		wrongSecret := []byte("wrong-secret")
		claims := &Claims{
			UserID: "user123",
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(wrongSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for wrong secret, got %d", w.Code)
		}
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns user ID from context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := req.Context()
		ctx = contextWithUserID(ctx, "test-user-123")
		req = req.WithContext(ctx)

		userID := GetUserID(req.Context())
		if userID != "test-user-123" {
			t.Errorf("Expected 'test-user-123', got '%s'", userID)
		}
	})

	t.Run("returns empty string when no user ID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		userID := GetUserID(req.Context())

		if userID != "" {
			t.Errorf("Expected empty string, got '%s'", userID)
		}
	})
}

// Helper function to add user ID to context (for testing GetUserID)
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
