-- Comprehensive Test Data Seed Script
-- This script creates diverse, realistic test data covering all field variations

-- Clean existing data
TRUNCATE TABLE 
    sale_items,
    payment_schedules,
    crate_ledger,
    wastage_log,
    expiry_alerts,
    price_history,
    transactions,
    inventory_items,
    customers,
    users
CASCADE;

-- Reset sequences if any
-- Note: We use UUIDs, so no sequences to reset

-- ============================================================================
-- USERS - Different roles and scenarios
-- ============================================================================
-- Password for all users: demo123 (bcrypt hash)
INSERT INTO users (id, name, email, password, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Admin User', 'admin@vms.com', '$2a$10$rN7YT4qX0mEXqjVYn3cJNOqXF8d1KzN7VqFx9hWqX7qJvQYX7xGZS', NOW() - INTERVAL '180 days'),
('550e8400-e29b-41d4-a716-446655440002', 'Manager User', 'manager@vms.com', '$2a$10$rN7YT4qX0mEXqjVYn3cJNOqXF8d1KzN7VqFx9hWqX7qJvQYX7xGZS', NOW() - INTERVAL '120 days'),
('550e8400-e29b-41d4-a716-446655440003', 'Sales User', 'sales@vms.com', '$2a$10$rN7YT4qX0mEXqjVYn3cJNOqXF8d1KzN7VqFx9hWqX7qJvQYX7xGZS', NOW() - INTERVAL '90 days'),
('550e8400-e29b-41d4-a716-446655440004', 'Demo User', 'demo@vms.com', '$2a$10$rN7YT4qX0mEXqjVYn3cJNOqXF8d1KzN7VqFx9hWqX7qJvQYX7xGZS', NOW() - INTERVAL '30 days');

-- ============================================================================
-- CUSTOMERS - All types, statuses, and field variations
-- ============================================================================

-- B2B Wholesale Customers (Large credit limits, payment terms)
INSERT INTO customers (id, name, email, address, contact_number, alternate_contact, whatsapp_number, 
    business_name, gstin, customer_type, credit_limit, current_balance, payment_terms_days, 
    interest_rate, status, kyc_document_type, kyc_document_number, aadhaar_verified, 
    notes, tags, total_purchases, loyalty_points, created_at, last_transaction_date) VALUES
('650e8400-e29b-41d4-a716-446655440001', 'Rajesh Kumar', 'rajesh@greenveg.com', 
    'Shop 12, Green Market, MG Road, Bangalore - 560001', '+919876543210', '+919876543211', '+919876543210',
    'Green Vegetables Wholesale', '29ABCDE1234F1Z5', 'b2b', 500000, 125000, 45, 
    1.5, 'active', 'GST Certificate', '29ABCDE1234F1Z5', true,
    'Premium wholesale customer. Always pays on time. Prefers early morning deliveries.',
    ARRAY['wholesale', 'vip', 'punctual'], 2500000, 1250, NOW() - INTERVAL '90 days', NOW() - INTERVAL '2 days'),

('650e8400-e29b-41d4-a716-446655440002', 'Priya Sharma', 'priya@freshmart.in',
    'Fresh Mart Supermarket, Brigade Road, Bangalore - 560025', '+919876543220', '+919876543221', '+919876543220',
    'Fresh Mart Retail Pvt Ltd', '29FGHIJ5678K1Z5', 'b2b', 750000, 85000, 30,
    1.2, 'active', 'GST Certificate', '29FGHIJ5678K1Z5', true,
    'Large retail chain. Orders 3 times a week. Needs invoice copies by email.',
    ARRAY['retail-chain', 'high-volume', 'regular'], 1800000, 900, NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day'),

('650e8400-e29b-41d4-a716-446655440003', 'Mohammed Salim', 'salim@organichub.net',
    'Organic Hub, Indiranagar, Bangalore - 560038', '+919876543230', NULL, '+919876543230',
    'Organic Hub India', '29KLMNO9012P1Z5', 'wholesale', 300000, 45000, 21,
    1.8, 'active', 'Trade License', 'TL/2024/BLR/1234', true,
    'Specializes in organic vegetables. Quality conscious. Premium pricing acceptable.',
    ARRAY['organic', 'premium', 'quality-focused'], 980000, 490, NOW() - INTERVAL '120 days', NOW() - INTERVAL '3 days'),

-- B2B Customer with Credit Issues
('650e8400-e29b-41d4-a716-446655440004', 'Suresh Patel', 'suresh@quickstore.com',
    'Quick Store, Whitefield, Bangalore - 560066', '+919876543240', '+919876543241', NULL,
    'Quick Store Retail', '29QRSTU3456V1Z5', 'b2b', 200000, 195000, 30,
    2.0, 'active', 'GST Certificate', '29QRSTU3456V1Z5', true,
    'CAUTION: Near credit limit. Request advance payment or reduce credit.',
    ARRAY['high-risk', 'slow-payer'], 450000, 50, NOW() - INTERVAL '45 days', NOW() - INTERVAL '25 days'),

-- B2C Retail Customers (Smaller purchases, cash/UPI)
('650e8400-e29b-41d4-a716-446655440005', 'Lakshmi Nair', 'lakshmi.nair@email.com',
    'Flat 204, Prestige Apartments, Koramangala, Bangalore - 560034', '+919876543250', NULL, '+919876543250',
    NULL, NULL, 'retail', 10000, 0, 0,
    0, 'active', 'Aadhaar', '1234 5678 9012', true,
    'Regular customer. Prefers UPI payments. Buys organic vegetables.',
    ARRAY['regular', 'organic-lover', 'upi'], 45000, 450, NOW() - INTERVAL '150 days', NOW() - INTERVAL '5 days'),

('650e8400-e29b-41d4-a716-446655440006', 'Amit Verma', 'amit.v@techmail.com',
    'HSR Layout, Sector 2, Bangalore - 560102', '+919876543260', '+919876543261', '+919876543260',
    NULL, NULL, 'retail', 5000, 1200, 7,
    0, 'active', 'Aadhaar', '9876 5432 1098', true,
    'Tech professional. Weekly orders. Prefers card payments.',
    ARRAY['weekly', 'card-payment'], 28000, 280, NOW() - INTERVAL '90 days', NOW() - INTERVAL '7 days'),

('650e8400-e29b-41d4-a716-446655440007', 'Sneha Reddy', NULL,
    'JP Nagar 3rd Phase, Bangalore - 560078', '+919876543270', NULL, '+919876543270',
    NULL, NULL, 'b2c', 3000, 500, 3,
    0, 'active', NULL, NULL, false,
    'New customer. Small orders. Cash only.',
    ARRAY['new', 'cash-only'], 2500, 25, NOW() - INTERVAL '15 days', NOW() - INTERVAL '10 days'),

-- Inactive/Blocked Customers
('650e8400-e29b-41d4-a716-446655440008', 'Ravi Shankar', 'ravi@oldstore.com',
    'Old Market Area, Bangalore - 560002', '+919876543280', NULL, NULL,
    'Ravi General Store', '29WXYZ7890A1Z5', 'wholesale', 100000, 25000, 15,
    0, 'inactive', 'GST Certificate', '29WXYZ7890A1Z5', false,
    'Store closed temporarily. Follow up after 3 months.',
    ARRAY['inactive', 'follow-up'], 180000, 0, NOW() - INTERVAL '200 days', NOW() - INTERVAL '65 days'),

('650e8400-e29b-41d4-a716-446655440009', 'Kiran Gupta', 'kiran@baddebt.com',
    'Unknown Address', '+919876543290', NULL, NULL,
    'Gupta Trading', '29BCDEF1111G1Z5', 'b2b', 150000, 148000, 30,
    5.0, 'blocked', 'GST Certificate', '29BCDEF1111G1Z5', true,
    'BLOCKED: Multiple payment defaults. Legal notice sent. Do not sell.',
    ARRAY['blocked', 'legal-issue', 'bad-debt'], 95000, 0, NOW() - INTERVAL '180 days', NOW() - INTERVAL '90 days'),

-- High Volume Customer
('650e8400-e29b-41d4-a716-446655440010', 'Metro Vegetables Ltd', 'orders@metroveg.com',
    'Metro Cash & Carry, Hebbal, Bangalore - 560024', '+919876543300', '+919876543301', '+919876543300',
    'Metro Vegetables Limited', '29HIJKL2222M1Z5', 'b2b', 1000000, 350000, 60,
    1.0, 'active', 'GST Certificate', '29HIJKL2222M1Z5', true,
    'VIP Customer. Largest account. Daily orders. Dedicated delivery slot 6-7 AM.',
    ARRAY['vip', 'high-volume', 'daily', 'priority'], 8500000, 8500, NOW() - INTERVAL '365 days', NOW());

-- ============================================================================
-- INVENTORY ITEMS - All categories, units, statuses, and variations
-- ============================================================================

-- Fresh Vegetables (High turnover, perishable)
INSERT INTO inventory_items (id, name, variant, category, sub_category, lot_number,
    quantity, unit, purchase_date, expiry_date, cost_price, selling_price,
    supplier_name, purchase_invoice, min_stock_level, reorder_point, shelf_life_days,
    storage_location, barcode, sku, hsn_code, gst_rate, status, is_perishable,
    weight_per_unit, packaging_type, notes, created_at) VALUES

-- Tomatoes (Multiple lots, different expiry dates)
('750e8400-e29b-41d4-a716-446655440001', 'Tomato', 'Local', 'Vegetables', 'Fruiting Vegetables', 'LOT-TOM-001-2025',
    250, 'kg', CURRENT_DATE - 2, CURRENT_DATE + 3, 35, 50,
    'Karnataka Farmers Co-op', 'INV-2025-001', 50, 100, 5,
    'Cold Storage A1', 'VEG001TOM', 'SKU-TOM-LOC', '07019000', 5, 'available', true,
    1, 'Plastic Crates', 'Fresh local tomatoes. Good quality.', NOW() - INTERVAL '2 days'),

('750e8400-e29b-41d4-a716-446655440002', 'Tomato', 'Organic', 'Vegetables', 'Fruiting Vegetables', 'LOT-TOM-002-2025',
    80, 'kg', CURRENT_DATE - 1, CURRENT_DATE + 4, 55, 75,
    'Organic Farms Bangalore', 'INV-2025-002', 30, 50, 5,
    'Cold Storage A2', 'VEG001TOM-ORG', 'SKU-TOM-ORG', '07019000', 5, 'available', true,
    1, 'Cardboard Boxes', 'Premium organic tomatoes. Higher margin.', NOW() - INTERVAL '1 day'),

('750e8400-e29b-41d4-a716-446655440003', 'Tomato', 'Hybrid', 'Vegetables', 'Fruiting Vegetables', 'LOT-TOM-003-2025',
    15, 'kg', CURRENT_DATE - 5, CURRENT_DATE + 1, 40, 55,
    'Hybrid Seeds Farm', 'INV-2025-003', 50, 100, 5,
    'Cold Storage A1', 'VEG001TOM-HYB', 'SKU-TOM-HYB', '07019000', 5, 'low_stock', true,
    1, 'Plastic Crates', 'Expiring soon - sell at discount.', NOW() - INTERVAL '5 days'),

-- Onions (Long shelf life)
('750e8400-e29b-41d4-a716-446655440004', 'Onion', 'Red', 'Vegetables', 'Bulb Vegetables', 'LOT-ONI-001-2025',
    500, 'kg', CURRENT_DATE - 10, CURRENT_DATE + 50, 28, 40,
    'Nashik Onion Traders', 'INV-2025-004', 100, 200, 60,
    'Dry Storage B1', 'VEG002ONI', 'SKU-ONI-RED', '07031000', 5, 'available', true,
    1, 'Mesh Bags', 'Good quality red onions from Nashik.', NOW() - INTERVAL '10 days'),

('750e8400-e29b-41d4-a716-446655440005', 'Onion', 'White', 'Vegetables', 'Bulb Vegetables', 'LOT-ONI-002-2025',
    150, 'kg', CURRENT_DATE - 8, CURRENT_DATE + 52, 32, 45,
    'Nashik Onion Traders', 'INV-2025-005', 50, 100, 60,
    'Dry Storage B1', 'VEG002ONI-WHT', 'SKU-ONI-WHT', '07031000', 5, 'available', true,
    1, 'Mesh Bags', 'Premium white onions.', NOW() - INTERVAL '8 days'),

-- Potatoes
('750e8400-e29b-41d4-a716-446655440006', 'Potato', 'Regular', 'Vegetables', 'Root Vegetables', 'LOT-POT-001-2025',
    600, 'kg', CURRENT_DATE - 15, CURRENT_DATE + 45, 22, 32,
    'Punjab Potato Suppliers', 'INV-2025-006', 150, 300, 60,
    'Dry Storage B2', 'VEG003POT', 'SKU-POT-REG', '07019000', 5, 'available', true,
    1, 'Jute Bags', 'Fresh potatoes from Punjab.', NOW() - INTERVAL '15 days'),

-- Leafy Vegetables (Very perishable)
('750e8400-e29b-41d4-a716-446655440007', 'Spinach', 'Fresh', 'Vegetables', 'Leafy Vegetables', 'LOT-SPI-001-2025',
    45, 'kg', CURRENT_DATE, CURRENT_DATE + 2, 38, 55,
    'Local Green Farms', 'INV-2025-007', 20, 40, 2,
    'Cold Storage A3', 'VEG004SPI', 'SKU-SPI-FRS', '07019000', 5, 'available', true,
    0.5, 'Bundles', 'Very fresh. Sell today.', NOW()),

('750e8400-e29b-41d4-a716-446655440008', 'Coriander', 'Fresh', 'Vegetables', 'Leafy Vegetables', 'LOT-COR-001-2025',
    12, 'kg', CURRENT_DATE - 1, CURRENT_DATE + 1, 45, 65,
    'Local Green Farms', 'INV-2025-008', 10, 20, 2,
    'Cold Storage A3', 'VEG005COR', 'SKU-COR-FRS', '07019000', 5, 'available', true,
    0.25, 'Bundles', 'Premium coriander.', NOW() - INTERVAL '1 day'),

-- Expired Item (Should trigger alerts)
('750e8400-e29b-41d4-a716-446655440009', 'Lettuce', 'Iceberg', 'Vegetables', 'Leafy Vegetables', 'LOT-LET-001-2025',
    8, 'kg', CURRENT_DATE - 4, CURRENT_DATE - 1, 50, 70,
    'Hydroponic Farms', 'INV-2025-009', 15, 30, 3,
    'Cold Storage A3', 'VEG006LET', 'SKU-LET-ICE', '07019000', 5, 'expired', true,
    0.5, 'Plastic Wrap', 'EXPIRED - Mark for wastage.', NOW() - INTERVAL '4 days'),

-- Fruits
('750e8400-e29b-41d4-a716-446655440010', 'Apple', 'Shimla', 'Fruits', 'Pome Fruits', 'LOT-APP-001-2025',
    200, 'kg', CURRENT_DATE - 5, CURRENT_DATE + 25, 85, 120,
    'Himachal Apple Suppliers', 'INV-2025-010', 50, 100, 30,
    'Cold Storage C1', 'FRT001APP', 'SKU-APP-SHI', '08081000', 5, 'available', true,
    0.15, 'Cardboard Boxes', 'Premium Shimla apples.', NOW() - INTERVAL '5 days'),

('750e8400-e29b-41d4-a716-446655440011', 'Banana', 'Robusta', 'Fruits', 'Tropical Fruits', 'LOT-BAN-001-2025',
    180, 'dozen', CURRENT_DATE - 2, CURRENT_DATE + 4, 35, 50,
    'Tamil Nadu Banana Board', 'INV-2025-011', 30, 60, 6,
    'Room Temperature D1', 'FRT002BAN', 'SKU-BAN-ROB', '08030000', 5, 'available', true,
    1.5, 'Bunches', 'Fresh robusta bananas.', NOW() - INTERVAL '2 days'),

('750e8400-e29b-41d4-a716-446655440012', 'Mango', 'Alphonso', 'Fruits', 'Tropical Fruits', 'LOT-MAN-001-2025',
    95, 'kg', CURRENT_DATE - 3, CURRENT_DATE + 7, 180, 250,
    'Ratnagiri Mango Exporters', 'INV-2025-012', 20, 40, 10,
    'Cold Storage C2', 'FRT003MAN', 'SKU-MAN-ALP', '08045000', 5, 'available', true,
    0.3, 'Wooden Crates', 'Premium Alphonso mangoes. High value.', NOW() - INTERVAL '3 days'),

-- Packaged Goods (Non-perishable or long shelf life)
('750e8400-e29b-41d4-a716-446655440013', 'Rice', 'Basmati Premium', 'Grains', 'Rice', 'LOT-RIC-001-2025',
    500, 'kg', CURRENT_DATE - 60, CURRENT_DATE + 305, 55, 75,
    'India Gate Foods', 'INV-2025-013', 100, 200, 365,
    'Dry Storage E1', 'GRN001RIC', 'SKU-RIC-BAS', '10063000', 5, 'available', false,
    1, '5kg Bags', 'Premium basmati rice.', NOW() - INTERVAL '60 days'),

('750e8400-e29b-41d4-a716-446655440014', 'Dal', 'Toor Dal', 'Pulses', 'Lentils', 'LOT-DAL-001-2025',
    300, 'kg', CURRENT_DATE - 45, CURRENT_DATE + 320, 85, 110,
    'Maharashtra Pulse Traders', 'INV-2025-014', 80, 150, 365,
    'Dry Storage E2', 'PLS001DAL', 'SKU-DAL-TOO', '07133000', 5, 'available', false,
    1, '1kg Packets', 'Quality toor dal.', NOW() - INTERVAL '45 days'),

-- Beverages (Different packaging)
('750e8400-e29b-41d4-a716-446655440015', 'Juice', 'Orange Fresh', 'Beverages', 'Fruit Juice', 'LOT-JUI-001-2025',
    144, 'bottles', CURRENT_DATE - 10, CURRENT_DATE + 80, 18, 30,
    'Fresh Juice Company', 'INV-2025-015', 50, 100, 90,
    'Cold Storage F1', 'BEV001JUI', 'SKU-JUI-ORA', '22029000', 12, 'available', true,
    1, 'Glass Bottles 500ml', 'Fresh orange juice. Refrigerated.', NOW() - INTERVAL '10 days'),

('750e8400-e29b-41d4-a716-446655440016', 'Milk', 'Full Cream', 'Dairy', 'Milk', 'LOT-MLK-001-2025',
    96, 'liters', CURRENT_DATE, CURRENT_DATE + 3, 52, 65,
    'Nandini Dairy', 'INV-2025-016', 100, 200, 3,
    'Cold Storage F2', 'DRY001MLK', 'SKU-MLK-FCR', '04012000', 5, 'available', true,
    1, '1L Packets', 'Fresh milk. Daily supply.', NOW()),

-- Damaged Item
('750e8400-e29b-41d4-a716-446655440017', 'Cooking Oil', 'Sunflower', 'Oils', 'Edible Oil', 'LOT-OIL-001-2025',
    25, 'liters', CURRENT_DATE - 30, CURRENT_DATE + 335, 125, 160,
    'Fortune Foods', 'INV-2025-017', 50, 100, 365,
    'Dry Storage E3', 'OIL001COO', 'SKU-OIL-SUN', '15121900', 5, 'damaged', false,
    1, '1L Bottles', 'DAMAGED - 5 bottles leaked during transport. Sell remaining at discount.', NOW() - INTERVAL '30 days'),

-- Low Stock Items
('750e8400-e29b-41d4-a716-446655440018', 'Carrot', 'Orange', 'Vegetables', 'Root Vegetables', 'LOT-CAR-001-2025',
    18, 'kg', CURRENT_DATE - 3, CURRENT_DATE + 12, 32, 48,
    'Local Farms', 'INV-2025-018', 40, 80, 15,
    'Cold Storage A4', 'VEG007CAR', 'SKU-CAR-ORA', '07019000', 5, 'low_stock', true,
    0.2, 'Plastic Bags', 'REORDER URGENTLY - Below reorder point.', NOW() - INTERVAL '3 days'),

('750e8400-e29b-41d4-a716-446655440019', 'Cucumber', 'Fresh', 'Vegetables', 'Fruiting Vegetables', 'LOT-CUC-001-2025',
    22, 'kg', CURRENT_DATE - 2, CURRENT_DATE + 6, 28, 42,
    'Local Farms', 'INV-2025-019', 35, 70, 8,
    'Cold Storage A4', 'VEG008CUC', 'SKU-CUC-FRS', '07019000', 5, 'low_stock', true,
    0.3, 'Plastic Crates', 'Running low. Order placed.', NOW() - INTERVAL '2 days'),

-- Out of Stock Item
('750e8400-e29b-41d4-a716-446655440020', 'Cauliflower', 'White', 'Vegetables', 'Flowering Vegetables', 'LOT-CAU-001-2025',
    0, 'kg', CURRENT_DATE - 7, CURRENT_DATE + 1, 35, 52,
    'Punjab Vegetable Suppliers', 'INV-2025-020', 50, 100, 8,
    'Cold Storage A5', 'VEG009CAU', 'SKU-CAU-WHT', '07019000', 5, 'out_of_stock', true,
    0.8, 'Plastic Crates', 'OUT OF STOCK - High demand. Rush order placed.', NOW() - INTERVAL '7 days');

-- ============================================================================
-- TRANSACTIONS - All types, payment methods, and scenarios
-- ============================================================================

-- Large B2B Credit Sales (Multiple items)
INSERT INTO transactions (id, customer_id, date, type, payment_amount, total_amount, payment_method,
    payment_reference, due_date, is_overdue, discount_amount, tax_amount, notes,
    invoice_number, balance_after, sale_type, delivery_status, delivery_date, status) VALUES

-- Recent sale to Metro Vegetables (VIP Customer)
('850e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440010',
    NOW(), 'sale', 0, 45500, 'credit', NULL,
    CURRENT_DATE + 60, false, 2500, 2285.71, 
    'Large order for weekend stock. Delivered at 6 AM as requested.',
    NULL, 395500, 'wholesale', 'delivered', NOW(), 'completed'),

-- Sale to Fresh Mart (Regular B2B)
('850e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440002',
    NOW() - INTERVAL '1 day', 'sale', 0, 28750, 'credit', NULL,
    CURRENT_DATE + 29, false, 750, 1333.33,
    'Weekly regular order. 3 varieties of tomatoes.',
    NULL, 113750, 'wholesale', 'delivered', NOW() - INTERVAL '1 day', 'completed'),

-- Sale to Green Vegetables Wholesale
('850e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440001',
    NOW() - INTERVAL '2 days', 'sale', 0, 32400, 'credit', NULL,
    CURRENT_DATE + 43, false, 900, 1504.76,
    'Mixed vegetables order. Premium organic items included.',
    NULL, 157400, 'wholesale', 'delivered', NOW() - INTERVAL '2 days', 'completed'),

-- Cash Sale to Retail Customer
('850e8400-e29b-41d4-a716-446655440004', '650e8400-e29b-41d4-a716-446655440005',
    NOW() - INTERVAL '5 days', 'sale', 2650, 2650, 'upi', 'UPI/123456789/2025',
    NULL, false, 100, 123.81,
    'Regular weekly purchase. Organic vegetables as usual.',
    NULL, 0, 'regular', 'delivered', NOW() - INTERVAL '5 days', 'completed'),

-- Card Payment Sale
('850e8400-e29b-41d4-a716-446655440005', '650e8400-e29b-41d4-a716-446655440006',
    NOW() - INTERVAL '7 days', 'sale', 1850, 1850, 'card', 'CARD/987654321/2025',
    NULL, false, 0, 88.10,
    'Weekly order. Paid by card at delivery.',
    NULL, 0, 'regular', 'delivered', NOW() - INTERVAL '7 days', 'completed'),

-- Small Credit Sale
('850e8400-e29b-41d4-a716-446655440006', '650e8400-e29b-41d4-a716-446655440007',
    NOW() - INTERVAL '10 days', 'sale', 0, 875, 'credit', NULL,
    CURRENT_DATE - 7, true, 0, 41.67,
    'Small order. Customer requested 3 days credit.',
    NULL, 875, 'credit', 'delivered', NOW() - INTERVAL '10 days', 'completed'),

-- Overdue Sale (Payment pending)
('850e8400-e29b-41d4-a716-446655440007', '650e8400-e29b-41d4-a716-446655440004',
    NOW() - INTERVAL '25 days', 'sale', 0, 15200, 'credit', NULL,
    CURRENT_DATE - 10, true, 0, 723.81,
    'OVERDUE PAYMENT. Follow up urgently.',
    NULL, 15200, 'credit', 'delivered', NOW() - INTERVAL '25 days', 'completed'),

-- Partial Payment Transaction
('850e8400-e29b-41d4-a716-446655440008', '650e8400-e29b-41d4-a716-446655440001',
    NOW() - INTERVAL '15 days', 'payment', 50000, 0, 'bank_transfer', 'NEFT/TXN/2025/001',
    NULL, false, 0, 0,
    'Partial payment received against outstanding balance.',
    NULL, 107400, NULL, NULL, NULL, 'completed'),

-- Full Payment Transaction
('850e8400-e29b-41d4-a716-446655440009', '650e8400-e29b-41d4-a716-446655440003',
    NOW() - INTERVAL '10 days', 'payment', 35000, 0, 'cheque', 'CHQ/123456/2025',
    NULL, false, 0, 0,
    'Full payment for previous month purchases. Cheque cleared.',
    NULL, 10000, NULL, NULL, NULL, 'completed'),

-- Recent Large Cash Sale
('850e8400-e29b-41d4-a716-446655440010', '650e8400-e29b-41d4-a716-446655440002',
    NOW() - INTERVAL '3 days', 'sale', 18900, 18900, 'bank_transfer', 'IMPS/2025/456',
    NULL, false, 600, 876.19,
    'Special order for store opening event. Immediate payment.',
    NULL, 0, 'wholesale', 'delivered', NOW() - INTERVAL '3 days', 'completed'),

-- Mixed Vegetable Sale
('850e8400-e29b-41d4-a716-446655440011', '650e8400-e29b-41d4-a716-446655440010',
    NOW() - INTERVAL '4 days', 'sale', 0, 52300, 'credit', NULL,
    CURRENT_DATE + 56, false, 3200, 2414.29,
    'Weekly stock replenishment. All items delivered fresh.',
    NULL, 447800, 'wholesale', 'delivered', NOW() - INTERVAL '4 days', 'completed'),

-- Return Transaction (Negative sale)
('850e8400-e29b-41d4-a716-446655440012', '650e8400-e29b-41d4-a716-446655440002',
    NOW() - INTERVAL '6 days', 'sale', 0, -1250, 'credit', NULL,
    NULL, false, 0, -59.52,
    'Return of damaged items. Credit note issued.',
    NULL, 115000, 'return', 'delivered', NOW() - INTERVAL '6 days', 'completed');

-- ============================================================================
-- SALE ITEMS - Detailed line items for transactions
-- ============================================================================

-- Sale Items for transaction 001 (Metro Vegetables - Large order)
INSERT INTO sale_items (transaction_id, inventory_item_id, quantity, unit_price, total_price) VALUES
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 100, 50, 5000),  -- Tomato Local
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440004', 150, 40, 6000),  -- Onion Red
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440006', 200, 32, 6400),  -- Potato
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440010', 80, 120, 9600),  -- Apple
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440011', 60, 50, 3000),   -- Banana
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440012', 50, 250, 12500); -- Mango Alphonso

-- Sale Items for transaction 002 (Fresh Mart)
INSERT INTO sale_items (transaction_id, inventory_item_id, quantity, unit_price, total_price) VALUES
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 80, 50, 4000),   -- Tomato Local
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440002', 50, 75, 3750),   -- Tomato Organic
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440007', 35, 55, 1925),   -- Spinach
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440008', 10, 65, 650),    -- Coriander
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440016', 80, 65, 5200);   -- Milk

