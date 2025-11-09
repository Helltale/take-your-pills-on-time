package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder *entities.Reminder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Reminder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error)
	GetDueReminders(ctx context.Context) ([]*entities.Reminder, error)
	Update(ctx context.Context, reminder *entities.Reminder) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateNextSendAt(ctx context.Context, id uuid.UUID, nextSendAt time.Time) error
	UpdateLastSentAt(ctx context.Context, id uuid.UUID, lastSentAt time.Time) error
}

type reminderRepository struct {
	db *sqlx.DB
}

func NewReminderRepository(db *sqlx.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(ctx context.Context, reminder *entities.Reminder) error {
	query := `
		INSERT INTO reminders (id, user_id, title, comment, image_url, type, interval_hours, 
		                      time_of_day, is_active, last_sent_at, next_send_at, created_at, updated_at)
		VALUES (:id, :user_id, :title, :comment, :image_url, :type, :interval_hours, 
		        :time_of_day, :is_active, :last_sent_at, :next_send_at, :created_at, :updated_at)
	`

	now := time.Now()
	reminder.ID = uuid.New()
	reminder.CreatedAt = now
	reminder.UpdatedAt = now

	_, err := r.db.NamedExecContext(ctx, query, reminder)
	return err
}

func (r *reminderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Reminder, error) {
	var reminder entities.Reminder
	query := `SELECT * FROM reminders WHERE id = $1`

	err := r.db.GetContext(ctx, &reminder, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (r *reminderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	query := `SELECT * FROM reminders WHERE user_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &reminders, query, userID)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	query := `SELECT * FROM reminders WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &reminders, query, userID)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) GetDueReminders(ctx context.Context) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	now := time.Now()
	query := `
		SELECT r.* FROM reminders r
		INNER JOIN users u ON r.user_id = u.id
		WHERE r.is_active = true 
		  AND u.is_active = true
		  AND (r.next_send_at IS NULL OR r.next_send_at <= $1)
		ORDER BY r.next_send_at ASC NULLS LAST
	`

	err := r.db.SelectContext(ctx, &reminders, query, now)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) Update(ctx context.Context, reminder *entities.Reminder) error {
	query := `
		UPDATE reminders 
		SET title = :title, comment = :comment, image_url = :image_url, type = :type,
		    interval_hours = :interval_hours, time_of_day = :time_of_day, is_active = :is_active,
		    last_sent_at = :last_sent_at, next_send_at = :next_send_at, updated_at = :updated_at
		WHERE id = :id
	`

	reminder.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(ctx, query, reminder)
	return err
}

func (r *reminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reminders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *reminderRepository) UpdateNextSendAt(ctx context.Context, id uuid.UUID, nextSendAt time.Time) error {
	query := `UPDATE reminders SET next_send_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, nextSendAt, id)
	return err
}

func (r *reminderRepository) UpdateLastSentAt(ctx context.Context, id uuid.UUID, lastSentAt time.Time) error {
	query := `UPDATE reminders SET last_sent_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, lastSentAt, id)
	return err
}
