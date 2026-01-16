DELETE FROM members_waiting
WHERE member_id = $1 AND facility_type_id = $2
RETURNING id, queued_at, notes;
