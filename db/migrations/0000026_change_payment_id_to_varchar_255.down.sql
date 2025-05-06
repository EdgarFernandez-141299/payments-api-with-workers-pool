-- First, revert foreign key column
-- Revert payment.order_id from varchar(255) to varchar(36)
ALTER TABLE "payment" ALTER COLUMN "order_id" TYPE VARCHAR(36);

-- Then, revert primary key column
-- Revert payment.id from varchar(255) to varchar(36)
ALTER TABLE "payment" ALTER COLUMN "id" TYPE VARCHAR(36);