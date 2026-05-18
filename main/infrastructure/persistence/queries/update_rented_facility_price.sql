UPDATE rented_facilities
SET price = $2
WHERE id = $1 AND deleted_at IS NULL;
