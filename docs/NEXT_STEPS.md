# PGVMS - Next Steps & Feature Roadmap

**Last Updated**: November 16, 2025  
**Current Version**: v1.0 (Base Implementation Complete)

---

## 📊 Current System Status

### ✅ Implemented Features (v2.0 - Audit Compliant)
- **Authentication**: User registration, login with JWT
- **Dashboard**: KPIs, quick stats, recent activity
- **Inventory Management**: Full CRUD with soft delete, FEFO sorting, expiry tracking, status badges
- **Customer Management**: Full CRUD with audit trail, contact info, KYC status
- **Transaction Ledger**: Sales and payments with edit/delete, multi-item support, balance tracking
- **Crate Management**: Issue/return tracking with delete, balance per customer
- **AI Forecasting**: Google Gemini integration (stub mode available)
- **Reports**: Sales, inventory, customer financial statements
- **Wastage Tracking** ✨ NEW: Log and track spoiled/damaged inventory
- **Expiry Alerts** ✨ NEW: Proactive monitoring with acknowledgment workflow
- **Audit Trail System** ✨ NEW: Comprehensive soft delete with attestation
- **Delete Confirmation** ✨ NEW: Attestation-required delete operations across all entities

### 🔒 Audit & Compliance Features (v2.0)
- ✅ **Soft Delete System**: No data permanently deleted, full audit trail
- ✅ **Attestation Requirement**: Must type "I CONFIRM DELETE" for all deletions
- ✅ **Audit Fields**: deleted_at, deleted_by, deletion_reason on all entities
- ✅ **Edit History**: Track who updated what and when
- ✅ **Reason Tracking**: Required reason for all delete operations
- ✅ **Data Recovery**: Soft-deleted records can be restored

### 📈 Current Data
- **Customers**: 3
- **Inventory Items**: 2
- **Transactions**: 3 (2 sales, 1 payment)
- **Total Sales**: ₹3,915
- **Total Payments**: ₹2,000
- **Wastage Entries**: Trackable with photo evidence
- **Expiry Alerts**: Auto-generated for items expiring soon

---

## 🎯 TIER 1: Critical Business Features (High Impact - Build Next)

### 1. Multi-Business/Multi-Tenant Support
**Business Problem**: System currently supports only one business - limits scalability  
**Impact**: High - Enables SaaS model, multiple franchises  
**Effort**: Medium (3-4 weeks)

**Features**:
- Business registration and onboarding workflow
- Tenant isolation (schema per business or tenant_id column)
- Business profile management (logo, GST number, address, contact)
- Staff/employee management with role-based permissions (Admin, Manager, Staff, Viewer)
- Business switcher for users managing multiple locations
- Subscription/pricing tier management
- Business settings and preferences

**Technical Requirements**:
- Add `businesses` table with tenant_id
- Add `tenant_id` to all existing tables
- Middleware to filter queries by tenant
- Row-level security policies
- Business-specific JWT claims

**Database Changes**:
```sql
CREATE TABLE businesses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  owner_email TEXT NOT NULL,
  gst_number TEXT,
  address TEXT,
  phone TEXT,
  logo_url TEXT,
  subscription_tier TEXT DEFAULT 'free',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE business_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  business_id UUID REFERENCES businesses(id),
  user_id UUID REFERENCES users(id),
  role TEXT CHECK (role IN ('admin','manager','staff','viewer')),
  created_at TIMESTAMPTZ DEFAULT now(),
  UNIQUE(business_id, user_id)
);
```

---

### 2. Advanced Payment Tracking & Credit Management
**Business Problem**: Credit management is manual, no alerts for overdue payments  
**Impact**: Very High - Directly improves cash flow and reduces bad debt  
**Effort**: Medium (2-3 weeks)

**Features**:
- **Credit Limits**: Set maximum credit per customer, auto-block sales when exceeded
- **Payment Plans**: Schedule installments, track payment schedules
- **Partial Payments**: Record multiple payments against a single sale
- **Payment History Timeline**: Visual timeline with all payment activities
- **Overdue Alerts**: Automatic SMS/WhatsApp reminders for overdue payments
- **Interest Calculation**: Optional interest on overdue amounts (configurable rate)
- **Payment Methods**: Differentiate UPI, Cash, Card, Bank Transfer
- **Payment Receipts**: Generate payment acknowledgment receipts
- **Collection Dashboard**: View all overdue accounts in one place
- **Payment Terms**: Net 15, Net 30, custom terms per customer

