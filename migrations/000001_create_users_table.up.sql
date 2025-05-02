-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ, -- Added to match gorm.Model expectation
    updated_at TIMESTAMPTZ, -- Added to match gorm.Model expectation
    deleted_at TIMESTAMPTZ, -- Added to match gorm.Model expectation
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    age INT NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- Optional: Add index for soft delete performance if using gorm.Model
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);