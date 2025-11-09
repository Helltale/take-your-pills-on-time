-- Создание таблицы выполнения напоминаний (статистика)
CREATE TABLE IF NOT EXISTS reminder_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reminder_id UUID NOT NULL REFERENCES reminders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('sent', 'confirmed', 'skipped')),
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска и аналитики
CREATE INDEX idx_executions_reminder_id ON reminder_executions(reminder_id);
CREATE INDEX idx_executions_user_id ON reminder_executions(user_id);
CREATE INDEX idx_executions_status ON reminder_executions(status);
CREATE INDEX idx_executions_sent_at ON reminder_executions(sent_at);
CREATE INDEX idx_executions_user_sent_at ON reminder_executions(user_id, sent_at);

-- Составной индекс для статистики по пользователю и напоминанию
CREATE INDEX idx_executions_user_reminder_status ON reminder_executions(user_id, reminder_id, status);

