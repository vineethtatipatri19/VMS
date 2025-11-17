package domain

import (
	"errors"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantMsg  string
		checkErr error
	}{
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			wantMsg:  "resource not found",
			checkErr: ErrNotFound,
		},
		{
			name:     "ErrAlreadyExists",
			err:      ErrAlreadyExists,
			wantMsg:  "resource already exists",
			checkErr: ErrAlreadyExists,
		},
		{
			name:     "ErrUnauthorized",
			err:      ErrUnauthorized,
			wantMsg:  "unauthorized",
			checkErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.wantMsg)
			}
			if !errors.Is(tt.err, tt.checkErr) {
				t.Errorf("errors.Is(%v, %v) = false, want true", tt.err, tt.checkErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "simple validation error",
			err:     ErrInvalidInput("name is required"),
			wantMsg: "name is required",
		},
		{
			name:    "field validation error",
			err:     ErrInvalidField("email", "invalid format"),
			wantMsg: "email: invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestDeleteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *DeleteRequest
		wantErr error
	}{
		{
			name: "valid delete request",
			req: &DeleteRequest{
				Reason:      "No longer needed",
				Attestation: "I CONFIRM DELETE",
			},
			wantErr: nil,
		},
		{
			name: "missing reason",
			req: &DeleteRequest{
				Reason:      "",
				Attestation: "I CONFIRM DELETE",
			},
			wantErr: ErrMissingReason,
		},
		{
			name: "invalid attestation",
			req: &DeleteRequest{
				Reason:      "Test",
				Attestation: "I confirm delete",
			},
			wantErr: ErrInvalidAttestation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
