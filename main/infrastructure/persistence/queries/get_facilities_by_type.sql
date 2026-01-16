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
    m.id AS rented_by_member_id,
    m.first_name AS rented_by_member_first_name,
    m.last_name AS rented_by_member_last_name
FROM facilities f
INNER JOIN facilities_catalog fc ON f.facility_type_id = fc.id
LEFT JOIN rented_facilities rf ON f.id = rf.facility_id AND rf.season_id = $2
LEFT JOIN seasons s ON rf.season_id = s.id
LEFT JOIN members m ON rf.member_id = m.id
WHERE f.facility_type_id = $1
ORDER BY f.id
