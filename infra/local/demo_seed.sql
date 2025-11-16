-- Demo Seed Data for VMS
-- Clear existing data
TRUNCATE TABLE transaction_items, transactions, inventory_lots, customers, crates, users, wastage_log, expiry_alerts, payment_schedules, pricing_tiers, price_history CASCADE;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE customers_id_seq RESTART WITH 1;
ALTER SEQUENCE inventory_lots_id_seq RESTART WITH 1;
ALTER SEQUENCE crates_id_seq RESTART WITH 1;
ALTER SEQUENCE transactions_id_seq RESTART WITH 1;
ALTER SEQUENCE transaction_items_id_seq RESTART WITH 1;
ALTER SEQUENCE wastage_log_id_seq RESTART WITH 1;
ALTER SEQUENCE expiry_alerts_id_seq RESTART WITH 1;
ALTER SEQUENCE payment_schedules_id_seq RESTART WITH 1;
ALTER SEQUENCE pricing_tiers_id_seq RESTART WITH 1;
ALTER SEQUENCE price_history_id_seq RESTART WITH 1;

-- Insert demo user
INSERT INTO users (email, password_hash, name, role) VALUES
('demo@vms.com', '$2a$10$rVxJZ9Y8W5vD5M8Q3QY8FuO8QhH5Z6X8Y9W5vD5M8Q3QY8FuO8QhH', 'Demo Admin', 'admin');

-- Insert diverse customers (20 customers with various types and credit statuses)
INSERT INTO customers (
    name, address, contact_number, photo_url, aadhaar_verified,
    email, whatsapp_number, alternate_contact, customer_type, business_name,
    gstin, credit_limit, current_balance, payment_terms_days, status,
    kyc_document_type, kyc_document_number, notes, created_at
) VALUES
-- B2B Customers (Premium)
('Organic Grocers Pvt Ltd', 'Shop 15, MG Road, Commercial Complex, Bangalore', '9876543210', NULL, true,
 'contact@organicgrocers.com', '9876543210', '9876543211', 'b2b', 'Organic Grocers Private Limited',
 '29ABCDE1234F1Z5', 50000.00, 12500.75, 30, 'active',
 'gstin', '29ABCDE1234F1Z5', 'Premium B2B client, monthly billing', NOW() - INTERVAL '6 months'),

('FreshMart Supermarket', '45 Brigade Road, Retail District, Bangalore', '9876543220', NULL, true,
 'orders@freshmart.in', '9876543220', '9876543221', 'b2b', 'FreshMart Retail Pvt Ltd',
 '29XYZAB5678C1D3', 75000.00, 28900.50, 45, 'active',
 'gstin', '29XYZAB5678C1D3', 'Large chain, bulk orders weekly', NOW() - INTERVAL '1 year'),

('Green Valley Restaurants', '12 Church Street, CBD, Bangalore', '9876543230', NULL, true,
 'procurement@greenvalley.com', '9876543230', NULL, 'b2b', 'Green Valley Food Services',
 '29PQRST9012E2F4', 35000.00, 0.00, 15, 'active',
 'gstin', '29PQRST9012E2F4', 'Restaurant chain, daily deliveries required', NOW() - INTERVAL '3 months'),

('Healthy Bites Cafe', '78 Indiranagar, Food Street, Bangalore', '9876543240', NULL, false,
 'cafe@healthybites.in', '9876543240', '9876543241', 'b2b', 'Healthy Bites Food & Beverages',
 NULL, 15000.00, 8200.00, 7, 'active',
 'pan', 'ABCDE1234F', 'New B2B client, weekly orders', NOW() - INTERVAL '1 month'),

-- B2B with issues
('Budget Bazaar', '90 KR Market, Wholesale Area, Bangalore', '9876543250', NULL, true,
 'accounts@budgetbazaar.com', '9876543250', NULL, 'b2b', 'Budget Bazaar Trading Co',
 '29LMNOP3456G3H5', 25000.00, 26500.00, 30, 'blocked',
 'gstin', '29LMNOP3456G3H5', 'BLOCKED: Credit limit exceeded, payment pending', NOW() - INTERVAL '8 months'),

-- Wholesale Customers
('Ravi Trading Company', '23 Avenue Road, Market Area, Bangalore', '9876543260', NULL, true,
 'ravi@ravitrading.com', '9876543260', '9876543261', 'wholesale', 'Ravi & Sons Trading',
 '29QWERT7890I4J6', 40000.00, 15600.00, 21, 'active',
 'gstin', '29QWERT7890I4J6', 'Wholesale dealer, bi-weekly orders', NOW() - INTERVAL '2 years'),

('Sri Venkateshwara Wholesale', '67 KR Puram, Wholesale District, Bangalore', '9876543270', NULL, true,
 'sri.venkat@svwholesale.com', '9876543270', NULL, 'wholesale', 'Sri Venkateshwara Traders',
 '29ASDFG4567K5L7', 60000.00, 34200.25, 30, 'active',
 'gstin', '29ASDFG4567K5L7', 'Long-term wholesale partner', NOW() - INTERVAL '5 years'),

('Metro Cash & Carry', '100 Whitefield, Industrial Area, Bangalore', '9876543280', NULL, true,
 'metro@cashcarry.in', '9876543280', '9876543281', 'wholesale', 'Metro Cash & Carry India',
 '29ZXCVB8901M6N8', 100000.00, 42000.00, 60, 'active',
 'gstin', '29ZXCVB8901M6N8', 'Large wholesale distributor', NOW() - INTERVAL '3 years'),

-- Retail Customers (Walk-ins)
('Rajesh Kumar', 'Flat 204, Prestige Apartments, Koramangala, Bangalore', '9876543290', NULL, true,
 'rajesh.kumar@email.com', '9876543290', NULL, 'retail', NULL,
 NULL, 5000.00, 0.00, 7, 'active',
 'aadhaar', '1234-5678-9012', 'Regular retail customer, prefers home delivery', NOW() - INTERVAL '1 year'),

