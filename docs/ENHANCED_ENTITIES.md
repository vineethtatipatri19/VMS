# Enhanced Entity Documentation

**Date**: November 16, 2025  
**Migration**: 004_enhance_entities.sql

---

## Overview

All core entities have been significantly enhanced with comprehensive business fields to support advanced features like credit management, profit tracking, supplier management, and wastage monitoring.

---

## 📋 Enhanced Entities

### 1. **CUSTOMERS** - Comprehensive Credit & KYC Management

#### Basic Information
- `id` - UUID primary key
- `name` - Customer name (required)
- `email` - Email address
- `business_name` - Business/company name for B2B customers

#### Contact Details
- `contact_number` - Primary phone
- `alternate_contact` - Secondary phone
- `whatsapp_number` - WhatsApp for automated messaging
- `address` - Delivery/billing address

#### KYC & Identity
- `photo_url` - Customer photo URL
- `aadhaar_verified` - Boolean flag
- `kyc_document_type` - Document type (Aadhaar, PAN, Voter ID, etc.)
- `kyc_document_number` - Document number
- `gstin` - GST Identification Number for B2B

#### Customer Classification
- `customer_type` - `b2b`, `b2c`, `retail`, `wholesale`
- `status` - `active`, `inactive`, `blocked`
- `tags` - Array of custom tags

#### Credit Management (Tier 1 Priority)
- `credit_limit` - Maximum credit allowed (₹)
- `current_balance` - Outstanding amount (₹)
- `payment_terms_days` - Net payment terms (default: 30)
- `interest_rate` - Interest on overdue (%)

#### Business Intelligence
- `last_transaction_date` - Last purchase date
- `total_purchases` - Lifetime purchase value (₹)
- `loyalty_points` - Reward points
- `notes` - Internal notes/remarks

#### Timestamps
- `created_at` - Account creation
- `updated_at` - Last modified

**Indexes**: contact_number, status, current_balance, email

---

### 2. **INVENTORY_ITEMS** - Advanced Stock & Profit Tracking

#### Basic Information
- `id` - UUID primary key
- `name` - Item name (required)
- `variant` - Size/grade variant
- `category` - Main category (vegetables, fruits, dairy, etc.)
- `sub_category` - Sub-category

#### Lot Management (FEFO)
- `lot_number` - Unique batch/lot identifier
- `purchase_date` - Date of purchase
- `expiry_date` - Expiration date (FEFO sorting)
- `shelf_life_days` - Expected shelf life

#### Quantity & Units
- `quantity` - Current stock quantity
- `unit` - `kg` or `lot`
- `weight_per_unit` - Weight per unit (for lots)
- `packaging_type` - Box, crate, bag, etc.

#### Pricing & Profitability (Tier 1 Priority)
- `cost_price` - Purchase price per unit (₹)
- `selling_price` - Retail selling price (₹)
- `margin_percentage` - **Auto-calculated** profit margin %
- `hsn_code` - HSN code for GST
- `gst_rate` - GST rate (%)

#### Supplier Tracking (Tier 1 Priority)
- `supplier_id` - Reference to supplier
- `supplier_name` - Supplier name
- `purchase_invoice` - Invoice/bill number

#### Stock Management
- `min_stock_level` - Minimum quantity threshold
- `reorder_point` - When to reorder
- `status` - **Auto-calculated**: `available`, `low_stock`, `out_of_stock`, `expired`, `damaged`

#### Warehouse Management
- `storage_location` - Warehouse location/section
- `barcode` - Barcode for scanning
- `sku` - Stock Keeping Unit code
- `image_url` - Product image

#### Performance Metrics
- `total_sold` - Total quantity sold (lifetime)
- `total_wasted` - Total quantity wasted
- `is_perishable` - Boolean flag

#### Timestamps
- `created_at`, `updated_at`

**Auto-Triggers**:
- Status automatically updates based on quantity and expiry
- Margin percentage calculated from cost/selling price

**Indexes**: expiry_date, status, category, barcode

---

### 3. **TRANSACTIONS** - Complete Sales & Payment Tracking

