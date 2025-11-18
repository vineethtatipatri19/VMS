-- Migration: Enhance all entities with comprehensive business fields
-- Date: 2025-11-16

-- ==========================================
-- CUSTOMERS - Enhanced for Credit Management
-- ==========================================

ALTER TABLE customers ADD COLUMN IF NOT EXISTS email TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS gstin TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customer_type TEXT CHECK (customer_type IN ('b2b', 'b2c', 'retail', 'wholesale')) DEFAULT 'retail';
ALTER TABLE customers ADD COLUMN IF NOT EXISTS credit_limit NUMERIC DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS current_balance NUMERIC DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS payment_terms_days INTEGER DEFAULT 30;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS interest_rate NUMERIC DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS status TEXT CHECK (status IN ('active', 'inactive', 'blocked')) DEFAULT 'active';
ALTER TABLE customers ADD COLUMN IF NOT EXISTS kyc_document_type TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS kyc_document_number TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS business_name TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tags TEXT[];
ALTER TABLE customers ADD COLUMN IF NOT EXISTS alternate_contact TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS whatsapp_number TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS last_transaction_date TIMESTAMPTZ;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS total_purchases NUMERIC DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_points INTEGER DEFAULT 0;

-- Index for faster customer searches
CREATE INDEX IF NOT EXISTS idx_customers_contact ON customers(contact_number);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status);
CREATE INDEX IF NOT EXISTS idx_customers_balance ON customers(current_balance);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);

-- ==========================================
-- INVENTORY - Enhanced for Profit & Supplier Tracking
-- ==========================================

ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS category TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS sub_category TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS cost_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS selling_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS margin_percentage NUMERIC GENERATED ALWAYS AS (
  CASE 
    WHEN selling_price > 0 THEN ((selling_price - COALESCE(cost_price, 0)) / selling_price * 100)
    ELSE 0
  END
) STORED;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS supplier_id UUID;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS supplier_name TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS purchase_invoice TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS min_stock_level NUMERIC DEFAULT 0;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS reorder_point NUMERIC DEFAULT 0;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS shelf_life_days INTEGER;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS storage_location TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS barcode TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS sku TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS hsn_code TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS gst_rate NUMERIC DEFAULT 0;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS status TEXT CHECK (status IN ('available', 'low_stock', 'out_of_stock', 'expired', 'damaged')) DEFAULT 'available';
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS is_perishable BOOLEAN DEFAULT true;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS weight_per_unit NUMERIC;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS packaging_type TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS image_url TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS total_sold NUMERIC DEFAULT 0;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS total_wasted NUMERIC DEFAULT 0;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS last_restock_date TIMESTAMPTZ;

-- Auto-update status based on quantity
CREATE OR REPLACE FUNCTION update_inventory_status()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.expiry_date < CURRENT_DATE THEN
    NEW.status = 'expired';
  ELSIF NEW.quantity = 0 THEN
    NEW.status = 'out_of_stock';
  ELSIF NEW.quantity <= NEW.min_stock_level THEN
    NEW.status = 'low_stock';
  ELSE
    NEW.status = 'available';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER inventory_status_trigger
BEFORE INSERT OR UPDATE ON inventory_items
FOR EACH ROW
EXECUTE FUNCTION update_inventory_status();

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_inventory_expiry ON inventory_items(expiry_date);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory_items(status);
CREATE INDEX IF NOT EXISTS idx_inventory_category ON inventory_items(category);
CREATE INDEX IF NOT EXISTS idx_inventory_barcode ON inventory_items(barcode);

-- ==========================================
-- TRANSACTIONS - Enhanced for Payment Tracking
-- ==========================================

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_method TEXT CHECK (payment_method IN ('cash', 'upi', 'card', 'bank_transfer', 'cheque', 'credit')) DEFAULT 'cash';
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_reference TEXT;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS due_date DATE;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS is_overdue BOOLEAN DEFAULT false;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS discount_amount NUMERIC DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS tax_amount NUMERIC DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS invoice_number TEXT;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS receipt_sent BOOLEAN DEFAULT false;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS balance_after NUMERIC;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS sale_type TEXT CHECK (sale_type IN ('regular', 'wholesale', 'credit', 'return'));
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS delivery_status TEXT CHECK (delivery_status IN ('pending', 'packed', 'dispatched', 'delivered', 'cancelled'));
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS delivery_date TIMESTAMPTZ;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS delivery_address TEXT;

-- Auto-generate invoice number
CREATE SEQUENCE IF NOT EXISTS invoice_number_seq START 1000;

-- Function to generate invoice number
CREATE OR REPLACE FUNCTION generate_invoice_number()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.type = 'sale' AND NEW.invoice_number IS NULL THEN
    NEW.invoice_number = 'INV-' || TO_CHAR(NEW.date, 'YYYYMMDD') || '-' || LPAD(nextval('invoice_number_seq')::TEXT, 4, '0');
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_invoice_trigger
BEFORE INSERT ON transactions
FOR EACH ROW
EXECUTE FUNCTION generate_invoice_number();

