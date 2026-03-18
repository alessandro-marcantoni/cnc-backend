-- Update the facility_id for an existing rental
-- This allows changing which specific facility is assigned to a rental
UPDATE rented_facilities
SET facility_id = $1
WHERE id = $2
  AND deleted_at IS NULL;
