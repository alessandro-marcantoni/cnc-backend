-- =========================
-- SAMPLE PRICING RULES DATA
-- =========================
-- This file contains example pricing rules that can be inserted after running the schema migration.
-- Adjust the facility_type_id values to match your actual facilities_catalog IDs.

-- First, let's check what facility types exist (this is just for reference, comment out when running)
-- SELECT id, name, description, suggested_price FROM facilities_catalog ORDER BY id;

-- Example pricing rules:
-- These assume the following facility type IDs (adjust based on your actual data):
-- 1 = Box (Posto Barca)
-- 2 = Canoe (Canoa)
-- 3 = Locker (Armadietto)
-- 4 = Mooring (Ormeggio)

-- Rule 1: If member has a Box, Canoe costs 80 EUR
INSERT INTO facility_pricing_rules (
    facility_type_id,
    required_facility_type_id,
    special_price,
    currency,
    description,
    active
)
VALUES (
    2,      -- Canoe gets the special price
    1,      -- When member has Box
    80.00,  -- Special price
    'EUR',
    'Canoe special price for Box holders',
    TRUE
)
ON CONFLICT (facility_type_id, required_facility_type_id) DO UPDATE
SET
    special_price = EXCLUDED.special_price,
    description = EXCLUDED.description,
    active = EXCLUDED.active,
    updated_at = now();

-- Rule 2: If member has a Box, Locker costs 50 EUR
INSERT INTO facility_pricing_rules (
    facility_type_id,
    required_facility_type_id,
    special_price,
    currency,
    description,
    active
)
VALUES (
    3,      -- Locker gets the special price
    1,      -- When member has Box
    50.00,  -- Special price
    'EUR',
    'Locker special price for Box holders',
    TRUE
)
ON CONFLICT (facility_type_id, required_facility_type_id) DO UPDATE
SET
    special_price = EXCLUDED.special_price,
    description = EXCLUDED.description,
    active = EXCLUDED.active,
    updated_at = now();

-- Rule 3: If member has a Mooring, Locker costs 60 EUR
INSERT INTO facility_pricing_rules (
    facility_type_id,
    required_facility_type_id,
    special_price,
    currency,
    description,
    active
)
VALUES (
    3,      -- Locker gets the special price
    4,      -- When member has Mooring
    60.00,  -- Special price
    'EUR',
    'Locker special price for Mooring holders',
    TRUE
)
ON CONFLICT (facility_type_id, required_facility_type_id) DO UPDATE
SET
    special_price = EXCLUDED.special_price,
    description = EXCLUDED.description,
    active = EXCLUDED.active,
    updated_at = now();

-- Add more pricing rules as needed...
-- Template for new rules:
-- INSERT INTO facility_pricing_rules (
--     facility_type_id,           -- The facility that gets the special price
--     required_facility_type_id,  -- The facility the member must already have
--     special_price,              -- The absolute special price
--     currency,
--     description,
--     active
-- )
-- VALUES (?, ?, ?, 'EUR', 'Description here', TRUE)
-- ON CONFLICT (facility_type_id, required_facility_type_id) DO UPDATE
-- SET special_price = EXCLUDED.special_price, updated_at = now();

-- Verify the inserted rules
-- SELECT
--     fpr.id,
--     ft1.name as facility_type,
--     ft2.name as required_facility_type,
--     fpr.special_price,
--     fpr.description,
--     fpr.active
-- FROM facility_pricing_rules fpr
-- JOIN facilities_catalog ft1 ON fpr.facility_type_id = ft1.id
-- JOIN facilities_catalog ft2 ON fpr.required_facility_type_id = ft2.id
-- ORDER BY ft1.name, ft2.name;
