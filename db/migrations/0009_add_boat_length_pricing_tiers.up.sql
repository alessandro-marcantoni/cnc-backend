-- =========================
-- BOAT LENGTH PRICING TIERS
-- =========================
-- This table stores pricing tiers based on boat length ranges
-- for facility types that require boats.
CREATE TABLE IF NOT EXISTS boat_length_pricing_tiers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_type_id BIGINT NOT NULL REFERENCES facilities_catalog(id) ON DELETE CASCADE,
    min_length_meters NUMERIC(10,2) NOT NULL CHECK (min_length_meters >= 0),
    max_length_meters NUMERIC(10,2) CHECK (max_length_meters IS NULL OR max_length_meters > min_length_meters),
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_length_range CHECK (
        max_length_meters IS NULL OR max_length_meters > min_length_meters
    )
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_boat_pricing_facility_type
    ON boat_length_pricing_tiers(facility_type_id);

CREATE INDEX IF NOT EXISTS idx_boat_pricing_active
    ON boat_length_pricing_tiers(active) WHERE active = TRUE;


INSERT INTO facilities_catalog (name, description, suggested_price)
VALUES
    ('Ormeggi Tavollo', 'Ormeggi Tavollo', 0.00),
    ('Ormeggi Ventena', 'Ormeggi Ventena', 0.00),
    ('Piazzale Ventena (Invernale)', 'Piazzale Ventena (Invernale)', 0.00);

INSERT INTO boat_length_pricing_tiers (facility_type_id, min_length_meters, max_length_meters, price)
VALUES
    (2, 0.00, 6.00, 170.00),
    (2, 6.00, 8.00, 220.00),
    (2, 8.00, 10.00, 270.00),
    (2, 10.00, NULL, 340.00);
