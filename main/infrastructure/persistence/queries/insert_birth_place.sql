INSERT INTO birth_places (member_id, country, city, zip_code, street, street_number)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
