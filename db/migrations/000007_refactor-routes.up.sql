DROP TABLE collection_center_payment_method;
DROP TABLE payment_method;
DROP TABLE payment_concept;

ALTER TABLE collection_center DROP COLUMN country_code;
ALTER TABLE collection_center ADD COLUMN available_currencies TEXT[];


ALTER TABLE collection_account DROP COLUMN country_code;
ALTER TABLE collection_account DROP COLUMN branch_code;


ALTER TABLE collection_account DROP COLUMN status;
ALTER TABLE collection_account ADD COLUMN status BOOLEAN DEFAULT TRUE;



ALTER TABLE collection_account_route ADD COLUMN status VARCHAR DEFAULT 'PENDING';
ALTER TABLE collection_account_route ADD COLUMN associated_origin VARCHAR;
ALTER TABLE collection_account_route DROP COLUMN payment_concept_code;


