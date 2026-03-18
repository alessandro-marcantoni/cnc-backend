-- Add discount_applied flag to track if a rental received a discount
ALTER TABLE rented_facilities
ADD COLUMN discount_applied BOOLEAN NOT NULL DEFAULT FALSE;

-- Add comment explaining the column
COMMENT ON COLUMN rented_facilities.discount_applied IS
'Indicates whether this rental received a discount from facility_pricing_rules. Used to enforce one-discount-per-season policy.';

-- Add index for efficient queries checking if member has already used discount in a season
CREATE INDEX idx_rented_facilities_member_season_discount
ON rented_facilities(member_id, season_id, discount_applied)
WHERE discount_applied = TRUE AND deleted_at IS NULL;