#### Basic Information
- `id` - UUID primary key
- `customer_id` - Foreign key to customers
- `date` - Transaction date/time
- `type` - `sale` or `payment`

#### Invoice & Reference
- `invoice_number` - **Auto-generated** (e.g., INV-20251116-1001)
- `payment_reference` - UPI ref, cheque no, etc.
- `notes` - Transaction notes/remarks

#### Payment Details (Tier 1 Priority)
- `payment_method` - `cash`, `upi`, `card`, `bank_transfer`, `cheque`, `credit`
- `payment_amount` - For payment transactions (₹)
- `due_date` - Payment due date (for credit sales)
- `is_overdue` - **Auto-calculated** boolean

#### Financial Breakdown
- `total_amount` - Total transaction amount (₹)
- `discount_amount` - Discount given (₹)
- `tax_amount` - GST/tax amount (₹)
- `balance_after` - Customer balance after transaction

#### Sale Type & Classification
- `sale_type` - `regular`, `wholesale`, `credit`, `return`
- `receipt_sent` - Boolean (SMS/WhatsApp sent?)

#### Delivery Management
- `delivery_status` - `pending`, `packed`, `dispatched`, `delivered`, `cancelled`
- `delivery_date` - Actual delivery date
- `delivery_address` - Delivery location

#### Details
- `details` - JSONB field with line items

#### Timestamps
- `created_at`

**Auto-Triggers**:
- Invoice number auto-generated on insert
- Overdue flag auto-updated based on due_date

**Indexes**: date, customer_id, type, is_overdue, invoice_number

---

### 4. **SALE_ITEMS** - Line Item Details with Profit

#### Item Details
- `id` - UUID primary key
- `transaction_id` - Foreign key to transactions
- `inventory_lot_id` - Foreign key to inventory
- `item_name` - Item name snapshot
- `batch_number` - Batch/lot number snapshot
- `expiry_date` - Expiry date snapshot

#### Quantity & Pricing
- `quantity` - Quantity sold
- `unit` - Unit of measurement
- `price_per_unit` - Selling price per unit (₹)
- `cost_per_unit` - Cost price per unit (₹)
- `profit` - **Auto-calculated**: (price - cost) × quantity

#### Discounts & Tax
- `discount_percentage` - Item-level discount %
- `tax_percentage` - GST rate %
- `hsn_code` - HSN code

#### Financial
- `total` - Line total (₹)

---

### 5. **CRATE_LEDGER** - Enhanced Asset Tracking

#### Basic Information
- `id` - UUID primary key
- `customer_id` - Foreign key to customers
- `date` - Transaction date
- `transaction_id` - Link to sale transaction

#### Crate Movement
- `crates_issued` - Number issued
- `crates_returned` - Number returned
- `balance` - Running balance

#### Crate Details
- `crate_type` - Type/size of crate
- `crate_value` - Replacement value per crate (₹)
- `notes` - Remarks

#### Timestamps
- `updated_at`

**Indexes**: customer_id, balance (where > 0)

---

## 🆕 New Tables

### 6. **WASTAGE_LOG** - Track Spoilage & Loss

**Purpose**: Record all wasted inventory for analysis and reduction

#### Fields
- `id` - UUID
- `inventory_item_id` - Reference to inventory
- `item_name` - Item name
- `quantity` - Quantity wasted
- `unit` - Unit of measurement
- `reason` - `expired`, `damaged`, `spoiled`, `pest`, `other`
- `reason_details` - Additional details
- `cost_value` - Financial impact (₹)
- `logged_by` - User who logged
- `logged_at` - Timestamp
- `photo_url` - Evidence photo

**Indexes**: logged_at, reason

---

### 7. **EXPIRY_ALERTS** - Proactive Expiry Monitoring

**Purpose**: Alert system for items approaching expiry (3-day window)

#### Fields
- `id` - UUID
- `inventory_item_id` - Reference to inventory
- `alert_date` - Date alert was generated
- `expiry_date` - Item expiry date
- `days_until_expiry` - Countdown
- `acknowledged` - Boolean
- `acknowledged_at` - When acknowledged
- `acknowledged_by` - User who acknowledged
- `created_at` - Timestamp

**Indexes**: alert_date, acknowledged

---

