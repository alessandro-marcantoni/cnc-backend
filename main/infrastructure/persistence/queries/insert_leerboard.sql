-- Insert leerboard information for a rented facility
INSERT INTO leeboards (rented_facility_id, color, type, length_meters)
VALUES ($1, $2, $3, $4)
RETURNING id;
