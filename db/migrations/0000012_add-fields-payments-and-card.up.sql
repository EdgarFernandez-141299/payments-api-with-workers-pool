ALTER TABLE card ADD COLUMN is_default BOOLEAN DEFAULT FALSE;

ALTER TABLE "payment" ADD COLUMN authorization_code VARCHAR(100);

ALTER TABLE "order" ADD COLUMN metadata json;