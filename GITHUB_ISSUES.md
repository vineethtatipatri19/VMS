# GitHub Issues for VMS Project

## Completed Today ✅

### Issue 1: Database Schema Enhancement
**Title:** Enhance database entities with comprehensive business fields

**Labels:** `enhancement`, `database`, `completed`

**Description:**
Enhanced all database tables from basic fields to production-ready schema supporting credit management, profit tracking, wastage monitoring, and supplier management.

**Changes:**
- Added 60+ new fields across 5 existing tables (customers, inventory_items, transactions, sale_items, crate_ledger)
- Created 5 new tables: wastage_log, expiry_alerts, payment_schedules, pricing_tiers, price_history
- Implemented auto-calculation triggers for margin%, profit, overdue status, inventory status
- Created 5 reporting views: customer_summary, inventory_summary, sales_performance, overdue_payments, wastage_summary
- Added 15+ performance indexes

**Migration File:** `infra/migrations/004_enhance_entities.sql`

**Status:** ✅ Completed and applied to database

---

### Issue 2: Backend API Enhancement - Entity Structs
**Title:** Update backend Go structs to match enhanced database schema

**Labels:** `enhancement`, `backend`, `completed`

**Description:**
Updated all backend entity structs and handlers to support enhanced database fields, enabling comprehensive business operations through the API.

**Changes:**
- **inventory.go**: InventoryItem struct expanded from 10 to 35 fields
  - Added cost tracking, supplier management, stock control, GST fields
  - Updated all CRUD handlers with new fields
  
- **transaction_service.go**: Transaction struct expanded from 8 to 21 fields
  - Added payment tracking, financial fields, delivery tracking
  - Enhanced TransactionRequest with 9 new fields
  
- **crates.go**: CrateLedgerEntry struct expanded from 7 to 9 fields
  - Added crate type, value, and transaction linking
  
- **customers.go**: Customer struct expanded from 10 to 26 fields (completed earlier)
  - Added credit management, KYC, business customer support

**Status:** ✅ Completed and tested

---

### Issue 3: New Entity Handlers - Wastage & Alerts
**Title:** Implement handlers for wastage tracking and expiry alerts

**Labels:** `feature`, `backend`, `completed`

**Description:**
Created new backend handlers and API endpoints for wastage tracking, expiry alerts, payment schedules, and reporting views.

**New File:** `backend/enhanced_entities.go`

**New Endpoints:**
- `GET/POST /api/v1/wastage` - Wastage log management
- `GET/PUT /api/v1/expiry-alerts` - Expiry alert tracking
- `GET /api/v1/payment-schedules` - Installment tracking
- `GET /api/v1/reports/overdue-payments` - Overdue payments report
- `GET /api/v1/reports/wastage-summary` - Wastage analytics

**Status:** ✅ Completed and deployed

---

### Issue 4: Frontend UI/UX Modernization
**Title:** Redesign frontend with modern 2025 design patterns

**Labels:** `enhancement`, `frontend`, `design`, `completed`

**Description:**
Complete frontend overhaul with modern design system including Inter font family, fluid typography, glassmorphism effects, and improved color schemes for better mobile app readiness.

**Changes:**
- Implemented Inter font family with fluid typography scale
- Added glassmorphism card effects with backdrop blur
- Enhanced color palette with gradients and modern accents
- Fixed dark mode visibility issues
- Improved spacing, shadows, and interactive states
- Updated all UI components (Button, Card, Badge, Modal, etc.)

**Status:** ✅ Completed and approved

---

## High Priority - Frontend Integration 🔥

### Issue 5: Update Inventory Form with Enhanced Fields
**Title:** Add cost tracking, supplier, and GST fields to inventory form

**Labels:** `enhancement`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Update the inventory creation/edit form to capture all new fields from the enhanced database schema.

**Required Fields to Add:**
1. **Cost & Pricing**
   - Cost Price (number input)
   - Selling Price (number input)
   - Margin % (read-only, auto-calculated)

2. **Supplier Information**
   - Supplier Name (text input with autocomplete)
   - Supplier ID (optional, hidden)
   - Purchase Invoice Number (text input)

