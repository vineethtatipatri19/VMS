package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/example/pgvms/internal/httputil"
)

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

// Handler for dashboard KPIs
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := DashboardStats{}

	// Total customers
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers`).Scan(&stats.TotalCustomers)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Expiring soon items (within 3 days)
	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inventory_items 
		WHERE expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days'
		AND quantity > 0`).Scan(&stats.ExpiringSoonItems)
	if err != nil {
		stats.ExpiringSoonItems = 0
	}

	// Expired items
	err = db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inventory_items 
		WHERE expiry_date < CURRENT_DATE
		AND quantity > 0`).Scan(&stats.ExpiredItems)
	if err != nil {
		stats.ExpiredItems = 0
	}

	// Unreturned crates (sum of all customer balances)
	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(balance), 0) 
		FROM (
			SELECT DISTINCT ON (customer_id) balance 
			FROM crate_ledger 
			ORDER BY customer_id, date DESC
		) as latest_balances`).Scan(&stats.UnreturnedCrates)
	if err != nil {
		stats.UnreturnedCrates = 0
	}

	// Outstanding balance (total sales - total payments)
	var totalSales, totalPayments float64
	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale'`).Scan(&totalSales)
	if err != nil {
		totalSales = 0
	}

	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(payment_amount), 0) 
		FROM transactions 
		WHERE type='payment'`).Scan(&totalPayments)
	if err != nil {
		totalPayments = 0
	}

	stats.OutstandingBalance = totalSales - totalPayments

	// Today's sales
	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale' 
		AND date >= CURRENT_DATE`).Scan(&stats.TodaysSales)
	if err != nil {
		stats.TodaysSales = 0
	}

	// Month sales
	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0) 
		FROM transactions 
		WHERE type='sale' 
		AND date >= date_trunc('month', CURRENT_DATE)`).Scan(&stats.MonthSales)
	if err != nil {
		stats.MonthSales = 0
	}

	// Total inventory value (quantity * cost_price for all available items)
	err = db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(quantity * COALESCE(cost_price, 0)), 0) 
		FROM inventory_items 
		WHERE status IN ('available', 'low_stock') 
		AND quantity > 0`).Scan(&stats.TotalInventoryValue)
	if err != nil {
		stats.TotalInventoryValue = 0
	}

	httputil.SendJSON(w, http.StatusOK, stats)
}

// RecentActivity represents recent system activity
type RecentActivity struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// Handler for recent activity
func recentActivityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get recent transactions
	rows, err := db.QueryContext(ctx, `
		SELECT t.type, t.total_amount, t.created_at, c.name
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		ORDER BY t.created_at DESC
		LIMIT 10`)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
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
			description = customerName + " - Sale: ₹" + formatFloat(amount)
		} else {
			description = customerName + " - Payment: ₹" + formatFloat(amount)
		}

		activities = append(activities, RecentActivity{
			Type:        txType,
			Description: description,
			Timestamp:   timestamp,
		})
	}

	httputil.SendJSON(w, http.StatusOK, activities)
}

// Helper function to format floats as currency
func formatFloat(f float64) string {
	// Simple formatting - in production use proper currency formatter
	return fmt.Sprintf("%.2f", f)
}