('Priya Sharma', '56 Jayanagar 4th Block, Bangalore', '9876543300', NULL, true,
 'priya.s@email.com', '9876543300', NULL, 'retail', NULL,
 NULL, 3000.00, 450.00, 0, 'active',
 'aadhaar', '2345-6789-0123', 'Weekly grocery shopping', NOW() - INTERVAL '6 months'),

('Mohammed Ali', '89 RT Nagar, Near Bus Stand, Bangalore', '9876543310', NULL, false,
 'ali.m@email.com', '9876543310', '9876543311', 'retail', NULL,
 NULL, 2000.00, 0.00, 0, 'active',
 'driving_license', 'KA0320190012345', 'Occasional buyer', NOW() - INTERVAL '2 months'),

('Lakshmi Devi', '34 Malleshwaram, Temple Street, Bangalore', '9876543320', NULL, true,
 'lakshmi.d@email.com', '9876543320', NULL, 'retail', NULL,
 NULL, 4000.00, 1200.00, 7, 'active',
 'aadhaar', '3456-7890-1234', 'Senior citizen, home delivery preferred', NOW() - INTERVAL '2 years'),

('Arjun Reddy', '12 HSR Layout, Tech Park Road, Bangalore', '9876543330', NULL, true,
 'arjun.reddy@email.com', '9876543330', NULL, 'retail', NULL,
 NULL, 2500.00, 0.00, 0, 'active',
 'pan', 'FGHIJ5678K', 'Young professional, weekend shopper', NOW() - INTERVAL '3 months'),

-- Inactive/Blocked retail customers
('Suresh Patel', '78 Yeshwanthpur, Old Layout, Bangalore', '9876543340', NULL, false,
 NULL, '9876543340', NULL, 'retail', NULL,
 NULL, 1500.00, 1650.00, 0, 'blocked',
 'voter_id', 'ABC1234567', 'BLOCKED: Payment overdue for 45 days', NOW() - INTERVAL '1 year'),

('Divya Menon', '45 Electronic City, Phase 1, Bangalore', '9876543350', NULL, true,
 'divya.m@email.com', NULL, NULL, 'retail', NULL,
 NULL, 2000.00, 0.00, 0, 'inactive',
 'aadhaar', '4567-8901-2345', 'Inactive: No purchases in 6 months', NOW() - INTERVAL '1.5 years'),

-- B2C Customers (Online/App)
('Ananya Iyer', '23 Bellandur, Lake View Apartments, Bangalore', '9876543360', NULL, true,
 'ananya.iyer@email.com', '9876543360', NULL, 'b2c', NULL,
 NULL, 8000.00, 2400.00, 14, 'active',
 'aadhaar', '5678-9012-3456', 'Premium online customer, subscription plan', NOW() - INTERVAL '8 months'),

('Vikram Singh', '67 Marathahalli, Tech Hub, Bangalore', '9876543370', NULL, true,
 'vikram.s@email.com', '9876543370', '9876543371', 'b2c', NULL,
 NULL, 6000.00, 0.00, 10, 'active',
 'pan', 'KLMNO6789P', 'App user, weekly online orders', NOW() - INTERVAL '4 months'),

('Meera Krishnan', '90 JP Nagar, Garden City, Bangalore', '9876543380', NULL, true,
 'meera.k@email.com', '9876543380', NULL, 'b2c', NULL,
 NULL, 5000.00, 1800.50, 7, 'active',
 'aadhaar', '6789-0123-4567', 'Regular online shopper', NOW() - INTERVAL '10 months'),

('Karthik Raj', '11 Banashankari, Temple Road, Bangalore', '9876543390', NULL, false,
 'karthik.raj@email.com', '9876543390', NULL, 'b2c', NULL,
 NULL, 3000.00, 0.00, 0, 'active',
 NULL, NULL, 'New B2C customer', NOW() - INTERVAL '2 weeks'),

('Sneha Gupta', '88 Yelahanka, New Town, Bangalore', '9876543400', NULL, true,
 'sneha.g@email.com', '9876543400', NULL, 'b2c', NULL,
 NULL, 7000.00, 3500.00, 15, 'active',
 'aadhaar', '7890-1234-5678', 'High-value B2C customer', NOW() - INTERVAL '1 year');

-- Insert diverse inventory items (40+ items across categories)
INSERT INTO inventory_lots (
    name, variant, quantity, unit, purchase_date, expiry_date,
    category, sub_category, cost_price, selling_price, supplier_name,
    purchase_invoice_number, min_stock_level, reorder_point, shelf_life_days,
    storage_location, barcode, sku, hsn_code, gst_rate,
    packaging_type, notes, status, created_at
) VALUES
-- Vegetables (Fresh)
('Tomato', 'Organic', 150.50, 'kg', '2025-11-14', '2025-11-18', 'vegetables', 'organic', 35.00, 50.00,
 'Green Farms Organic', 'GFO-2025-1234', 50.0, 75.0, 4, 'Cold Storage A1', '8901234567801', 'VEG-TOM-ORG-001',
 '07019000', 0, 'crate', 'Premium organic tomatoes from Ooty', 'available', NOW() - INTERVAL '2 days'),

('Potato', 'Regular', 280.00, 'kg', '2025-11-10', '2025-12-10', 'vegetables', 'root', 18.00, 28.00,
 'Karnataka Farmers Coop', 'KFC-2025-5678', 100.0, 150.0, 30, 'Warehouse Section B', '8901234567802', 'VEG-POT-REG-001',
 '07019000', 0, 'bag', 'Local farm fresh potatoes', 'available', NOW() - INTERVAL '6 days'),

('Onion', 'Red', 220.00, 'kg', '2025-11-12', '2025-12-12', 'vegetables', 'bulb', 25.00, 40.00,
 'Nashik Traders', 'NSH-2025-9012', 80.0, 120.0, 30, 'Warehouse Section B', '8901234567803', 'VEG-ONI-RED-001',
 '07019000', 0, 'bag', 'Premium red onions from Nashik', 'available', NOW() - INTERVAL '4 days'),

