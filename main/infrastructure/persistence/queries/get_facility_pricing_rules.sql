-- Get all active facility pricing rules
SELECT
    id,
    facility_type_id,
    required_facility_type_id,
    special_price,
    currency,
    description,
    active,
    created_at,
    updated_at
FROM facility_pricing_rules
WHERE active = TRUE
ORDER BY facility_type_id, special_price ASC;