3. **Product Details**
   - Category (dropdown: Vegetables, Fruits, Dairy, Grains, etc.)
   - Sub-category (text input)
   - Barcode (text input with scanner support)
   - SKU (text input)
   - HSN Code (text input for GST)
   - GST Rate (dropdown: 0%, 5%, 12%, 18%, 28%)

4. **Stock Management**
   - Min Stock Level (number input)
   - Reorder Point (number input)
   - Shelf Life (days, number input)
   - Storage Location (text input)
   - Packaging Type (dropdown: Box, Bag, Crate, Loose)

5. **Additional**
   - Weight per Unit (number input, optional)
   - Image Upload (file input)
   - Notes (textarea)

**API Endpoint:** `POST/PUT /api/v1/inventory`

**Acceptance Criteria:**
- [ ] Form validates all required fields
- [ ] Margin % auto-calculates when cost/selling prices change
- [ ] Category dropdown populated from config
- [ ] GST rate dropdown with standard rates
- [ ] Image upload preview before submit
- [ ] Form shows existing data on edit mode
- [ ] Mobile-responsive layout

**Files to Update:**
- `frontend/src/pages/Inventory.js`
- `frontend/src/components/InventoryForm.js` (if separate)

**Priority:** High
**Estimated Effort:** 4-6 hours

---

### Issue 6: Update Customer Form with Credit Management
**Title:** Add credit limits, payment terms, and business fields to customer form

**Labels:** `enhancement`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Update customer creation/edit form to capture enhanced fields for credit management, KYC, and business customer support.

**Required Fields to Add:**
1. **Contact Information**
   - Email (email input)
   - WhatsApp Number (phone input)
   - Alternate Contact (phone input)

2. **Business Details**
   - Customer Type (dropdown: B2B, B2C, Retail, Wholesale)
   - Business Name (text input, required for B2B)
   - GSTIN (text input with validation, 15 chars)

3. **Credit Management**
   - Credit Limit (number input with currency)
   - Current Balance (read-only display)
   - Payment Terms (number input, days)
   - Interest Rate (number input, %)
   - Status (dropdown: Active, Inactive, Blocked)

4. **KYC Documents**
   - KYC Document Type (dropdown: Aadhaar, PAN, Passport, Driving License)
   - KYC Document Number (text input)
   - Aadhaar Verified (checkbox)

5. **Additional**
   - Loyalty Points (read-only display)
   - Total Purchases (read-only display)
   - Last Transaction Date (read-only display)
   - Tags (multi-select or text input)
   - Notes (textarea)

**API Endpoint:** `POST/PUT /api/v1/customers`

**Acceptance Criteria:**
- [ ] Customer type selection shows/hides business fields
- [ ] GSTIN validation for 15-character format
- [ ] Credit limit with visual indicator of usage
- [ ] Current balance displays with color coding (red if overdue)
- [ ] Read-only fields disabled but visible
- [ ] Form validation for required fields based on customer type
- [ ] Mobile-responsive layout

**Files to Update:**
- `frontend/src/pages/Customers.js`

**Priority:** High
**Estimated Effort:** 5-7 hours

---

### Issue 7: Update Transaction Form with Payment Tracking
**Title:** Add payment methods, discounts, and delivery tracking to transaction form

**Labels:** `enhancement`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Enhance transaction creation form to capture payment details, discounts, taxes, and delivery information.

**Required Fields to Add:**
1. **Payment Information**
   - Payment Method (dropdown: Cash, UPI, Card, Bank Transfer, Credit)
   - Payment Reference (text input, optional for cash)
   - Due Date (date picker, for credit sales)

2. **Financial Details**
   - Discount Amount (number input)
   - Discount % (number input, auto-calculates amount)
   - Tax Amount (number input)
   - Tax % (number input, auto-calculates amount)
   - Total with Discount/Tax (read-only, auto-calculated)

3. **Sale Type**
   - Sale Type (dropdown: Cash, Credit, Wholesale, Retail)

4. **Delivery Tracking**
   - Delivery Status (dropdown: Pending, In Transit, Delivered)
   - Delivery Date (date picker)
   - Delivery Address (textarea)

