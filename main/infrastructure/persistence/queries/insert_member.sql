-- Insert a new member and return the generated ID
INSERT INTO members (first_name, last_name, date_of_birth, email)
VALUES ($1, $2, $3, $4)
RETURNING id;
