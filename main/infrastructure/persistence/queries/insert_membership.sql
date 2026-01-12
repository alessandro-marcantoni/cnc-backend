-- Insert a new membership for a member
INSERT INTO memberships (member_id, number, created_at)
VALUES ($1, $2, now())
RETURNING id;
