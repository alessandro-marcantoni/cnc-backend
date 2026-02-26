-- Make email optional and remove unique constraint
ALTER TABLE members ALTER COLUMN email DROP NOT NULL;
ALTER TABLE members DROP CONSTRAINT IF EXISTS members_email_key;

-- Add unique constraint that allows multiple NULLs
-- (NULL values are not considered equal in UNIQUE constraints, so multiple NULLs are allowed)
CREATE UNIQUE INDEX IF NOT EXISTS members_email_unique
ON members(email)
WHERE email IS NOT NULL;
