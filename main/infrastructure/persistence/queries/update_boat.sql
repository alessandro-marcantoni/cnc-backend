-- Update boat information for a rented facility
UPDATE boats
SET name = $2,
    length_meters = $3,
    width_meters = $4,
    engine_info = $5,
    type = $6
WHERE rented_facility_id = $1
RETURNING id;
