package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/example/pgvms/internal/domain"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// ErrorData represents error details in API response
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: statusCode < 400,
		Data:    data,
	})
}

// SendError sends an error JSON response
func SendError(w http.ResponseWriter, err error) {
	statusCode, errData := mapError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   errData,
	})
}

// mapError maps domain errors to HTTP status codes and error data
func mapError(err error) (int, *ErrorData) {
	// Check for validation errors
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest, &ErrorData{
			Code:    "VALIDATION_ERROR",
			Message: validationErr.Message,
			Field:   validationErr.Field,
		}
	}

	// Check for business errors
	var businessErr *domain.BusinessError
	if errors.As(err, &businessErr) {
		return http.StatusUnprocessableEntity, &ErrorData{
			Code:    businessErr.Code,
			Message: businessErr.Message,
		}
	}

	// Check for common domain errors
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, &ErrorData{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict, &ErrorData{
			Code:    "ALREADY_EXISTS",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, &ErrorData{
			Code:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, &ErrorData{
			Code:    "FORBIDDEN",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrInvalidAttestation):
		return http.StatusBadRequest, &ErrorData{
			Code:    "INVALID_ATTESTATION",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrMissingReason):
		return http.StatusBadRequest, &ErrorData{
			Code:    "MISSING_REASON",
			Message: err.Error(),
		}
	default:
		// Generic internal server error
		// TODO: Add proper logging system to capture actual errors
		log.Printf("Unhandled error: %v (type: %T)", err, err)
		return http.StatusInternalServerError, &ErrorData{
			Code:    "INTERNAL_ERROR",
			Message: "An internal error occurred",
		}
	}
}

// DecodeJSON decodes JSON request body
func DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return domain.ErrInvalidInput("invalid JSON body")
	}
	return nil
}

// RespondSuccess sends a successful response with data
func RespondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	SendJSON(w, statusCode, data)
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
		},
	})
}
