package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
)

type ReminderExecutionRepository interface {
	Create(ctx context.Context, execution *entities.ReminderExecution) error
	GetByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error)
	GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*ExecutionStatistics, error)
	GetStatisticsByReminderID(ctx context.Context, reminderID uuid.UUID, fromDate, toDate time.Time) (*ExecutionStatistics, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ExecutionStatus) error
}

type ExecutionStatistics struct {
	TotalSent        int     `json:"total_sent"`
	TotalConfirmed   int     `json:"total_confirmed"`
	TotalSkipped     int     `json:"total_skipped"`
	ConfirmationRate float64 `json:"confirmation_rate"`
}

type reminderExecutionRepository struct {
	db *gorm.DB
}

func NewReminderExecutionRepository(db *gorm.DB) ReminderExecutionRepository {
	return &reminderExecutionRepository{db: db}
}

func (r *reminderExecutionRepository) Create(ctx context.Context, execution *entities.ReminderExecution) error {
	execution.ID = uuid.New()
	execution.CreatedAt = time.Now()
	if execution.SentAt.IsZero() {
		execution.SentAt = time.Now()
	}

	return r.db.WithContext(ctx).Create(execution).Error
}

func (r *reminderExecutionRepository) GetByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	var executions []*entities.ReminderExecution
	query := r.db.WithContext(ctx).
		Where("reminder_id = ?", reminderID).
		Order("sent_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&executions).Error
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (r *reminderExecutionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	var executions []*entities.ReminderExecution
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("sent_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&executions).Error
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (r *reminderExecutionRepository) GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*ExecutionStatistics, error) {
	var stats ExecutionStatistics

	err := r.db.WithContext(ctx).
		Model(&entities.ReminderExecution{}).
		Select(`
			COUNT(*) FILTER (WHERE status = 'sent') as total_sent,
			COUNT(*) FILTER (WHERE status = 'confirmed') as total_confirmed,
			COUNT(*) FILTER (WHERE status = 'skipped') as total_skipped
		`).
		Where("user_id = ? AND sent_at >= ? AND sent_at <= ?", userID, fromDate, toDate).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	if stats.TotalSent > 0 {
		stats.ConfirmationRate = float64(stats.TotalConfirmed) / float64(stats.TotalSent) * 100
	}

	return &stats, nil
}

func (r *reminderExecutionRepository) GetStatisticsByReminderID(ctx context.Context, reminderID uuid.UUID, fromDate, toDate time.Time) (*ExecutionStatistics, error) {
	var stats ExecutionStatistics

	err := r.db.WithContext(ctx).
		Model(&entities.ReminderExecution{}).
		Select(`
			COUNT(*) FILTER (WHERE status = 'sent') as total_sent,
			COUNT(*) FILTER (WHERE status = 'confirmed') as total_confirmed,
			COUNT(*) FILTER (WHERE status = 'skipped') as total_skipped
		`).
		Where("reminder_id = ? AND sent_at >= ? AND sent_at <= ?", reminderID, fromDate, toDate).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	if stats.TotalSent > 0 {
		stats.ConfirmationRate = float64(stats.TotalConfirmed) / float64(stats.TotalSent) * 100
	}

	return &stats, nil
}

func (r *reminderExecutionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ExecutionStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == entities.ExecutionStatusConfirmed {
		now := time.Now()
		updates["confirmed_at"] = &now
	}

	return r.db.WithContext(ctx).Model(&entities.ReminderExecution{}).
		Where("id = ?", id).
		Updates(updates).Error
}
