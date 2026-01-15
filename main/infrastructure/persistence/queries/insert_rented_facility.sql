-- Insert a rental for a facility
INSERT INTO rented_facilities (facility_id, member_id, rented_at, expires_at, season_id, price)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
