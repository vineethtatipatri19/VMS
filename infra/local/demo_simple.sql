-- Simple Demo Data for VMS - Vegetables & Fruits Focus
-- Clean data first
TRUNCATE TABLE sale_items, transactions, inventory_items, customers CASCADE;

-- Insert diverse customers (15 customers)
INSERT INTO customers (
    name, address, contact_number, aadhaar_verified,
    email, whatsapp_number, alternate_contact, customer_type, business_name,
    gstin, credit_limit, current_balance, payment_terms_days, status,
    kyc_document_type, kyc_document_number, notes
) VALUES
-- B2B Customers
('Fresh Mart Supermarket', '45 Brigade Road, Bangalore', '9876543210', true,
 'orders@freshmart.in', '9876543210', '9876543211', 'b2b', 'Fresh Mart Retail Pvt Ltd',
 '29XYZAB5678C1D3', 75000.00, 18500.00, 45, 'active',
 'gstin', '29XYZAB5678C1D3', 'Large retail chain, weekly bulk orders'),

('Green Valley Restaurant', '12 Church Street, Bangalore', '9876543220', true,
 'chef@greenvalley.com', '9876543220', NULL, 'b2b', 'Green Valley Food Services',
 '29PQRST9012E2F4', 35000.00, 8200.00, 15, 'active',
 'gstin', '29PQRST9012E2F4', 'Restaurant chain, daily fresh produce'),

('Organic Hub', '78 Koramangala, Bangalore', '9876543230', true,
 'buy@organichub.com', '9876543230', NULL, 'b2b', 'Organic Hub Pvt Ltd',
 '29ABCDE1234F1Z5', 50000.00, 0.00, 30, 'active',
 'gstin', '29ABCDE1234F1Z5', 'Organic specialty store'),

-- Wholesale Customers
('Ravi Wholesale Trading', '23 Avenue Road, Bangalore', '9876543240', true,
 'ravi@trading.com', '9876543240', '9876543241', 'wholesale', 'Ravi & Sons Trading Co',
 '29QWERT7890I4J6', 60000.00, 24600.00, 30, 'active',
 'gstin', '29QWERT7890I4J6', 'Wholesale dealer for local markets'),

('Metro Vegetables', '100 KR Market, Bangalore', '9876543250', true,
 'metro@veggies.com', '9876543250', NULL, 'wholesale', 'Metro Vegetables & Fruits',
 '29ASDFG4567K5L7', 80000.00, 12000.00, 21, 'active',
 'gstin', '29ASDFG4567K5L7', 'Large wholesale distributor'),

-- Retail Customers
('Rajesh Kumar', 'Flat 204, Koramangala, Bangalore', '9876543260', true,
 'rajesh.k@email.com', '9876543260', NULL, 'retail', NULL,
 NULL, 5000.00, 850.00, 7, 'active',
 'aadhaar', '1234-5678-9012', 'Regular customer, weekly orders'),

('Priya Sharma', '56 Jayanagar, Bangalore', '9876543270', true,
 'priya.s@email.com', '9876543270', NULL, 'retail', NULL,
 NULL, 3000.00, 0.00, 0, 'active',
 'aadhaar', '2345-6789-0123', 'Prefers organic produce'),

('Mohammed Ali', '89 RT Nagar, Bangalore', '9876543280', false,
 'ali.m@email.com', '9876543280', NULL, 'retail', NULL,
 NULL, 2000.00, 450.00, 0, 'active',
 'pan', 'ABCDE1234F', 'Walk-in customer'),

('Lakshmi Devi', '34 Malleshwaram, Bangalore', '9876543290', true,
 'lakshmi.d@email.com', '9876543290', NULL, 'retail', NULL,
 NULL, 4000.00, 1200.00, 7, 'active',
 'aadhaar', '3456-7890-1234', 'Senior citizen, home delivery'),

('Arjun Reddy', '12 HSR Layout, Bangalore', '9876543300', true,
 'arjun.r@email.com', '9876543300', NULL, 'retail', NULL,
 NULL, 2500.00, 0.00, 0, 'active',
 'pan', 'FGHIJ5678K', 'Young professional'),

-- B2C Customers  
('Ananya Iyer', '23 Bellandur, Bangalore', '9876543310', true,
 'ananya.i@email.com', '9876543310', NULL, 'b2c', NULL,
 NULL, 8000.00, 2400.00, 14, 'active',
 'aadhaar', '5678-9012-3456', 'Premium online customer'),