**Technical Requirements**:
```sql
-- Add to customers table
ALTER TABLE customers ADD COLUMN credit_limit NUMERIC DEFAULT 0;
ALTER TABLE customers ADD COLUMN payment_terms_days INTEGER DEFAULT 30;
ALTER TABLE customers ADD COLUMN interest_rate NUMERIC DEFAULT 0;

-- Add payment_method to transactions
ALTER TABLE transactions ADD COLUMN payment_method TEXT;
ALTER TABLE transactions ADD COLUMN due_date DATE;
ALTER TABLE transactions ADD COLUMN is_overdue BOOLEAN DEFAULT false;

-- Create payment schedule table
CREATE TABLE payment_schedules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id UUID REFERENCES customers(id),
  transaction_id UUID REFERENCES transactions(id),
  due_date DATE NOT NULL,
  amount_due NUMERIC NOT NULL,
  amount_paid NUMERIC DEFAULT 0,
  status TEXT CHECK (status IN ('pending','paid','overdue','cancelled')),
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Create payment reminders log
CREATE TABLE payment_reminders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id UUID REFERENCES customers(id),
  transaction_id UUID REFERENCES transactions(id),
  reminder_type TEXT CHECK (reminder_type IN ('sms','whatsapp','email')),
  sent_at TIMESTAMPTZ DEFAULT now(),
  status TEXT
);
```

**API Endpoints**:
- `POST /api/v1/customers/{id}/credit-limit` - Set credit limit
- `GET /api/v1/customers/{id}/payment-schedule` - Get payment schedule
- `POST /api/v1/payments/partial` - Record partial payment
- `GET /api/v1/payments/overdue` - List all overdue payments
- `POST /api/v1/payments/send-reminder` - Send payment reminder

**UI Components**:
- Credit limit indicator on customer profile
- Payment schedule calendar view
- Overdue payments dashboard (red alerts)
- Payment reminder settings page
- Payment method selector on transaction form

**Expected ROI**:
- Reduce bad debt by 30-40%
- Improve collection time by 25%
- Reduce manual follow-up effort by 60%

---

### 3. ✅ Real Expiry & Wastage Tracking (COMPLETED)
**Business Problem**: Perishable goods expire causing losses, no proactive monitoring  
**Impact**: Very High - Directly reduces waste and improves profitability  
**Effort**: Medium (2 weeks) - **STATUS: IMPLEMENTED**

**Features** (✅ = Implemented, ⏳ = Planned):
- ✅ **Wastage Log**: Daily entry for spoiled/expired items with reason codes
- ✅ **Expiry Alerts**: Automatic alerts for items expiring within 3 days
- ✅ **Expiry Alert Acknowledgment**: Users can acknowledge and dismiss alerts
- ✅ **Photo Evidence**: Upload photos of wasted items
- ✅ **Reason Categorization**: expired, damaged, spoiled, pest, other
- ✅ **Cost Impact Tracking**: Financial loss calculation per wastage entry
- ✅ **Wastage Reports**: View and filter wastage entries
- ✅ **Soft Delete Support**: Wastage and expiry entries can be safely deleted
- ⏳ **Auto-mark Expired**: Daily background job to flag expired inventory (planned)
- ⏳ **Waste vs Sales Ratio**: KPI showing waste percentage per item (planned)
- ⏳ **Smart Reorder Points**: Calculate based on actual shelf life and sales velocity (planned)
- ⏳ **Near-Expiry Pricing**: Suggest discounts for items approaching expiry (planned)
- ⏳ **Supplier Quality Metrics**: Link wastage to suppliers (planned)