-- Sale Items for transaction 003 (Green Vegetables)
INSERT INTO sale_items (transaction_id, inventory_item_id, quantity, unit_price, total_price) VALUES
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440002', 60, 75, 4500),   -- Tomato Organic
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440005', 100, 45, 4500),  -- Onion White
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440007', 40, 55, 2200),   -- Spinach
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440010', 60, 120, 7200);  -- Apple

-- Sale Items for transaction 004 (Retail - Lakshmi)
INSERT INTO sale_items (transaction_id, inventory_item_id, quantity, unit_price, total_price) VALUES
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440002', 8, 75, 600),     -- Tomato Organic
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440007', 5, 55, 275),     -- Spinach
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440008', 3, 65, 195),     -- Coriander
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440010', 10, 120, 1200);  -- Apple

-- Sale Items for transaction 005 (Retail - Amit)
INSERT INTO sale_items (transaction_id, inventory_item_id, quantity, unit_price, total_price) VALUES
('850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440001', 15, 50, 750),    -- Tomato
('850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440006', 20, 32, 640),    -- Potato
('850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440011', 8, 50, 400);     -- Banana

-- ============================================================================
-- CRATE LEDGER - Crate tracking transactions
-- ============================================================================

INSERT INTO crate_ledger (customer_id, transaction_id, crates_given, crates_returned, 
    transaction_type, notes, created_at) VALUES

