-- Update insurance information for a boat
UPDATE insurances
SET provider = $2,
    number = $3,
    expires_at = $4
WHERE boat_id = $1
RETURNING id;
