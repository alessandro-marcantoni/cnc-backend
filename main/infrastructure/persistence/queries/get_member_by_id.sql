WITH membership_details AS (
    SELECT
        m.id AS membership_id,
        m.number AS membership_number,
        s.starts_at AS valid_from,
        s.ends_at AS expires_at,
        mp.exclusion_deliberated_at,
        mp.status_id,
        mp.price,
        p.id AS payment_id,
        p.amount AS payment_amount,
        p.currency AS payment_currency,
        p.paid_at AS payment_date,
        p.payment_method,
        p.transaction_ref
    FROM memberships m
    JOIN membership_periods mp ON m.id = mp.membership_id
    LEFT JOIN seasons s ON s.id = mp.season_id
    LEFT JOIN payments p ON mp.id = p.membership_period_id
    WHERE m.member_id = $1
    AND s.id = $2
)
SELECT
    m.id AS member_id,
    m.first_name,
    m.last_name,
    m.date_of_birth,
    m.email,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object(
            'prefix', pn.description,
            'number', pn.number
        )) FILTER (WHERE pn.id IS NOT NULL),
        '[]'::json
    ) AS phone_numbers,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object(
            'country', a.country,
            'city', a.city,
            'street', a.street,
            'street_number', a.street_number,
            'zip_code', a.zip_code
        )) FILTER (WHERE a.id IS NOT NULL),
        '[]'::json
    ) AS addresses,
    COALESCE(
        json_agg(DISTINCT jsonb_build_object(
            'membership_id', md.membership_id,
            'membership_number', md.membership_number,
            'valid_from', md.valid_from,
            'expires_at', md.expires_at,
            'exclusion_deliberated_at', md.exclusion_deliberated_at,
            'status', ms.status,
            'price', md.price,
            'payment', CASE
                WHEN md.payment_amount IS NOT NULL THEN
                    jsonb_build_object(
                        'id', md.payment_id,
                        'amount', md.payment_amount,
                        'currency', md.payment_currency,
                        'paid_at', md.payment_date,
                        'payment_method', md.payment_method,
                        'transaction_ref', md.transaction_ref
                    )
                ELSE NULL
            END
        )) FILTER (WHERE md.membership_id IS NOT NULL),
        '[]'::json
    ) AS memberships
FROM members m
LEFT JOIN phone_numbers pn ON m.id = pn.member_id
LEFT JOIN addresses a ON m.id = a.member_id
LEFT JOIN membership_details md ON m.id = md.membership_id
LEFT JOIN membership_statuses ms ON ms.id = md.status_id
WHERE m.id = $1
GROUP BY m.id;
