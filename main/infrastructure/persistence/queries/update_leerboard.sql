-- Update leerboard information for a rented facility
UPDATE leeboards
SET color = $2,
    type = $3,
    length_meters = $4
WHERE rented_facility_id = $1
RETURNING id;
