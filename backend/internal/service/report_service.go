package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/example/pgvms/internal/domain"
)

// ReportService handles report generation
type ReportService struct {
	db *sql.DB
}

// NewReportService creates a new report service
func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

// ReportType represents different report types
type ReportType string

const (
	ReportTypeSales     ReportType = "sales"
	ReportTypeInventory ReportType = "inventory"
	ReportTypeCustomer  ReportType = "customer"
)

// TopItem represents top selling items
type TopItem struct {
	ItemName      string  `json:"itemName"`
	TotalQuantity float64 `json:"totalQuantity"`
	TotalRevenue  float64 `json:"totalRevenue"`
}

// SalesReport represents a sales report
type SalesReport struct {
	TotalSales   float64               `json:"totalSales"`
	TotalItems   int                   `json:"totalItems"`
	Transactions []domain.Transaction  `json:"transactions"`
	TopItems     []TopItem             `json:"topItems"`
	GeneratedAt  time.Time             `json:"generatedAt"`
}

// InventoryReport represents an inventory report
type InventoryReport struct {
	TotalItems   int                    `json:"totalItems"`
	ExpiringSoon int                    `json:"expiringSoon"`
	Expired      int                    `json:"expired"`
	Items        []domain.InventoryItem `json:"items"`
	GeneratedAt  time.Time              `json:"generatedAt"`
}

// CustomerReport represents a customer financial report
type CustomerReport struct {
	CustomerID         string               `json:"customerId"`
	CustomerName       string               `json:"customerName"`
	TotalSales         float64              `json:"totalSales"`
	TotalPayments      float64              `json:"totalPayments"`
	OutstandingBalance float64              `json:"outstandingBalance"`
	CrateBalance       int                  `json:"crateBalance"`
	Transactions       []domain.Transaction `json:"transactions"`
	GeneratedAt        time.Time            `json:"generatedAt"`
}

// GenerateSalesReport generates a sales report
func (s *ReportService) GenerateSalesReport(ctx context.Context, startDate, endDate string) (*SalesReport, error) {
	report := &SalesReport{
		GeneratedAt: time.Now(),
	}

	// Build query with date filters
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount, details, created_at 
		FROM transactions WHERE type='sale' AND deleted_at IS NULL`

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

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []domain.Transaction{}
	totalSales := 0.0

	for rows.Next() {
		var tx domain.Transaction
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
	topItems, err := s.getTopSellingItems(ctx, startDate, endDate)
	if err == nil {
		report.TopItems = topItems
	}

	return report, nil
}

// getTopSellingItems gets top selling items
func (s *ReportService) getTopSellingItems(ctx context.Context, startDate, endDate string) ([]TopItem, error) {
	query := `
		SELECT si.item_name, SUM(si.quantity) as total_qty, SUM(si.total) as total_revenue
		FROM sale_items si
		JOIN transactions t ON si.transaction_id = t.id
		WHERE t.deleted_at IS NULL AND si.deleted_at IS NULL`

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

	rows, err := s.db.QueryContext(ctx, query, args...)
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

// GenerateInventoryReport generates an inventory report
func (s *ReportService) GenerateInventoryReport(ctx context.Context) (*InventoryReport, error) {
	report := &InventoryReport{
		GeneratedAt: time.Now(),
	}

	// Get all inventory items
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, created_at, updated_at 
		FROM inventory_items 
		WHERE deleted_at IS NULL
		ORDER BY expiry_date ASC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.InventoryItem{}
	expiringSoon := 0
	expired := 0

	for rows.Next() {
		var item domain.InventoryItem
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

// GenerateCustomerReport generates a customer-specific report
func (s *ReportService) GenerateCustomerReport(ctx context.Context, customerID, startDate, endDate string) (*CustomerReport, error) {
	report := &CustomerReport{
		CustomerID:  customerID,
		GeneratedAt: time.Now(),
	}

	// Get customer name
	err := s.db.QueryRowContext(ctx, `SELECT name FROM customers WHERE id=$1 AND deleted_at IS NULL`, customerID).Scan(&report.CustomerName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Get transactions
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount, details, created_at 
		FROM transactions WHERE customer_id=$1 AND deleted_at IS NULL`

	args := []interface{}{customerID}
	argIdx := 2

	if startDate != "" {
		query += ` AND date >= $2`
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		if len(args) > 1 {
			query += ` AND date <= $3`
		} else {
			query += ` AND date <= $2`
		}
		args = append(args, endDate)
	}

	query += ` ORDER BY date DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []domain.Transaction{}
	totalSales := 0.0
	totalPayments := 0.0

	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount, &tx.Details, &tx.CreatedAt); err != nil {
			continue
		}
		transactions = append(transactions, tx)

		if tx.Type == "sale" {
			totalSales += tx.TotalAmount
		} else if tx.Type == "payment" && tx.PaymentAmount > 0 {
			totalPayments += tx.PaymentAmount
		}
	}

	report.Transactions = transactions
	report.TotalSales = totalSales
	report.TotalPayments = totalPayments
	report.OutstandingBalance = totalSales - totalPayments

	// Get crate balance
	var crateBalance int
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(balance, 0) 
		FROM crate_ledger 
		WHERE customer_id=$1 AND deleted_at IS NULL
		ORDER BY date DESC 
		LIMIT 1`, customerID).Scan(&crateBalance)

	if err == nil {
		report.CrateBalance = crateBalance
	}

	return report, nil
}