-- Check overdue payments
CREATE OR REPLACE FUNCTION check_overdue_payments()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.due_date IS NOT NULL AND NEW.due_date < CURRENT_DATE AND NEW.type = 'sale' THEN
    NEW.is_overdue = true;
  ELSE
    NEW.is_overdue = false;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_overdue_trigger
BEFORE INSERT OR UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION check_overdue_payments();

-- Indexes
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_customer ON transactions(customer_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_overdue ON transactions(is_overdue) WHERE is_overdue = true;
CREATE INDEX IF NOT EXISTS idx_transactions_invoice ON transactions(invoice_number);

-- ==========================================
-- SALE ITEMS - Enhanced for Profit Tracking
-- ==========================================

ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS cost_per_unit NUMERIC;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS profit NUMERIC GENERATED ALWAYS AS (
  (price_per_unit - COALESCE(cost_per_unit, 0)) * quantity
) STORED;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS discount_percentage NUMERIC DEFAULT 0;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS tax_percentage NUMERIC DEFAULT 0;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS hsn_code TEXT;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS batch_number TEXT;
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS expiry_date DATE;

-- ==========================================
-- CRATE LEDGER - Enhanced
-- ==========================================

ALTER TABLE crate_ledger ADD COLUMN IF NOT EXISTS crate_type TEXT DEFAULT 'standard';
ALTER TABLE crate_ledger ADD COLUMN IF NOT EXISTS crate_value NUMERIC DEFAULT 50;
ALTER TABLE crate_ledger ADD COLUMN IF NOT EXISTS transaction_id UUID REFERENCES transactions(id);
ALTER TABLE crate_ledger ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_crates_customer ON crate_ledger(customer_id);
CREATE INDEX IF NOT EXISTS idx_crates_balance ON crate_ledger(balance) WHERE balance > 0;

-- ==========================================
-- NEW TABLE: WASTAGE LOG
-- ==========================================

CREATE TABLE IF NOT EXISTS wastage_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id) ON DELETE SET NULL,
  item_name TEXT NOT NULL,
  quantity NUMERIC NOT NULL,
  unit TEXT NOT NULL,
  reason TEXT CHECK (reason IN ('expired', 'damaged', 'spoiled', 'pest', 'other')) NOT NULL,
  reason_details TEXT,
  cost_value NUMERIC,
  logged_by TEXT,
  logged_at TIMESTAMPTZ DEFAULT now(),
  photo_url TEXT
);

CREATE INDEX IF NOT EXISTS idx_wastage_date ON wastage_log(logged_at DESC);
CREATE INDEX IF NOT EXISTS idx_wastage_reason ON wastage_log(reason);

-- ==========================================
-- NEW TABLE: EXPIRY ALERTS
-- ==========================================

CREATE TABLE IF NOT EXISTS expiry_alerts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id) ON DELETE CASCADE,
  alert_date DATE NOT NULL,
  expiry_date DATE NOT NULL,
  days_until_expiry INTEGER NOT NULL,
  acknowledged BOOLEAN DEFAULT false,
  acknowledged_at TIMESTAMPTZ,
  acknowledged_by TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_expiry_alerts_date ON expiry_alerts(alert_date);
CREATE INDEX IF NOT EXISTS idx_expiry_alerts_ack ON expiry_alerts(acknowledged);

-- ==========================================
-- NEW TABLE: PAYMENT SCHEDULES
-- ==========================================

CREATE TABLE IF NOT EXISTS payment_schedules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
  customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
  installment_number INTEGER NOT NULL,
  due_date DATE NOT NULL,
  amount_due NUMERIC NOT NULL,
  amount_paid NUMERIC DEFAULT 0,
  status TEXT CHECK (status IN ('pending', 'partial', 'paid', 'overdue')) DEFAULT 'pending',
  payment_date TIMESTAMPTZ,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payment_schedule_customer ON payment_schedules(customer_id);
CREATE INDEX IF NOT EXISTS idx_payment_schedule_status ON payment_schedules(status);
CREATE INDEX IF NOT EXISTS idx_payment_schedule_due ON payment_schedules(due_date);

-- ==========================================
-- NEW TABLE: PRICING TIERS (for bulk pricing)
-- ==========================================

CREATE TABLE IF NOT EXISTS pricing_tiers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id) ON DELETE CASCADE,
  min_quantity NUMERIC NOT NULL,
  max_quantity NUMERIC,
  price_per_unit NUMERIC NOT NULL,
  tier_name TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- ==========================================
-- NEW TABLE: PRICE HISTORY
-- ==========================================

CREATE TABLE IF NOT EXISTS price_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id) ON DELETE CASCADE,
  old_price NUMERIC,
  new_price NUMERIC NOT NULL,
  changed_by TEXT,
  reason TEXT,
  changed_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_price_history_item ON price_history(inventory_item_id, changed_at DESC);

-- ==========================================
-- VIEWS FOR REPORTING
-- ==========================================

