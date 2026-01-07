SELECT
    m.id AS member_id,
    m.first_name,
    m.last_name,
    m.date_of_birth,
    mem.number AS membership_number,
    s.starts_at AS season_starts_at,
    s.ends_at AS season_ends_at,
    mp.exclusion_deliberated_at,
    p.amount AS amount_paid,
    p.paid_at AS paid_at,
    p.currency AS currency
FROM members m
LEFT JOIN memberships mem ON m.id = mem.member_id
LEFT JOIN membership_periods mp ON mem.id = mp.membership_id
LEFT JOIN payments p ON mp.id = p.membership_period_id
LEFT JOIN seasons s ON mp.season_id = s.id
WHERE s.code = $1
ORDER BY m.last_name, m.first_name