('Vikram Singh', '67 Marathahalli, Bangalore', '9876543320', true,
 'vikram.s@email.com', '9876543320', NULL, 'b2c', NULL,
 NULL, 6000.00, 1800.00, 10, 'active',
 'pan', 'KLMNO6789P', 'App user, weekly orders'),

-- Blocked/Inactive
('Budget Store', '90 KR Market, Bangalore', '9876543330', true,
 'accounts@budgetstore.com', '9876543330', NULL, 'b2b', 'Budget Store Trading',
 '29LMNOP3456G3H5', 25000.00, 26500.00, 30, 'blocked',
 'gstin', '29LMNOP3456G3H5', 'BLOCKED: Payment overdue'),

('Suresh Patel', '78 Yeshwanthpur, Bangalore', '9876543340', false,
 NULL, '9876543340', NULL, 'retail', NULL,
 NULL, 1500.00, 1650.00, 0, 'blocked',
 'voter_id', 'ABC1234567', 'BLOCKED: Overdue payment'),

('Meera Krishnan', '45 Electronic City, Bangalore', '9876543350', true,
 'meera.k@email.com', '9876543350', NULL, 'retail', NULL,
 NULL, 3000.00, 0.00, 0, 'inactive',
 'aadhaar', '6789-0123-4567', 'Inactive: No orders in 6 months');

-- Insert vegetables inventory (25 items)
INSERT INTO inventory_items (
    name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
    category, sub_category, cost_price, selling_price, supplier_name,
    purchase_invoice, min_stock_level, reorder_point, shelf_life_days,
    storage_location, barcode, sku, hsn_code, gst_rate,
    packaging_type, status
) VALUES
-- Fresh Vegetables
('Tomato', 'Organic', 'LOT-TOM-001', 150.50, 'kg', '2025-11-14', '2025-11-19',
 'vegetables', 'organic', 35.00, 55.00, 'Green Farms Organic',
 'GFO-2025-1234', 50.0, 75.0, 5, 'Cold Storage A1', '8901234567801', 'VEG-TOM-ORG-001',
 '07019000', 0, 'crate', 'available'),

('Tomato', 'Regular', 'LOT-TOM-002', 200.00, 'kg', '2025-11-15', '2025-11-20',
 'vegetables', 'regular', 28.00, 45.00, 'Local Farmers Market',
 'LFM-2025-7890', 80.0, 120.0, 5, 'Cold Storage A1', '8901234567802', 'VEG-TOM-REG-001',
 '07019000', 0, 'crate', 'available'),

('Potato', 'Regular', 'LOT-POT-001', 300.00, 'kg', '2025-11-10', '2025-12-10',
 'vegetables', 'root', 18.00, 30.00, 'Karnataka Farmers Coop',
 'KFC-2025-5678', 100.0, 150.0, 30, 'Warehouse B1', '8901234567803', 'VEG-POT-REG-001',
 '07019000', 0, 'bag', 'available'),

('Onion', 'Red', 'LOT-ONI-001', 250.00, 'kg', '2025-11-12', '2025-12-12',
 'vegetables', 'bulb', 25.00, 42.00, 'Nashik Traders',
 'NSH-2025-9012', 80.0, 120.0, 30, 'Warehouse B1', '8901234567804', 'VEG-ONI-RED-001',
 '07019000', 0, 'bag', 'available'),

('Onion', 'White', 'LOT-ONI-002', 120.00, 'kg', '2025-11-13', '2025-12-13',
 'vegetables', 'bulb', 30.00, 48.00, 'Nashik Traders',
 'NSH-2025-9013', 50.0, 80.0, 30, 'Warehouse B1', '8901234567805', 'VEG-ONI-WHI-001',
 '07019000', 0, 'bag', 'available'),

('Carrot', 'Organic', 'LOT-CAR-001', 95.00, 'kg', '2025-11-13', '2025-11-21',
 'vegetables', 'root', 40.00, 65.00, 'Green Farms Organic',
 'GFO-2025-1235', 30.0, 50.0, 8, 'Cold Storage A1', '8901234567806', 'VEG-CAR-ORG-001',
 '07019000', 0, 'crate', 'available'),

