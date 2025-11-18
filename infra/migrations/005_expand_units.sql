-- Migration to expand allowed inventory units
-- This allows common units of measurement for dairy products and general retail

-- Drop the old constraint
ALTER TABLE inventory_items DROP CONSTRAINT IF EXISTS inventory_items_unit_check;

-- Add new constraint with expanded units
ALTER TABLE inventory_items ADD CONSTRAINT inventory_items_unit_check 
  CHECK (unit IN (
    'kg',          -- Kilograms
    'g',           -- Grams  
    'lot',         -- Lot/Batch
    'pieces',      -- Individual pieces
    'bottles',     -- Bottles
    'cartons',     -- Cartons
    'liters',      -- Liters
    'ml',          -- Milliliters
    'boxes',       -- Boxes
    'packets',     -- Packets
    'cans',        -- Cans
    'jars',        -- Jars
    'bags',        -- Bags
    'units',       -- Generic units
    'dozen',       -- Dozen
    'pairs'        -- Pairs
  ));
