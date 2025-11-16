
# User Stories

## Authentication & Access Control

**US-1.1**: As a **user**, I want to **register with email and password**, so that **I can create my account securely**.  
**Acceptance Criteria**:
- Password must be hashed before storage
- Email validation required
- Success message displayed after registration
- Automatic redirect to login page

**US-1.2**: As a **user**, I want to **log in with my credentials**, so that **I can access the system**.  
**Acceptance Criteria**:
- JWT token generated on successful login
- Token stored in localStorage
- Redirect to dashboard after login
- Error message for invalid credentials

---

## Inventory Management

**US-2.1**: As a **vendor**, I want to **add new inventory items with complete details**, so that **I can track stock, costs, and expiry dates**.  
**Acceptance Criteria**:
- Form includes: name, category, quantity, unit, cost price, selling price, expiry date, supplier, lot number
- Margin percentage auto-calculated
- Status auto-assigned based on quantity
- Success toast notification
- Item appears in inventory list immediately

**US-2.2**: As a **vendor**, I want to **edit inventory item details**, so that **I can correct mistakes or update information**.  
**Acceptance Criteria**:
- All fields editable except ID
- Changes saved to database
- Audit trail recorded (who updated, when)
- Success confirmation message

**US-2.3**: As a **vendor**, I want to **delete inventory items with attestation**, so that **I can remove items while maintaining audit trail**.  
**Acceptance Criteria**:
- Delete button opens confirmation modal
- Must provide deletion reason
- Must type "I CONFIRM DELETE" exactly
- Record soft-deleted (not permanently removed)
- Audit fields populated: deleted_at, deleted_by, deletion_reason
- Success message: "Inventory item deleted successfully"

**US-2.4**: As a **vendor**, I want to **see items sorted by expiry date (FEFO)**, so that **I can sell older stock first**.  
**Acceptance Criteria**:
- Items automatically sorted by expiry_date ascending
- Near-expiry items highlighted in yellow
- Expired items highlighted in red
- Days until expiry displayed

**US-2.5**: As a **vendor**, I want to **view stock status at a glance**, so that **I know which items need reordering**.  
**Acceptance Criteria**:
- Color-coded status badges
- Green = Available, Yellow = Low Stock, Red = Out of Stock, Black = Expired
- Low stock threshold: quantity <= min_stock_level

---

## Customer Management

**US-3.1**: As a **vendor**, I want to **add new customers with complete profiles**, so that **I can manage credit, contact details, and KYC information**.  
**Acceptance Criteria**:
- Form includes: name, contact, business name, customer type, credit limit, payment terms
- Optional fields: KYC documents, GSTIN, photo, address
- Customer added to database
- Appears in customer list immediately

**US-3.2**: As a **vendor**, I want to **edit customer information**, so that **I can keep records up to date**.  
**Acceptance Criteria**:
- All fields editable
- Changes saved with audit trail
- Current balance recalculated if needed
- Success notification

**US-3.3**: As a **vendor**, I want to **delete customers with attestation**, so that **I can remove inactive customers while maintaining compliance**.  
**Acceptance Criteria**:
- Delete confirmation modal appears
- Must provide deletion reason
- Must type "I CONFIRM DELETE"
- Customer soft-deleted (remains in database)
- Audit fields populated
- Removed from active customer list

**US-3.4**: As a **vendor**, I want to **view customer transaction history**, so that **I can track purchases and payments**.  
**Acceptance Criteria**:
- Customer detail page shows all transactions
- Running balance displayed
- Payment history visible
- Outstanding crates shown

---

## Transaction Management

**US-4.1**: As a **vendor**, I want to **create sales transactions with multiple items**, so that **I can record customer purchases**.  
**Acceptance Criteria**:
- Select customer from dropdown
- Add multiple line items (item, quantity, price)
- Discount and tax fields available
- Total calculated automatically
- Invoice number auto-generated
- Inventory quantities decremented
- Customer balance updated
- Success message with invoice number

**US-4.2**: As a **vendor**, I want to **record payments**, so that **I can track customer balances**.  
**Acceptance Criteria**:
- Select customer
- Enter payment amount
- Select payment method (cash, UPI, card, etc.)
- Optional payment reference (UPI ID, cheque number)
- Customer balance decremented
- Payment recorded with timestamp
- Receipt generated (optional)

