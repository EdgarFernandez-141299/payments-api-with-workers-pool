ALTER TABLE user_payment DROP COLUMN token;

ALTER TABLE credit_card ADD COLUMN status VARCHAR(20);
ALTER TABLE credit_card ADD COLUMN billing_address_id VARCHAR(14);

CREATE TABLE billing_address (
  id VARCHAR(14) PRIMARY KEY,
  address VARCHAR(50) NOT NULL,
  zip VARCHAR(10) NOT NULL,
  city VARCHAR(25) NOT NULL,
  state VARCHAR(3) NOT NULL,
  country_code VARCHAR(3) NOT NULL,
  phone VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

ALTER TABLE credit_card 
ADD CONSTRAINT fk_billing_address 
FOREIGN KEY (billing_address_id) 
REFERENCES billing_address(id) 
ON DELETE CASCADE;