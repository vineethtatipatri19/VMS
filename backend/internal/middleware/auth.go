package middleware

import (
"context"
"net/http"
"strings"

"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT token claims
type Claims struct {
UserID string `json:"userId"`
Email  string `json:"email"`
jwt.RegisteredClaims
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const userIDKey contextKey = "userId"

// Auth validates JWT tokens and adds user ID to request context
func Auth(jwtSecret []byte) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
http.Error(w, "missing authorization header", http.StatusUnauthorized)
return
}

// Extract token from "Bearer <token>"
parts := strings.Split(authHeader, " ")
if len(parts) != 2 || parts[0] != "Bearer" {
http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
return
}
tokenString := parts[1]

// Parse and validate token
claims := &Claims{}
token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
return jwtSecret, nil
})

if err != nil || !token.Valid {
http.Error(w, "invalid or expired token", http.StatusUnauthorized)
return
}

// Add user ID to context
ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
next.ServeHTTP(w, r.WithContext(ctx))
})
}
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) string {
if userID, ok := ctx.Value(userIDKey).(string); ok {
return userID
}
return ""
}
