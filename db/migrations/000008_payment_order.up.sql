CREATE TABLE IF NOT EXISTS payment_order (
    id VARCHAR(14) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL,
    origin_category VARCHAR(50) NOT NULL,
    collection_account_id VARCHAR(14) NOT NULL,
    reference VARCHAR(255),
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_collection_account FOREIGN KEY (collection_account_id) REFERENCES collection_account(id)
);

CREATE TABLE IF NOT EXISTS payment_order_item (
    id VARCHAR(14) PRIMARY KEY,
    payment_order_id VARCHAR(14) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL,
    reference VARCHAR(255),
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_payment_order FOREIGN KEY (payment_order_id) REFERENCES payment_order(id)
);