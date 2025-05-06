-- Revert user.id from varchar(255) to varchar(36)
ALTER TABLE "user" ALTER COLUMN "id" TYPE VARCHAR(36);