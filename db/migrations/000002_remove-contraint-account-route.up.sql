ALTER TABLE collection_account_route DROP COLUMN payment_concept_id;
ALTER TABLE collection_account_route ADD COLUMN payment_concept_code VARCHAR(50);
ALTER TABLE collection_account ADD COLUMN interbank_account_number VARCHAR(30);

