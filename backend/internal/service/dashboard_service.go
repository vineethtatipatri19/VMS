package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DashboardService handles dashboard statistics
type DashboardService struct {
	db *sql.DB
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(db *sql.DB) *DashboardService {
	return &DashboardService{db: db}
}

// DashboardStats represents KPI data for the dashboard
type DashboardStats struct {
	TotalCustomers      int     `json:"totalCustomers"`
	ExpiringSoonItems   int     `json:"expiringSoonItems"`
	ExpiredItems        int     `json:"expiredItems"`
	TotalInventoryValue float64 `json:"totalInventoryValue"`
	UnreturnedCrates    int     `json:"unreturnedCrates"`
	OutstandingBalance  float64 `json:"outstandingBalance"`
	TodaysSales         float64 `json:"todaysSales"`
	MonthSales          float64 `json:"monthSales"`
}

// RecentActivity represents recent system activity
type RecentActivity struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetStats retrieves dashboard statistics
func (s *DashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Total customers
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL`).Scan(&stats.TotalCustomers)
	if err != nil {
		return nil, err
	}

	// Expiring soon items (within 3 days)
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inventory_items 
		WHERE expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days'
		AND quantity > 0
		AND deleted_at IS NULL`).Scan(&stats.ExpiringSoonItems)
	if err != nil {
		stats.ExpiringSoonItems = 0
	}

	// Expired items
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inventory_items 
		WHERE expiry_date < CURRENT_DATE
		AND quantity > 0
		AND deleted_at IS NULL`).Scan(&stats.ExpiredItems)
	if err != nil {
		stats.ExpiredItems = 0
	}

	// Unreturned crates (sum of all customer balances)
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(balance), 0) 
		FROM (
			SELECT DISTINCT ON (customer_id) balance 
			FROM crate_ledger 
			WHERE deleted_at IS NULL
			ORDER BY customer_id, date DESC
		) as latest_balances`).Scan(&stats.UnreturnedCrates)
	if err != nil {
		stats.UnreturnedCrates = 0
	}

	// Outstanding balance (total sales - total payments)
	var totalSales, totalPayments float64
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale' AND deleted_at IS NULL`).Scan(&totalSales)
	if err != nil {
		totalSales = 0
	}

	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(payment_amount), 0) 
		FROM transactions 
		WHERE type='payment' AND deleted_at IS NULL`).Scan(&totalPayments)
	if err != nil {
		totalPayments = 0
	}

	stats.OutstandingBalance = totalSales - totalPayments

	// Today's sales
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale' 
		AND date >= CURRENT_DATE
		AND deleted_at IS NULL`).Scan(&stats.TodaysSales)
	if err != nil {
		stats.TodaysSales = 0
	}

	// Month sales
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale' 
		AND date >= date_trunc('month', CURRENT_DATE)
		AND deleted_at IS NULL`).Scan(&stats.MonthSales)
	if err != nil {
		stats.MonthSales = 0
	}

	// Total inventory value
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(quantity * COALESCE(cost_price, 0)), 0) 
		FROM inventory_items 
		WHERE status IN ('available', 'low_stock') 
		AND quantity > 0
		AND deleted_at IS NULL`).Scan(&stats.TotalInventoryValue)
	if err != nil {
		stats.TotalInventoryValue = 0
	}

	return stats, nil
}

// GetRecentActivity retrieves recent activity
func (s *DashboardService) GetRecentActivity(ctx context.Context, limit int) ([]RecentActivity, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT t.type, t.total_amount, t.created_at, c.name
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		WHERE t.deleted_at IS NULL
		ORDER BY t.created_at DESC
		LIMIT $1`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []RecentActivity{}
	for rows.Next() {
		var txType string
		var amount float64
		var timestamp time.Time
		var customerName string

		if err := rows.Scan(&txType, &amount, &timestamp, &customerName); err != nil {
			continue
		}

		var description string
		if txType == "sale" {
			description = customerName + " - Sale: ₹" + formatCurrency(amount)
		} else {
			description = customerName + " - Payment: ₹" + formatCurrency(amount)
		}

		activities = append(activities, RecentActivity{
			Type:        txType,
			Description: description,
			Timestamp:   timestamp,
		})
	}

	return activities, nil
}

// Helper function to format currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
