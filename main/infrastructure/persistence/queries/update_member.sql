UPDATE members
SET
    first_name = $1,
    last_name = $2,
    date_of_birth = $3,
    email = $4,
    tax_code = $5
WHERE id = $6
RETURNING id;