**US-4.3**: As a **vendor**, I want to **edit transactions**, so that **I can correct errors after submission**.  
**Acceptance Criteria**:
- Edit button on transaction details
- All fields editable (except invoice number)
- Audit trail recorded
- Inventory and balances recalculated
- Updated_at timestamp refreshed

**US-4.4**: As a **vendor**, I want to **delete transactions with attestation**, so that **I can remove incorrect entries while maintaining audit compliance**.  
**Acceptance Criteria**:
- Delete confirmation modal
- Must provide reason
- Must type "I CONFIRM DELETE"
- Transaction soft-deleted
- Inventory and customer balances rolled back
- Audit trail complete

---

## Crate Management

**US-5.1**: As a **vendor**, I want to **issue crates to customers**, so that **I can track returnable assets**.  
**Acceptance Criteria**:
- Issue crates form: customer, quantity, crate type, date
- Balance incremented for customer
- Transaction ID linked (if part of sale)
- Crate value tracked

**US-5.2**: As a **vendor**, I want to **record crate returns**, so that **I can decrease customer crate balances**.  
**Acceptance Criteria**:
- Return form: customer, quantity
- Balance decremented
- Cannot return more than issued
- Return date recorded

**US-5.3**: As a **vendor**, I want to **delete crate entries with attestation**, so that **I can correct mistakes**.  
**Acceptance Criteria**:
- Delete modal with reason and attestation
- Soft delete applied
- Balance recalculated
- Audit trail maintained

---

## Wastage Tracking

**US-6.1**: As a **vendor**, I want to **log wasted inventory**, so that **I can track losses and identify patterns**.  
**Acceptance Criteria**:
- Wastage form: item, quantity, reason (expired, damaged, spoiled, pest, other)
- Optional: photo upload, detailed notes
- Cost impact calculated automatically
- Wastage logged with timestamp and user ID
- Item's total_wasted incremented

**US-6.2**: As a **vendor**, I want to **view wastage reports**, so that **I can analyze loss patterns and reduce waste**.  
**Acceptance Criteria**:
- Wastage list page shows all entries
- Filter by date range, reason, item
- Total wastage cost displayed
- Export to Excel (future)

**US-6.3**: As a **vendor**, I want to **delete wastage entries with attestation**, so that **I can remove incorrect logs**.  
**Acceptance Criteria**:
- Delete confirmation with reason and attestation
- Soft delete applied
- Audit trail maintained

---

## Expiry Alert System

**US-7.1**: As a **vendor**, I want to **receive alerts for items expiring soon**, so that **I can take action before they expire**.  
**Acceptance Criteria**:
- Alerts generated for items expiring within 3 days
- Alert dashboard shows all pending alerts
- Days until expiry displayed
- Item details (name, lot number, quantity) visible

**US-7.2**: As a **vendor**, I want to **acknowledge expiry alerts**, so that **I can mark them as seen and handled**.  
**Acceptance Criteria**:
- Acknowledge button on each alert
- Alert marked as acknowledged with timestamp and user ID
- Acknowledged alerts removed from pending list
- Acknowledged alerts viewable in history (optional)

**US-7.3**: As a **vendor**, I want to **delete expiry alerts with attestation**, so that **I can remove false alerts**.  
**Acceptance Criteria**:
- Delete confirmation modal
- Must provide reason and attestation
- Alert soft-deleted with audit trail

---

## Dashboard & Reporting

**US-8.1**: As a **vendor**, I want to **see key metrics on the dashboard**, so that **I can monitor business performance at a glance**.  
**Acceptance Criteria**:
- Total customers (active only)
- Total inventory value (cost_price × quantity)
- Today's sales amount
- Outstanding receivables (sum of customer balances)
- Low stock items count
- Items expiring within 3 days count
- All metrics update in real-time

**US-8.2**: As a **vendor**, I want to **generate sales reports**, so that **I can analyze revenue trends**.  
**Acceptance Criteria**:
- Date range filter (from, to)
- Display: total sales, transaction count, average transaction value
- Filter by customer (optional)
- Export to Excel (future)

