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
    rf.expires_at
FROM facilities f
INNER JOIN facilities_catalog fc ON f.facility_type_id = fc.id
LEFT JOIN rented_facilities rf ON f.id = rf.facility_id
    AND rf.expires_at > NOW() AND rf.rented_at <= NOW()
WHERE f.facility_type_id = $1
ORDER BY f.identifier