('Carrot', 'Regular', 'LOT-CAR-002', 140.00, 'kg', '2025-11-14', '2025-11-22',
 'vegetables', 'root', 32.00, 52.00, 'Local Farmers Market',
 'LFM-2025-7891', 50.0, 80.0, 8, 'Cold Storage A1', '8901234567807', 'VEG-CAR-REG-001',
 '07019000', 0, 'crate', 'available'),

('Capsicum', 'Green', 'LOT-CAP-001', 55.00, 'kg', '2025-11-15', '2025-11-20',
 'vegetables', 'bell_pepper', 70.00, 110.00, 'Hydro Gardens Ltd',
 'HGL-2025-3456', 20.0, 30.0, 5, 'Cold Storage A2', '8901234567808', 'VEG-CAP-GRE-001',
 '07019000', 5, 'box', 'available'),

('Capsicum', 'Red', 'LOT-CAP-002', 42.00, 'kg', '2025-11-15', '2025-11-20',
 'vegetables', 'bell_pepper', 90.00, 135.00, 'Hydro Gardens Ltd',
 'HGL-2025-3457', 20.0, 30.0, 5, 'Cold Storage A2', '8901234567809', 'VEG-CAP-RED-001',
 '07019000', 5, 'box', 'available'),

('Capsicum', 'Yellow', 'LOT-CAP-003', 38.00, 'kg', '2025-11-15', '2025-11-20',
 'vegetables', 'bell_pepper', 85.00, 130.00, 'Hydro Gardens Ltd',
 'HGL-2025-3458', 20.0, 30.0, 5, 'Cold Storage A2', '8901234567810', 'VEG-CAP-YEL-001',
 '07019000', 5, 'box', 'available'),

('Spinach', 'Fresh', 'LOT-SPI-001', 28.00, 'kg', '2025-11-16', '2025-11-18',
 'vegetables', 'leafy', 30.00, 52.00, 'Local Farmers Market',
 'LFM-2025-7892', 15.0, 25.0, 2, 'Cold Storage A2', '8901234567811', 'VEG-SPI-FRE-001',
 '07019000', 0, 'bundle', 'available'),

('Cabbage', 'Regular', 'LOT-CAB-001', 80.00, 'kg', '2025-11-14', '2025-11-22',
 'vegetables', 'cruciferous', 25.00, 40.00, 'Hill Station Produce',
 'HSP-2025-4567', 30.0, 50.0, 8, 'Cold Storage A1', '8901234567812', 'VEG-CAB-REG-001',
 '07019000', 0, 'crate', 'available'),

('Cauliflower', 'Large', 'LOT-CAU-001', 65.00, 'kg', '2025-11-14', '2025-11-22',
 'vegetables', 'cruciferous', 35.00, 58.00, 'Hill Station Produce',
 'HSP-2025-4568', 25.0, 40.0, 8, 'Cold Storage A1', '8901234567813', 'VEG-CAU-LAR-001',
 '07019000', 0, 'crate', 'available'),

('Broccoli', 'Fresh', 'LOT-BRO-001', 32.00, 'kg', '2025-11-15', '2025-11-21',
 'vegetables', 'cruciferous', 80.00, 125.00, 'Hydro Gardens Ltd',
 'HGL-2025-3459', 15.0, 25.0, 6, 'Cold Storage A2', '8901234567814', 'VEG-BRO-FRE-001',
 '07019000', 5, 'box', 'available'),

('Cucumber', 'Fresh', 'LOT-CUC-001', 72.00, 'kg', '2025-11-15', '2025-11-21',
 'vegetables', 'gourd', 22.00, 38.00, 'Local Farmers Market',
 'LFM-2025-7893', 30.0, 50.0, 6, 'Cold Storage A2', '8901234567815', 'VEG-CUC-FRE-001',
 '07019000', 0, 'crate', 'available'),

('Beetroot', 'Fresh', 'LOT-BEE-001', 58.00, 'kg', '2025-11-13', '2025-11-23',
 'vegetables', 'root', 35.00, 55.00, 'Local Farmers Market',
 'LFM-2025-7894', 25.0, 40.0, 10, 'Warehouse B2', '8901234567816', 'VEG-BEE-FRE-001',
 '07019000', 0, 'bag', 'available'),

