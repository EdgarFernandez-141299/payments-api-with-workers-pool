-- Change refund.id from varchar(36) to varchar(255)
ALTER TABLE "refund" ALTER COLUMN "id" TYPE VARCHAR(255);

-- Change refund.payment_id from varchar(36) to varchar(255)
ALTER TABLE "refund" ALTER COLUMN "payment_id" TYPE VARCHAR(255);