**US-8.3**: As a **vendor**, I want to **generate inventory reports**, so that **I can understand stock levels and value**.  
**Acceptance Criteria**:
- Filter by category, status
- Display: item name, quantity, value, status, expiry date
- Total inventory value calculated
- Export functionality (future)

**US-8.4**: As a **vendor**, I want to **generate customer financial statements**, so that **I can send account summaries**.  
**Acceptance Criteria**:
- Select customer
- Display: all transactions (sales, payments), running balance
- Date range filter
- Print-friendly format
- PDF export (future)

---

## Audit & Compliance

**US-9.1**: As an **admin**, I want to **view audit trail for all deletions**, so that **I can ensure compliance and accountability**.  
**Acceptance Criteria**:
- Audit log page shows all soft-deleted records
- Columns: entity type, entity ID, deleted_at, deleted_by, deletion_reason
- Filter by entity type, date, user
- Export audit log (future)

**US-9.2**: As an **admin**, I want to **require attestation for all deletions**, so that **users cannot accidentally delete data**.  
**Acceptance Criteria**:
- All delete operations show confirmation modal
- Must type "I CONFIRM DELETE" exactly (case-sensitive)
- Invalid attestation shows error message
- Delete blocked until correct attestation provided

**US-9.3**: As an **admin**, I want to **restore soft-deleted records**, so that **I can recover from mistakes**.  
**Acceptance Criteria** (Future):
- Restore button on audit log page
- Clears deleted_at, deleted_by, deletion_reason
- Record appears in active list again
- Restoration logged with timestamp and user

---

## AI Forecasting

**US-10.1**: As a **vendor**, I want to **get demand forecasts for items**, so that **I can optimize inventory levels**.  
**Acceptance Criteria**:
- Forecasting page with item selector
- Input: item name
- Output: predicted demand (quantity), suggested reorder quantity, confidence level
- Forecast based on historical sales data
- Fallback to stub mode if API unavailable

---

## Security & Audit

**US-11.1**: As a **user**, I want **my password to be securely stored**, so that **my account is protected**.  
**Acceptance Criteria**:
- Password hashed with bcrypt
- Plain text password never stored
- Salt applied automatically

**US-11.2**: As a **vendor**, I want **my data to never be permanently deleted**, so that **I can maintain compliance and recover from mistakes**.  
**Acceptance Criteria**:
- All delete operations are soft deletes
- Records remain in database with deleted_at timestamp
- Default queries exclude soft-deleted records
- Audit fields track deletion metadata

**US-11.3**: As an **admin**, I want **complete audit trail for all critical operations**, so that **I can track who did what and when**.  
**Acceptance Criteria**:
- All deletes tracked: entity, user, timestamp, reason
- All updates tracked: updated_at, updated_by (future)
- All creates tracked: created_at
- Audit data exportable for compliance reports

---

# Roadmap

## Phase 1: MVP - ✅ COMPLETED
- ✅ Inventory, Customers, Transactions basic CRUD
- ✅ Print-friendly customer statements
- ✅ Crate ledger tracking
- ✅ FEFO visualization
- ✅ Authentication with JWT

## Phase 2: Audit Compliance - ✅ COMPLETED
- ✅ Enhanced entity fields (credit, supplier, profit tracking)
- ✅ Soft delete system with attestation
- ✅ Wastage tracking with photo evidence
- ✅ Expiry alert system with acknowledgment
- ✅ Full edit/delete UI across all entities
- ✅ Audit trail for all operations

## Phase 3: Documentation & Stabilization - ⏳ IN PROGRESS
- ⏳ Update all documentation
- ⏳ Performance testing
- ⏳ User acceptance testing
- ⏳ Production deployment

## Phase 4: Business Features - 📅 PLANNED
- 📅 Credit management enforcement
- 📅 GST billing and tax compliance
- 📅 WhatsApp integration for receipts and reminders
- 📅 Supplier management with purchase orders
- 📅 Payment plans and installments
- 📅 Mobile-first POS interface

## Phase 5: Advanced Analytics - 🔮 FUTURE
- 🔮 Customer segmentation and loyalty
- 🔮 Enhanced AI forecasting with ML models
- 🔮 Route optimization for delivery
- 🔮 Customer-facing marketplace

---

**Version**: 2.0  
**Last Updated**: November 16, 2025  
**Status**: Production-ready with audit compliance