-- Crates given with sales
('650e8400-e29b-41d4-a716-446655440010', '850e8400-e29b-41d4-a716-446655440001', 
    25, 0, 'sale', 'Large order - 25 crates delivered with vegetables', NOW()),

('650e8400-e29b-41d4-a716-446655440002', '850e8400-e29b-41d4-a716-446655440002',
    12, 0, 'sale', '12 crates for mixed vegetable order', NOW() - INTERVAL '1 day'),

('650e8400-e29b-41d4-a716-446655440001', '850e8400-e29b-41d4-a716-446655440003',
    15, 0, 'sale', '15 crates with organic vegetables', NOW() - INTERVAL '2 days'),

-- Crates returned
('650e8400-e29b-41d4-a716-446655440010', NULL,
    0, 20, 'return', 'Returned 20 crates from previous deliveries', NOW() - INTERVAL '3 days'),

('650e8400-e29b-41d4-a716-446655440002', NULL,
    0, 8, 'return', 'Partial return of crates', NOW() - INTERVAL '5 days'),

('650e8400-e29b-41d4-a716-446655440001', NULL,
    0, 10, 'return', 'Returned crates', NOW() - INTERVAL '7 days'),

-- More crate movements
('650e8400-e29b-41d4-a716-446655440010', '850e8400-e29b-41d4-a716-446655440011',
    22, 0, 'sale', 'Weekly stock - 22 crates', NOW() - INTERVAL '4 days'),

