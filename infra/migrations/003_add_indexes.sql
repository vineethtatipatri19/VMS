-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_inventory_expiry ON inventory_items(expiry_date);
CREATE INDEX IF NOT EXISTS idx_inventory_name ON inventory_items(name);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);
CREATE INDEX IF NOT EXISTS idx_transactions_customer ON transactions(customer_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_sale_items_transaction ON sale_items(transaction_id);
CREATE INDEX IF NOT EXISTS idx_crate_ledger_customer ON crate_ledger(customer_id);
CREATE INDEX IF NOT EXISTS idx_crate_ledger_date ON crate_ledger(date);
