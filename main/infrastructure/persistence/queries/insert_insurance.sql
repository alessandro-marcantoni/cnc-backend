-- Insert insurance information for a boat
INSERT INTO insurances (boat_id, provider, number, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id;