### 8. **PAYMENT_SCHEDULES** - Installment Tracking

**Purpose**: Support payment plans for large credit sales

#### Fields
- `id` - UUID
- `transaction_id` - Original sale transaction
- `customer_id` - Customer reference
- `installment_number` - 1, 2, 3, etc.
- `due_date` - Payment due date
- `amount_due` - Expected amount (₹)
- `amount_paid` - Amount received (₹)
- `status` - `pending`, `partial`, `paid`, `overdue`
- `payment_date` - Actual payment date
- `notes` - Remarks
- `created_at`, `updated_at`

**Indexes**: customer_id, status, due_date

---

### 9. **PRICING_TIERS** - Bulk Pricing Support

**Purpose**: Different prices for different quantities (wholesale vs retail)

#### Fields
- `id` - UUID
- `inventory_item_id` - Item reference
- `min_quantity` - Minimum quantity for tier
- `max_quantity` - Maximum quantity (null = unlimited)
- `price_per_unit` - Price at this tier (₹)
- `tier_name` - Description (e.g., "Retail", "Wholesale")
- `created_at`

---

### 10. **PRICE_HISTORY** - Audit Trail

**Purpose**: Track all price changes for compliance and analysis

#### Fields
- `id` - UUID
- `inventory_item_id` - Item reference
- `old_price` - Previous price (₹)
- `new_price` - New price (₹)
- `changed_by` - User who changed
- `reason` - Why price changed
- `changed_at` - Timestamp

**Index**: inventory_item_id + changed_at

---

## 📊 New Views for Reporting

### 1. **customer_summary**
Aggregated customer data with balances, transaction counts, and crate status.

**Columns**:
- Customer basic info
- `total_transactions` - Count of sales
- `total_sales` - Lifetime sales amount (₹)
- `total_payments` - Payments received (₹)
- `last_transaction_date`
- `crate_balance` - Outstanding crates

---

### 2. **inventory_summary**
Enhanced inventory view with profitability and expiry status.

**Columns**:
- All inventory fields
- `expiry_status` - `expired`, `expiring_soon`, `fresh`
- `days_until_expiry` - Countdown
- `total_value` - cost_price × quantity

---

### 3. **sales_performance**
Daily sales metrics with profit analysis.

**Columns**:
- `sale_date`
- `transaction_count`
- `unique_customers`
- `total_sales`, `total_tax`, `total_discounts`
- `avg_transaction_value`
- `total_profit`
- `profit_margin_percentage`

---

### 4. **overdue_payments**
List of all overdue credit sales with customer contact info.

**Columns**:
- Transaction details
- Customer name, contact, WhatsApp
- `days_overdue`
- `balance_due`

---

### 5. **wastage_summary**
Daily wastage analysis by reason.

**Columns**:
- `wastage_date`, `reason`
- `item_count`, `total_quantity`, `total_cost`
- `items_wasted` - Aggregated list

---

## 🎯 Benefits of Enhanced Entities

### For Credit Management (Tier 1)
✅ Enforce credit limits automatically  
✅ Track payment terms per customer  
✅ Monitor overdue payments with alerts  
✅ Support installment plans  
✅ Calculate interest on overdue amounts

### For Profit Tracking (Tier 1)
✅ Know cost price and selling price for every item  
✅ Auto-calculate profit margins  
✅ Track profitability per transaction and per item  
✅ Identify high-margin vs low-margin products

### For Wastage Reduction (Tier 1)
✅ Log all wastage with reasons  
✅ Get alerts before items expire  
✅ Analyze wastage patterns  
✅ Calculate financial impact of waste

### For Supplier Management (Tier 1)
✅ Track which supplier provides which items  
✅ Record purchase invoices  
✅ Link inventory to suppliers

### For GST Compliance (Tier 2)
✅ Store GST IN for B2B customers  
✅ HSN codes on inventory items  
✅ Tax calculation fields ready

### For Customer Segmentation (Tier 3)
✅ Customer types (B2B, B2C, wholesale, retail)  
✅ Loyalty points system ready  
✅ Tags for custom segmentation  
✅ Track lifetime purchase value

---

## 🗑️ Delete Operations & Audit Trail