5. **Documentation**
   - Invoice Number (read-only, auto-generated)
   - Receipt Sent (checkbox)
   - Notes (textarea)

**API Endpoint:** `POST /api/v1/transactions`

**Acceptance Criteria:**
- [ ] Payment method changes form requirements (e.g., UPI requires reference)
- [ ] Discount can be entered as amount OR percentage
- [ ] Tax can be entered as amount OR percentage
- [ ] Total recalculates on any change
- [ ] Credit sales require due date
- [ ] Due date validation (must be future date)
- [ ] Invoice number displays after creation
- [ ] Form shows customer credit limit warning if applicable
- [ ] Mobile-responsive layout

**Files to Update:**
- `frontend/src/pages/Transactions.js`

**Priority:** High
**Estimated Effort:** 6-8 hours

---

### Issue 8: Update List Views with Enhanced Fields
**Title:** Display priority fields in inventory, customer, and transaction lists

**Labels:** `enhancement`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Update list/table views to show enhanced fields based on display priorities (HIGH/MEDIUM/LOW) from `docs/ENTITY_FIELDS_SUMMARY.md`.

**Inventory List Updates:**
- Add columns: Category, Cost Price, Selling Price, Margin %, Status, Supplier Name
- Add status badge with colors (available=green, low_stock=yellow, out_of_stock=red, expired=gray)
- Add margin % with color coding (green >30%, yellow 10-30%, red <10%)
- Add filter by: Category, Status, Supplier
- Add sort by: Margin %, Expiry Date, Quantity

**Customer List Updates:**
- Add columns: Type, Credit Limit, Current Balance, Status, Last Transaction
- Add type badge (B2B, B2C, Retail, Wholesale)
- Add status badge (Active=green, Inactive=gray, Blocked=red)
- Add credit usage indicator (progress bar or %)
- Add filter by: Type, Status, Overdue Balance
- Highlight overdue customers in red

**Transaction List Updates:**
- Add columns: Payment Method, Invoice Number, Discount, Tax, Status, Overdue
- Add payment method badge
- Add overdue indicator (red flag icon)
- Add delivery status badge
- Add filter by: Payment Method, Sale Type, Overdue
- Add total amount with discount/tax breakdown on hover

**Acceptance Criteria:**
- [ ] Tables show new columns with proper formatting
- [ ] Status badges use consistent color coding
- [ ] Filters work for all new fields
- [ ] Sort works for all new columns
- [ ] Mobile view collapses to card layout showing priority fields only
- [ ] Loading states for all data fetching
- [ ] Pagination works with filters

**Files to Update:**
- `frontend/src/pages/Inventory.js`
- `frontend/src/pages/Customers.js`
- `frontend/src/pages/Transactions.js`

**Priority:** High
**Estimated Effort:** 6-8 hours

---

## High Priority - New Features 🚀

### Issue 9: Create Wastage Log Page
**Title:** Build wastage tracking UI with log entry and analytics

**Labels:** `feature`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Create a new page for tracking inventory wastage with entry form and analytics dashboard.

**Features Required:**
1. **Wastage Log List**
   - Table showing: Item Name, Lot Number, Quantity, Reason, Value Lost, Date
   - Filter by: Reason Code, Date Range, Item Name
   - Sort by: Date, Value Lost
   - Export to CSV functionality

2. **Wastage Entry Form**
   - Inventory Item selector (dropdown with search)
   - Auto-populate: Item Name, Lot Number, Current Quantity
   - Quantity Wasted (number input with max validation)
   - Reason Code (dropdown: Expired, Damaged, Spoiled, Recalled, Other)
   - Reason Details (textarea)
   - Estimated Value (auto-calculated: quantity × cost_price)
   - Reported By (text input, default to current user)

3. **Analytics Dashboard**
   - Total value lost (this month/year)
   - Most wasted items (bar chart)
   - Wastage by reason (pie chart)
   - Trend over time (line chart)
   - Top 5 items by waste count

**API Endpoints:**
- `GET /api/v1/wastage`
- `POST /api/v1/wastage`
- `GET /api/v1/reports/wastage-summary`

