DROP TABLE IF EXISTS deuna_payment;

CREATE TABLE IF NOT EXISTS deuna_payment (
    payment_id VARCHAR(255) NOT NULL,
    order_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (payment_id)
);