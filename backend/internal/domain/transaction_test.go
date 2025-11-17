package domain

import (
	"testing"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		txn     *Transaction
		wantErr bool
	}{
		{
			name: "valid sale transaction",
			txn: &Transaction{
				CustomerID:  "cust123",
				Type:        "sale",
				TotalAmount: 100.0,
			},
			wantErr: false,
		},
		{
			name: "valid payment transaction",
			txn: &Transaction{
				CustomerID:    "cust123",
				Type:          "payment",
				PaymentAmount: 200.0,
			},
			wantErr: false,
		},
		{
			name: "missing customer ID",
			txn: &Transaction{
				CustomerID:  "",
				Type:        "sale",
				TotalAmount: 100.0,
			},
			wantErr: true,
		},
		{
			name: "invalid transaction type",
			txn: &Transaction{
				CustomerID:  "cust123",
				Type:        "invalid",
				TotalAmount: 100.0,
			},
			wantErr: true,
		},
		{
			name: "negative total amount for sale",
			txn: &Transaction{
				CustomerID:  "cust123",
				Type:        "sale",
				TotalAmount: -50.0,
			},
			wantErr: true,
		},
		{
			name: "zero payment amount",
			txn: &Transaction{
				CustomerID:    "cust123",
				Type:          "payment",
				PaymentAmount: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.txn.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaleItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    *SaleItem
		wantErr bool
	}{
		{
			name: "valid sale item",
			item: &SaleItem{
				ItemName:     "Test Item",
				Quantity:     5,
				PricePerUnit: 100,
			},
			wantErr: false,
		},
		{
			name: "missing item name",
			item: &SaleItem{
				ItemName:     "",
				Quantity:     5,
				PricePerUnit: 100,
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			item: &SaleItem{
				ItemName:     "Test",
				Quantity:     0,
				PricePerUnit: 100,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			item: &SaleItem{
				ItemName:     "Test",
				Quantity:     5,
				PricePerUnit: -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SaleItem.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaleItem_CalculateProfit(t *testing.T) {
	item := &SaleItem{
		Quantity:     10,
		PricePerUnit: 150,
		CostPerUnit:  100,
	}
	item.CalculateProfit()

	want := 500.0
	if item.Profit != want {
		t.Errorf("SaleItem.CalculateProfit() = %.2f, want %.2f", item.Profit, want)
	}
}

func TestSaleItem_CalculateTotal(t *testing.T) {
	tests := []struct {
		name string
		item *SaleItem
		want float64
	}{
		{
			name: "no discount or tax",
			item: &SaleItem{
				Quantity:        10,
				PricePerUnit:    100,
				DiscountPercent: 0,
				TaxPercent:      0,
			},
			want: 1000.0,
		},
		{
			name: "with 10% discount",
			item: &SaleItem{
				Quantity:        10,
				PricePerUnit:    100,
				DiscountPercent: 10,
				TaxPercent:      0,
			},
			want: 900.0,
		},
		{
			name: "with 18% tax",
			item: &SaleItem{
				Quantity:        10,
				PricePerUnit:    100,
				DiscountPercent: 0,
				TaxPercent:      18,
			},
			want: 1180.0,
		},
		{
			name: "with discount and tax",
			item: &SaleItem{
				Quantity:        10,
				PricePerUnit:    100,
				DiscountPercent: 10,
				TaxPercent:      18,
			},
			want: 1062.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.CalculateTotal()
			if tt.item.Total != tt.want {
				t.Errorf("SaleItem.CalculateTotal() = %.2f, want %.2f", tt.item.Total, tt.want)
			}
		})
	}
}
