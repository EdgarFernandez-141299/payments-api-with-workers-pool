DO $$ 
DECLARE 
    r RECORD;
BEGIN --SELECT tablename FROM pg_tables WHERE schemaname = 'public'
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public' and tablename != 'schema_migrations') 
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;

CREATE TABLE "user" (
  "id" varchar(36) PRIMARY KEY,
  "user_type" varchar(50) NOT NULL,
  "email" varchar(100) NOT NULL,
  "address" varchar(100),
  "zip" varchar(10),
  "city" varchar(50),
  "state" varchar(50),
  "country_code" varchar(20),
  "phone" varchar(15),
  "enterprise_id" varchar(36) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "card" (
  "id" varchar(36) PRIMARY KEY,
  "external_card_id" varchar(36) UNIQUE,
  "user_id" varchar(36) NOT NULL,
  "card_holder" varchar(150),
  "alias" varchar(50),
  "bin" char(6),
  "last_four" char(4),
  "brand" varchar(25),
  "expiration_date" DATE,
  "card_type" varchar(30),
  "status" varchar(25) NOT NULL DEFAULT 'ACTIVE',
  "is_recurrent" bool DEFAULT false,
  "retry_attempts" int DEFAULT 0,
  "enterprise_id" varchar(36) NOT NULL,
  "card_failure_reason" varchar(100),
  "card_failure_code" varchar(20),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "collection_center" (
  "id" varchar(36) PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "country_code" varchar(4),
  "description" varchar(200),
  "available_currencies" varchar(100),
  "enterprise_id" varchar(36) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "collection_account_route" (
  "id" varchar(36) PRIMARY KEY,
  "collection_account_id" varchar(36),
  "associated_origin" varchar(50),
  "country_code" varchar(4),
  "currency_code" varchar(4),
  "status" varchar(25),
  "enterprise_id" varchar(36),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "disabled_at" timestamp
);

CREATE TABLE "collection_account" (
  "id" varchar(36) PRIMARY KEY,
  "collection_center_id" varchar(36) NOT NULL,
  "bank_name" varchar(50),
  "account_type" varchar(50),
  "account_number" varchar(25),
  "interbank_account_number" varchar(20),
  "country_code" varchar(4),
  "currency_code" varchar(4),
  "status" bool DEFAULT false,
  "enterprise_id" varchar(36) NOT NULL,
  "billing_address" json,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "payment" (
  "id" varchar(36) PRIMARY KEY,
  "order_id" varchar(36) NOT NULL,
  "associated_origin" varchar(50),
  "currency_code" varchar(4),
  "country_code" varchar(4),
  "card_id" varchar(36),
  "card_detail" json,
  "payment_method" varchar(50),
  "collection_account_id" varchar(36) NOT NULL,
  "metadata" json,
  "status" varchar(20),
  "total_amount" decimal,
  "reference" varchar(100),
  "failure_reason" varchar(100),
  "failure_code" varchar(20),
  "enterprise_id" varchar(36) NOT NULL,
  "ip_address" varchar(45),
  "device_fingerprint" varchar(100),
  "transaction_date" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "order" (
  "id" varchar(36) PRIMARY KEY,
  "user_id" varchar(36),
  "reference_order_id" varchar(36),
  "total_amount" decimal,
  "country_code" varchar(4),
  "currency_code" varchar(4),
  "status" varchar(20) NOT NULL DEFAULT 'pending',
  "expires_at" timestamp,
  "enterprise_id" varchar(36),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "refund" (
  "id" varchar(36) PRIMARY KEY,
  "payment_id" varchar(36) NOT NULL,
  "amount" decimal,
  "reason" varchar(200),
  "status" varchar(20),
  "enterprise_id" varchar(36) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

ALTER TABLE "refund" ADD FOREIGN KEY ("payment_id") REFERENCES "payment" ("id");

ALTER TABLE "collection_account" ADD FOREIGN KEY ("collection_center_id") REFERENCES "collection_center" ("id");

ALTER TABLE "collection_account_route" ADD FOREIGN KEY ("collection_account_id") REFERENCES "collection_account" ("id");

ALTER TABLE "payment" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");

ALTER TABLE "payment" ADD FOREIGN KEY ("collection_account_id") REFERENCES "collection_account" ("id");