('Pumpkin', 'Regular', 'LOT-PUM-001', 145.00, 'kg', '2025-11-12', '2025-12-02',
 'vegetables', 'gourd', 15.00, 28.00, 'Karnataka Farmers Coop',
 'KFC-2025-5679', 60.0, 90.0, 20, 'Warehouse B2', '8901234567817', 'VEG-PUM-REG-001',
 '07019000', 0, 'loose', 'available'),

('Bottle Gourd', 'Fresh', 'LOT-BOT-001', 85.00, 'kg', '2025-11-14', '2025-11-24',
 'vegetables', 'gourd', 18.00, 32.00, 'Local Farmers Market',
 'LFM-2025-7895', 40.0, 60.0, 10, 'Warehouse B2', '8901234567818', 'VEG-BOT-FRE-001',
 '07019000', 0, 'loose', 'available'),

('Bitter Gourd', 'Fresh', 'LOT-BIT-001', 42.00, 'kg', '2025-11-15', '2025-11-21',
 'vegetables', 'gourd', 40.00, 62.00, 'Local Farmers Market',
 'LFM-2025-7896', 20.0, 35.0, 6, 'Cold Storage A2', '8901234567819', 'VEG-BIT-FRE-001',
 '07019000', 0, 'crate', 'available'),

('Ladies Finger', 'Fresh', 'LOT-LAD-001', 38.00, 'kg', '2025-11-16', '2025-11-19',
 'vegetables', 'pod', 45.00, 72.00, 'Local Farmers Market',
 'LFM-2025-7897', 20.0, 30.0, 3, 'Cold Storage A2', '8901234567820', 'VEG-LAD-FRE-001',
 '07019000', 0, 'bundle', 'available'),

('Green Beans', 'Fresh', 'LOT-BEA-001', 22.00, 'kg', '2025-11-15', '2025-11-19',
 'vegetables', 'legume', 50.00, 78.00, 'Local Farmers Market',
 'LFM-2025-7898', 20.0, 30.0, 4, 'Cold Storage A2', '8901234567821', 'VEG-BEA-FRE-001',
 '07019000', 0, 'bundle', 'available'),

('Radish', 'White', 'LOT-RAD-001', 65.00, 'kg', '2025-11-14', '2025-11-24',
 'vegetables', 'root', 20.00, 35.00, 'Local Farmers Market',
 'LFM-2025-7899', 30.0, 50.0, 10, 'Cold Storage A1', '8901234567822', 'VEG-RAD-WHI-001',
 '07019000', 0, 'bundle', 'available'),

('Ginger', 'Fresh', 'LOT-GIN-001', 48.00, 'kg', '2025-11-10', '2025-12-10',
 'vegetables', 'rhizome', 80.00, 120.00, 'Kerala Spices',
 'KS-2025-1111', 25.0, 40.0, 30, 'Dry Storage C1', '8901234567823', 'VEG-GIN-FRE-001',
 '09101100', 5, 'bag', 'available'),

('Garlic', 'Fresh', 'LOT-GAR-001', 52.00, 'kg', '2025-11-11', '2025-12-11',
 'vegetables', 'bulb', 120.00, 175.00, 'Local Farmers Market',
 'LFM-2025-7900', 30.0, 45.0, 30, 'Dry Storage C1', '8901234567824', 'VEG-GAR-FRE-001',
 '07032000', 5, 'bag', 'available'),

('Green Chili', 'Hot', 'LOT-CHI-001', 35.00, 'kg', '2025-11-14', '2025-11-24',
 'vegetables', 'pepper', 60.00, 95.00, 'Andhra Traders',
 'AT-2025-2222', 20.0, 30.0, 10, 'Cold Storage A2', '8901234567825', 'VEG-CHI-HOT-001',
 '07096000', 5, 'bag', 'available');

-- Insert fruits inventory (20 items)
INSERT INTO inventory_items (
    name, variant, lot_number, quantity, unit, purchase_date, expiry_date,
    category, sub_category, cost_price, selling_price, supplier_name,
    purchase_invoice, min_stock_level, reorder_point, shelf_life_days,
    storage_location, barcode, sku, hsn_code, gst_rate,
    packaging_type, status
) VALUES
('Apple', 'Shimla Deluxe', 'LOT-APP-001', 125.00, 'kg', '2025-11-10', '2025-12-10',
 'fruits', 'pome', 120.00, 185.00, 'Kashmir Fruit Traders',
 'KFT-2025-1111', 40.0, 60.0, 30, 'Cold Storage B1', '8901234567826', 'FRU-APP-SHI-001',
 '08081000', 0, 'box', 'available'),