('Carrot', 'Organic', 95.00, 'kg', '2025-11-13', '2025-11-20', 'vegetables', 'root', 40.00, 60.00,
 'Green Farms Organic', 'GFO-2025-1235', 30.0, 50.0, 7, 'Cold Storage A1', '8901234567804', 'VEG-CAR-ORG-001',
 '07019000', 0, 'crate', 'Fresh organic carrots', 'available', NOW() - INTERVAL '3 days'),

('Capsicum', 'Mixed', 45.00, 'kg', '2025-11-15', '2025-11-19', 'vegetables', 'bell_pepper', 80.00, 120.00,
 'Hydro Gardens Ltd', 'HGL-2025-3456', 20.0, 30.0, 4, 'Cold Storage A2', '8901234567805', 'VEG-CAP-MIX-001',
 '07019000', 5, 'box', 'Red, yellow, green capsicum mix', 'available', NOW() - INTERVAL '1 day'),

('Spinach', 'Fresh', 25.00, 'kg', '2025-11-16', '2025-11-18', 'vegetables', 'leafy', 30.00, 50.00,
 'Local Farmers Market', 'LFM-2025-7890', 15.0, 25.0, 2, 'Cold Storage A2', '8901234567806', 'VEG-SPI-FRE-001',
 '07019000', 0, 'bundle', 'Fresh morning spinach', 'available', NOW()),

('Cauliflower', 'Large', 65.00, 'kg', '2025-11-14', '2025-11-21', 'vegetables', 'cruciferous', 35.00, 55.00,
 'Hill Station Produce', 'HSP-2025-4567', 25.0, 40.0, 7, 'Cold Storage A1', '8901234567807', 'VEG-CAU-LAR-001',
 '07019000', 0, 'crate', 'Fresh cauliflower from Ooty', 'available', NOW() - INTERVAL '2 days'),

-- Vegetables (Near Expiry)
('Green Beans', 'Regular', 18.00, 'kg', '2025-11-12', '2025-11-17', 'vegetables', 'legume', 50.00, 75.00,
 'Local Farmers Market', 'LFM-2025-7891', 20.0, 30.0, 5, 'Cold Storage A2', '8901234567808', 'VEG-BEA-REG-001',
 '07019000', 0, 'bundle', 'Near expiry - discount price', 'near_expiry', NOW() - INTERVAL '4 days'),

-- Fruits (Fresh)
('Apple', 'Shimla Deluxe', 120.00, 'kg', '2025-11-10', '2025-12-10', 'fruits', 'pome', 120.00, 180.00,
 'Kashmir Fruit Traders', 'KFT-2025-1111', 40.0, 60.0, 30, 'Cold Storage B1', '8901234567809', 'FRU-APP-SHI-001',
 '08081000', 0, 'box', 'Premium Shimla apples, Grade A', 'available', NOW() - INTERVAL '6 days'),

('Banana', 'Robusta', 180.00, 'dozen', '2025-11-14', '2025-11-19', 'fruits', 'tropical', 30.00, 48.00,
 'Tamil Nadu Banana Co', 'TNBC-2025-2222', 100.0, 150.0, 5, 'Ambient Storage C', '8901234567810', 'FRU-BAN-ROB-001',
 '08030000', 0, 'bunch', 'Fresh Robusta bananas', 'available', NOW() - INTERVAL '2 days'),

('Mango', 'Alphonso', 85.00, 'kg', '2025-11-15', '2025-11-22', 'fruits', 'tropical', 200.00, 320.00,
 'Ratnagiri Mango Suppliers', 'RMS-2025-3333', 30.0, 50.0, 7, 'Cool Storage D1', '8901234567811', 'FRU-MAN-ALP-001',
 '08045000', 5, 'crate', 'Premium Alphonso mangoes', 'available', NOW() - INTERVAL '1 day'),

('Orange', 'Nagpur', 95.00, 'kg', '2025-11-13', '2025-11-27', 'fruits', 'citrus', 45.00, 70.00,
 'Citrus Fresh Ltd', 'CFL-2025-4444', 50.0, 80.0, 14, 'Ambient Storage C', '8901234567812', 'FRU-ORA-NAG-001',
 '08051000', 5, 'net_bag', 'Juicy Nagpur oranges', 'available', NOW() - INTERVAL '3 days'),

('Grapes', 'Black Seedless', 55.00, 'kg', '2025-11-14', '2025-11-21', 'fruits', 'berry', 90.00, 140.00,
 'Maharashtra Grape Growers', 'MGG-2025-5555', 25.0, 40.0, 7, 'Cold Storage B1', '8901234567813', 'FRU-GRA-BLA-001',
 '08061000', 5, 'box', 'Premium black seedless grapes', 'available', NOW() - INTERVAL '2 days'),

('Watermelon', 'Large', 145.00, 'kg', '2025-11-15', '2025-11-25', 'fruits', 'melon', 15.00, 25.00,
 'Karnataka Farmers Coop', 'KFC-2025-5679', 80.0, 120.0, 10, 'Ambient Storage C', '8901234567814', 'FRU-WAT-LAR-001',
 '08071000', 0, 'loose', 'Fresh watermelons', 'available', NOW() - INTERVAL '1 day'),

-- Dairy Products
('Milk', 'Full Cream', 95.00, 'liter', '2025-11-16', '2025-11-18', 'dairy', 'liquid', 48.00, 60.00,
 'Nandini Dairy', 'ND-2025-6666', 100.0, 150.0, 2, 'Refrigerator E1', '8901234567815', 'DAI-MIL-FUL-001',
 '04011000', 5, 'packet', 'Pasteurized full cream milk', 'available', NOW()),

('Curd', 'Fresh', 42.00, 'kg', '2025-11-16', '2025-11-19', 'dairy', 'fermented', 55.00, 75.00,
 'Nandini Dairy', 'ND-2025-6667', 40.0, 60.0, 3, 'Refrigerator E1', '8901234567816', 'DAI-CUR-FRE-001',
 '04031000', 5, 'box', 'Fresh homestyle curd', 'available', NOW()),

