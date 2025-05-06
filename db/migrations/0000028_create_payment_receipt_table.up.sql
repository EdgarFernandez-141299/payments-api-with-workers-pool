CREATE TABLE IF NOT EXISTS payment_receipt (
    id VARCHAR(255) PRIMARY KEY,
    payment_receipt_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    enterprise_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    reference_order_id VARCHAR(255) NOT NULL,
    payment_id VARCHAR(255) NOT NULL,
    payment_status VARCHAR(255) NOT NULL,
    payment_amount DECIMAL(19, 2) NOT NULL,
    payment_country_code VARCHAR(10) NOT NULL,
    payment_currency_code VARCHAR(10) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    payment_date VARCHAR(255) NOT NULL,
    file_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);