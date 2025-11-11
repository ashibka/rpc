CREATE TABLE orders (
                        id VARCHAR(36) PRIMARY KEY,
                        item VARCHAR(255) NOT NULL,
                        quantity INTEGER NOT NULL,
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_created_at ON orders(created_at);