('Paneer', 'Premium', 28.00, 'kg', '2025-11-15', '2025-11-20', 'dairy', 'cheese', 280.00, 400.00,
 'Amul Dairy', 'AD-2025-7777', 15.0, 25.0, 5, 'Refrigerator E2', '8901234567817', 'DAI-PAN-PRE-001',
 '04061000', 12, 'vacuum_pack', 'Fresh paneer cubes', 'available', NOW() - INTERVAL '1 day'),

('Butter', 'Salted', 35.00, 'kg', '2025-11-10', '2025-12-25', 'dairy', 'spread', 380.00, 500.00,
 'Amul Dairy', 'AD-2025-7778', 20.0, 30.0, 45, 'Refrigerator E2', '8901234567818', 'DAI-BUT-SAL-001',
 '04051000', 12, 'box', 'Premium salted butter', 'available', NOW() - INTERVAL '6 days'),

-- Grains & Pulses
('Rice', 'Basmati Premium', 450.00, 'kg', '2025-11-01', '2026-05-01', 'grains', 'rice', 65.00, 90.00,
 'India Gate Foods', 'IGF-2025-8888', 200.0, 300.0, 180, 'Dry Storage F1', '8901234567819', 'GRA-RIC-BAS-001',
 '10063000', 5, 'bag', '1121 Basmati rice, aged', 'available', NOW() - INTERVAL '15 days'),

('Wheat Flour', 'Whole Grain', 320.00, 'kg', '2025-11-05', '2026-02-05', 'grains', 'flour', 35.00, 48.00,
 'Aashirvaad', 'ASH-2025-9999', 150.0, 250.0, 90, 'Dry Storage F1', '8901234567820', 'GRA-WHE-WHO-001',
 '11010000', 5, 'bag', 'Whole wheat atta', 'available', NOW() - INTERVAL '11 days'),

('Toor Dal', 'Yellow', 180.00, 'kg', '2025-10-20', '2026-04-20', 'grains', 'pulses', 110.00, 145.00,
 'Karnataka Pulses', 'KP-2025-1010', 100.0, 150.0, 180, 'Dry Storage F2', '8901234567821', 'GRA-TOO-YEL-001',
 '07131000', 5, 'bag', 'Premium toor dal', 'available', NOW() - INTERVAL '27 days'),

('Moong Dal', 'Green Split', 95.00, 'kg', '2025-10-25', '2026-04-25', 'grains', 'pulses', 95.00, 125.00,
 'Karnataka Pulses', 'KP-2025-1011', 60.0, 90.0, 180, 'Dry Storage F2', '8901234567822', 'GRA-MOO-GRE-001',
 '07133000', 5, 'bag', 'Split green moong dal', 'available', NOW() - INTERVAL '22 days'),

-- Meat & Seafood
('Chicken', 'Broiler Fresh', 85.00, 'kg', '2025-11-16', '2025-11-17', 'meat', 'poultry', 160.00, 220.00,
 'Venky''s Poultry', 'VP-2025-1212', 50.0, 75.0, 1, 'Deep Freezer G1', '8901234567823', 'MEA-CHI-BRO-001',
 '02071100', 0, 'tray', 'Fresh broiler chicken, antibiotic-free', 'available', NOW()),

('Mutton', 'Goat Curry Cut', 42.00, 'kg', '2025-11-15', '2025-11-17', 'meat', 'red_meat', 520.00, 700.00,
 'Al Kabeer Meats', 'AKM-2025-1313', 30.0, 50.0, 2, 'Deep Freezer G1', '8901234567824', 'MEA-MUT-GOA-001',
 '02041000', 0, 'vacuum_pack', 'Fresh goat meat, curry cut', 'available', NOW() - INTERVAL '1 day'),

('Fish', 'Pomfret', 35.00, 'kg', '2025-11-16', '2025-11-17', 'seafood', 'fish', 280.00, 400.00,
 'Coastal Seafoods', 'CS-2025-1414', 20.0, 35.0, 1, 'Deep Freezer G2', '8901234567825', 'SEA-FIS-POM-001',
 '03038900', 5, 'ice_pack', 'Fresh pomfret, cleaned', 'available', NOW()),

('Prawns', 'Jumbo', 22.00, 'kg', '2025-11-16', '2025-11-17', 'seafood', 'shellfish', 480.00, 650.00,
 'Coastal Seafoods', 'CS-2025-1415', 15.0, 25.0, 1, 'Deep Freezer G2', '8901234567826', 'SEA-PRA-JUM-001',
 '03061700', 5, 'ice_pack', 'Jumbo prawns, de-veined', 'available', NOW()),

-- Beverages
('Orange Juice', 'Fresh Squeezed', 48.00, 'liter', '2025-11-16', '2025-11-19', 'beverages', 'juice', 85.00, 120.00,
 'Tropicana Fresh', 'TF-2025-1515', 30.0, 50.0, 3, 'Refrigerator E3', '8901234567827', 'BEV-ORA-FRE-001',
 '20091100', 12, 'bottle', 'No preservatives, fresh squeezed', 'available', NOW()),

('Mineral Water', '1L Bottles', 240.00, 'unit', '2025-11-01', '2026-05-01', 'beverages', 'water', 15.00, 20.00,
 'Bisleri', 'BIS-2025-1616', 300.0, 500.0, 180, 'Ambient Storage H', '8901234567828', 'BEV-WAT-MIN-001',
 '22011000', 12, 'case', 'Packaged drinking water', 'available', NOW() - INTERVAL '15 days'),

-- Spices & Condiments
('Turmeric Powder', 'Organic', 45.00, 'kg', '2025-10-15', '2026-10-15', 'spices', 'powder', 280.00, 380.00,
 'Everest Spices', 'ES-2025-1717', 25.0, 40.0, 365, 'Dry Storage I1', '8901234567829', 'SPI-TUR-ORG-001',
 '09103000', 5, 'packet', 'Organic turmeric powder', 'available', NOW() - INTERVAL '32 days'),

