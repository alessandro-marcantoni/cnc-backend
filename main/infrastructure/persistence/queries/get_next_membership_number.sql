-- Get the next membership number (max + 1)
SELECT COALESCE(MAX(number), 0) + 1 AS next_number
FROM memberships;
