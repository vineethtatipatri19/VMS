-- Add status column to transactions table
ALTER TABLE transactions 
ADD COLUMN status TEXT CHECK (status IN ('pending', 'completed', 'cancelled')) DEFAULT 'completed';

-- Add comment explaining the column
COMMENT ON COLUMN transactions.status IS 'Transaction processing status: pending, completed, cancelled';