-- Customer Summary View
CREATE OR REPLACE VIEW customer_summary AS
SELECT 
  c.id,
  c.name,
  c.contact_number,
  c.customer_type,
  c.current_balance,
  c.credit_limit,
  c.status,
  COUNT(DISTINCT t.id) FILTER (WHERE t.type = 'sale') as total_transactions,
  COALESCE(SUM(t.total_amount) FILTER (WHERE t.type = 'sale'), 0) as total_sales,
  COALESCE(SUM(t.payment_amount) FILTER (WHERE t.type = 'payment'), 0) as total_payments,
  MAX(t.date) as last_transaction_date,
  COALESCE(cl.balance, 0) as crate_balance
FROM customers c
LEFT JOIN transactions t ON c.id = t.customer_id
LEFT JOIN LATERAL (
  SELECT balance 
  FROM crate_ledger 
  WHERE customer_id = c.id 
  ORDER BY date DESC 
  LIMIT 1
) cl ON true
GROUP BY c.id, c.name, c.contact_number, c.customer_type, c.current_balance, 
         c.credit_limit, c.status, cl.balance;

-- Inventory Summary View
CREATE OR REPLACE VIEW inventory_summary AS
SELECT 
  i.id,
  i.name,
  i.category,
  i.quantity,
  i.unit,
  i.cost_price,
  i.selling_price,
  i.margin_percentage,
  i.expiry_date,
  i.status,
  CASE 
    WHEN i.expiry_date < CURRENT_DATE THEN 'expired'
    WHEN i.expiry_date <= CURRENT_DATE + INTERVAL '3 days' THEN 'expiring_soon'
    ELSE 'fresh'
  END as expiry_status,
  i.expiry_date - CURRENT_DATE as days_until_expiry,
  i.total_sold,
  i.total_wasted,
  COALESCE(i.cost_price * i.quantity, 0) as total_value
FROM inventory_items i;

-- Sales Performance View
CREATE OR REPLACE VIEW sales_performance AS
SELECT 
  DATE(t.date) as sale_date,
  COUNT(DISTINCT t.id) as transaction_count,
  COUNT(DISTINCT t.customer_id) as unique_customers,
  SUM(t.total_amount) as total_sales,
  SUM(t.tax_amount) as total_tax,
  SUM(t.discount_amount) as total_discounts,
  AVG(t.total_amount) as avg_transaction_value,
  SUM(si.profit) as total_profit,
  SUM(si.profit) / NULLIF(SUM(si.cost_per_unit * si.quantity), 0) * 100 as profit_margin_percentage
FROM transactions t
LEFT JOIN sale_items si ON t.id = si.transaction_id
WHERE t.type = 'sale'
GROUP BY DATE(t.date)
ORDER BY sale_date DESC;

-- Overdue Payments View
CREATE OR REPLACE VIEW overdue_payments AS
SELECT 
  t.id as transaction_id,
  t.invoice_number,
  t.customer_id,
  c.name as customer_name,
  c.contact_number,
  c.whatsapp_number,
  t.date as transaction_date,
  t.due_date,
  CURRENT_DATE - t.due_date as days_overdue,
  t.total_amount,
  COALESCE(SUM(t2.payment_amount) FILTER (WHERE t2.type = 'payment'), 0) as amount_paid,
  t.total_amount - COALESCE(SUM(t2.payment_amount) FILTER (WHERE t2.type = 'payment'), 0) as balance_due
FROM transactions t
JOIN customers c ON t.customer_id = c.id
LEFT JOIN transactions t2 ON t2.customer_id = c.id AND t2.date >= t.date AND t2.type = 'payment'
WHERE t.type = 'sale' 
  AND t.is_overdue = true
  AND t.total_amount > COALESCE((SELECT SUM(payment_amount) FROM transactions WHERE customer_id = t.customer_id AND type = 'payment' AND date >= t.date), 0)
GROUP BY t.id, t.invoice_number, t.customer_id, c.name, c.contact_number, 
         c.whatsapp_number, t.date, t.due_date, t.total_amount
ORDER BY days_overdue DESC;

-- Wastage Summary View
CREATE OR REPLACE VIEW wastage_summary AS
SELECT 
  DATE(w.logged_at) as wastage_date,
  w.reason,
  COUNT(*) as item_count,
  SUM(w.quantity) as total_quantity,
  SUM(w.cost_value) as total_cost,
  STRING_AGG(DISTINCT w.item_name, ', ') as items_wasted
FROM wastage_log w
GROUP BY DATE(w.logged_at), w.reason
ORDER BY wastage_date DESC;

-- Comment on tables
COMMENT ON TABLE customers IS 'Enhanced customer table with credit management, KYC, and business details';
COMMENT ON TABLE inventory_items IS 'Enhanced inventory with cost/pricing, supplier tracking, and profit margins';
COMMENT ON TABLE transactions IS 'Enhanced transactions with payment methods, invoice generation, and delivery tracking';
COMMENT ON TABLE wastage_log IS 'Track all wasted/spoiled inventory items';
COMMENT ON TABLE expiry_alerts IS 'Alert system for items approaching expiry';
COMMENT ON TABLE payment_schedules IS 'Installment/payment plan tracking for credit sales';
COMMENT ON TABLE pricing_tiers IS 'Bulk pricing and tiered pricing support';
COMMENT ON TABLE price_history IS 'Audit trail of all price changes';
