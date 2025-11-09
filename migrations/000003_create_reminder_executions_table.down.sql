-- Удаление индексов
DROP INDEX IF EXISTS idx_executions_user_reminder_status;
DROP INDEX IF EXISTS idx_executions_user_sent_at;
DROP INDEX IF EXISTS idx_executions_sent_at;
DROP INDEX IF EXISTS idx_executions_status;
DROP INDEX IF EXISTS idx_executions_user_id;
DROP INDEX IF EXISTS idx_executions_reminder_id;

-- Удаление таблицы выполнения напоминаний
DROP TABLE IF EXISTS reminder_executions;

