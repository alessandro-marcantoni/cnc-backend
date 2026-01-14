SELECT
    rf.id                 AS rented_facility_id,
    rf.rented_at,
    rf.expires_at,
    rf.price,

    f.id                  AS facility_id,
    f.identifier          AS facility_identifier,

    fc.id                 AS facility_type_id,
    fc.name               AS facility_type,
    fc.description        AS facility_type_description,
    fc.suggested_price,

    b.id                  AS boat_id,
    b.name                AS boat_name,
    b.length_meters,
    b.width_meters,

    p.id                  AS payment_id,
    p.amount              AS payment_amount,
    p.currency            AS payment_currency,
    p.paid_at             AS payment_paid_at,
    p.payment_method,
    p.transaction_ref
FROM rented_facilities rf
JOIN facilities f
    ON f.id = rf.facility_id
JOIN facilities_catalog fc
    ON fc.id = f.facility_type_id
LEFT JOIN boats b
    ON b.rented_facility_id = rf.id
LEFT JOIN payments p
    ON p.rented_facility_id = rf.id
LEFT JOIN seasons s
    ON s.id = rf.season_id
WHERE rf.member_id = $1
AND s.code = $2
ORDER BY rf.rented_at DESC;