('650e8400-e29b-41d4-a716-446655440003', '850e8400-e29b-41d4-a716-446655440009',
    8, 0, 'sale', '8 crates for organic order', NOW() - INTERVAL '10 days');

-- ============================================================================
-- EXPIRY ALERTS - Alerts for items expiring soon
-- ============================================================================

INSERT INTO expiry_alerts (inventory_item_id, alert_date, severity, status, notes) VALUES

-- Critical alerts (Expiring today or tomorrow)
('750e8400-e29b-41d4-a716-446655440003', CURRENT_DATE, 'high', 'pending',
    'Tomato Hybrid lot expiring in 1 day - 15 kg remaining. Offer discount sale.'),

('750e8400-e29b-41d4-a716-446655440007', CURRENT_DATE, 'high', 'pending',
    'Spinach expiring in 2 days - 45 kg. Very perishable. Sell urgently.'),

('750e8400-e29b-41d4-a716-446655440008', CURRENT_DATE, 'high', 'pending',
    'Coriander expiring in 1 day - 12 kg. Contact regular customers.'),

-- Medium priority (3-5 days)
('750e8400-e29b-41d4-a716-446655440001', CURRENT_DATE, 'medium', 'pending',
    'Tomato Local expiring in 3 days - 250 kg. High stock level.'),