**Acceptance Criteria:**
- [ ] Wastage entry form validates quantity against current stock
- [ ] Form shows estimated value calculation
- [ ] List view shows all wastage logs with filters
- [ ] Analytics dashboard displays charts correctly
- [ ] Export to CSV includes all fields
- [ ] Mobile-responsive design
- [ ] Success toast on wastage entry
- [ ] Confirmation dialog before submitting

**New Files to Create:**
- `frontend/src/pages/Wastage.js`
- `frontend/src/components/WastageForm.js`
- `frontend/src/components/WastageChart.js`

**Priority:** High
**Estimated Effort:** 8-10 hours

---

### Issue 10: Create Expiry Alerts Dashboard
**Title:** Build expiry alert monitoring with acknowledgment system

**Labels:** `feature`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Create a dashboard widget and dedicated page for monitoring items approaching expiry with alert management.

**Features Required:**
1. **Dashboard Widget** (on main dashboard)
   - Count of pending expiry alerts
   - Top 3 items expiring soonest
   - "View All" link to dedicated page
   - Color-coded by urgency (red ≤1 day, yellow 2-3 days)

2. **Expiry Alerts Page**
   - Table: Item Name, Lot Number, Expiry Date, Days to Expiry, Quantity, Status
   - Filter by: Status (Pending, Acknowledged, Resolved)
   - Sort by: Days to Expiry, Quantity
   - Bulk action: Acknowledge selected alerts
   - Action buttons: Acknowledge, Mark Resolved, View Item

3. **Alert Details Modal**
   - Full item information
   - Recommended actions:
     - Create wastage log (if expired)
     - Apply discount (to sell quickly)
     - Move to wastage
   - Action buttons directly in modal

**API Endpoints:**
- `GET /api/v1/expiry-alerts`
- `PUT /api/v1/expiry-alerts/{id}`

**Acceptance Criteria:**
- [ ] Dashboard widget shows real-time count
- [ ] Alerts table color-codes by urgency
- [ ] Acknowledge button updates status immediately
- [ ] Resolved alerts disappear from pending view
- [ ] Modal provides quick actions
- [ ] Mobile-responsive design
- [ ] Auto-refresh every 5 minutes
- [ ] Browser notification for new critical alerts (≤1 day)

**Files to Update:**
- `frontend/src/pages/Dashboard.js` (add widget)

**New Files to Create:**
- `frontend/src/pages/ExpiryAlerts.js`
- `frontend/src/components/ExpiryAlertWidget.js`
- `frontend/src/components/ExpiryAlertModal.js`

**Priority:** High
**Estimated Effort:** 6-8 hours

---

### Issue 11: Create Overdue Payments Report
**Title:** Build overdue payments tracking and reminder system

**Labels:** `feature`, `frontend`, `high-priority`

**Assignee:** [To be assigned]

**Description:**
Create a dedicated page for tracking overdue payments with customer contact and reminder features.

**Features Required:**
1. **Overdue Payments List**
   - Table: Customer Name, Invoice Number, Amount, Due Date, Days Overdue, Balance
   - Sort by: Days Overdue, Amount
   - Filter by: Days Overdue Range, Customer
   - Color-code by severity (yellow 1-7 days, orange 8-30 days, red >30 days)

2. **Summary Cards**
   - Total overdue amount
   - Number of overdue customers
   - Average days overdue
   - Trend vs last month

3. **Customer Action Buttons**
   - Send Reminder (WhatsApp icon if number exists)
   - View Customer Details
   - View Transaction Details
   - Record Payment

4. **Quick Actions**
   - Export to CSV
   - Send bulk reminders
   - Generate collection report

**API Endpoints:**
- `GET /api/v1/reports/overdue-payments`
- `GET /api/v1/customers/{id}` (for details)
- `GET /api/v1/transactions/{id}` (for details)

**Acceptance Criteria:**
- [ ] List displays all overdue payments correctly
- [ ] Color coding by overdue severity
- [ ] Summary cards calculate totals accurately
- [ ] WhatsApp link works (if number exists)
- [ ] Customer details modal shows credit history
- [ ] Export includes all relevant fields
- [ ] Mobile-responsive design
- [ ] Auto-refresh data