('Chili Powder', 'Hot Red', 38.00, 'kg', '2025-10-20', '2026-10-20', 'spices', 'powder', 220.00, 300.00,
 'Everest Spices', 'ES-2025-1718', 20.0, 35.0, 365, 'Dry Storage I1', '8901234567830', 'SPI-CHI-HOT-001',
 '09042200', 5, 'packet', 'Hot red chili powder', 'available', NOW() - INTERVAL '27 days'),

('Salt', 'Iodized', 280.00, 'kg', '2025-09-01', '2027-09-01', 'condiments', 'basic', 12.00, 18.00,
 'Tata Salt', 'TS-2025-1819', 200.0, 300.0, 730, 'Dry Storage I2', '8901234567831', 'CON-SAL-IOD-001',
 '25010000', 5, 'bag', 'Iodized table salt', 'available', NOW() - INTERVAL '76 days'),

-- Cooking Oil
('Sunflower Oil', 'Refined', 125.00, 'liter', '2025-10-01', '2026-10-01', 'oil', 'cooking', 145.00, 185.00,
 'Fortune Foods', 'FF-2025-1919', 80.0, 120.0, 365, 'Dry Storage J1', '8901234567832', 'OIL-SUN-REF-001',
 '15121100', 5, 'can', 'Refined sunflower oil', 'available', NOW() - INTERVAL '46 days'),

('Olive Oil', 'Extra Virgin', 28.00, 'liter', '2025-09-15', '2026-09-15', 'oil', 'premium', 650.00, 850.00,
 'Figaro', 'FIG-2025-2020', 20.0, 30.0, 365, 'Dry Storage J1', '8901234567833', 'OIL-OLI-EXT-001',
 '15091000', 12, 'bottle', 'Extra virgin olive oil, imported', 'available', NOW() - INTERVAL '62 days'),

-- Snacks & Packaged
('Potato Chips', 'Salted', 180.00, 'packet', '2025-10-10', '2026-02-10', 'snacks', 'chips', 25.00, 40.00,
 'Lay''s', 'LAYS-2025-2121', 150.0, 250.0, 120, 'Dry Storage K1', '8901234567834', 'SNA-CHI-SAL-001',
 '20052000', 12, 'box', '50g packets', 'available', NOW() - INTERVAL '37 days'),

-- Bakery
('Bread', 'Whole Wheat', 65.00, 'unit', '2025-11-16', '2025-11-19', 'bakery', 'bread', 32.00, 45.00,
 'Britannia', 'BR-2025-2222', 80.0, 120.0, 3, 'Ambient Storage L', '8901234567835', 'BAK-BRE-WHO-001',
 '19051000', 5, 'packet', 'Whole wheat bread loaf', 'available', NOW()),

-- Expired/Damaged Items
('Mixed Vegetables', 'Frozen', 8.50, 'kg', '2025-10-15', '2025-11-15', 'vegetables', 'frozen', 90.00, 130.00,
 'McCain Foods', 'MC-2025-2323', 25.0, 40.0, 30, 'Deep Freezer G3', '8901234567836', 'VEG-MIX-FRO-001',
 '07108000', 12, 'bag', 'EXPIRED - Move to wastage', 'expired', NOW() - INTERVAL '31 days'),

('Yogurt', 'Flavored', 0.00, 'kg', '2025-11-10', '2025-11-15', 'dairy', 'fermented', 60.00, 85.00,
 'Nestle', 'NES-2025-2424', 30.0, 50.0, 5, 'Refrigerator E1', '8901234567837', 'DAI-YOG-FLA-001',
 '04031000', 12, 'cup', 'DAMAGED - Container leaked', 'damaged', NOW() - INTERVAL '6 days');

-- Insert sample transactions (30+ transactions showing various scenarios)
-- B2B Credit Sales
INSERT INTO transactions (customer_id, type, total_amount, date, payment_method, payment_reference, 
    due_date, discount_amount, tax_amount, sale_type, delivery_status, delivery_date, delivery_address, 
    invoice_number, notes, created_at)
VALUES
(2, 'sale', 24500.00, '2025-11-01', 'bank_transfer', 'NEFT202511011234', '2025-12-16', 500.00, 3000.00, 'credit', 'delivered', '2025-11-02',
 '45 Brigade Road, Retail District, Bangalore', 'INV-2025-0001', 'Bulk order - monthly supply', NOW() - INTERVAL '15 days'),

(1, 'sale', 18750.00, '2025-11-05', 'cheque', 'CHQ456789', '2025-12-05', 250.00, 2500.00, 'credit', 'delivered', '2025-11-06',
 'Shop 15, MG Road, Commercial Complex, Bangalore', 'INV-2025-0002', 'Weekly order', NOW() - INTERVAL '11 days'),

(3, 'sale', 12400.00, '2025-11-10', 'upi', 'UPI202511101234567', '2025-11-25', 0.00, 1600.00, 'credit', 'delivered', '2025-11-11',
 '12 Church Street, CBD, Bangalore', 'INV-2025-0003', 'Restaurant supplies', NOW() - INTERVAL '6 days'),

-- Payments received
(2, 'payment', 28900.50, '2025-11-14', 'bank_transfer', 'NEFT202511141234', NULL, 0.00, 0.00, NULL, NULL, NULL, NULL,
 NULL, 'Payment for pending invoices', NOW() - INTERVAL '2 days'),

(1, 'payment', 10000.00, '2025-11-12', 'cheque', 'CHQ789012', NULL, 0.00, 0.00, NULL, NULL, NULL, NULL,
 NULL, 'Partial payment', NOW() - INTERVAL '4 days'),

-- Wholesale transactions
(6, 'sale', 32500.00, '2025-11-03', 'bank_transfer', 'RTGS202511031234', '2025-11-24', 1500.00, 4000.00, 'credit', 'delivered', '2025-11-04',
 '23 Avenue Road, Market Area, Bangalore', 'INV-2025-0004', 'Bi-weekly bulk order', NOW() - INTERVAL '13 days'),

