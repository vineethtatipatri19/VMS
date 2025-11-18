package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/repository"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (
		id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, status, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.CustomerID, tx.Date, tx.Type, tx.PaymentAmount, tx.TotalAmount,
		toNullString(tx.PaymentMethod), toNullString(tx.PaymentRef),
		toNullString(tx.Notes), toNullString(tx.Status), tx.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, status, created_at
	FROM transactions 
	WHERE id = $1 AND deleted_at IS NULL`

	var tx domain.Transaction
	var paymentMethod, paymentRef, notes, status sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
		&paymentMethod, &paymentRef, &notes, &status, &tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	tx.PaymentMethod = fromNullString(paymentMethod)
	tx.PaymentRef = fromNullString(paymentRef)
	tx.Notes = fromNullString(notes)
	tx.Status = fromNullString(status)

	return &tx, nil
}

func (r *transactionRepository) ListByCustomer(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	query := `
		SELECT id, customer_id, date, type, payment_amount, total_amount,
			payment_method, payment_reference, notes, status, invoice_number,
			sale_type, delivery_status, delivery_date, delivery_address,
			discount_amount, tax_amount, balance_after, receipt_sent,
			due_date, is_overdue, created_at
		FROM transactions
		WHERE customer_id = $1 AND deleted_at IS NULL
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		var tx domain.Transaction
		var paymentMethod, paymentRef, notes, status, invoiceNumber, saleType, deliveryStatus, deliveryAddress sql.NullString
		var dueDate, deliveryDate sql.NullTime

		if err := rows.Scan(
			&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
			&paymentMethod, &paymentRef, &notes, &status, &invoiceNumber,
			&saleType, &deliveryStatus, &deliveryDate, &deliveryAddress,
			&tx.DiscountAmount, &tx.TaxAmount, &tx.BalanceAfter, &tx.ReceiptSent,
			&dueDate, &tx.IsOverdue, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		tx.PaymentMethod = fromNullString(paymentMethod)
		tx.PaymentRef = fromNullString(paymentRef)
		tx.Notes = fromNullString(notes)
		tx.Status = fromNullString(status)
		tx.InvoiceNumber = fromNullString(invoiceNumber)
		tx.SaleType = fromNullString(saleType)
		tx.DeliveryStatus = fromNullString(deliveryStatus)
		tx.DeliveryAddress = fromNullString(deliveryAddress)
		if dueDate.Valid {
			tx.DueDate = &dueDate.Time
		}
		if deliveryDate.Valid {
			tx.DeliveryDate = &deliveryDate.Time
		}

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (r *transactionRepository) List(ctx context.Context, txType string, startDate, endDate time.Time) ([]*domain.Transaction, error) {
	query := `SELECT id, customer_id, date, type, payment_amount, total_amount,
		payment_method, payment_reference, notes, status, invoice_number,
		sale_type, delivery_status, delivery_date, delivery_address,
		discount_amount, tax_amount, balance_after, receipt_sent,
		due_date, is_overdue, created_at
	FROM transactions 
	WHERE deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	if txType != "" {
		query += fmt.Sprintf(` AND type = $%d`, argCount)
		args = append(args, txType)
		argCount++
	}

	if !startDate.IsZero() {
		query += fmt.Sprintf(` AND date >= $%d`, argCount)
		args = append(args, startDate)
		argCount++
	}

	if !endDate.IsZero() {
		query += fmt.Sprintf(` AND date <= $%d`, argCount)
		args = append(args, endDate)
		argCount++
	}

	query += ` ORDER BY date DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		var tx domain.Transaction
		var paymentMethod, paymentRef, notes, status, invoiceNumber, saleType, deliveryStatus, deliveryAddress sql.NullString
		var dueDate, deliveryDate sql.NullTime

		if err := rows.Scan(
			&tx.ID, &tx.CustomerID, &tx.Date, &tx.Type, &tx.PaymentAmount, &tx.TotalAmount,
			&paymentMethod, &paymentRef, &notes, &status, &invoiceNumber,
			&saleType, &deliveryStatus, &deliveryDate, &deliveryAddress,
			&tx.DiscountAmount, &tx.TaxAmount, &tx.BalanceAfter, &tx.ReceiptSent,
			&dueDate, &tx.IsOverdue, &tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		tx.PaymentMethod = fromNullString(paymentMethod)
		tx.PaymentRef = fromNullString(paymentRef)
		tx.Notes = fromNullString(notes)
		tx.Status = fromNullString(status)
		tx.InvoiceNumber = fromNullString(invoiceNumber)
		tx.SaleType = fromNullString(saleType)
		tx.DeliveryStatus = fromNullString(deliveryStatus)
		tx.DeliveryAddress = fromNullString(deliveryAddress)
		if dueDate.Valid {
			tx.DueDate = &dueDate.Time
		}
		if deliveryDate.Valid {
			tx.DeliveryDate = &deliveryDate.Time
		}

		transactions = append(transactions, &tx)
	}

	// Populate sale items for all transactions
	if err := r.populateSaleItems(ctx, transactions); err != nil {
		return nil, fmt.Errorf("failed to populate sale items: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `UPDATE transactions SET
		payment_amount = $2, total_amount = $3, payment_method = $4,
		payment_reference = $5, notes = $6
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.PaymentAmount, tx.TotalAmount,
		toNullString(tx.PaymentMethod), toNullString(tx.PaymentRef),
		toNullString(tx.Notes),
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id string, req *domain.DeleteRequest) error {
	query := `UPDATE transactions 
	SET deleted_at = $2, deleted_by = $3, deletion_reason = $4
	WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now(), "system", req.Reason)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// populateSaleItems fetches and populates sale items for transactions
func (r *transactionRepository) populateSaleItems(ctx context.Context, transactions []*domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	// Build list of transaction IDs
	txIDs := make([]string, len(transactions))
	txMap := make(map[string]*domain.Transaction)
	for i, tx := range transactions {
		txIDs[i] = tx.ID
		txMap[tx.ID] = tx
		// Initialize Details map
		tx.Details = make(map[string]interface{})
	}

	// Query sale items for all transactions
	query := `SELECT transaction_id, item_name, quantity, unit, price_per_unit, total
		FROM sale_items
		WHERE transaction_id = ANY($1)
		ORDER BY transaction_id, item_name`

	rows, err := r.db.QueryContext(ctx, query, txIDs)
	if err != nil {
		return fmt.Errorf("failed to query sale items: %w", err)
	}
	defer rows.Close()

	// Group items by transaction
	itemsByTx := make(map[string][]map[string]interface{})
	for rows.Next() {
		var txID, itemName, unit string
		var quantity, pricePerUnit, total float64

		if err := rows.Scan(&txID, &itemName, &quantity, &unit, &pricePerUnit, &total); err != nil {
			return fmt.Errorf("failed to scan sale item: %w", err)
		}

		item := map[string]interface{}{
			"item_name":      itemName,
			"ItemName":       itemName, // Add both formats for compatibility
			"quantity":       quantity,
			"Quantity":       quantity,
			"unit":           unit,
			"price_per_unit": pricePerUnit,
			"total":          total,
		}

		itemsByTx[txID] = append(itemsByTx[txID], item)
	}

	// Populate Details for each transaction
	for txID, items := range itemsByTx {
		if tx, ok := txMap[txID]; ok {
			tx.Details["items"] = items
		}
	}

	return nil
}
