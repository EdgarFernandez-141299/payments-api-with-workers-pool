-- First, revert foreign key columns
-- Revert order.reference_order_id from varchar(255) to varchar(36)
ALTER TABLE "order" ALTER COLUMN "reference_order_id" TYPE VARCHAR(36);

-- Revert order.user_id from varchar(255) to varchar(36)
ALTER TABLE "order" ALTER COLUMN "user_id" TYPE VARCHAR(36);

-- Then, revert primary key column
-- Revert order.id from varchar(255) to varchar(36)
ALTER TABLE "order" ALTER COLUMN "id" TYPE VARCHAR(36);