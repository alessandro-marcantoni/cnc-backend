-- Get a single facility by ID with its type information
SELECT
    f.id,
    f.facility_type_id,
    f.identifier,
    fc.name AS facility_type_name,
    fc.description AS facility_type_description,
    fc.suggested_price,
    CASE
        WHEN rf.id IS NOT NULL THEN TRUE
        ELSE FALSE
    END AS is_rented,
    s.ends_at AS expires_at,
    rf.member_id AS rented_by_member_id,
    m.first_name AS rented_by_member_first_name,
    m.last_name AS rented_by_member_last_name
FROM facilities f
JOIN facilities_catalog fc
    ON fc.id = f.facility_type_id
LEFT JOIN rented_facilities rf
    ON rf.facility_id = f.id
    AND rf.deleted_at IS NULL
LEFT JOIN seasons s
    ON s.id = rf.season_id
LEFT JOIN members m
    ON m.id = rf.member_id
WHERE f.id = $1
LIMIT 1;
