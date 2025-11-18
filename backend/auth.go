package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "default-secret-change-in-production"
	}
	return secret
}

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "INVALID_INPUT",
				"message": "Invalid request format",
			},
		})
		return
	}

	// Validate required fields
	if u.Email == "" || u.Password == "" || u.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "MISSING_FIELDS",
				"message": "Name, email, and password are required",
			},
		})
		return
	}

	// Validate password length
	if len(u.Password) < 6 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "WEAK_PASSWORD",
				"message": "Password must be at least 6 characters long",
			},
		})
		return
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "HASH_ERROR",
				"message": "Failed to process password",
			},
		})
		return
	}

	id := uuid.New().String()
	_, err = db.ExecContext(r.Context(), `INSERT INTO users(id,name,email,password_hash,created_at) VALUES ($1,$2,$3,$4,now())`, id, u.Name, u.Email, hash)
	if err != nil {
		// Check if it's a duplicate email error
		errorMsg := "Registration failed"
		if contains(err.Error(), "duplicate") || contains(err.Error(), "unique") {
			errorMsg = "An account with this email already exists"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DUPLICATE_EMAIL",
				"message": errorMsg,
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"id":      id,
			"message": "Registration successful",
		},
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct{ Email, Password string }

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "INVALID_INPUT",
				"message": "Invalid request format",
			},
		})
		return
	}

	// Validate required fields
	if creds.Email == "" || creds.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "MISSING_FIELDS",
				"message": "Email and password are required",
			},
		})
		return
	}

	var id, hash string
	err := db.QueryRowContext(r.Context(), `SELECT id,password_hash FROM users WHERE email=$1`, creds.Email).Scan(&id, &hash)

	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DATABASE_ERROR",
				"message": "An error occurred during login",
			},
		})
		return
	}

	if err := checkPasswordHash(hash, creds.Password); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Create token
	exp := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		Email:  creds.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(jwtKey)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate authentication token",
			},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"token": s,
		},
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		log.Printf("Auth header: %s", auth)
		if auth == "" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"missing authorization header"}`, 401)
			return
		}
		var tokenString string
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenString = auth[7:]
			log.Printf("Token string length: %d", len(tokenString))
		} else {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"invalid auth header format, use: Bearer <token>"}`, 401)
			return
		}
		claims := &Claims{}
		log.Printf("Using JWT secret length: %d", len(jwtKey))
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			log.Printf("In key func, returning key")
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			errMsg := "invalid or expired token"
			if err != nil {
				errMsg = err.Error()
				log.Printf("JWT Parse Error: %v", err)
			}
			if token != nil && !token.Valid {
				log.Printf("Token invalid but parsed. Valid=%v", token.Valid)
			}
			http.Error(w, `{"error":"`+errMsg+`"}`, 401)
			return
		}
		log.Printf("Token valid for user: %s", claims.UserID)
		// Store user id in context
		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
