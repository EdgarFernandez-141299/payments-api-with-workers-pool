-- Change order.id from varchar(36) to varchar(255)
ALTER TABLE "order" ALTER COLUMN "id" TYPE VARCHAR(255);

-- Change order.user_id from varchar(36) to varchar(255)
ALTER TABLE "order" ALTER COLUMN "user_id" TYPE VARCHAR(255);

-- Change order.reference_order_id from varchar(36) to varchar(255)
ALTER TABLE "order" ALTER COLUMN "reference_order_id" TYPE VARCHAR(255);