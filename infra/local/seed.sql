
-- Seed sample data for local development
INSERT INTO customers (id, name, address, contact_number, photo_url, aadhaar_verified, created_at)
VALUES (gen_random_uuid(), 'Local Customer', '100 Local St', '9999999999', '', false, now())
ON CONFLICT DO NOTHING;