(7, 'sale', 45800.00, '2025-11-08', 'upi', 'UPI202511082345678', '2025-12-08', 2200.00, 6000.00, 'credit', 'in_transit', NULL,
 '67 KR Puram, Wholesale District, Bangalore', 'INV-2025-0005', 'Large wholesale order - in transit', NOW() - INTERVAL '8 days'),

(8, 'sale', 67500.00, '2025-11-12', 'bank_transfer', 'NEFT202511121234', '2026-01-12', 2500.00, 8500.00, 'credit', 'delivered', '2025-11-13',
 '100 Whitefield, Industrial Area, Bangalore', 'INV-2025-0006', 'Monthly bulk supply', NOW() - INTERVAL '4 days'),

(6, 'payment', 15000.00, '2025-11-15', 'upi', 'UPI202511153456789', NULL, 0.00, 0.00, NULL, NULL, NULL, NULL,
 NULL, 'Advance payment for next order', NOW() - INTERVAL '1 day'),

-- Retail cash sales
(9, 'sale', 850.00, '2025-11-14', 'cash', NULL, NULL, 0.00, 50.00, 'cash', 'delivered', '2025-11-14',
 'Flat 204, Prestige Apartments, Koramangala, Bangalore', 'INV-2025-0007', 'Weekly grocery delivery', NOW() - INTERVAL '2 days'),

(10, 'sale', 1250.00, '2025-11-15', 'upi', 'UPI202511154567890', NULL, 0.00, 100.00, 'cash', 'delivered', '2025-11-15',
 '56 Jayanagar 4th Block, Bangalore', 'INV-2025-0008', 'Weekly shopping', NOW() - INTERVAL '1 day'),

(11, 'sale', 620.00, '2025-11-16', 'cash', NULL, NULL, 20.00, 40.00, 'cash', 'pending', NULL,
 '89 RT Nagar, Near Bus Stand, Bangalore', 'INV-2025-0009', 'Walk-in purchase - pickup', NOW()),

(12, 'sale', 2150.00, '2025-11-13', 'upi', 'UPI202511135678901', NULL, 50.00, 200.00, 'cash', 'delivered', '2025-11-13',
 '34 Malleshwaram, Temple Street, Bangalore', 'INV-2025-0010', 'Home delivery', NOW() - INTERVAL '3 days'),

-- Retail credit sales
(12, 'sale', 1500.00, '2025-11-08', 'upi', 'UPI202511086789012', '2025-11-15', 0.00, 150.00, 'credit', 'delivered', '2025-11-08',
 '34 Malleshwaram, Temple Street, Bangalore', 'INV-2025-0011', '7-day credit sale', NOW() - INTERVAL '8 days'),

(10, 'sale', 950.00, '2025-11-10', 'card', 'CARD202511107890123', '2025-11-17', 0.00, 80.00, 'credit', 'delivered', '2025-11-10',
 '56 Jayanagar 4th Block, Bangalore', 'INV-2025-0012', 'Credit sale - due soon', NOW() - INTERVAL '6 days'),

-- B2C online orders
(16, 'sale', 3200.00, '2025-11-14', 'upi', 'UPI202511148901234', '2025-11-28', 200.00, 400.00, 'credit', 'delivered', '2025-11-15',
 '23 Bellandur, Lake View Apartments, Bangalore', 'INV-2025-0013', 'Online app order - subscription', NOW() - INTERVAL '2 days'),

(17, 'sale', 1850.00, '2025-11-15', 'card', 'CARD202511159012345', NULL, 50.00, 200.00, 'cash', 'in_transit', NULL,
 '67 Marathahalli, Tech Hub, Bangalore', 'INV-2025-0014', 'Online order - in transit', NOW() - INTERVAL '1 day'),

(18, 'sale', 2650.00, '2025-11-12', 'upi', 'UPI202511120123456', '2025-11-19', 150.00, 300.00, 'credit', 'delivered', '2025-11-13',
 '90 JP Nagar, Garden City, Bangalore', 'INV-2025-0015', 'Weekly online order', NOW() - INTERVAL '4 days'),

(19, 'sale', 890.00, '2025-11-16', 'upi', 'UPI202511161234567', NULL, 0.00, 90.00, 'cash', 'pending', NULL,
 '11 Banashankari, Temple Road, Bangalore', 'INV-2025-0016', 'First order - new customer', NOW()),

(20, 'sale', 4200.00, '2025-11-11', 'card', 'CARD202511112345678', '2025-11-26', 300.00, 500.00, 'credit', 'delivered', '2025-11-12',
 '88 Yelahanka, New Town, Bangalore', 'INV-2025-0017', 'Premium B2C order', NOW() - INTERVAL '5 days'),

-- More B2B transactions
(4, 'sale', 8900.00, '2025-11-15', 'bank_transfer', 'NEFT202511153456', '2025-11-22', 100.00, 1100.00, 'credit', 'delivered', '2025-11-16',
 '78 Indiranagar, Food Street, Bangalore', 'INV-2025-0018', 'Weekly cafe supply', NOW() - INTERVAL '1 day'),

(4, 'payment', 5000.00, '2025-11-16', 'upi', 'UPI202511164567890', NULL, 0.00, 0.00, NULL, NULL, NULL, NULL,
 NULL, 'Partial payment on account', NOW()),

-- Overdue credit sale (blocked customer)
(5, 'sale', 3500.00, '2025-10-01', 'cheque', 'CHQ345678', '2025-10-31', 0.00, 450.00, 'credit', 'delivered', '2025-10-02',
 '90 KR Market, Wholesale Area, Bangalore', 'INV-2025-0019', 'OVERDUE - Customer blocked', NOW() - INTERVAL '46 days'),

-- Recent mixed transactions
(13, 'sale', 750.00, '2025-11-16', 'cash', NULL, NULL, 0.00, 50.00, 'cash', 'delivered', '2025-11-16',
 '12 HSR Layout, Tech Park Road, Bangalore', 'INV-2025-0020', 'Weekend shopping', NOW()),

