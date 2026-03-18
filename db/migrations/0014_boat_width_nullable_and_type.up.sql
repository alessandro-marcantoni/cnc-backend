-- Make boat width_meters nullable and add type field
ALTER TABLE boats
ALTER COLUMN width_meters DROP NOT NULL;

ALTER TABLE boats
ADD COLUMN type VARCHAR(100);

-- Add comment explaining the fields
COMMENT ON COLUMN boats.width_meters IS
'Width of the boat in meters. Can be null if not measured or not applicable.';

COMMENT ON COLUMN boats.type IS
'Type or category of the boat (e.g., Sailing, Motor, Inflatable, Kayak, etc.)';
