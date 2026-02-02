UPDATE payments
SET
    amount = $1,
    currency = $2,
    payment_method = $3,
    notes = $4
WHERE id = $5;