(9, 'sale', 1450.00, '2025-11-16', 'upi', 'UPI202511165678901', NULL, 50.00, 120.00, 'cash', 'pending', NULL,
 'Flat 204, Prestige Apartments, Koramangala, Bangalore', 'INV-2025-0021', 'Order for delivery tomorrow', NOW()),

(7, 'sale', 38500.00, '2025-11-16', 'bank_transfer', 'RTGS202511163456', '2025-12-16', 1500.00, 5000.00, 'credit', 'pending', NULL,
 '67 KR Puram, Wholesale District, Bangalore', 'INV-2025-0022', 'New bulk order - processing', NOW()),

(16, 'payment', 1500.00, '2025-11-16', 'upi', 'UPI202511166789012', NULL, 0.00, 0.00, NULL, NULL, NULL, NULL,
 NULL, 'Payment for previous order', NOW()),

-- High value transactions
(8, 'sale', 95000.00, '2025-11-14', 'bank_transfer', 'RTGS202511144567', '2026-01-14', 5000.00, 12000.00, 'credit', 'in_transit', NULL,
 '100 Whitefield, Industrial Area, Bangalore', 'INV-2025-0023', 'Large wholesale order - premium items', NOW() - INTERVAL '2 days'),

(2, 'sale', 42500.00, '2025-11-15', 'cheque', 'CHQ901234', '2025-12-30', 2500.00, 5500.00, 'credit', 'delivered', '2025-11-16',
 '45 Brigade Road, Retail District, Bangalore', 'INV-2025-0024', 'Monthly bulk supply - delivered', NOW() - INTERVAL '1 day');

-- Insert transaction items for the sales (sample items for first few transactions)
INSERT INTO transaction_items (transaction_id, inventory_lot_id, item_name, quantity, price_per_unit, unit)
VALUES
-- Transaction 1 (B2B - FreshMart)
(1, 1, 'Tomato - Organic', 50.0, 50.00, 'kg'),
(1, 2, 'Potato - Regular', 100.0, 28.00, 'kg'),
(1, 3, 'Onion - Red', 80.0, 40.00, 'kg'),
(1, 9, 'Apple - Shimla Deluxe', 40.0, 180.00, 'kg'),
(1, 10, 'Banana - Robusta', 60.0, 48.00, 'dozen'),

-- Transaction 3 (B2B - Green Valley Restaurant)
(3, 1, 'Tomato - Organic', 30.0, 50.00, 'kg'),
(3, 4, 'Carrot - Organic', 25.0, 60.00, 'kg'),
(3, 5, 'Capsicum - Mixed', 15.0, 120.00, 'kg'),
(3, 17, 'Milk - Full Cream', 40.0, 60.00, 'liter'),
(3, 19, 'Paneer - Premium', 10.0, 400.00, 'kg'),

-- Transaction 6 (Wholesale - Ravi Trading)
(6, 2, 'Potato - Regular', 150.0, 28.00, 'kg'),
(6, 3, 'Onion - Red', 120.0, 40.00, 'kg'),
(6, 21, 'Rice - Basmati Premium', 80.0, 90.00, 'kg'),
(6, 22, 'Wheat Flour - Whole Grain', 60.0, 48.00, 'kg'),

-- Transaction 10 (Retail - Rajesh Kumar)
(10, 1, 'Tomato - Organic', 3.0, 50.00, 'kg'),
(10, 2, 'Potato - Regular', 5.0, 28.00, 'kg'),
(10, 3, 'Onion - Red', 3.0, 40.00, 'kg'),
(10, 10, 'Banana - Robusta', 2.0, 48.00, 'dozen'),
(10, 17, 'Milk - Full Cream', 4.0, 60.00, 'liter'),

-- Transaction 11 (Retail - Priya Sharma)
(11, 2, 'Potato - Regular', 4.0, 28.00, 'kg'),
(11, 6, 'Spinach - Fresh', 2.0, 50.00, 'kg'),
(11, 9, 'Apple - Shimla Deluxe', 2.0, 180.00, 'kg'),
(11, 18, 'Curd - Fresh', 2.0, 75.00, 'kg'),

-- Transaction 13 (Retail - Mohammed Ali)
(13, 10, 'Banana - Robusta', 3.0, 48.00, 'dozen'),
(13, 17, 'Milk - Full Cream', 3.0, 60.00, 'liter'),
(13, 37, 'Bread - Whole Wheat', 4.0, 45.00, 'unit');

-- Insert wastage log entries
INSERT INTO wastage_log (inventory_lot_id, item_name, quantity_wasted, unit, wastage_reason, 
    wastage_date, reported_by, estimated_loss, notes)
VALUES
(38, 'Mixed Vegetables - Frozen', 8.50, 'kg', 'expired', '2025-11-15', 1, 765.00, 
 'Expired frozen vegetables - moved from inventory'),
(39, 'Yogurt - Flavored', 12.00, 'kg', 'damaged', '2025-11-15', 1, 720.00, 
 'Container leaked during transport - full batch damaged'),
(8, 'Green Beans - Regular', 2.00, 'kg', 'spoilage', '2025-11-16', 1, 100.00, 
 'Started rotting - removed before full spoilage'),
(6, 'Spinach - Fresh', 3.00, 'kg', 'spoilage', '2025-11-17', 1, 90.00, 
 'Wilted leaves - quality compromised'),
(26, 'Chicken - Broiler Fresh', 5.50, 'kg', 'spoilage', '2025-11-16', 1, 880.00, 
 'Temperature fluctuation in freezer - had to discard');

-- Insert expiry alerts
INSERT INTO expiry_alerts (inventory_lot_id, item_name, expiry_date, days_until_expiry, 
    alert_level, quantity_remaining, unit, created_at)