('Apple', 'Kashmiri', 'LOT-APP-002', 95.00, 'kg', '2025-11-11', '2025-12-11',
 'fruits', 'pome', 150.00, 220.00, 'Kashmir Fruit Traders',
 'KFT-2025-1112', 30.0, 50.0, 30, 'Cold Storage B1', '8901234567827', 'FRU-APP-KAS-001',
 '08081000', 0, 'box', 'available'),

('Banana', 'Robusta', 'LOT-BAN-001', 200.00, 'kg', '2025-11-14', '2025-11-20',
 'fruits', 'tropical', 30.00, 50.00, 'Tamil Nadu Banana Co',
 'TNBC-2025-2222', 100.0, 150.0, 6, 'Ambient Storage C', '8901234567828', 'FRU-BAN-ROB-001',
 '08030000', 0, 'bunch', 'available'),

('Banana', 'Red', 'LOT-BAN-002', 85.00, 'kg', '2025-11-15', '2025-11-21',
 'fruits', 'tropical', 45.00, 72.00, 'Tamil Nadu Banana Co',
 'TNBC-2025-2223', 50.0, 80.0, 6, 'Ambient Storage C', '8901234567829', 'FRU-BAN-RED-001',
 '08030000', 0, 'bunch', 'available'),

('Mango', 'Alphonso', 'LOT-MAN-001', 72.00, 'kg', '2025-11-15', '2025-11-23',
 'fruits', 'tropical', 200.00, 325.00, 'Ratnagiri Mango Suppliers',
 'RMS-2025-3333', 30.0, 50.0, 8, 'Cool Storage D1', '8901234567830', 'FRU-MAN-ALP-001',
 '08045000', 5, 'crate', 'available'),

('Orange', 'Nagpur', 'LOT-ORA-001', 110.00, 'kg', '2025-11-13', '2025-11-28',
 'fruits', 'citrus', 45.00, 72.00, 'Citrus Fresh Ltd',
 'CFL-2025-4444', 50.0, 80.0, 15, 'Ambient Storage C', '8901234567831', 'FRU-ORA-NAG-001',
 '08051000', 5, 'net_bag', 'available'),

('Orange', 'Sweet', 'LOT-ORA-002', 88.00, 'kg', '2025-11-14', '2025-11-29',
 'fruits', 'citrus', 50.00, 78.00, 'Citrus Fresh Ltd',
 'CFL-2025-4445', 40.0, 60.0, 15, 'Ambient Storage C', '8901234567832', 'FRU-ORA-SWE-001',
 '08051000', 5, 'net_bag', 'available'),

('Grapes', 'Black Seedless', 'LOT-GRA-001', 62.00, 'kg', '2025-11-14', '2025-11-22',
 'fruits', 'berry', 90.00, 145.00, 'Maharashtra Grape Growers',
 'MGG-2025-5555', 25.0, 40.0, 8, 'Cold Storage B1', '8901234567833', 'FRU-GRA-BLA-001',
 '08061000', 5, 'box', 'available'),

('Grapes', 'Green Seedless', 'LOT-GRA-002', 55.00, 'kg', '2025-11-15', '2025-11-23',
 'fruits', 'berry', 85.00, 135.00, 'Maharashtra Grape Growers',
 'MGG-2025-5556', 25.0, 40.0, 8, 'Cold Storage B1', '8901234567834', 'FRU-GRA-GRE-001',
 '08061000', 5, 'box', 'available'),

('Watermelon', 'Large', 'LOT-WAT-001', 180.00, 'kg', '2025-11-15', '2025-11-26',
 'fruits', 'melon', 15.00, 28.00, 'Karnataka Farmers Coop',
 'KFC-2025-5680', 80.0, 120.0, 11, 'Ambient Storage C', '8901234567835', 'FRU-WAT-LAR-001',
 '08071000', 0, 'loose', 'available'),

('Papaya', 'Ripe', 'LOT-PAP-001', 95.00, 'kg', '2025-11-14', '2025-11-22',
 'fruits', 'tropical', 28.00, 48.00, 'Local Farmers Market',
 'LFM-2025-8001', 40.0, 60.0, 8, 'Ambient Storage C', '8901234567836', 'FRU-PAP-RIP-001',
 '08072000', 0, 'loose', 'available'),

