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
),
rented_services AS (
    SELECT
        rf.id AS rented_facility_id,
        f.identifier AS facility_identifier,
        fc.name AS facility_name,
        rf.rented_at,
        rf.expires_at,
        p.amount AS payment_amount,
        p.currency AS payment_currency,
        p.paid_at AS payment_date,
        p.payment_method,
        p.transaction_ref
    FROM rented_facilities rf
    JOIN facilities f ON rf.facility_id = f.id
    JOIN facilities_catalog fc ON f.facility_type_id = fc.id
    LEFT JOIN payments p ON rf.id = p.rented_facility_id
    WHERE rf.member_id = $1
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
    )) AS memberships,
    json_agg(DISTINCT jsonb_build_object(
        'rented_facility_id', rs.rented_facility_id,
        'facility_identifier', rs.facility_identifier,
        'facility_name', rs.facility_name,
        'rented_at', rs.rented_at,
        'expires_at', rs.expires_at,
        'payment', jsonb_build_object(
            'amount', rs.payment_amount,
            'currency', rs.payment_currency,
            'paid_at', rs.payment_date,
            'payment_method', rs.payment_method,
            'transaction_ref', rs.transaction_ref
        )
    )) AS rented_services
FROM members m
LEFT JOIN phone_numbers pn ON m.id = pn.member_id
LEFT JOIN addresses a ON m.id = a.member_id
LEFT JOIN membership_details md ON m.id = md.membership_id
LEFT JOIN rented_services rs ON m.id = rs.rented_facility_id
WHERE m.id = $1
GROUP BY m.id;
