-- Migration: Add soft delete support to rented_facilities
-- Date: 2024-01-XX
-- Description: Adds deleted_at column and updates unique constraint to allow re-renting

-- =========================
-- Step 1: Add deleted_at column
-- =========================
ALTER TABLE rented_facilities
ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

-- Add index for better query performance on non-deleted records
CREATE INDEX IF NOT EXISTS idx_rented_facilities_deleted_at
ON rented_facilities(deleted_at) WHERE deleted_at IS NULL;

-- =========================
-- Step 2: Update UNIQUE constraint
-- =========================
-- Drop existing constraint that prevents re-renting
ALTER TABLE rented_facilities
DROP CONSTRAINT IF EXISTS rented_facilities_facility_id_season_id_key;

-- Create partial unique index that only applies to non-deleted records
-- This allows the same facility to be rented again in the same season after soft delete
CREATE UNIQUE INDEX IF NOT EXISTS idx_rented_facilities_active_unique
ON rented_facilities(facility_id, season_id)
WHERE deleted_at IS NULL;

-- =========================
-- Step 3: Add comment for documentation
-- =========================
COMMENT ON COLUMN rented_facilities.deleted_at IS
'Timestamp when the rental was soft-deleted. NULL means the rental is active. Soft delete preserves payment history and audit trail.';
