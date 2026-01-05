WITH membership_details AS (
    SELECT
        m.id AS membership_id,
        mp.valid_from,
        mp.expires_at,
        ms.status AS membership_status,
        mp.exclusion_deliberated_at,
        mp.excluded_at,
        p.amount AS payment_amount,
        p.currency AS payment_currency,
        p.paid_at AS payment_date,
        p.payment_method,
        p.transaction_ref
    FROM memberships m
    JOIN membership_periods mp ON m.id = mp.membership_id
    JOIN membership_statuses ms ON mp.status_id = ms.id
    LEFT JOIN payments p ON mp.id = p.membership_period_id
    WHERE m.member_id = $1
)
SELECT
    m.id AS member_id,
    m.first_name,
    m.last_name,
    m.date_of_birth,
    m.email,
    json_agg(DISTINCT jsonb_build_object(
        'prefix', pn.description,
        'number', pn.number
    )) AS phone_numbers,
    json_agg(DISTINCT jsonb_build_object(
        'country', a.country,
        'city', a.city,
        'street', a.street,
        'street_number', a.street_number
    )) AS addresses,
    json_agg(DISTINCT jsonb_build_object(
        'membership_id', md.membership_id,
        'valid_from', md.valid_from,
        'expires_at', md.expires_at,
        'status', md.membership_status,
        'exclusion_deliberated_at', md.exclusion_deliberated_at,
        'excluded_at', md.excluded_at,
        'payment', jsonb_build_object(
            'amount', md.payment_amount,
            'currency', md.payment_currency,
            'paid_at', md.payment_date,
            'payment_method', md.payment_method,
            'transaction_ref', md.transaction_ref
        )
    )) AS memberships
FROM members m
LEFT JOIN phone_numbers pn ON m.id = pn.member_id
LEFT JOIN addresses a ON m.id = a.member_id
LEFT JOIN membership_details md ON m.id = md.membership_id
WHERE m.id = $1
GROUP BY m.id;