**Technical Requirements**:
```sql
-- Create wastage log table
CREATE TABLE wastage_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id),
  quantity_wasted NUMERIC NOT NULL,
  unit TEXT NOT NULL,
  reason TEXT CHECK (reason IN ('expired','spoiled','damaged','pest','other')),
  notes TEXT,
  cost_impact NUMERIC,
  logged_by UUID REFERENCES users(id),
  logged_at TIMESTAMPTZ DEFAULT now()
);

-- Create expiry alerts table
CREATE TABLE expiry_alerts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id),
  alert_date DATE NOT NULL,
  days_to_expiry INTEGER,
  status TEXT CHECK (status IN ('pending','acknowledged','resolved')),
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Add fields to inventory_items
ALTER TABLE inventory_items ADD COLUMN cost_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN is_expired BOOLEAN DEFAULT false;
ALTER TABLE inventory_items ADD COLUMN shelf_life_days INTEGER;
```

**Background Jobs**:
- Daily job (runs at midnight): Check expiry dates, mark expired items
- Daily job (runs at 8 AM): Generate expiry alerts for items expiring in 3 days
- Weekly job: Calculate wastage statistics and trends

**API Endpoints**:
- `POST /api/v1/wastage/log` - Log wasted item
- `GET /api/v1/wastage/report` - Get wastage report
- `GET /api/v1/inventory/expiring-soon` - Get items expiring soon
- `POST /api/v1/inventory/{id}/acknowledge-expiry` - Acknowledge expiry alert
- `GET /api/v1/analytics/waste-metrics` - Get wastage KPIs

**UI Components**:
- Wastage entry form (quick add from inventory page)
- Expiry alert dashboard widget (red for expired, yellow for expiring soon)
- Wastage report with charts
- Near-expiry items list with suggested discounts
- Wastage cost impact summary

**Expected ROI**:
- Reduce wastage by 15-20%
- Improve profit margins by 10-12%
- Better inventory planning

---

### 4. Supplier/Vendor Management
**Business Problem**: No visibility into where stock comes from, can't track quality  
**Impact**: High - Enables quality control, better pricing negotiations  
**Effort**: Medium (2-3 weeks)

**Features**:
- **Supplier Master**: Name, contact, GST, payment terms, bank details
- **Purchase Orders**: Create POs, track order status, fulfillment
- **Goods Receipt**: Record incoming stock with invoice numbers
- **Purchase Invoices**: Link invoices to payments and inventory
- **Supplier Payments**: Track payables (accounts payable)
- **Supplier Performance**: Quality ratings, on-time delivery, defect rates
- **Price History**: Track price changes over time per supplier per item
- **Supplier Comparison**: Compare prices across suppliers
- **Preferred Suppliers**: Mark preferred suppliers per item
- **Supplier Ledger**: Running balance of what you owe suppliers

**Technical Requirements**:
```sql
-- Create suppliers table
CREATE TABLE suppliers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  contact_person TEXT,
  phone TEXT,
  email TEXT,
  address TEXT,
  gst_number TEXT,
  payment_terms_days INTEGER DEFAULT 30,
  bank_account_number TEXT,
  ifsc_code TEXT,
  is_active BOOLEAN DEFAULT true,
  rating NUMERIC CHECK (rating BETWEEN 0 AND 5),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Create purchase orders table
CREATE TABLE purchase_orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  po_number TEXT UNIQUE NOT NULL,
  supplier_id UUID REFERENCES suppliers(id),
  order_date DATE NOT NULL,
  expected_delivery_date DATE,
  status TEXT CHECK (status IN ('draft','sent','confirmed','received','cancelled')),
  total_amount NUMERIC,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Create PO items
CREATE TABLE po_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  po_id UUID REFERENCES purchase_orders(id),
  item_name TEXT NOT NULL,
  quantity NUMERIC NOT NULL,
  unit TEXT NOT NULL,
  unit_price NUMERIC NOT NULL,
  total NUMERIC NOT NULL
);

-- Create goods receipts
CREATE TABLE goods_receipts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  po_id UUID REFERENCES purchase_orders(id),
  supplier_id UUID REFERENCES suppliers(id),
  receipt_date DATE NOT NULL,
  invoice_number TEXT,
  invoice_date DATE,
  total_amount NUMERIC,
  quality_rating INTEGER CHECK (quality_rating BETWEEN 1 AND 5),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Create supplier payments
CREATE TABLE supplier_payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  supplier_id UUID REFERENCES suppliers(id),
  payment_date DATE NOT NULL,
  amount NUMERIC NOT NULL,
  payment_method TEXT,
  reference_number TEXT,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Add supplier tracking to inventory
ALTER TABLE inventory_items ADD COLUMN supplier_id UUID REFERENCES suppliers(id);
ALTER TABLE inventory_items ADD COLUMN purchase_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN purchase_invoice TEXT;
```

