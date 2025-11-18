-- Add related data (crates, wastage, alerts, sale items)
-- This assumes customers, inventory, and transactions already exist

BEGIN;

-- Get some customer and inventory IDs to use
DO $$ 
DECLARE
    cust1_id UUID;
    cust2_id UUID;
    cust3_id UUID;
    inv1_id UUID;
    inv2_id UUID;
    inv3_id UUID;
    tx1_id UUID;
    tx2_id UUID;
BEGIN
    -- Get customer IDs
    SELECT id INTO cust1_id FROM customers WHERE name LIKE '%Metro%' LIMIT 1;
    SELECT id INTO cust2_id FROM customers WHERE name LIKE '%Fresh Mart%' LIMIT 1;
    SELECT id INTO cust3_id FROM customers WHERE name LIKE '%Green%' LIMIT 1;
    
    -- Get inventory IDs
    SELECT id INTO inv1_id FROM inventory_items WHERE name = 'Tomato' AND variant = 'Local' LIMIT 1;
    SELECT id INTO inv2_id FROM inventory_items WHERE name = 'Lettuce' LIMIT 1;
    SELECT id INTO inv3_id FROM inventory_items WHERE name = 'Onion' LIMIT 1;
    
    -- Get transaction IDs
    SELECT id INTO tx1_id FROM transactions WHERE customer_id = cust1_id ORDER BY date DESC LIMIT 1;
    SELECT id INTO tx2_id FROM transactions WHERE customer_id = cust2_id ORDER BY date DESC LIMIT 1;

    -- Insert Crate Ledger entries (note: columns are crates_issued and crates_returned)
    IF cust1_id IS NOT NULL AND tx1_id IS NOT NULL THEN
        INSERT INTO crate_ledger (customer_id, transaction_id, crates_issued, crates_returned, balance, notes, date) VALUES
        (cust1_id, tx1_id, 25, 0, 25, 'Large order - 25 crates delivered', NOW()),
        (cust1_id, NULL, 0, 20, 5, 'Returned 20 crates from previous delivery', NOW() - INTERVAL '3 days');
    END IF;

    IF cust2_id IS NOT NULL AND tx2_id IS NOT NULL THEN
        INSERT INTO crate_ledger (customer_id, transaction_id, crates_issued, crates_returned, balance, notes, date) VALUES
        (cust2_id, tx2_id, 12, 0, 12, '12 crates for mixed vegetable order', NOW() - INTERVAL '1 day'),
        (cust2_id, NULL, 0, 8, 4, 'Partial crate return', NOW() - INTERVAL '5 days');
    END IF;

    IF cust3_id IS NOT NULL THEN
        INSERT INTO crate_ledger (customer_id, crates_issued, crates_returned, balance, notes, date) VALUES
        (cust3_id, 15, 0, 15, '15 crates with organic vegetables', NOW() - INTERVAL '2 days'),
        (cust3_id, 0, 10, 5, 'Returned crates from last week', NOW() - INTERVAL '7 days');
    END IF;

    -- Insert Wastage Log entries (columns: inventory_item_id, item_name, quantity, unit, reason, cost_value, logged_by)
    IF inv2_id IS NOT NULL THEN
        INSERT INTO wastage_log (inventory_item_id, item_name, quantity, unit, reason, reason_details, cost_value, logged_by, logged_at) VALUES
        (inv2_id, 'Lettuce Iceberg', 8, 'kg', 'expired', 'Could not sell in time due to slow demand', 560, 'Manager User', NOW() - INTERVAL '1 day');
    END IF;

    IF inv1_id IS NOT NULL THEN
        INSERT INTO wastage_log (inventory_item_id, item_name, quantity, unit, reason, reason_details, cost_value, logged_by, logged_at) VALUES
        (inv1_id, 'Tomato Local', 12, 'kg', 'spoiled', 'Started rotting in cold storage', 600, 'Sales User', NOW() - INTERVAL '3 days');
    END IF;

    -- Insert more wastage entries
    INSERT INTO wastage_log (inventory_item_id, item_name, quantity, unit, reason, reason_details, cost_value, logged_by, logged_at)
    SELECT id, name || ' ' || COALESCE(variant, ''), 5, unit, 'damaged', 'Damaged during transport - leakage', 250, 'Admin User', NOW() - INTERVAL '5 days'
    FROM inventory_items WHERE name = 'Cooking Oil' LIMIT 1;

    INSERT INTO wastage_log (inventory_item_id, item_name, quantity, unit, reason, reason_details, cost_value, logged_by, logged_at)
    SELECT id, name || ' ' || COALESCE(variant, ''), 3, unit, 'damaged', 'Bruised during delivery to customer', 360, 'Manager User', NOW() - INTERVAL '6 days'
    FROM inventory_items WHERE name = 'Apple' LIMIT 1;

    INSERT INTO wastage_log (inventory_item_id, item_name, quantity, unit, reason, reason_details, cost_value, logged_by, logged_at)
    SELECT id, name || ' ' || COALESCE(variant, ''), 8, unit, 'pest', 'Rodent damage in dry storage area', 256, 'Admin User', NOW() - INTERVAL '15 days'
    FROM inventory_items WHERE name = 'Potato' LIMIT 1;

    -- Insert Expiry Alerts for items expiring soon (columns: inventory_item_id, alert_date, expiry_date, days_until_expiry)
    INSERT INTO expiry_alerts (inventory_item_id, alert_date, expiry_date, days_until_expiry, acknowledged)
    SELECT id, CURRENT_DATE, expiry_date,
        (expiry_date - CURRENT_DATE),
        false
    FROM inventory_items
    WHERE expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
    AND quantity > 0;

    -- Mark expired items
    INSERT INTO expiry_alerts (inventory_item_id, alert_date, expiry_date, days_until_expiry, acknowledged, acknowledged_at, acknowledged_by)
    SELECT id, CURRENT_DATE - 1, expiry_date,
        (expiry_date - CURRENT_DATE),
        true, NOW() - INTERVAL '1 day', 'System'
    FROM inventory_items
    WHERE expiry_date < CURRENT_DATE
    AND quantity > 0;

    -- Insert Sale Items for existing transactions
    -- Get some transaction IDs and add items to them
    IF tx1_id IS NOT NULL THEN
        INSERT INTO sale_items (transaction_id, inventory_lot_id, item_name, quantity, unit, price_per_unit, total, cost_per_unit)
        SELECT tx1_id, id, name || ' ' || COALESCE(variant, ''), 50, unit, selling_price, 50 * selling_price, cost_price
        FROM inventory_items WHERE name = 'Tomato' AND variant = 'Local' LIMIT 1;
        
        INSERT INTO sale_items (transaction_id, inventory_lot_id, item_name, quantity, unit, price_per_unit, total, cost_per_unit)
        SELECT tx1_id, id, name || ' ' || COALESCE(variant, ''), 100, unit, selling_price, 100 * selling_price, cost_price
        FROM inventory_items WHERE name = 'Onion' LIMIT 1;
    END IF;

    IF tx2_id IS NOT NULL THEN
        INSERT INTO sale_items (transaction_id, inventory_lot_id, item_name, quantity, unit, price_per_unit, total, cost_per_unit)
        SELECT tx2_id, id, name || ' ' || COALESCE(variant, ''), 30, unit, selling_price, 30 * selling_price, cost_price
        FROM inventory_items WHERE name = 'Tomato' LIMIT 1;
        
        INSERT INTO sale_items (transaction_id, inventory_lot_id, item_name, quantity, unit, price_per_unit, total, cost_per_unit)
        SELECT tx2_id, id, name || ' ' || COALESCE(variant, ''), 25, unit, selling_price, 25 * selling_price, cost_price
        FROM inventory_items WHERE name = 'Spinach' LIMIT 1;
    END IF;

    -- Add price history
    INSERT INTO price_history (inventory_item_id, old_price, new_price, changed_by, reason, changed_at)
    SELECT id, 45, selling_price, 'Admin User', 'Seasonal price adjustment', NOW() - INTERVAL '5 days'
    FROM inventory_items WHERE name = 'Tomato' AND variant = 'Local' LIMIT 1;

    INSERT INTO price_history (inventory_item_id, old_price, new_price, changed_by, reason, changed_at)
    SELECT id, selling_price + 20, selling_price, 'Manager User', 'Clearance discount - expiring soon', NOW() - INTERVAL '1 day'
    FROM inventory_items WHERE expiry_date <= CURRENT_DATE + 2 LIMIT 1;

END $$;

-- Summary
SELECT 'Crate Ledger' as table_name, COUNT(*) as count FROM crate_ledger
UNION ALL
SELECT 'Wastage Log', COUNT(*) FROM wastage_log
UNION ALL
SELECT 'Expiry Alerts', COUNT(*) FROM expiry_alerts
UNION ALL
SELECT 'Sale Items', COUNT(*) FROM sale_items
UNION ALL
SELECT 'Price History', COUNT(*) FROM price_history;

COMMIT;
