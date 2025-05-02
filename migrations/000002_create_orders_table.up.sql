-- Create the orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, -- Changed from user's specific definition to match potential gorm.Model
    updated_at TIMESTAMPTZ, -- Added for potential gorm.Model
    deleted_at TIMESTAMPTZ, -- Added for potential gorm.Model
    user_id INT NOT NULL, -- Make NOT NULL explicit
    product VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    -- Add foreign key constraint separately for clarity and control
    CONSTRAINT fk_orders_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE -- If an user is deleted, their orders are also deleted
);

-- Optional: Add index for soft delete performance if using gorm.Model for Order
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
-- Optional: Add index on user_id for faster lookup of user's orders
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);