SELECT
    m.id                AS member_id,
    m.first_name,
    m.last_name,
    m.date_of_birth,
    m.email,

    mem.id              AS membership_id,
    mem.number          AS membership_number,

    mp.id               AS membership_period_id,
    mp.valid_from,
    mp.expires_at,
    ms.status           AS membership_status,
    mp.exclusion_deliberated_at,
    mp.excluded_at,

    COALESCE(
        (
            SELECT json_agg(
                jsonb_build_object(
                    'prefix', pn.description,
                    'number', pn.number
                )
            )
            FROM phone_numbers pn
            WHERE pn.member_id = m.id
        ),
        '[]'::json
    ) AS phone_numbers,

    COALESCE(
        (
            SELECT json_agg(
                jsonb_build_object(
                    'country', a.country,
                    'city', a.city,
                    'street', a.street,
                    'street_number', a.street_number
                )
            )
            FROM addresses a
            WHERE a.member_id = m.id
        ),
        '[]'::json
    ) AS addresses
FROM members m
LEFT JOIN (
    SELECT
        mem.member_id,
        mem.id,
        mem.number,
        mp.id AS period_id,
        mp.valid_from,
        mp.expires_at,
        mp.status_id,
        mp.exclusion_deliberated_at,
        mp.excluded_at,
        ROW_NUMBER() OVER (
            PARTITION BY mem.member_id
            ORDER BY mp.valid_from DESC
        ) AS rn
    FROM memberships mem
    JOIN membership_periods mp
        ON mp.membership_id = mem.id
) latest
    ON latest.member_id = m.id
   AND latest.rn = 1
LEFT JOIN memberships mem
    ON mem.id = latest.id
LEFT JOIN membership_periods mp
    ON mp.id = latest.period_id
LEFT JOIN membership_statuses ms
    ON ms.id = mp.status_id
ORDER BY m.last_name, m.first_name;