### Soft Delete System

**All entities support soft deletes** with audit trail preservation. No data is permanently removed.

#### Audit Fields (All Tables)
- `deleted_at` - Timestamp when deleted (NULL = not deleted)
- `deleted_by` - User who performed deletion
- `deletion_reason` - Required reason for deletion

#### Attestation Requirement

**All DELETE operations require user attestation** to prevent accidental data loss:

```json
DELETE /api/v1/inventory/{id}
{
  "reason": "Duplicate entry created by mistake",
  "attestation": "I CONFIRM DELETE"
}
```

The attestation string **must exactly match** `"I CONFIRM DELETE"` or the operation will be rejected.

#### Soft Delete Behavior

**Database Layer**:
- Records are marked as deleted with `deleted_at = CURRENT_TIMESTAMP`
- `deleted_by` is set to the authenticated user ID
- `deletion_reason` stores the provided reason
- Original data remains intact in the database

**API Layer**:
- Default queries exclude soft-deleted records (`WHERE deleted_at IS NULL`)
- Special endpoints can retrieve deleted records for audit purposes
- Restore functionality can be implemented by clearing `deleted_at`

**Example**: Soft delete a wastage entry
```sql
UPDATE wastage_log
SET deleted_at = CURRENT_TIMESTAMP,
    deleted_by = 'user-uuid',
    deletion_reason = 'Logged in error'
WHERE id = 'wastage-uuid';
```

#### Entities with Soft Delete Support

✅ **customers** - Customer accounts  
✅ **inventory_items** - Inventory records  
✅ **transactions** - Sales and payment records  
✅ **sale_items** - Transaction line items  
✅ **crate_ledger** - Crate tracking entries  
✅ **wastage_log** - Wastage entries  
✅ **expiry_alerts** - Alert records  
✅ **payment_schedules** - Installment plans  

#### Benefits of Soft Deletes

🔒 **Compliance**: Full audit trail for regulatory requirements  
🔍 **Forensics**: Investigate issues with complete historical data  
↩️ **Recovery**: Restore accidentally deleted records  
📊 **Reporting**: Historical analysis includes deleted records  
🛡️ **Security**: Track who deleted what and why  

---

## 🔧 Next Steps

1. ✅ **Backend Handlers Updated** - Go API handlers support all enhanced fields
2. ✅ **Frontend Forms Updated** - React forms include all new fields
3. ✅ **Features Implemented**:
   - ✅ Credit limit enforcement
   - ✅ Wastage log entry with photo upload
   - ✅ Expiry alert dashboard with acknowledgment
   - ✅ Profit analysis in reports
   - ✅ Overdue payment tracking
   - ✅ Soft delete with attestation
   - ✅ Audit trail system
   - ⏳ Payment plan creation (planned)

4. ✅ **Seed Test Data** - Demo data includes enhanced fields

---

## 📝 Migration Status

✅ **Migration Applied**: November 16, 2025  
✅ **Database**: All new columns, tables, indexes, views created  
✅ **Triggers**: Auto-calculations active  
✅ **Backend**: All handlers updated with CRUD operations  
✅ **Frontend**: Full UI implementation complete  
✅ **Audit System**: Soft deletes and attestation active  
✅ **Delete Handlers**: Attestation-required delete endpoints deployed

---

## 💡 Usage Examples

### Query Customer with Enhanced Fields
```sql
SELECT id, name, customer_type, credit_limit, current_balance, 
       payment_terms_days, status
FROM customers 
WHERE status = 'active' 
  AND current_balance > 0
ORDER BY current_balance DESC;
```

### Get Items Near Expiry
```sql
SELECT * FROM inventory_summary
WHERE expiry_status = 'expiring_soon'
ORDER BY days_until_expiry;
```

### Calculate Daily Profit
```sql
SELECT sale_date, total_sales, total_profit, profit_margin_percentage
FROM sales_performance
WHERE sale_date >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY sale_date DESC;
```

### Find Overdue Payments
```sql
SELECT * FROM overdue_payments
WHERE days_overdue > 7
ORDER BY days_overdue DESC;
```

---

**End of Enhanced Entities Documentation**