('750e8400-e29b-41d4-a716-446655440011', CURRENT_DATE, 'medium', 'pending',
    'Banana expiring in 4 days - 180 dozen. Plan promotions.'),

-- Already expired
('750e8400-e29b-41d4-a716-446655440009', CURRENT_DATE - 1, 'high', 'acknowledged',
    'Lettuce EXPIRED. Moved to wastage. 8 kg loss.'),

-- Low priority (6-10 days)
('750e8400-e29b-41d4-a716-446655440012', CURRENT_DATE, 'low', 'pending',
    'Alphonso Mango expiring in 7 days - 95 kg. Premium item, monitor closely.'),

('750e8400-e29b-41d4-a716-446655440002', CURRENT_DATE, 'low', 'pending',
    'Organic Tomato expiring in 4 days - 80 kg. Premium pricing maintained.');

-- ============================================================================
-- WASTAGE LOG - Track damaged/expired inventory
-- ============================================================================

INSERT INTO wastage_log (inventory_item_id, quantity, reason, recorded_by, 
    estimated_loss, notes, created_at) VALUES

-- Expired wastage
('750e8400-e29b-41d4-a716-446655440009', 8, 'expired', 'Manager User',
    560, 'Lettuce expired. Could not sell in time due to slow demand.', NOW() - INTERVAL '1 day'),

