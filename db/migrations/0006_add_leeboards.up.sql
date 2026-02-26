CREATE TABLE IF NOT EXISTS leeboards (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    rented_facility_id BIGINT NOT NULL UNIQUE
        REFERENCES rented_facilities(id) ON DELETE CASCADE,
    color VARCHAR(255),
    type VARCHAR(255),
    length_meters NUMERIC(10,2) NOT NULL CHECK (length_meters > 0)
);

ALTER TABLE facilities_catalog ADD COLUMN has_leerboard BOOLEAN NOT NULL DEFAULT FALSE;
