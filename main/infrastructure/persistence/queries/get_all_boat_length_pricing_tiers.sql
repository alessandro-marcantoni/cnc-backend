-- Get all active boat length pricing tiers for all facility types
SELECT
    id,
    facility_type_id,
    min_length_meters,
    max_length_meters,
    price,
    currency,
    active,
    created_at,
    updated_at
FROM boat_length_pricing_tiers
WHERE active = TRUE
ORDER BY facility_type_id ASC, min_length_meters ASC;