-- Damaged during transport
('750e8400-e29b-41d4-a716-446655440017', 5, 'damaged', 'Admin User',
    625, 'Cooking oil bottles leaked during transport. 5 liters lost.', NOW() - INTERVAL '30 days'),

-- Quality issues
('750e8400-e29b-41d4-a716-446655440001', 12, 'spoiled', 'Sales User',
    600, 'Tomatoes started rotting. Removed from cold storage.', NOW() - INTERVAL '3 days'),

-- Customer return - damaged
('750e8400-e29b-41d4-a716-446655440010', 3, 'damaged', 'Manager User',
    360, 'Apples returned by customer - bruised during delivery.', NOW() - INTERVAL '6 days'),

-- Pest damage
('750e8400-e29b-41d4-a716-446655440006', 8, 'other', 'Admin User',
    256, 'Potato bag damaged by rodents in storage. Improved pest control.', NOW() - INTERVAL '15 days');

-- ============================================================================
-- PRICE HISTORY - Track price changes
-- ============================================================================

INSERT INTO price_history (inventory_item_id, old_price, new_price, changed_by, 
    reason, created_at) VALUES

-- Seasonal price adjustments
('750e8400-e29b-41d4-a716-446655440001', 45, 50, 'Admin User',
    'Seasonal increase - Summer demand high', NOW() - INTERVAL '5 days'),

