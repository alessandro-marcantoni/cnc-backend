-- Insert a membership period for a membership
INSERT INTO membership_periods (membership_id, valid_from, expires_at, status_id, season_id, price)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
