package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

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
	db *sqlx.DB
}

func NewReminderExecutionRepository(db *sqlx.DB) ReminderExecutionRepository {
	return &reminderExecutionRepository{db: db}
}

func (r *reminderExecutionRepository) Create(ctx context.Context, execution *entities.ReminderExecution) error {
	query := `
		INSERT INTO reminder_executions (id, reminder_id, user_id, status, sent_at, confirmed_at, created_at)
		VALUES (:id, :reminder_id, :user_id, :status, :sent_at, :confirmed_at, :created_at)
	`

	execution.ID = uuid.New()
	execution.CreatedAt = time.Now()
	if execution.SentAt.IsZero() {
		execution.SentAt = time.Now()
	}

	_, err := r.db.NamedExecContext(ctx, query, execution)
	return err
}

func (r *reminderExecutionRepository) GetByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	var executions []*entities.ReminderExecution
	query := `SELECT * FROM reminder_executions WHERE reminder_id = $1 ORDER BY sent_at DESC LIMIT $2`

	err := r.db.SelectContext(ctx, &executions, query, reminderID, limit)
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (r *reminderExecutionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	var executions []*entities.ReminderExecution
	query := `SELECT * FROM reminder_executions WHERE user_id = $1 ORDER BY sent_at DESC LIMIT $2`

	err := r.db.SelectContext(ctx, &executions, query, userID, limit)
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (r *reminderExecutionRepository) GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*ExecutionStatistics, error) {
	var stats ExecutionStatistics
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'sent') as total_sent,
			COUNT(*) FILTER (WHERE status = 'confirmed') as total_confirmed,
			COUNT(*) FILTER (WHERE status = 'skipped') as total_skipped
		FROM reminder_executions
		WHERE user_id = $1 AND sent_at >= $2 AND sent_at <= $3
	`

	err := r.db.GetContext(ctx, &stats, query, userID, fromDate, toDate)
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
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'sent') as total_sent,
			COUNT(*) FILTER (WHERE status = 'confirmed') as total_confirmed,
			COUNT(*) FILTER (WHERE status = 'skipped') as total_skipped
		FROM reminder_executions
		WHERE reminder_id = $1 AND sent_at >= $2 AND sent_at <= $3
	`

	err := r.db.GetContext(ctx, &stats, query, reminderID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	if stats.TotalSent > 0 {
		stats.ConfirmationRate = float64(stats.TotalConfirmed) / float64(stats.TotalSent) * 100
	}

	return &stats, nil
}

func (r *reminderExecutionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ExecutionStatus) error {
	query := `UPDATE reminder_executions SET status = $1, confirmed_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
