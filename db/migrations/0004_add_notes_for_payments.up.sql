ALTER TABLE payments
DROP CONSTRAINT IF EXISTS payments_transaction_ref_key;

ALTER TABLE payments
RENAME COLUMN transaction_ref TO notes;