**New Files to Create:**
- `frontend/src/pages/OverduePayments.js`
- `frontend/src/components/OverduePaymentCard.js`

**Priority:** High
**Estimated Effort:** 5-7 hours

---

## Medium Priority - Feature Enhancements 🎯

### Issue 12: Implement Credit Limit Enforcement
**Title:** Add credit limit validation in transaction creation

**Labels:** `enhancement`, `backend`, `frontend`, `medium-priority`

**Assignee:** [To be assigned]

**Description:**
Add validation logic to prevent credit sales exceeding customer credit limits, with override capability for admins.

**Backend Changes:**
- In `createTransaction` handler, check customer credit limit before creating credit sale
- Query: `SELECT credit_limit, current_balance FROM customers WHERE id = ?`
- Validation: `new_balance = current_balance + transaction_total`
- If `new_balance > credit_limit`, return 400 error with message
- Add override parameter for admin users

**Frontend Changes:**
- Show credit limit warning in transaction form
- Display: "Credit Available: ₹X / ₹Y (Z% used)"
- Progress bar showing credit utilization
- Warning modal if transaction exceeds limit
- Admin users see "Override Credit Limit" checkbox

**API Changes:**
```json
POST /api/v1/transactions
{
  "customerId": "uuid",
  "type": "sale",
  "items": [...],
  "saleType": "credit",
  "overrideCreditLimit": false
}

Error Response (400):
{
  "error": "Credit limit exceeded",
  "creditLimit": 50000,
  "currentBalance": 45000,
  "transactionAmount": 8000,
  "availableCredit": 5000
}
```

**Acceptance Criteria:**
- [ ] Backend validates credit limit on credit sales
- [ ] Cash sales bypass credit check
- [ ] Error response includes all relevant amounts
- [ ] Frontend shows credit usage in transaction form
- [ ] Warning modal appears before exceeding limit
- [ ] Admin override works correctly
- [ ] Current balance updates after transaction

**Files to Update:**
- `backend/transaction_service.go`
- `frontend/src/pages/Transactions.js`

**Priority:** Medium
**Estimated Effort:** 4-6 hours

---

### Issue 13: Build Profit Margin Dashboard
**Title:** Create profit analysis dashboard with item and category insights

**Labels:** `feature`, `frontend`, `analytics`, `medium-priority`

**Assignee:** [To be assigned]

**Description:**
Build analytics dashboard showing profit margins across items, categories, and time periods.

**Features Required:**
1. **Overview Cards**
   - Total profit (today/week/month/year)
   - Average margin %
   - Best performing item
   - Worst performing item

2. **Profit by Item Chart**
   - Bar chart: Top 10 items by profit amount
   - Toggle: Profit Amount vs Margin %
   - Filter by date range

3. **Profit by Category Chart**
   - Pie chart: Profit distribution by category
   - Filter by date range

4. **Margin Analysis Table**
   - Table: Item Name, Units Sold, Revenue, Cost, Profit, Margin %
   - Sort by any column
   - Color-code margin % (green >30%, yellow 10-30%, red <10%)
   - Filter by: Category, Date Range, Margin Range

5. **Trend Analysis**
   - Line chart: Profit over time (daily/weekly/monthly)
   - Compare multiple time periods

**Data Source:**
- Query from `sales_performance` view
- Calculate from `sale_items` table with cost_per_unit and profit fields

**Acceptance Criteria:**
- [ ] All charts render correctly with real data
- [ ] Date range filters work across all widgets
- [ ] Margin % color coding consistent
- [ ] Export to PDF functionality
- [ ] Mobile-responsive design
- [ ] Loading states for all charts
- [ ] Empty states for no data

**New Files to Create:**
- `frontend/src/pages/ProfitAnalysis.js`
- `frontend/src/components/ProfitChart.js`
- `frontend/src/components/MarginTable.js`

**Priority:** Medium
**Estimated Effort:** 8-10 hours

---

### Issue 14: Add Barcode Scanner Support
**Title:** Integrate barcode scanning for inventory and sales

