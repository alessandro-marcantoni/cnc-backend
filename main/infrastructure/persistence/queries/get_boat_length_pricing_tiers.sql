-- Get all active boat length pricing tiers for a specific facility type
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
WHERE facility_type_id = $1
  AND active = TRUE
ORDER BY min_length_meters ASC;
