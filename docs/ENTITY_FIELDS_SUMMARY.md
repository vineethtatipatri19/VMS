# Entity Fields Quick Reference

**Last Updated**: November 16, 2025

---

## 👥 CUSTOMERS (26 fields)

### Display Priority: HIGH
- ✅ **name** - Customer name
- ✅ **contact_number** - Primary phone
- ✅ **customer_type** - b2b/b2c/retail/wholesale
- ✅ **credit_limit** - Max credit (₹)
- ✅ **current_balance** - Outstanding (₹)
- ✅ **status** - active/inactive/blocked

### Display Priority: MEDIUM
- email, business_name, gstin
- whatsapp_number, alternate_contact
- payment_terms_days, interest_rate
- total_purchases, loyalty_points
- last_transaction_date

### Display Priority: LOW (Admin Only)
- address, photo_url
- aadhaar_verified, kyc_document_type, kyc_document_number
- notes, tags
- created_at, updated_at

---

## 📦 INVENTORY ITEMS (35 fields)

### Display Priority: HIGH
- ✅ **name** - Item name
- ✅ **category** - Main category
- ✅ **quantity** - Stock quantity
- ✅ **unit** - kg/lot
- ✅ **selling_price** - Retail price (₹)
- ✅ **expiry_date** - Expiration date
- ✅ **status** - available/low_stock/expired

### Display Priority: MEDIUM
- cost_price, margin_percentage
- supplier_name, purchase_invoice
- min_stock_level, reorder_point
- lot_number, purchase_date
- shelf_life_days, storage_location

### Display Priority: LOW
- variant, sub_category
- barcode, sku, hsn_code, gst_rate
- weight_per_unit, packaging_type
- image_url, supplier_id
- is_perishable
- total_sold, total_wasted
- created_at, updated_at

---

## 💰 TRANSACTIONS (21 fields)

### Display Priority: HIGH
- ✅ **invoice_number** - Auto-generated
- ✅ **customer_id** - Customer reference
- ✅ **date** - Transaction date
- ✅ **type** - sale/payment
- ✅ **total_amount** - Total (₹)
- ✅ **payment_method** - cash/upi/card/credit
- ✅ **is_overdue** - Overdue flag

### Display Priority: MEDIUM
- payment_amount, due_date
- discount_amount, tax_amount
- balance_after
- sale_type, delivery_status

### Display Priority: LOW
- payment_reference, notes
- receipt_sent
- delivery_date, delivery_address
- details (JSONB)
- created_at

---

## 📋 SALE ITEMS (13 fields)

### Display Priority: HIGH
- ✅ **item_name** - Product name
- ✅ **quantity** - Qty sold
- ✅ **unit** - kg/lot
- ✅ **price_per_unit** - Selling price (₹)
- ✅ **total** - Line total (₹)

### Display Priority: MEDIUM
- cost_per_unit, profit (auto-calc)
- discount_percentage, tax_percentage

### Display Priority: LOW
- transaction_id, inventory_lot_id
- batch_number, expiry_date, hsn_code

---

## 📦 CRATES (9 fields)

### Display Priority: HIGH
- ✅ **customer_id** - Customer
- ✅ **crates_issued** - Given out
- ✅ **crates_returned** - Returned
- ✅ **balance** - Outstanding
- ✅ **date** - Transaction date

### Display Priority: MEDIUM
- crate_type, crate_value

### Display Priority: LOW
- transaction_id, notes, updated_at

---

## ��️ WASTAGE LOG (10 fields)

### Display Priority: HIGH
- ✅ **item_name** - Wasted item
- ✅ **quantity** - Amount wasted
- ✅ **reason** - expired/damaged/spoiled/pest
- ✅ **cost_value** - Financial loss (₹)
- ✅ **logged_at** - Date

### Display Priority: MEDIUM
- reason_details, logged_by

### Display Priority: LOW
- inventory_item_id, unit, photo_url

---

## ⏰ EXPIRY ALERTS (9 fields)

### Display Priority: HIGH
- ✅ **inventory_item_id** - Item reference
- ✅ **expiry_date** - When expires
- ✅ **days_until_expiry** - Countdown
- ✅ **alert_date** - Alert generated
- ✅ **acknowledged** - Seen?

### Display Priority: LOW
- acknowledged_at, acknowledged_by, created_at

---

## 💳 PAYMENT SCHEDULES (12 fields)

### Display Priority: HIGH
- ✅ **customer_id** - Customer
- ✅ **installment_number** - #1, #2, etc.
- ✅ **due_date** - When due
- ✅ **amount_due** - Expected (₹)
- ✅ **amount_paid** - Received (₹)
- ✅ **status** - pending/partial/paid/overdue

### Display Priority: MEDIUM
- transaction_id, payment_date, notes

### Display Priority: LOW
- created_at, updated_at

---

## 💵 PRICING TIERS (6 fields)

### Display Priority: HIGH
- ✅ **inventory_item_id** - Item
- ✅ **min_quantity** - Min qty
- ✅ **max_quantity** - Max qty
- ✅ **price_per_unit** - Price (₹)
- ✅ **tier_name** - Description

---

## 📈 PRICE HISTORY (6 fields)

### Display Priority: MEDIUM
- inventory_item_id, old_price, new_price
- changed_by, reason, changed_at

---

## 🎯 Recommended Display Fields by Page

### Customer List View
Show: name, contact_number, customer_type, current_balance, status, last_transaction_date

### Customer Detail View
Show: All HIGH + MEDIUM priority fields

### Inventory List View
Show: name, category, quantity, unit, selling_price, expiry_date, status

### Inventory Detail View
Show: All fields except internal IDs

### Transaction List View
Show: invoice_number, date, customer_id (with name), type, total_amount, payment_method, is_overdue

### Transaction Detail View
Show: All transaction fields + sale_items breakdown

### Dashboard KPIs
- Total customers (active)
- Total inventory value (cost_price × quantity)
- Today's sales amount
- Outstanding receivables (sum of current_balance)
- Overdue payments count
- Items expiring soon (next 3 days)
- Low stock items count
- Today's profit

---

**Note**: LOW priority fields can be hidden by default and shown in "More Details" expandable sections or admin-only views.
