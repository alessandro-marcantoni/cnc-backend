-- Insert boat information for a rented facility
INSERT INTO boats (rented_facility_id, name, length_meters, width_meters)
VALUES ($1, $2, $3, $4)
RETURNING id;
