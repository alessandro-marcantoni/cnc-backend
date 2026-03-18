-- Insert a rental for a facility
INSERT INTO rented_facilities (facility_id, member_id, season_id, price, discount_applied)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