**API Endpoints**:
- `GET /api/v1/suppliers` - List suppliers
- `POST /api/v1/suppliers` - Create supplier
- `GET /api/v1/purchase-orders` - List POs
- `POST /api/v1/purchase-orders` - Create PO
- `POST /api/v1/goods-receipts` - Record goods receipt
- `GET /api/v1/suppliers/{id}/ledger` - Supplier account statement
- `GET /api/v1/suppliers/{id}/performance` - Performance metrics
- `POST /api/v1/supplier-payments` - Record payment to supplier

**UI Pages**:
- Suppliers list and detail page
- Create/edit supplier form
- Purchase order creation wizard
- Goods receipt entry form
- Supplier payments page
- Supplier performance dashboard

---

### 5. Smart Pricing & Profit Margin Tracking
**Business Problem**: No visibility into profitability per item or transaction  
**Impact**: High - Understand what's profitable, optimize pricing  
**Effort**: Low-Medium (1-2 weeks)

**Features**:
- **Cost Price Tracking**: Record purchase/cost price for every item
- **Selling Price Management**: Track selling price separately
- **Margin Calculation**: Auto-calculate margin % and absolute profit
- **Dynamic Pricing**: Suggest price reductions for near-expiry items
- **Profit per Transaction**: Show margin on each sale line item
- **Profit Reports**: Daily/weekly/monthly profit analysis
- **Item Profitability Ranking**: Sort items by profit margin
- **Bulk Pricing Tiers**: Set different prices for wholesale quantities
- **Competitive Pricing**: Track competitor prices (manual entry)
- **Price Change History**: Audit trail of all price changes

**Technical Requirements**:
```sql
-- Already have in inventory_items (add if missing)
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS cost_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN selling_price NUMERIC;
ALTER TABLE inventory_items ADD COLUMN margin_percentage NUMERIC;

-- Add profit tracking to sale_items
ALTER TABLE sale_items ADD COLUMN cost_per_unit NUMERIC;
ALTER TABLE sale_items ADD COLUMN profit NUMERIC;

-- Create pricing tiers table
CREATE TABLE pricing_tiers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id),
  min_quantity NUMERIC NOT NULL,
  max_quantity NUMERIC,
  price_per_unit NUMERIC NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Create price history
CREATE TABLE price_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_item_id UUID REFERENCES inventory_items(id),
  old_price NUMERIC,
  new_price NUMERIC,
  change_reason TEXT,
  changed_by UUID REFERENCES users(id),
  changed_at TIMESTAMPTZ DEFAULT now()
);
```

**Calculations**:
- Margin % = ((Selling Price - Cost Price) / Selling Price) × 100
- Profit = (Selling Price - Cost Price) × Quantity
- ROI = (Profit / Cost) × 100

**API Endpoints**:
- `PUT /api/v1/inventory/{id}/pricing` - Update prices
- `GET /api/v1/analytics/profit-margins` - Get profit analysis
- `GET /api/v1/inventory/profitability-ranking` - Rank items by profit
- `POST /api/v1/pricing-tiers` - Set bulk pricing
- `GET /api/v1/inventory/{id}/price-history` - View price changes

**UI Enhancements**:
- Add cost price field to inventory form
- Show margin % on inventory list
- Profit column in transaction details
- Profit dashboard widget
- Low-margin items alert
- Price change approval workflow (optional)

---

## 🚀 TIER 2: Operational Efficiency (Medium-High Impact)

### 6. Mobile-First POS Interface
**Business Problem**: Current web UI not optimized for mobile speed  
**Impact**: Medium-High - Faster checkout, better user experience  
**Effort**: Medium (3 weeks)

**Features**:
- Touch-optimized large buttons and controls
- Barcode/QR scanning for items
- Full offline mode with local storage
- Auto-sync when connection restored
- Quick customer search (type-ahead)
- Favorite/frequently sold items section
- Swipe gestures for navigation
- Voice input for quantities
- Calculator integration
- Receipt preview before print

