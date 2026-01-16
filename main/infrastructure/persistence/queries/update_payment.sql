UPDATE payments
SET
    amount = $1,
    currency = $2,
    payment_method = $3,
    transaction_ref = $4
WHERE id = $5;
