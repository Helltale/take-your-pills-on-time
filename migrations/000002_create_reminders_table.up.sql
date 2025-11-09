-- Создание таблицы напоминаний
CREATE TABLE IF NOT EXISTS reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    image_url TEXT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('daily', 'weekly', 'custom', 'specific')),
    interval_hours INTEGER CHECK (interval_hours > 0),
    time_of_day VARCHAR(5) CHECK (time_of_day ~ '^([0-1][0-9]|2[0-3]):[0-5][0-9]$'),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_sent_at TIMESTAMP,
    next_send_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_reminders_user_id ON reminders(user_id);
CREATE INDEX idx_reminders_is_active ON reminders(is_active);
CREATE INDEX idx_reminders_next_send_at ON reminders(next_send_at) WHERE is_active = true;
CREATE INDEX idx_reminders_type ON reminders(type);

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_reminders_updated_at BEFORE UPDATE ON reminders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

