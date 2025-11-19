package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request format")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.Name == "" {
		httputil.RespondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Name, email, and password are required")
		return
	}

	// Validate password length
	if len(req.Password) < 6 {
		httputil.RespondError(w, http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 6 characters long")
		return
	}

	userID, err := h.authService.Register(r.Context(), service.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Check if it's a validation error (duplicate email)
		if errStr := err.Error(); errStr == "email already exists" || contains(errStr, "duplicate") {
			httputil.RespondError(w, http.StatusConflict, "DUPLICATE_EMAIL", "An account with this email already exists")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Registration failed")
		return
	}

	httputil.RespondSuccess(w, http.StatusCreated, map[string]string{
		"id":      userID,
		"message": "Registration successful",
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request format")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		httputil.RespondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required")
		return
	}

	authResp, err := h.authService.Login(r.Context(), service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// For security, don't reveal if email exists or password is wrong
		httputil.RespondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	httputil.RespondSuccess(w, http.StatusOK, map[string]string{
		"token": authResp.Token,
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