**Technical Stack**:
- Progressive Web App (PWA) for offline
- IndexedDB for local storage
- Service workers for background sync
- ZXing library for barcode scanning
- Responsive touch-friendly CSS

---

### 7. WhatsApp Business API Integration
**Business Problem**: No automated customer communication  
**Impact**: High - Better engagement, faster collections  
**Effort**: Medium (2-3 weeks including approval)

**Features**:
- Auto-send receipts via WhatsApp after transaction
- Send outstanding balance notifications
- Payment confirmation messages
- Expiry alerts for customers (for their purchases)
- Broadcast messages (price updates, offers)
- Order taking via WhatsApp bot
- Two-way communication (balance inquiry)
- Scheduled reminders for payments
- Customizable message templates
- Delivery notifications

**Requirements**:
- WhatsApp Business API account
- Facebook Business Manager verification
- Message template approval
- Webhook server for incoming messages

**Integration**:
- Use Twilio WhatsApp API or official WhatsApp Business API
- Template messages for transactional notifications
- Session messages for conversations
- Media messages for receipts (PDF)

---

### 8. GST Billing & Tax Compliance
**Business Problem**: No tax handling, compliance issues  
**Impact**: High - Legal requirement for most businesses  
**Effort**: Medium (2-3 weeks)

**Features**:
- GST invoice generation with GSTIN, HSN codes
- Tax calculation (CGST, SGST, IGST)
- Input and output tax tracking
- GSTR-1, GSTR-3B report generation
- E-invoice integration (for B2B > ₹50L)
- E-way bill generation
- TDS tracking and deduction
- Tax summary reports
- Reverse charge mechanism
- Composition scheme support

**Technical Requirements**:
```sql
-- Add to businesses table
ALTER TABLE businesses ADD COLUMN gstin TEXT;
ALTER TABLE businesses ADD COLUMN gst_scheme TEXT CHECK (gst_scheme IN ('regular','composition'));

-- Add to customers
ALTER TABLE customers ADD COLUMN gstin TEXT;
ALTER TABLE customers ADD COLUMN customer_type TEXT CHECK (customer_type IN ('b2b','b2c'));

-- Add to inventory_items
ALTER TABLE inventory_items ADD COLUMN hsn_code TEXT;
ALTER TABLE inventory_items ADD COLUMN gst_rate NUMERIC;

-- Tax tracking
CREATE TABLE tax_invoices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_number TEXT UNIQUE NOT NULL,
  transaction_id UUID REFERENCES transactions(id),
  invoice_date DATE NOT NULL,
  gstin TEXT,
  place_of_supply TEXT,
  taxable_amount NUMERIC,
  cgst_amount NUMERIC,
  sgst_amount NUMERIC,
  igst_amount NUMERIC,
  total_tax NUMERIC,
  invoice_value NUMERIC,
  irn TEXT, -- for e-invoice
  qr_code TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
```

---

### 9. Weighing Scale Integration
**Business Problem**: Manual weight entry causes errors  
**Impact**: Medium - Accuracy and speed improvement  
**Effort**: Medium (device dependent)

**Features**:
- Direct USB/Bluetooth scale integration
- Auto-populate weight in transaction
- Support multiple scale brands (Avery, Essae, etc.)
- Tare weight calculation
- Weight stability indicator
- Unit conversion (kg to gm)
- Weight history tracking

**Technical Approach**:
- USB HID device communication
- Serial port communication
- Bluetooth LE for wireless scales
- WebUSB API for browser access
- Fallback to manual entry

---

### 10. Smart Dashboard with Alerts
**Business Problem**: Reactive management, no proactive alerts  
**Impact**: Medium - Better decision making  
**Effort**: Low-Medium (1-2 weeks)

**Features**:
- Daily morning summary email/SMS
- Low stock alerts
- High debt customer alerts
- Expiry warnings (3 days)
- Unusual activity detection
- Target vs actual tracking
- Customizable alert thresholds
- Alert history log
- Snooze/dismiss alerts
- Alert priority levels

---

## 🔮 TIER 3: Advanced Analytics (Medium Impact, Future)

### 11. Customer Segmentation & Loyalty
**Effort**: Medium (3-4 weeks)

