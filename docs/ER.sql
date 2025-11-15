
-- PostgreSQL schema for PGVMS (core tables)
CREATE TABLE customers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  address TEXT,
  contact_number TEXT,
  photo_url TEXT,
  aadhaar_verified BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  variant TEXT,
  lot_number TEXT UNIQUE NOT NULL,
  quantity NUMERIC NOT NULL,
  unit TEXT CHECK (unit IN ('kg','lot')) NOT NULL,
  purchase_date DATE NOT NULL,
  expiry_date DATE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
  date TIMESTAMPTZ NOT NULL DEFAULT now(),
  type TEXT CHECK (type IN ('sale','payment')) NOT NULL,
  payment_amount NUMERIC,
  total_amount NUMERIC NOT NULL,
  details JSONB,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE sale_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
  inventory_lot_id UUID REFERENCES inventory_items(id),
  item_name TEXT,
  quantity NUMERIC,
  unit TEXT,
  price_per_unit NUMERIC,
  total NUMERIC
);

CREATE TABLE crate_ledger (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
  date TIMESTAMPTZ NOT NULL DEFAULT now(),
  crates_issued INTEGER DEFAULT 0,
  crates_returned INTEGER DEFAULT 0,
  balance INTEGER NOT NULL,
  notes TEXT
);