VALUES
(6, 'Spinach - Fresh', '2025-11-18', 2, 'critical', 25.00, 'kg', NOW()),
(1, 'Tomato - Organic', '2025-11-18', 2, 'critical', 150.50, 'kg', NOW()),
(5, 'Capsicum - Mixed', '2025-11-19', 3, 'warning', 45.00, 'kg', NOW()),
(17, 'Milk - Full Cream', '2025-11-18', 2, 'critical', 95.00, 'liter', NOW()),
(18, 'Curd - Fresh', '2025-11-19', 3, 'warning', 42.00, 'kg', NOW()),
(10, 'Banana - Robusta', '2025-11-19', 3, 'warning', 180.00, 'dozen', NOW()),
(19, 'Paneer - Premium', '2025-11-20', 4, 'info', 28.00, 'kg', NOW()),
(26, 'Chicken - Broiler Fresh', '2025-11-17', 1, 'critical', 85.00, 'kg', NOW()),
(27, 'Mutton - Goat Curry Cut', '2025-11-17', 1, 'critical', 42.00, 'kg', NOW()),
(28, 'Fish - Pomfret', '2025-11-17', 1, 'critical', 35.00, 'kg', NOW()),
(29, 'Prawns - Jumbo', '2025-11-17', 1, 'critical', 22.00, 'kg', NOW()),
(37, 'Bread - Whole Wheat', '2025-11-19', 3, 'warning', 65.00, 'unit', NOW());

-- Insert payment schedules (for overdue and upcoming payments)
INSERT INTO payment_schedules (customer_id, transaction_id, invoice_number, amount_due, 
    due_date, payment_status, reminder_sent, notes, created_at)
VALUES
(5, 25, 'INV-2025-0019', 3500.00, '2025-10-31', 'overdue', true, 
 'OVERDUE by 17 days - Customer blocked', NOW() - INTERVAL '17 days'),
(12, 15, 'INV-2025-0011', 1500.00, '2025-11-15', 'overdue', true, 
 'OVERDUE by 1 day - Follow up needed', NOW() - INTERVAL '8 days'),
(10, 16, 'INV-2025-0012', 950.00, '2025-11-17', 'pending', false, 
 'Due tomorrow', NOW() - INTERVAL '6 days'),
(3, 3, 'INV-2025-0003', 12400.00, '2025-11-25', 'pending', false, 
 'Due in 9 days', NOW() - INTERVAL '6 days'),
(16, 13, 'INV-2025-0013', 3200.00, '2025-11-28', 'pending', false, 
 'Subscription customer - usually pays on time', NOW() - INTERVAL '2 days'),
(18, 19, 'INV-2025-0015', 2650.00, '2025-11-19', 'pending', false, 
 'Due in 3 days', NOW() - INTERVAL '4 days'),
(20, 21, 'INV-2025-0017', 4200.00, '2025-11-26', 'pending', false, 
 'Premium customer - credit sale', NOW() - INTERVAL '5 days'),
(4, 22, 'INV-2025-0018', 8900.00, '2025-11-22', 'pending', false, 
 'Weekly cafe order', NOW() - INTERVAL '1 day');

-- Add some pricing tier examples
INSERT INTO pricing_tiers (customer_id, item_category, min_quantity, max_quantity, 
    discount_percentage, effective_from, effective_to, is_active)
VALUES
(2, 'vegetables', 100.0, NULL, 10.0, '2025-01-01', '2025-12-31', true),
(2, 'fruits', 50.0, NULL, 8.0, '2025-01-01', '2025-12-31', true),
(7, 'grains', 200.0, NULL, 15.0, '2025-01-01', '2025-12-31', true),
(8, 'vegetables', 150.0, NULL, 12.0, '2025-01-01', '2025-12-31', true),
(8, 'dairy', 50.0, NULL, 10.0, '2025-01-01', '2025-12-31', true);

-- Add price history for trending
INSERT INTO price_history (inventory_lot_id, item_name, cost_price, selling_price, 
    effective_date, recorded_by)
VALUES
(1, 'Tomato - Organic', 30.00, 45.00, '2025-10-01', 1),
(1, 'Tomato - Organic', 32.00, 48.00, '2025-10-15', 1),
(1, 'Tomato - Organic', 35.00, 50.00, '2025-11-01', 1),
(2, 'Potato - Regular', 15.00, 25.00, '2025-10-01', 1),
(2, 'Potato - Regular', 16.00, 26.00, '2025-10-20', 1),
(2, 'Potato - Regular', 18.00, 28.00, '2025-11-01', 1),
(3, 'Onion - Red', 20.00, 35.00, '2025-10-01', 1),
(3, 'Onion - Red', 23.00, 38.00, '2025-10-20', 1),
(3, 'Onion - Red', 25.00, 40.00, '2025-11-01', 1);

-- Update some balances to reflect transactions
UPDATE customers SET current_balance = 12500.75 WHERE id = 1;
UPDATE customers SET current_balance = 0.00 WHERE id = 2;
UPDATE customers SET current_balance = 0.00 WHERE id = 3;
UPDATE customers SET current_balance = 8200.00 - 5000.00 WHERE id = 4;
UPDATE customers SET current_balance = 26500.00 WHERE id = 5;
UPDATE customers SET current_balance = 15600.00 + 15000.00 WHERE id = 6;
UPDATE customers SET current_balance = 34200.25 WHERE id = 7;
UPDATE customers SET current_balance = 42000.00 WHERE id = 8;
UPDATE customers SET current_balance = 0.00 WHERE id = 9;
UPDATE customers SET current_balance = 450.00 + 950.00 WHERE id = 10;
UPDATE customers SET current_balance = 0.00 WHERE id = 11;
UPDATE customers SET current_balance = 1200.00 + 1500.00 WHERE id = 12;
UPDATE customers SET current_balance = 0.00 WHERE id = 13;
UPDATE customers SET current_balance = 1650.00 WHERE id = 14;
UPDATE customers SET current_balance = 0.00 WHERE id = 15;
UPDATE customers SET current_balance = 2400.00 - 1500.00 WHERE id = 16;
UPDATE customers SET current_balance = 0.00 WHERE id = 17;
UPDATE customers SET current_balance = 1800.50 + 2650.00 WHERE id = 18;
UPDATE customers SET current_balance = 0.00 WHERE id = 19;
UPDATE customers SET current_balance = 3500.00 + 4200.00 WHERE id = 20;

ANALYZE;
