-- Удаление триггера
DROP TRIGGER IF EXISTS update_reminders_updated_at ON reminders;

-- Удаление индексов
DROP INDEX IF EXISTS idx_reminders_type;
DROP INDEX IF EXISTS idx_reminders_next_send_at;
DROP INDEX IF EXISTS idx_reminders_is_active;
DROP INDEX IF EXISTS idx_reminders_user_id;

-- Удаление таблицы напоминаний
DROP TABLE IF EXISTS reminders;

