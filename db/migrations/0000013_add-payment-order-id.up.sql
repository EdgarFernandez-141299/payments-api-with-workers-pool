ALTER TABLE payment ADD COLUMN payment_order_id VARCHAR(255);
--  add authorization_code column
ALTER TABLE payment RENAME COLUMN authorization_code TO order_token;