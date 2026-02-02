INSERT INTO payments (
    rented_facility_id,
    membership_period_id,
    amount,
    currency,
    paid_at,
    payment_method,
    notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