**Features**:
- RFM Analysis (Recency, Frequency, Monetary value)
- Customer lifetime value (CLV) calculation
- Loyalty points program
- Tier-based benefits (Silver, Gold, Platinum)
- VIP customer identification
- Churn prediction (customers who stopped buying)
- Customer engagement scoring
- Personalized offers
- Birthday/anniversary wishes
- Referral tracking

---

### 12. AI/ML Demand Forecasting
**Effort**: High (6-8 weeks including model training)

**Features**:
- Per-item demand prediction based on:
  - Historical sales patterns
  - Day of week trends
  - Seasonal patterns
  - Weather data integration
  - Festival/holiday calendar
  - Local events
- Purchase quantity recommendations
- Optimal stock levels
- Reorder point calculation
- Safety stock recommendations
- Slow-moving item detection

**Technical Stack**:
- Python backend service
- TensorFlow or PyTorch
- LSTM or Transformer models
- Historical data: minimum 1 year required
- External APIs: Weather, calendar

---

### 13. Route Optimization for Delivery
**Effort**: Medium-High (4 weeks)

**Features**:
- Delivery scheduling
- Multi-stop route planning
- Google Maps integration
- Real-time delivery tracking
- Delivery slot booking
- Customer delivery preferences
- Delivery cost calculation
- Driver assignment
- Proof of delivery (photo, signature)
- Failed delivery handling

---

### 14. Customer-Facing Marketplace
**Effort**: High (8+ weeks)

**Features**:
- Customer web portal
- Mobile app for customers
- Browse available products
- Online ordering
- Delivery slot selection
- Digital payments (UPI, cards, wallets)
- Order history
- Balance checking
- Dispute resolution
- Customer ratings and reviews

---

## 📊 Quick Wins (Low Effort, High Value - Build Immediately)

### 15. Immediate Enhancements (1 day each)

**A. Export to Excel**
- Add export button to all reports
- Generate XLSX files with formatting
- Include charts and summaries

**B. Backup & Restore**
- One-click database backup
- Download backup file
- Restore from backup file
- Automated daily backups to cloud storage

**C. Transaction Notes**
- Add notes/remarks field to transactions
- Search by notes
- Display notes in receipts (optional)

**D. Customer Photos Display**
- Show customer photo in transaction form
- Quick visual identification
- Photo upload/update on customer profile

**E. Quick Stats Widget**
- Today's sales count and amount
- Pending payments total
- Crates outstanding
- Items expiring today
- Refresh every 5 minutes

**F. Global Search**
- Search across customers, items, transactions
- Keyboard shortcut (Ctrl+K)
- Quick results dropdown
- Recent searches history

**G. Bulk Actions**
- Bulk price update
- Bulk SMS to customers
- Bulk delete transactions (admin only)
- Bulk export selected items

**H. Print Customization**
- Customize receipt header/footer
- Add business logo
- Configure paper size (thermal, A4)
- Multi-language support

---

## 💡 Recommended Implementation Priority

### **Phase 1: Foundation** (PARTIALLY COMPLETE)
**Focus**: Critical business needs