**Labels:** `feature`, `frontend`, `medium-priority`

**Assignee:** [To be assigned]

**Description:**
Add barcode scanning capability for quick inventory lookup and sale item selection.

**Features Required:**
1. **Barcode Scanner Component**
   - Use device camera or external USB scanner
   - Library: `react-webcam` or `quagga2`
   - Scan formats: EAN-13, UPC-A, Code 128

2. **Inventory Page Integration**
   - Scanner button in inventory list
   - On scan, search inventory by barcode
   - Open item details if found
   - Show "not found" message with option to add

3. **Transaction Page Integration**
   - Scanner button in add item section
   - On scan, auto-add item to cart
   - If multiple lots exist, show selection modal
   - If not found, show search form

4. **Settings**
   - Enable/disable scanner
   - Select scanner type (camera/USB)
   - Test scanner functionality

**Acceptance Criteria:**
- [ ] Scanner detects barcodes accurately (>95%)
- [ ] Works with both camera and USB scanners
- [ ] Fast lookup (<500ms after scan)
- [ ] Clear error messages for invalid/not found codes
- [ ] Works on mobile devices
- [ ] Scanner respects permissions
- [ ] Can toggle scanner on/off

**New Files to Create:**
- `frontend/src/components/BarcodeScanner.js`
- `frontend/src/hooks/useBarcodeScanner.js`

**Dependencies:**
- `npm install react-webcam` or `npm install quagga2`

**Priority:** Medium
**Estimated Effort:** 6-8 hours

---

### Issue 15: Create Payment Schedule Management
**Title:** Build installment plan creation and tracking UI

**Labels:** `feature`, `frontend`, `medium-priority`

**Assignee:** [To be assigned]

**Description:**
Create UI for managing payment schedules and installment plans for credit sales.

**Features Required:**
1. **Create Payment Schedule** (in Transaction Form)
   - Toggle: "Enable Installments"
   - Number of installments (dropdown: 2, 3, 4, 6, 12)
   - Auto-calculate: Amount per installment
   - Auto-generate due dates (monthly intervals)
   - Preview schedule before saving

2. **Payment Schedule List Page**
   - Table: Customer, Invoice, Installment, Due Date, Amount, Paid, Status
   - Filter by: Customer, Status, Due Date Range
   - Sort by: Due Date, Amount
   - Status badges: Pending, Partial, Paid, Overdue

3. **Record Payment Modal**
   - Input: Amount Paid (validates ≤ installment amount)
   - Payment Date (defaults to today)
   - Payment Method
   - Payment Reference
   - Notes
   - Auto-update status to "Partial" or "Paid"

4. **Customer Payment History**
   - Show all installments for a customer
   - Payment completion timeline
   - Next due date highlighted

**API Endpoints:**
- `GET /api/v1/payment-schedules`
- `POST /api/v1/payment-schedules` (created with transaction)
- `PUT /api/v1/payment-schedules/{id}` (record payment)

**Acceptance Criteria:**
- [ ] Installment creation validates total equals transaction amount
- [ ] Due dates generate correctly (monthly intervals)
- [ ] Payment recording updates status correctly
- [ ] Overdue installments auto-flagged
- [ ] Customer view shows all their installments
- [ ] Mobile-responsive design
- [ ] Export to CSV functionality

**New Files to Create:**
- `frontend/src/pages/PaymentSchedules.js`
- `frontend/src/components/PaymentScheduleForm.js`
- `frontend/src/components/RecordPaymentModal.js`

**Priority:** Medium
**Estimated Effort:** 8-10 hours

---

## Low Priority - Nice to Have ✨

### Issue 16: Add Bulk Import for Inventory
**Title:** CSV import functionality for bulk inventory upload

**Labels:** `feature`, `backend`, `frontend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Add ability to import multiple inventory items via CSV file upload.

**Features:**
- CSV template download
- File upload with validation
- Preview before import
- Error handling for invalid rows
- Success/failure summary report

**Priority:** Low
**Estimated Effort:** 6-8 hours

---

### Issue 17: Add Multi-language Support
**Title:** Implement i18n for English and Hindi

**Labels:** `enhancement`, `frontend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Add internationalization support for English and Hindi languages.

