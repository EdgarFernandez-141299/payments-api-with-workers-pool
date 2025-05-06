-- First, revert foreign key column
-- Revert refund.payment_id from varchar(255) to varchar(36)
ALTER TABLE "refund" ALTER COLUMN "payment_id" TYPE VARCHAR(36);

-- Then, revert primary key column
-- Revert refund.id from varchar(255) to varchar(36)
ALTER TABLE "refund" ALTER COLUMN "id" TYPE VARCHAR(36);