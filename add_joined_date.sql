-- Add joined_date column to users table
ALTER TABLE users ADD COLUMN joined_date DATE DEFAULT NULL AFTER created_at;

-- Update existing users to have their created_at date as joined_date
UPDATE users SET joined_date = DATE(created_at) WHERE joined_date IS NULL;

-- Show updated structure
DESCRIBE users;
