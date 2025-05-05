-- Создать таблицу заказов
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INT NOT NULL,
    product VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    -- Добавить ограничение внешнего ключа отдельно для ясности и контроля
    CONSTRAINT fk_orders_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE -- Если пользователь удалён, его заказы также удаляются
);

-- добавлен индекс для повышения производительности soft delete при использовании gorm.Model для Order
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
-- добавлен индекс по user_id для ускорения поиска заказов пользователя
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);