-- Удаление триггера и функции
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_telegram_id;

-- Удаление таблицы пользователей
DROP TABLE IF EXISTS users;

