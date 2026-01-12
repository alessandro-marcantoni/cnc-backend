-- Insert a phone number for a member
INSERT INTO phone_numbers (member_id, number, description)
VALUES ($1, $2, $3)
RETURNING id;
