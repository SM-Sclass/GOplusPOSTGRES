CREATE TABLE IF NOT EXISTS stocks (
    stockid BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price BIGINT NOT NULL,
    company VARCHAR(100) NOT NULL
);