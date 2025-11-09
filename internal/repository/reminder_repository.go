package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

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
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(ctx context.Context, reminder *entities.Reminder) error {
	now := time.Now()
	reminder.ID = uuid.New()
	reminder.CreatedAt = now
	reminder.UpdatedAt = now

	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *reminderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Reminder, error) {
	var reminder entities.Reminder
	err := r.db.WithContext(ctx).First(&reminder, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (r *reminderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) GetDueReminders(ctx context.Context) ([]*entities.Reminder, error) {
	var reminders []*entities.Reminder
	now := time.Now()

	err := r.db.WithContext(ctx).
		Joins("INNER JOIN users ON reminders.user_id = users.id").
		Where("reminders.is_active = ? AND users.is_active = ? AND (reminders.next_send_at IS NULL OR reminders.next_send_at <= ?)", true, true, now).
		Order("reminders.next_send_at ASC NULLS LAST").
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *reminderRepository) Update(ctx context.Context, reminder *entities.Reminder) error {
	reminder.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(reminder).Error
}

func (r *reminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Reminder{}, "id = ?", id).Error
}

func (r *reminderRepository) UpdateNextSendAt(ctx context.Context, id uuid.UUID, nextSendAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"next_send_at": nextSendAt,
			"updated_at":   time.Now(),
		}).Error
}

func (r *reminderRepository) UpdateLastSentAt(ctx context.Context, id uuid.UUID, lastSentAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_sent_at": lastSentAt,
			"updated_at":   time.Now(),
		}).Error
}
