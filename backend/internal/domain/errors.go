package domain

import (
	"errors"
	"fmt"
)

// Common domain errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation error")
	ErrInvalidAttestation = errors.New("invalid attestation - must type 'I CONFIRM DELETE' exactly")
	ErrMissingReason      = errors.New("deletion reason is required")
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ErrInvalidInput creates a validation error
func ErrInvalidInput(message string) error {
	return &ValidationError{Message: message}
}

// NewValidationError creates a validation error with a message
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

// ErrInvalidField creates a field-specific validation error
func ErrInvalidField(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

// BusinessError represents a business logic error
type BusinessError struct {
	Code    string
	Message string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewBusinessError creates a new business error
func NewBusinessError(code, message string) error {
	return &BusinessError{Code: code, Message: message}
}

// DeleteRequest represents a soft delete request with attestation
type DeleteRequest struct {
	Reason      string `json:"reason"`
	Attestation string `json:"attestation"`
}

// Validate validates delete request
func (d *DeleteRequest) Validate() error {
	if d.Reason == "" {
		return ErrMissingReason
	}
	if d.Attestation != "I CONFIRM DELETE" {
		return ErrInvalidAttestation
	}
	return nil
}
