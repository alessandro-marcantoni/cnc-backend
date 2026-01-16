INSERT INTO members_waiting (member_id, facility_type_id, notes)
VALUES ($1, $2, $3)
RETURNING id, queued_at;
