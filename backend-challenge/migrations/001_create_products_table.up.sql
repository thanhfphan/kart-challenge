CREATE TABLE IF NOT EXISTS products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert 5 example products
INSERT INTO products (sku, name, created_at, updated_at) VALUES
('SKU001', 'Product 1', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('SKU002', 'Product 2', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('SKU003', 'Product 3', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('SKU004', 'Product 4', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
('SKU005', 'Product 5', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
