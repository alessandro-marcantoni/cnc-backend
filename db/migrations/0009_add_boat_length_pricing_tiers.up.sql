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


INSERT INTO facilities_catalog (name, description, suggested_price, has_boat)
VALUES
    ('Ormeggi Tavollo', 'Ormeggi Tavollo', 0.00, TRUE),
    ('Ormeggi Piazzale Ventena (Estivo)', 'Ormeggi Piazzale Ventena (Estivo)', 0.00, TRUE),
    ('Ormeggi Piazzale Ventena (Invernale)', 'Ormeggi Piazzale Ventena (Invernale)', 0.00, TRUE),
    ('Box Kite', 'Box Kite', 150.00, FALSE),
    ('Rastrelliera Verticale Enel', 'Rastrelliera Verticale Enel', 40.00, FALSE);

INSERT INTO boat_length_pricing_tiers (facility_type_id, min_length_meters, max_length_meters, price)
VALUES
    (12, 0.00, 4.00, 98.00),
    (12, 4.00, 5.00, 175.00),
    (12, 5.00, 6.00, 242.00),
    (12, 6.00, 7.50, 291.00),
    (12, 7.50, NULL, 500.00),
    (11, 0.00, 4.00, 450.00),
    (11, 4.00, 5.00, 550.00),
    (11, 5.00, 6.00, 750.00),
    (11, 6.00, 7.50, 950.00),
    (11, 7.50, NULL, 1080.00),
    (10, 0.00, 6.00, 825.00),
    (10, 6.00, 7.00, 1067.00),
    (10, 7.00, 8.00, 1250.00),
    (10, 8.00, 9.00, 1500.00),
    (10, 9.00, 10.00, 1625.00),
    (10, 10.00, 11.00, 1750.00),
    (10, 11.00, 12.00, 1875.00),
    (10, 12.00, 13.00, 2000.00),
    (10, 13.00, NULL, 2125.00);

UPDATE facilities_catalog
SET has_leerboard = TRUE
WHERE name = 'Posti Barca Piazzale Derive';

UPDATE facilities_catalog
SET suggested_price = 75.00
WHERE name = 'Rastrelliera Tavole Aperta';

UPDATE facilities_catalog
SET suggested_price = 100.00
WHERE name = 'Deposito SUP Chiuso';
