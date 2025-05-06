-- Change payment.id from varchar(36) to varchar(255)
ALTER TABLE "payment" ALTER COLUMN "id" TYPE VARCHAR(255);

-- Change payment.order_id from varchar(36) to varchar(255)
ALTER TABLE "payment" ALTER COLUMN "order_id" TYPE VARCHAR(255);