('750e8400-e29b-41d4-a716-446655440012', 220, 250, 'Manager User',
    'Premium Alphonso season peak - increased to market rate', NOW() - INTERVAL '10 days'),

-- Discount for near expiry
('750e8400-e29b-41d4-a716-446655440003', 55, 45, 'Sales User',
    'Clearance discount - expiring in 2 days', NOW() - INTERVAL '1 day'),

-- Damaged goods discount
('750e8400-e29b-41d4-a716-446655440017', 160, 120, 'Manager User',
    'Discount on remaining stock after transport damage', NOW() - INTERVAL '29 days'),

-- Market price adjustment
('750e8400-e29b-41d4-a716-446655440004', 38, 40, 'Admin User',
    'Market rate increased - adjusted wholesale price', NOW() - INTERVAL '12 days');

-- ============================================================================
-- PAYMENT SCHEDULES - Installment plans for credit customers
-- ============================================================================

INSERT INTO payment_schedules (transaction_id, customer_id, installment_number, 
    due_date, amount_due, amount_paid, status, notes) VALUES

-- Payment plan for Suresh Patel (near credit limit)
('850e8400-e29b-41d4-a716-446655440007', '650e8400-e29b-41d4-a716-446655440004',
    1, CURRENT_DATE - 10, 5000, 0, 'overdue',
    'First installment overdue. Follow up immediately.'),

