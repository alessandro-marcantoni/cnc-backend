SELECT
    m.id                AS member_id,
    m.first_name,
    m.last_name,
    m.date_of_birth,
    mem.number          AS membership_number,
    CASE WHEN mp.id IS NOT NULL THEN
        CASE WHEN CURRENT_DATE >= s.starts_at AND CURRENT_DATE < s.ends_at THEN 'CURRENT'
         WHEN CURRENT_DATE < s.starts_at THEN 'FUTURE'
         WHEN CURRENT_DATE >= s.ends_at THEN 'PAST'
        END
    END AS season,
    s.starts_at AS season_starts_at,
    s.ends_at   AS season_ends_at,
    ms.status AS membership_status,
    mp.exclusion_deliberated_at,
    mp.price   AS price,
    p.amount   AS amount_paid,
    p.paid_at  AS paid_at,
    p.currency AS currency
FROM members m
LEFT JOIN (
    SELECT
        mem.member_id,
        mem.id,
        mem.number,
        mp.id AS period_id,
        mp.season_id,
        mp.exclusion_deliberated_at,
        mp.status_id,
        ROW_NUMBER() OVER (
            PARTITION BY mem.member_id
            ORDER BY mp.season_id DESC NULLS LAST
        ) AS rn
    FROM memberships mem
    LEFT JOIN membership_periods mp
        ON mp.membership_id = mem.id
) latest_membership_period
    ON latest_membership_period.member_id = m.id
   AND latest_membership_period.rn = 1
LEFT JOIN memberships mem
    ON mem.id = latest_membership_period.id
LEFT JOIN membership_periods mp
    ON mp.id = latest_membership_period.period_id
LEFT JOIN membership_statuses ms
    ON ms.id = mp.status_id
LEFT JOIN payments p
    ON p.membership_period_id = mp.id
LEFT JOIN seasons s
    ON s.id = mp.season_id
ORDER BY m.last_name, m.first_name
