-- Insert an address for a member
INSERT INTO addresses (member_id, country, city, street, street_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
