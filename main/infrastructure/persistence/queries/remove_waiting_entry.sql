DELETE FROM members_waiting
WHERE id = $1
RETURNING id, member_id, facility_type_id, queued_at, notes;
