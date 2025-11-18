-- Migration: Add soft delete support to all main entities
-- Date: 2025-11-18

-- Add soft delete columns to customers
ALTER TABLE customers 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to inventory_items
ALTER TABLE inventory_items 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to transactions
ALTER TABLE transactions 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to sale_items
ALTER TABLE sale_items 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to crate_ledger
ALTER TABLE crate_ledger 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to wastage_log
ALTER TABLE wastage_log 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to expiry_alerts
ALTER TABLE expiry_alerts 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to payment_schedules
ALTER TABLE payment_schedules 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to pricing_tiers
ALTER TABLE pricing_tiers 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Add soft delete columns to price_history
ALTER TABLE price_history 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS deleted_by TEXT,
  ADD COLUMN IF NOT EXISTS deletion_reason TEXT;

-- Create indexes for soft delete queries (only on main tables for performance)
CREATE INDEX IF NOT EXISTS idx_customers_deleted ON customers(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_deleted ON inventory_items(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_deleted ON transactions(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sale_items_deleted ON sale_items(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crate_ledger_deleted ON crate_ledger(deleted_at) WHERE deleted_at IS NULL;

COMMENT ON COLUMN customers.deleted_at IS 'Soft delete timestamp - NULL means record is active';
COMMENT ON COLUMN inventory_items.deleted_at IS 'Soft delete timestamp - NULL means record is active';
COMMENT ON COLUMN transactions.deleted_at IS 'Soft delete timestamp - NULL means record is active';