('850e8400-e29b-41d4-a716-446655440007', '650e8400-e29b-41d4-a716-446655440004',
    2, CURRENT_DATE + 5, 5000, 0, 'pending',
    'Second installment due in 5 days.'),

('850e8400-e29b-41d4-a716-446655440007', '650e8400-e29b-41d4-a716-446655440004',
    3, CURRENT_DATE + 20, 5200, 0, 'pending',
    'Final installment includes late payment charges.');

-- ============================================================================
-- PRICING TIERS - Volume-based pricing (if applicable)
-- ============================================================================

INSERT INTO pricing_tiers (customer_id, item_category, min_quantity, max_quantity, 
    discount_percentage, effective_from, effective_to, notes) VALUES

-- VIP customer special pricing
('650e8400-e29b-41d4-a716-446655440010', 'Vegetables', 100, NULL, 8.5,
    CURRENT_DATE - 30, CURRENT_DATE + 335,
    'Metro Vegetables - VIP volume discount on all vegetables'),

('650e8400-e29b-41d4-a716-446655440010', 'Fruits', 50, NULL, 7.0,
    CURRENT_DATE - 30, CURRENT_DATE + 335,
    'Metro Vegetables - Volume discount on fruits'),

-- Wholesale customer tiers
('650e8400-e29b-41d4-a716-446655440001', 'Vegetables', 50, 200, 5.0,
    CURRENT_DATE - 90, CURRENT_DATE + 275,
    'Green Vegetables - Tier 1 discount'),

('650e8400-e29b-41d4-a716-446655440001', 'Vegetables', 200, NULL, 8.0,
    CURRENT_DATE - 90, CURRENT_DATE + 275,
    'Green Vegetables - Tier 2 discount for bulk orders'),

-- Seasonal promotion
('650e8400-e29b-41d4-a716-446655440002', 'Fruits', 30, NULL, 6.0,
    CURRENT_DATE - 15, CURRENT_DATE + 45,
    'Fresh Mart - Summer fruit promotion discount');

-- ============================================================================
-- Summary Statistics
-- ============================================================================

-- Update inventory quantities after sales
UPDATE inventory_items SET quantity = quantity - 100 WHERE id = '750e8400-e29b-41d4-a716-446655440001'; -- Tomato Local
UPDATE inventory_items SET quantity = quantity - 50 WHERE id = '750e8400-e29b-41d4-a716-446655440002';  -- Tomato Organic
UPDATE inventory_items SET quantity = quantity - 150 WHERE id = '750e8400-e29b-41d4-a716-446655440004'; -- Onion
UPDATE inventory_items SET quantity = quantity - 200 WHERE id = '750e8400-e29b-41d4-a716-446655440006'; -- Potato
UPDATE inventory_items SET quantity = quantity - 80 WHERE id = '750e8400-e29b-41d4-a716-446655440007';  -- Spinach
UPDATE inventory_items SET quantity = quantity - 13 WHERE id = '750e8400-e29b-41d4-a716-446655440008';  -- Coriander
UPDATE inventory_items SET quantity = quantity - 140 WHERE id = '750e8400-e29b-41d4-a716-446655440010'; -- Apple
UPDATE inventory_items SET quantity = quantity - 68 WHERE id = '750e8400-e29b-41d4-a716-446655440011';  -- Banana
UPDATE inventory_items SET quantity = quantity - 50 WHERE id = '750e8400-e29b-41d4-a716-446655440012';  -- Mango
UPDATE inventory_items SET quantity = quantity - 80 WHERE id = '750e8400-e29b-41d4-a716-446655440016';  -- Milk

-- Mark items as low_stock where quantity < reorder_point
UPDATE inventory_items SET status = 'low_stock' 
WHERE quantity < reorder_point AND status = 'available';

-- Verification Queries
SELECT 'Users Created' as entity, COUNT(*) as count FROM users
UNION ALL
SELECT 'Customers Created', COUNT(*) FROM customers
UNION ALL
SELECT 'Inventory Items', COUNT(*) FROM inventory_items
UNION ALL
SELECT 'Transactions', COUNT(*) FROM transactions
UNION ALL
SELECT 'Sale Items', COUNT(*) FROM sale_items
UNION ALL
SELECT 'Crate Ledger Entries', COUNT(*) FROM crate_ledger
UNION ALL
SELECT 'Expiry Alerts', COUNT(*) FROM expiry_alerts
UNION ALL
SELECT 'Wastage Logs', COUNT(*) FROM wastage_log
UNION ALL
SELECT 'Price History', COUNT(*) FROM price_history
UNION ALL
SELECT 'Payment Schedules', COUNT(*) FROM payment_schedules
UNION ALL
SELECT 'Pricing Tiers', COUNT(*) FROM pricing_tiers;

-- Summary by status
SELECT 'Inventory by Status' as summary, status, COUNT(*) as count, SUM(quantity) as total_qty
FROM inventory_items 
GROUP BY status
ORDER BY status;

SELECT 'Customers by Type' as summary, customer_type, COUNT(*) as count
FROM customers
GROUP BY customer_type;

SELECT 'Transactions by Type' as summary, type, COUNT(*) as count, SUM(total_amount) as total
FROM transactions
GROUP BY type;

COMMIT;