('Pomegranate', 'Fresh', 'LOT-POM-001', 68.00, 'kg', '2025-11-13', '2025-11-28',
 'fruits', 'berry', 110.00, 165.00, 'Maharashtra Grape Growers',
 'MGG-2025-5557', 30.0, 50.0, 15, 'Ambient Storage C', '8901234567837', 'FRU-POM-FRE-001',
 '08109000', 5, 'box', 'available'),

('Pineapple', 'Fresh', 'LOT-PIN-001', 85.00, 'kg', '2025-11-14', '2025-11-24',
 'fruits', 'tropical', 35.00, 58.00, 'Kerala Fruits',
 'KF-2025-6666', 40.0, 60.0, 10, 'Ambient Storage C', '8901234567838', 'FRU-PIN-FRE-001',
 '08043000', 5, 'crate', 'available'),

('Guava', 'Fresh', 'LOT-GUA-001', 72.00, 'kg', '2025-11-15', '2025-11-22',
 'fruits', 'tropical', 32.00, 52.00, 'Local Farmers Market',
 'LFM-2025-8002', 30.0, 50.0, 7, 'Ambient Storage C', '8901234567839', 'FRU-GUA-FRE-001',
 '08045000', 5, 'loose', 'available'),

('Sweet Lime', 'Fresh', 'LOT-LIM-001', 58.00, 'kg', '2025-11-14', '2025-11-29',
 'fruits', 'citrus', 40.00, 65.00, 'Citrus Fresh Ltd',
 'CFL-2025-4446', 25.0, 40.0, 15, 'Ambient Storage C', '8901234567840', 'FRU-LIM-FRE-001',
 '08054000', 5, 'net_bag', 'available'),

('Sapota', 'Ripe', 'LOT-SAP-001', 45.00, 'kg', '2025-11-15', '2025-11-20',
 'fruits', 'tropical', 38.00, 60.00, 'Karnataka Farmers Coop',
 'KFC-2025-5681', 20.0, 35.0, 5, 'Ambient Storage C', '8901234567841', 'FRU-SAP-RIP-001',
 '08109000', 5, 'box', 'available'),

('Dragon Fruit', 'White', 'LOT-DRA-001', 28.00, 'kg', '2025-11-15', '2025-11-23',
 'fruits', 'exotic', 180.00, 280.00, 'Exotic Fruits Import',
 'EFI-2025-7777', 15.0, 25.0, 8, 'Cold Storage B1', '8901234567842', 'FRU-DRA-WHI-001',
 '08109000', 5, 'box', 'available'),

('Kiwi', 'Green', 'LOT-KIW-001', 22.00, 'kg', '2025-11-12', '2025-12-02',
 'fruits', 'exotic', 220.00, 340.00, 'Exotic Fruits Import',
 'EFI-2025-7778', 10.0, 20.0, 20, 'Cold Storage B1', '8901234567843', 'FRU-KIW-GRE-001',
 '08105000', 5, 'box', 'available'),

('Strawberry', 'Fresh', 'LOT-STR-001', 18.00, 'kg', '2025-11-16', '2025-11-20',
 'fruits', 'berry', 180.00, 280.00, 'Hill Station Produce',
 'HSP-2025-4569', 10.0, 20.0, 4, 'Cold Storage B1', '8901234567844', 'FRU-STR-FRE-001',
 '08101000', 5, 'box', 'available'),

('Coconut', 'Fresh', 'LOT-COC-001', 150.00, 'kg', '2025-11-13', '2025-12-13',
 'fruits', 'tropical', 25.00, 40.00, 'Kerala Fruits',
 'KF-2025-6667', 80.0, 120.0, 30, 'Warehouse B2', '8901234567845', 'FRU-COC-FRE-001',
 '08011000', 5, 'bag', 'available');

ANALYZE;

-- Show summary
SELECT 'Customers created:' as info, COUNT(*)::text as count FROM customers
UNION ALL
SELECT 'Inventory items created:', COUNT(*)::text FROM inventory_items
UNION ALL
SELECT 'Vegetables:', COUNT(*)::text FROM inventory_items WHERE category = 'vegetables'
UNION ALL
SELECT 'Fruits:', COUNT(*)::text FROM inventory_items WHERE category = 'fruits';
