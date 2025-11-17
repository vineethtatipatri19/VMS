package domain

import (
	"testing"
)

func TestCustomer_Validate(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		wantErr  bool
	}{
		{
			name: "valid customer",
			customer: Customer{
				Name:         "John Doe",
				CustomerType: "b2b",
				Status:       "active",
				CreditLimit:  1000.0,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			customer: Customer{
				CustomerType: "b2b",
				Status:       "active",
			},
			wantErr: true,
		},
		{
			name: "invalid customer type",
			customer: Customer{
				Name:         "John Doe",
				CustomerType: "invalid",
				Status:       "active",
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			customer: Customer{
				Name:         "John Doe",
				CustomerType: "b2b",
				Status:       "invalid",
			},
			wantErr: true,
		},
		{
			name: "negative credit limit",
			customer: Customer{
				Name:        "John Doe",
				CreditLimit: -100.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.customer.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Customer.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomer_CanPurchase(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		amount   float64
		want     bool
	}{
		{
			name: "can purchase - within limit",
			customer: Customer{
				Status:         "active",
				CurrentBalance: 500.0,
				CreditLimit:    1000.0,
			},
			amount: 400.0,
			want:   true,
		},
		{
			name: "cannot purchase - exceeds limit",
			customer: Customer{
				Status:         "active",
				CurrentBalance: 800.0,
				CreditLimit:    1000.0,
			},
			amount: 300.0,
			want:   false,
		},
		{
			name: "cannot purchase - inactive status",
			customer: Customer{
				Status:         "inactive",
				CurrentBalance: 0,
				CreditLimit:    1000.0,
			},
			amount: 100.0,
			want:   false,
		},
		{
			name: "can purchase - zero balance",
			customer: Customer{
				Status:         "active",
				CurrentBalance: 0,
				CreditLimit:    1000.0,
			},
			amount: 500.0,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.customer.CanPurchase(tt.amount); got != tt.want {
				t.Errorf("Customer.CanPurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomer_IsOverdue(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		want     bool
	}{
		{
			name: "is overdue - positive balance",
			customer: Customer{
				CurrentBalance: 500.0,
			},
			want: true,
		},
		{
			name: "not overdue - zero balance",
			customer: Customer{
				CurrentBalance: 0,
			},
			want: false,
		},
		{
			name: "not overdue - negative balance",
			customer: Customer{
				CurrentBalance: -100.0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.customer.IsOverdue(); got != tt.want {
				t.Errorf("Customer.IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCustomerType(t *testing.T) {
	tests := []struct {
		name         string
		customerType string
		want         bool
	}{
		{"valid b2b", "b2b", true},
		{"valid b2c", "b2c", true},
		{"valid retail", "retail", true},
		{"valid wholesale", "wholesale", true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidCustomerType(tt.customerType); got != tt.want {
				t.Errorf("isValidCustomerType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"valid active", "active", true},
		{"valid inactive", "inactive", true},
		{"valid blocked", "blocked", true},
		{"invalid status", "invalid", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidStatus(tt.status); got != tt.want {
				t.Errorf("isValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
