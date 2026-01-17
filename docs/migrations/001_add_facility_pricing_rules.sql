-- =========================
-- FACILITY PRICING RULES
-- =========================
-- This table stores special pricing rules that apply when a member
-- already has a specific facility type rented.
-- Example: If member has a Box, then Canoe costs 80 EUR instead of base price

CREATE TABLE IF NOT EXISTS facility_pricing_rules (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id) ON DELETE CASCADE,
    required_facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id) ON DELETE CASCADE,
    special_price NUMERIC(10,2) NOT NULL CHECK (special_price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    -- Ensure we don't have duplicate rules for the same combination
    UNIQUE(facility_type_id, required_facility_type_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pricing_rules_facility_type
ON facility_pricing_rules(facility_type_id);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_required_facility_type
ON facility_pricing_rules(required_facility_type_id);

CREATE INDEX IF NOT EXISTS idx_pricing_rules_active
ON facility_pricing_rules(active);

-- Comments for documentation
COMMENT ON TABLE facility_pricing_rules IS 'Special pricing rules that apply when a member already has a specific facility type rented';
COMMENT ON COLUMN facility_pricing_rules.facility_type_id IS 'The facility type that gets the special price';
COMMENT ON COLUMN facility_pricing_rules.required_facility_type_id IS 'The facility type the member must already have to get the special price';
COMMENT ON COLUMN facility_pricing_rules.special_price IS 'The absolute special price to apply (e.g., 80.00 EUR)';
COMMENT ON COLUMN facility_pricing_rules.active IS 'Whether this pricing rule is currently active';
COMMENT ON COLUMN facility_pricing_rules.description IS 'Human-readable description of the rule (e.g., "Canoe discount for Box holders")';

-- Example data (can be inserted later)
-- INSERT INTO facility_pricing_rules (facility_type_id, required_facility_type_id, special_price, description)
-- VALUES
--   (2, 1, 80.00, 'Canoe special price for Box holders'),
--   (3, 1, 50.00, 'Locker special price for Box holders'),
--   (3, 4, 60.00, 'Locker special price for Mooring holders');
