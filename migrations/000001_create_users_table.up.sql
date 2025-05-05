-- Создать таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    age INT NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- добавлен индекс для повышения производительности soft delete при использовании gorm.Model
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);