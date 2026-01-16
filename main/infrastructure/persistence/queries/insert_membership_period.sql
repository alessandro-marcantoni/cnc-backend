-- Insert a membership period for a membership
INSERT INTO membership_periods (membership_id, status_id, season_id, price)
VALUES ($1, $2, $3, $4)
RETURNING id;
