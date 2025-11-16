package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// ReportType represents different report types
type ReportType string

const (
	ReportTypeSales      ReportType = "sales"
	ReportTypeInventory  ReportType = "inventory"
	ReportTypeCustomer   ReportType = "customer"
	ReportTypeCrates     ReportType = "crates"
	ReportTypeFinancial  ReportType = "financial"
)

// ReportRequest represents a report generation request
type ReportRequest struct {
	Type      ReportType `json:"type"`
	StartDate string     `json:"startDate,omitempty"`
	EndDate   string     `json:"endDate,omitempty"`
	CustomerID string    `json:"customerId,omitempty"`
}

// SalesReport represents a sales report
type SalesReport struct {
	TotalSales    float64           `json:"totalSales"`
	TotalItems    int               `json:"totalItems"`
	Transactions  []Transaction     `json:"transactions"`
	TopItems      []TopItem         `json:"topItems"`
	GeneratedAt   time.Time         `json:"generatedAt"`
}

// TopItem represents top selling items
type TopItem struct {
	ItemName     string  `json:"itemName"`
	TotalQuantity float64 `json:"totalQuantity"`
	TotalRevenue  float64 `json:"totalRevenue"`
}

// InventoryReport represents an inventory report
type InventoryReport struct {
	TotalItems      int             `json:"totalItems"`
	ExpiringSoon    int             `json:"expiringSoon"`
	Expired         int             `json:"expired"`
	Items           []InventoryItem `json:"items"`
	GeneratedAt     time.Time       `json:"generatedAt"`
}

// CustomerReport represents a customer financial report
type CustomerReport struct {
	CustomerID       string        `json:"customerId"`
	CustomerName     string        `json:"customerName"`
	TotalSales       float64       `json:"totalSales"`
	TotalPayments    float64       `json:"totalPayments"`
	OutstandingBalance float64     `json:"outstandingBalance"`
	CrateBalance     int           `json:"crateBalance"`
	Transactions     []Transaction `json:"transactions"`
	GeneratedAt      time.Time     `json:"generatedAt"`
}

// Handler for generating reports
func generateReportHandler(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	ctx := r.Context()
	
	switch req.Type {
	case ReportTypeSales:
		report, err := generateSalesReport(ctx, req.StartDate, req.EndDate)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		
	case ReportTypeInventory:
		report, err := generateInventoryReport(ctx)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		
	case ReportTypeCustomer:
		if req.CustomerID == "" {
			http.Error(w, "customerId is required for customer reports", 400)
			return
		}
		report, err := generateCustomerReport(ctx, req.CustomerID, req.StartDate, req.EndDate)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		
	default:
		http.Error(w, "unsupported report type", 400)
	}
}

