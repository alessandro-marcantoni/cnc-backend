SELECT
    id,
    member_id,
    facility_type_id,
    queued_at,
    notes
FROM members_waiting
WHERE facility_type_id = $1
ORDER BY queued_at ASC
LIMIT 1;
