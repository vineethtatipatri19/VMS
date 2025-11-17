package domain

import (
	"testing"
	"time"
)

func TestInventoryItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		inv     *InventoryItem
		wantErr bool
	}{
		{
			name: "valid inventory",
			inv: &InventoryItem{
				Name:         "Test Product",
				LotNumber:    "LOT001",
				Quantity:     10,
				Unit:         "kg",
				CostPrice:    100.0,
				SellingPrice: 150.0,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			inv: &InventoryItem{
				LotNumber: "LOT001",
				Quantity:  10,
				Unit:      "kg",
			},
			wantErr: true,
		},
		{
			name: "missing lot number",
			inv: &InventoryItem{
				Name:     "Test",
				Quantity: 10,
				Unit:     "kg",
			},
			wantErr: true,
		},
		{
			name: "negative quantity",
			inv: &InventoryItem{
				Name:      "Test",
				LotNumber: "LOT001",
				Quantity:  -1,
				Unit:      "kg",
			},
			wantErr: true,
		},
		{
			name: "missing unit",
			inv: &InventoryItem{
				Name:      "Test",
				LotNumber: "LOT001",
				Quantity:  10,
			},
			wantErr: true,
		},
		{
			name: "negative cost price",
			inv: &InventoryItem{
				Name:      "Test",
				LotNumber: "LOT001",
				Quantity:  10,
				Unit:      "kg",
				CostPrice: -50,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inv.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("InventoryItem.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInventoryItem_IsExpired(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")

	tests := []struct {
		name string
		inv  *InventoryItem
		want bool
	}{
		{
			name: "expired yesterday",
			inv: &InventoryItem{
				ExpiryDate: yesterday,
			},
			want: true,
		},
		{
			name: "expires tomorrow",
			inv: &InventoryItem{
				ExpiryDate: tomorrow,
			},
			want: false,
		},
		{
			name: "no expiry date",
			inv: &InventoryItem{
				ExpiryDate: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.inv.IsExpired(); got != tt.want {
				t.Errorf("InventoryItem.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInventoryItem_IsLowStock(t *testing.T) {
	tests := []struct {
		name string
		inv  *InventoryItem
		want bool
	}{
		{
			name: "low stock - below min level",
			inv: &InventoryItem{
				Quantity:      5,
				MinStockLevel: 10,
			},
			want: true,
		},
		{
			name: "sufficient stock - equal min level",
			inv: &InventoryItem{
				Quantity:      10,
				MinStockLevel: 10,
			},
			want: true,
		},
		{
			name: "sufficient stock - above min level",
			inv: &InventoryItem{
				Quantity:      15,
				MinStockLevel: 10,
			},
			want: false,
		},
		{
			name: "zero stock",
			inv: &InventoryItem{
				Quantity:      0,
				MinStockLevel: 10,
			},
			want: true,
		},
		{
			name: "no min level set",
			inv: &InventoryItem{
				Quantity:      5,
				MinStockLevel: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.inv.IsLowStock(); got != tt.want {
				t.Errorf("InventoryItem.IsLowStock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInventoryItem_CalculateMargin(t *testing.T) {
	tests := []struct {
		name string
		inv  *InventoryItem
		want float64
	}{
		{
			name: "normal profit margin",
			inv: &InventoryItem{
				CostPrice:    100,
				SellingPrice: 150,
			},
			want: 33.33,
		},
		{
			name: "zero cost",
			inv: &InventoryItem{
				CostPrice:    0,
				SellingPrice: 100,
			},
			want: 100,
		},
		{
			name: "zero selling price",
			inv: &InventoryItem{
				CostPrice:    100,
				SellingPrice: 0,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.inv.CalculateMargin()
			// Allow small floating point differences
			if diff := tt.inv.MarginPercentage - tt.want; diff > 0.01 || diff < -0.01 {
				t.Errorf("InventoryItem.CalculateMargin() = %.2f, want %.2f", tt.inv.MarginPercentage, tt.want)
			}
		})
	}
}
