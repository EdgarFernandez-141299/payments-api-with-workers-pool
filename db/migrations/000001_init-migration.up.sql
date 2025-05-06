CREATE TABLE IF NOT EXISTS payment_concept (
    id VARCHAR(14) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    enterprise_id VARCHAR(14) NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payment_method(
    id VARCHAR(14) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    enterprise_id VARCHAR(14) NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_center(
    id VARCHAR(14) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    description VARCHAR(255),
    enterprise_id VARCHAR(14) NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_center_payment_method(
    collection_center_id VARCHAR(14),
    payment_method_id VARCHAR(14),
    CONSTRAINT fk_collection_center FOREIGN KEY (collection_center_id) REFERENCES collection_center(id),
    CONSTRAINT fk_payment_method FOREIGN KEY (payment_method_id) REFERENCES payment_method(id)
);

CREATE TABLE IF NOT EXISTS collection_account(
    id VARCHAR(14) PRIMARY KEY,
    collection_center_id VARCHAR(14) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    status VARCHAR(50),
    bank_name VARCHAR(255) NOT NULL,
    branch_code VARCHAR(50) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    enterprise_id VARCHAR(14) NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_collection_center FOREIGN KEY (collection_center_id) REFERENCES collection_center(id)
);

CREATE TABLE IF NOT EXISTS collection_account_route (
    id VARCHAR(14) PRIMARY KEY,
    collection_account_id VARCHAR(14) NOT NULL,
    payment_concept_id VARCHAR(50) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    enterprise_id VARCHAR(14) NOT NULL,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payment_concept FOREIGN KEY (payment_concept_id) REFERENCES payment_concept(id),
    CONSTRAINT fk_collection_account FOREIGN KEY (collection_account_id) REFERENCES collection_account(id)
)