**Features:**
- Language switcher in header
- Translation files for all UI text
- Persist language preference
- RTL support preparation

**Priority:** Low
**Estimated Effort:** 10-12 hours

---

### Issue 18: Add Email Notifications
**Title:** Email alerts for expiry, overdue payments, and low stock

**Labels:** `feature`, `backend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Implement email notification system for important alerts.

**Features:**
- Email service integration (SendGrid/AWS SES)
- Templates for different alert types
- User notification preferences
- Daily summary emails

**Priority:** Low
**Estimated Effort:** 8-10 hours

---

### Issue 19: Add WhatsApp Integration
**Title:** WhatsApp message sending for payment reminders

**Labels:** `feature`, `backend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Integrate WhatsApp Business API for sending payment reminders and transaction receipts.

**Features:**
- WhatsApp Business API setup
- Message templates
- Send payment reminders
- Send transaction receipts
- Delivery status tracking

**Priority:** Low
**Estimated Effort:** 12-15 hours

---

### Issue 20: Create Mobile App (React Native)
**Title:** Build iOS and Android apps using React Native

**Labels:** `feature`, `mobile`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Create native mobile applications for iOS and Android using React Native, reusing frontend components and API.

**Scope:**
- React Native setup
- Navigation structure
- Core features: Inventory, Customers, Transactions
- Camera integration for barcode scanning
- Push notifications
- Offline mode with sync

**Priority:** Low
**Estimated Effort:** 40-60 hours

---

## Testing & Documentation 📝

### Issue 21: Write API Documentation
**Title:** Complete OpenAPI/Swagger documentation for all endpoints

**Labels:** `documentation`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Create comprehensive API documentation using OpenAPI/Swagger specification.

**Deliverables:**
- OpenAPI 3.0 spec file
- Swagger UI integration
- Example requests/responses
- Authentication documentation
- Error code reference

**Priority:** Low
**Estimated Effort:** 6-8 hours

---

### Issue 22: Write Unit Tests for Backend
**Title:** Add unit tests for all backend handlers

**Labels:** `testing`, `backend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Write comprehensive unit tests for all Go backend handlers and business logic.

**Coverage Target:** 80%

**Test Areas:**
- CRUD operations for all entities
- Credit limit validation
- Auto-calculation triggers
- Error handling
- Authentication/authorization

**Priority:** Low
**Estimated Effort:** 20-25 hours

---

### Issue 23: Write Frontend Tests
**Title:** Add React component tests and E2E tests

**Labels:** `testing`, `frontend`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Write unit tests for React components and E2E tests for critical user flows.

**Test Framework:** Jest + React Testing Library + Cypress

**Coverage:**
- Component unit tests
- Form validation tests
- API integration tests
- E2E: Create inventory → Create sale → Record payment

**Priority:** Low
**Estimated Effort:** 25-30 hours

---

### Issue 24: Create User Documentation
**Title:** Write user manual and video tutorials

**Labels:** `documentation`, `low-priority`

**Assignee:** [To be assigned]

**Description:**
Create comprehensive user documentation with screenshots and video tutorials.

**Deliverables:**
- User manual (PDF)
- Video tutorials for key workflows
- FAQ section
- Troubleshooting guide
- Admin guide

**Priority:** Low
**Estimated Effort:** 15-20 hours

---

## Summary

**Total Issues:** 24

**Completed Today:** 4 issues
- Database schema enhancement
- Backend entity struct updates
- New entity handlers
- Frontend UI modernization

**High Priority:** 7 issues (frontend integration + new features)
**Medium Priority:** 4 issues (feature enhancements)
**Low Priority:** 9 issues (nice to have + testing)

**Estimated Total Effort:** 300-400 hours

**Recommended Sprint Plan:**
- **Sprint 1 (2 weeks):** Issues 5-8 (frontend integration)
- **Sprint 2 (2 weeks):** Issues 9-11 (new feature pages)
- **Sprint 3 (1-2 weeks):** Issues 12-15 (enhancements)
- **Sprint 4+:** Low priority items as time permits
