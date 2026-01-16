SELECT
    id,
    facility_type_id,
    queued_at,
    notes
FROM members_waiting
WHERE member_id = $1 AND facility_type_id = $2;
