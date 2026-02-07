SELECT
    m.id AS member_id,
    m.first_name,
    m.last_name,
    m.email,
    m.date_of_birth,
    m.tax_code,
    mem.number AS membership_number,
    s.starts_at AS season_starts_at,
    s.ends_at AS season_ends_at,
    mp.exclusion_deliberated_at,
    mp.price AS price,
    ms.status AS membership_status,
    p.amount AS amount_paid,
    p.paid_at AS paid_at,
    p.currency AS currency,
    CASE
        WHEN EXISTS (
            SELECT 1
            FROM rented_facilities rf
            LEFT JOIN payments fp ON fp.rented_facility_id = rf.id
            WHERE rf.member_id = m.id
            AND rf.season_id = s.id
            AND fp.id IS NULL
        ) THEN true
        ELSE false
    END AS has_unpaid_facilities
FROM members m
LEFT JOIN memberships mem ON m.id = mem.member_id
LEFT JOIN membership_periods mp ON mem.id = mp.membership_id
LEFT JOIN membership_statuses ms ON mp.status_id = ms.id
LEFT JOIN payments p ON mp.id = p.membership_period_id
LEFT JOIN seasons s ON mp.season_id = s.id
WHERE s.id = $1
ORDER BY m.last_name, m.first_name