// generateSalesReport generates a sales report
func generateSalesReport(ctx interface{}, startDate, endDate string) (*SalesReport, error) {
	report := &SalesReport{
		GeneratedAt: time.Now(),
	}
	
	// Build query with date filters
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount, details, created_at 
		FROM transactions WHERE type='sale'`
	
	args := []interface{}{}
	if startDate != "" {
		query += ` AND date >= $1`
		args = append(args, startDate)
	}
	if endDate != "" {
		if len(args) > 0 {
			query += ` AND date <= $2`
		} else {
			query += ` AND date <= $1`
		}
		args = append(args, endDate)
	}
	
	query += ` ORDER BY date DESC`
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	transactions := []Transaction{}
	totalSales := 0.0
	
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount, &tx.Details, &tx.CreatedAt); err != nil {
			continue
		}
		transactions = append(transactions, tx)
		totalSales += tx.TotalAmount
	}
	
	report.Transactions = transactions
	report.TotalSales = totalSales
	report.TotalItems = len(transactions)
	
	// Get top selling items
	topItems, err := getTopSellingItems(ctx, startDate, endDate)
	if err == nil {
		report.TopItems = topItems
	}
	
	return report, nil
}

// getTopSellingItems gets top selling items
func getTopSellingItems(ctx interface{}, startDate, endDate string) ([]TopItem, error) {
	query := `
		SELECT si.item_name, SUM(si.quantity) as total_qty, SUM(si.total) as total_revenue
		FROM sale_items si
		JOIN transactions t ON si.transaction_id = t.id
		WHERE 1=1`
	
	args := []interface{}{}
	if startDate != "" {
		query += ` AND t.date >= $1`
		args = append(args, startDate)
	}
	if endDate != "" {
		if len(args) > 0 {
			query += ` AND t.date <= $2`
		} else {
			query += ` AND t.date <= $1`
		}
		args = append(args, endDate)
	}
	
	query += ` GROUP BY si.item_name ORDER BY total_revenue DESC LIMIT 10`
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	topItems := []TopItem{}
	for rows.Next() {
		var item TopItem
		if err := rows.Scan(&item.ItemName, &item.TotalQuantity, &item.TotalRevenue); err != nil {
			continue
		}
		topItems = append(topItems, item)
	}
	
	return topItems, nil
}

// generateInventoryReport generates an inventory report
func generateInventoryReport(ctx interface{}) (*InventoryReport, error) {
	report := &InventoryReport{
		GeneratedAt: time.Now(),
	}
	
	// Get all inventory items
	rows, err := db.Query(`
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, created_at, updated_at 
		FROM inventory_items 
		ORDER BY expiry_date ASC`)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	items := []InventoryItem{}
	expiringSoon := 0
	expired := 0
	
	for rows.Next() {
		var item InventoryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Variant, &item.LotNumber, &item.Quantity, &item.Unit, &item.PurchaseDate, &item.ExpiryDate, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		items = append(items, item)
		
		// Calculate status
		expiryDate, _ := time.Parse("2006-01-02", item.ExpiryDate)
		daysUntilExpiry := int(time.Until(expiryDate).Hours() / 24)
		
		if daysUntilExpiry < 0 {
			expired++
		} else if daysUntilExpiry <= 3 {
			expiringSoon++
		}
	}
	
	report.Items = items
	report.TotalItems = len(items)
	report.ExpiringSoon = expiringSoon
	report.Expired = expired
	
	return report, nil
}

// generateCustomerReport generates a customer-specific report
func generateCustomerReport(ctx interface{}, customerID, startDate, endDate string) (*CustomerReport, error) {
	report := &CustomerReport{
		CustomerID:  customerID,
		GeneratedAt: time.Now(),
	}
	
	// Get customer name
	err := db.QueryRow(`SELECT name FROM customers WHERE id=$1`, customerID).Scan(&report.CustomerName)
	if err != nil {
		return nil, err
	}
	
	// Get transactions
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount, details, created_at 
		FROM transactions WHERE customer_id=$1`
	
	args := []interface{}{customerID}
	argIdx := 2
	
	if startDate != "" {
		query += ` AND date >= $` + string(rune('0'+argIdx))
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		query += ` AND date <= $` + string(rune('0'+argIdx))
		args = append(args, endDate)
	}
	
	query += ` ORDER BY date DESC`
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	transactions := []Transaction{}
	totalSales := 0.0
	totalPayments := 0.0
	
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount, &tx.Details, &tx.CreatedAt); err != nil {
			continue
		}
		transactions = append(transactions, tx)
		
		if tx.Type == "sale" {
			totalSales += tx.TotalAmount
		} else if tx.Type == "payment" && tx.PaymentAmount != nil {
			totalPayments += *tx.PaymentAmount
		}
	}
	
	report.Transactions = transactions
	report.TotalSales = totalSales
	report.TotalPayments = totalPayments
	report.OutstandingBalance = totalSales - totalPayments
	
	// Get crate balance
	var crateBalance int
	err = db.QueryRow(`
		SELECT COALESCE(balance, 0) 
		FROM crate_ledger 
		WHERE customer_id=$1 
		ORDER BY date DESC 
		LIMIT 1`, customerID).Scan(&crateBalance)
	
	if err == nil {
		report.CrateBalance = crateBalance
	}
	
	return report, nil
}