1. ⏳ **Payment Tracking & Credit Management** (#2) - IN PROGRESS
   - ⏳ Week 1: Credit limits, payment plans
   - ⏳ Week 2: Overdue alerts, reminders
   - ✅ Database schema ready (credit_limit, payment_terms_days, interest_rate fields exist)
   - ⏳ Backend enforcement logic needed
   - ⏳ Frontend UI for credit management needed
   
2. ✅ **Expiry & Wastage Tracking** (#3) - COMPLETED
   - ✅ Week 3: Wastage log, expiry alerts
   - ✅ Week 4: Reports and analytics
   - ✅ Full CRUD operations for wastage
   - ✅ Expiry alert system with acknowledgment
   - ✅ Photo upload for wastage evidence

**Completed Deliverables**:
- ✅ Wastage tracking dashboard
- ✅ Expiry alert notifications and acknowledgment
- ✅ Soft delete system with attestation
- ✅ Audit trail for all operations

**Pending Deliverables**:
- ⏳ Credit limit enforcement logic
- ⏳ Payment reminder system
- ⏳ Automated expiry checks (background job)

---

### **Phase 2: Communication (Weeks 5-7)**
**Focus**: Customer engagement

3. ✅ **WhatsApp Integration** (#7)
   - Week 5: Setup WhatsApp Business API
   - Week 6: Receipt and reminder messages
   - Week 7: Two-way communication

**Deliverables**:
- Auto-send receipts via WhatsApp
- Payment reminders
- Balance inquiry bot

---

### **Phase 3: Visibility (Weeks 8-11)**
**Focus**: Supplier and profitability

4. ✅ **Supplier Management** (#4)
   - Week 8-9: Supplier CRUD, purchase orders
   - Week 10: Goods receipt, payments
   - Week 11: Performance metrics

5. ✅ **Profit Tracking** (#5)
   - Week 11: Cost price, margin calculation
   - Week 12: Profit reports

**Deliverables**:
- Complete supplier management
- Profitability dashboard
- Purchase order workflow

---

### **Phase 4: Compliance (Weeks 12-15)**
**Focus**: Legal requirements

6. ✅ **GST Billing** (#8)
   - Week 12-13: GST invoice generation
   - Week 14: Tax reports
   - Week 15: E-invoice integration

**Deliverables**:
- GST compliant invoices
- GSTR reports
- E-invoice support

---

### **Phase 5: Scale (Weeks 16-20)**
**Focus**: Growth enablers

7. ✅ **Multi-Tenant Support** (#1)
   - Week 16-17: Database restructuring
   - Week 18: Business management
   - Week 19: Role-based access
   - Week 20: Testing and migration

8. ✅ **Mobile POS** (#6)
   - Week 18-20: Parallel development
   
**Deliverables**:
- Support for multiple businesses
- Mobile-optimized POS
- Staff management

---

### **Phase 6: Intelligence (Month 6+)**
**Focus**: Advanced features

9. Customer segmentation
10. AI forecasting
11. Route optimization
12. Customer marketplace

---

## 🎯 Success Metrics

### Business KPIs to Track
- **Revenue Growth**: Month-over-month sales increase
- **Bad Debt Reduction**: Decrease in write-offs
- **Wastage Reduction**: Percentage decrease in spoilage
- **Collection Time**: Average days to collect payment
- **Customer Retention**: Repeat customer rate
- **Inventory Turnover**: How fast inventory moves
- **Gross Margin**: Average profit margin %
- **Customer Satisfaction**: NPS score

### Technical Metrics
- **System Uptime**: 99.9% target
- **Response Time**: < 200ms for most queries
- **Mobile Performance**: < 3s page load
- **Data Accuracy**: 99.99% transaction accuracy
- **Backup Success**: 100% daily backups
- **Security**: Zero data breaches

---

## 📝 Notes

### Dependencies
- **External APIs**: WhatsApp, Weather, Maps, Payment gateways
- **Hardware**: Weighing scales, barcode scanners, thermal printers
- **Third-party Services**: SMS gateway, Email service, Cloud storage

### Assumptions
- Users have stable internet for core features
- Mobile devices support modern browsers
- Basic smartphone capabilities (camera for scanning)
- Users willing to adopt digital processes

### Risks
- **Data Migration**: Moving from manual to digital
- **User Adoption**: Training and change management
- **Integration Complexity**: Multiple external systems
- **Compliance**: GST and privacy regulations
- **Scalability**: Database performance at scale

---

## 🚀 Getting Started

To begin implementing these features, the recommended approach is:

1. **Validate with Users**: Get feedback on priority features
2. **Start Small**: Implement Quick Wins first for immediate value
3. **Follow the Phases**: Stick to the recommended sequence
4. **Iterate Quickly**: Release MVPs, gather feedback, improve
5. **Measure Impact**: Track metrics before and after each feature

### Next Immediate Actions:
1. ✅ Review and prioritize this roadmap
2. ✅ Set up project tracking (GitHub Issues/Projects)
3. ✅ Begin Phase 1: Payment & Credit Management
4. ✅ Collect user feedback on existing features
5. ✅ Plan infrastructure for scaling (if needed)

---

**Document Version**: 1.0  
**Author**: Development Team  
**Review Date**: Every Quarter  
**Last Review**: November 16, 2025
