-- Migration: Add tax_code column to members table
-- Description: Adds the Italian tax code (codice fiscale) field to store member tax identification

ALTER TABLE members
ADD COLUMN tax_code VARCHAR(100);

-- Add comment to document the column
COMMENT ON COLUMN members.tax_code IS 'Italian tax code (Codice Fiscale) - alphanumeric code of 16 characters';
