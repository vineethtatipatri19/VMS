-- Additional Demo Data: Transactions, Crates, Wastage, and Expiry Alerts
-- Compatible with actual VMS schema

DO $$
DECLARE
    customer_1 UUID;
    customer_2 UUID;
    customer_3 UUID;
    customer_4 UUID;
    inv_tomato UUID;
    inv_potato UUID;
    inv_onion UUID;
    inv_carrot UUID;
    inv_capsicum UUID;
    inv_apple UUID;
    inv_banana UUID;
    inv_mango UUID;
    inv_grapes UUID;
    inv_spinach UUID;
    inv_strawberry UUID;
    trans_1 UUID := gen_random_uuid();
    trans_2 UUID := gen_random_uuid();
    trans_3 UUID := gen_random_uuid();
    trans_4 UUID := gen_random_uuid();
BEGIN
    -- Get customer IDs
    SELECT id INTO customer_1 FROM customers WHERE name = 'Fresh Mart Supermarket' LIMIT 1;
    SELECT id INTO customer_2 FROM customers WHERE name = 'Green Valley Restaurant' LIMIT 1;
    SELECT id INTO customer_3 FROM customers WHERE name = 'Rajesh Kumar' LIMIT 1;
    SELECT id INTO customer_4 FROM customers WHERE name = 'Ravi Wholesale Trading' LIMIT 1;

    -- Get inventory IDs
    SELECT id INTO inv_tomato FROM inventory_items WHERE name = 'Tomato' AND variant = 'Regular' LIMIT 1;
    SELECT id INTO inv_potato FROM inventory_items WHERE name = 'Potato' LIMIT 1;
    SELECT id INTO inv_onion FROM inventory_items WHERE name = 'Onion' AND variant = 'Red' LIMIT 1;
    SELECT id INTO inv_carrot FROM inventory_items WHERE name = 'Carrot' AND variant = 'Regular' LIMIT 1;
    SELECT id INTO inv_capsicum FROM inventory_items WHERE name = 'Capsicum' AND variant = 'Green' LIMIT 1;
    SELECT id INTO inv_apple FROM inventory_items WHERE name = 'Apple' AND variant = 'Shimla Deluxe' LIMIT 1;
    SELECT id INTO inv_banana FROM inventory_items WHERE name = 'Banana' AND variant = 'Robusta' LIMIT 1;
    SELECT id INTO inv_mango FROM inventory_items WHERE name = 'Mango' LIMIT 1;
    SELECT id INTO inv_grapes FROM inventory_items WHERE name = 'Grapes' AND variant = 'Black Seedless' LIMIT 1;
    SELECT id INTO inv_spinach FROM inventory_items WHERE name = 'Spinach' LIMIT 1;
    SELECT id INTO inv_strawberry FROM inventory_items WHERE name = 'Strawberry' LIMIT 1;

    -- Transaction 1: Fresh Mart - Large B2B Order (Paid with UPI)
    INSERT INTO transactions (
        id, invoice_number, date, type, sale_type, customer_id,
        total_amount, payment_amount, payment_method,
        delivery_status, status, notes
    ) VALUES
    (trans_1, NULL, '2025-11-15 10:30:00', 'sale', 'credit', customer_1,
     8610.00, 8610.00, 'upi',
     'delivered', 'completed', 'Weekly order for Fresh Mart Supermarket');

    INSERT INTO sale_items (transaction_id, inventory_lot_id, quantity, price_per_unit, total, hsn_code, tax_percentage)
    VALUES 
        (trans_1, inv_tomato, 50.00, 45.00, 2250.00, '07019000', 0),
        (trans_1, inv_potato, 100.00, 30.00, 3000.00, '07019000', 0),
        (trans_1, inv_onion, 80.00, 42.00, 3360.00, '07019000', 0);

    -- Transaction 2: Green Valley Restaurant (Credit - Unpaid)
    INSERT INTO transactions (
        id, invoice_number, date, type, sale_type, customer_id,
        total_amount, payment_amount, payment_method,
        delivery_status, status, notes, due_date
    ) VALUES
    (trans_2, NULL, '2025-11-16 08:15:00', 'sale', 'credit', customer_2,
     4850.00, 0.00, 'credit',
     'delivered', 'completed', 'Daily fresh produce - Payment pending', '2025-11-30');

    INSERT INTO sale_items (transaction_id, inventory_lot_id, quantity, price_per_unit, total, hsn_code, tax_percentage)
    VALUES 
        (trans_2, inv_tomato, 30.00, 45.00, 1350.00, '07019000', 0),
        (trans_2, inv_capsicum, 20.00, 110.00, 2200.00, '07019000', 5),
        (trans_2, inv_carrot, 25.00, 52.00, 1300.00, '07019000', 0);

    -- Transaction 3: Ravi Wholesale - Large Order with Discount
    INSERT INTO transactions (
        id, invoice_number, date, type, sale_type, customer_id,
        total_amount, payment_amount, payment_method, discount_amount,
        delivery_status, status, notes
    ) VALUES
    (trans_3, NULL, '2025-11-16 14:20:00', 'sale', 'wholesale', customer_4,
     12700.00, 12700.00, 'bank_transfer', 1400.00,
     'delivered', 'completed', 'Wholesale order - 10% discount applied');

    INSERT INTO sale_items (transaction_id, inventory_lot_id, quantity, price_per_unit, total, hsn_code, tax_percentage, discount_percentage)
    VALUES 
        (trans_3, inv_potato, 150.00, 30.00, 4500.00, '07019000', 0, 10.0),
        (trans_3, inv_onion, 100.00, 42.00, 4200.00, '07019000', 0, 10.0),
        (trans_3, inv_banana, 80.00, 50.00, 4000.00, '08030000', 0, 10.0);

    -- Transaction 4: Rajesh Kumar - Retail Purchase
    INSERT INTO transactions (
        id, invoice_number, date, type, sale_type, customer_id,
        total_amount, payment_amount, payment_method,
        delivery_status, status, notes
    ) VALUES
    (trans_4, NULL, '2025-11-17 16:45:00', 'sale', 'regular', customer_3,
     1180.00, 1180.00, 'cash',
     'delivered', 'completed', 'Weekend grocery shopping');

    INSERT INTO sale_items (transaction_id, inventory_lot_id, quantity, price_per_unit, total, hsn_code, tax_percentage)
    VALUES 
        (trans_4, inv_tomato, 5.00, 45.00, 225.00, '07019000', 0),
        (trans_4, inv_banana, 10.00, 50.00, 500.00, '08030000', 0),
        (trans_4, inv_apple, 3.00, 185.00, 555.00, '08081000', 0);

    -- Crate Ledger: Fresh Mart - issued
    INSERT INTO crate_ledger (
        customer_id, transaction_id, date, crates_issued, crates_returned,
        balance, crate_type, crate_value, notes
    ) VALUES
    (customer_1, trans_1, '2025-11-15 10:30:00', 15, 0, 15, 'plastic', 50, 'Crates issued with delivery');
    
    -- Crate Ledger: Fresh Mart - returned
    INSERT INTO crate_ledger (
        customer_id, date, crates_issued, crates_returned,
        balance, crate_type, crate_value, notes
    ) VALUES
    (customer_1, '2025-11-17 09:00:00', 0, 12, 3, 'plastic', 50, 'Crates returned from previous delivery');
    
    -- Crate Ledger: Green Valley - issued
    INSERT INTO crate_ledger (
        customer_id, transaction_id, date, crates_issued, crates_returned,
        balance, crate_type, crate_value, notes
    ) VALUES
    (customer_2, trans_2, '2025-11-16 08:15:00', 8, 0, 8, 'plastic', 50, 'Crates for restaurant order');
    
    -- Crate Ledger: Ravi Wholesale - wooden crates issued
    INSERT INTO crate_ledger (
        customer_id, transaction_id, date, crates_issued, crates_returned,
        balance, crate_type, crate_value, notes
    ) VALUES
    (customer_4, trans_3, '2025-11-16 14:20:00', 20, 0, 20, 'wooden', 80, 'Wooden crates for wholesale order');
    
    -- Crate Ledger: Ravi Wholesale - returned
    INSERT INTO crate_ledger (
        customer_id, date, crates_issued, crates_returned,
        balance, crate_type, crate_value, notes
    ) VALUES
    (customer_4, '2025-11-18 10:00:00', 0, 18, 2, 'wooden', 80, 'Returned wooden crates - 2 damaged');

    -- Wastage Log entries
    INSERT INTO wastage_log (
        inventory_item_id, item_name, quantity, unit, reason, 
        reason_details, cost_value, logged_by, logged_at
    ) VALUES
    (inv_tomato, 'Tomato - Regular', 8.5, 'kg', 'spoiled',
     'Over-ripe tomatoes, unsellable quality', 238.00, 'warehouse_manager', '2025-11-14 18:00:00'),
    
    (inv_banana, 'Banana - Robusta', 12.0, 'kg', 'damaged',
     'Physical damage during transport', 360.00, 'warehouse_manager', '2025-11-15 09:30:00'),
    
    (inv_capsicum, 'Capsicum - Green', 3.2, 'kg', 'expired',
     'Past expiry date', 224.00, 'quality_inspector', '2025-11-17 07:00:00'),
    
    (inv_spinach, 'Spinach - Fresh', 5.5, 'kg', 'spoiled',
     'Wilted and yellowing', 165.00, 'warehouse_staff', '2025-11-18 06:00:00'),
    
    (inv_strawberry, 'Strawberry - Fresh', 2.8, 'kg', 'spoiled',
     'Mold formation detected', 504.00, 'quality_inspector', '2025-11-18 07:30:00'),
    
    (inv_mango, 'Mango - Alphonso', 4.0, 'kg', 'damaged',
     'Bruised during handling', 800.00, 'warehouse_manager', '2025-11-16 16:00:00'),
    
    ((SELECT id FROM inventory_items WHERE name = 'Ladies Finger' LIMIT 1), 
     'Ladies Finger - Fresh', 6.2, 'kg', 'expired',
     'Expired, quality degraded', 279.00, 'quality_inspector', '2025-11-17 08:00:00');

    -- Expiry Alerts
    INSERT INTO expiry_alerts (
        inventory_item_id, alert_date, expiry_date, days_until_expiry, acknowledged
    ) VALUES
    (inv_spinach, CURRENT_DATE, '2025-11-18', 0, false),
    ((SELECT id FROM inventory_items WHERE name = 'Ladies Finger' LIMIT 1), CURRENT_DATE, '2025-11-19', 1, false),
    (inv_strawberry, CURRENT_DATE, '2025-11-20', 2, false),
    ((SELECT id FROM inventory_items WHERE name = 'Green Beans' LIMIT 1), CURRENT_DATE, '2025-11-19', 1, false),
    (inv_tomato, CURRENT_DATE, '2025-11-20', 2, false),
    (inv_capsicum, CURRENT_DATE, '2025-11-20', 2, false),
    ((SELECT id FROM inventory_items WHERE name = 'Sapota' LIMIT 1), CURRENT_DATE, '2025-11-20', 2, false),
    ((SELECT id FROM inventory_items WHERE name = 'Cucumber' LIMIT 1), CURRENT_DATE, '2025-11-21', 3, false),
    ((SELECT id FROM inventory_items WHERE name = 'Broccoli' LIMIT 1), CURRENT_DATE, '2025-11-21', 3, false),
    (inv_carrot, CURRENT_DATE, '2025-11-22', 4, false),
    (inv_grapes, CURRENT_DATE, '2025-11-22', 4, false),
    ((SELECT id FROM inventory_items WHERE name = 'Guava' LIMIT 1), CURRENT_DATE, '2025-11-22', 4, false);

END $$;

ANALYZE;

-- Show summary
SELECT 'Summary of Additional Demo Data' as info, '' as count
UNION ALL
SELECT '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', '━━━━━━━━━━━━━━━━━━━━━━━━━'
UNION ALL
SELECT 'Transactions:', COUNT(*)::text FROM transactions
UNION ALL
SELECT '  - Sales:', COUNT(*)::text FROM transactions WHERE type = 'sale'
UNION ALL
SELECT 'Sale Items:', COUNT(*)::text FROM sale_items
UNION ALL
SELECT 'Crate Ledger:', COUNT(*)::text FROM crate_ledger
UNION ALL
SELECT '  - Issued:', SUM(crates_issued)::text FROM crate_ledger
UNION ALL
SELECT '  - Returned:', SUM(crates_returned)::text FROM crate_ledger
UNION ALL
SELECT 'Wastage Log:', COUNT(*)::text FROM wastage_log
UNION ALL
SELECT 'Expiry Alerts:', COUNT(*)::text FROM expiry_alerts
UNION ALL
SELECT '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', '━━━━━━━━━━━━━━━━━━━━━━━━━';
