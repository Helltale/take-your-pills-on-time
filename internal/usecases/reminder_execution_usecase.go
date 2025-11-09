package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
)

type ReminderExecutionUsecase interface {
	RecordSent(ctx context.Context, reminderID, userID uuid.UUID) (*entities.ReminderExecution, error)
	RecordConfirmed(ctx context.Context, executionID uuid.UUID) error
	RecordSkipped(ctx context.Context, executionID uuid.UUID) error
	GetHistoryByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error)
	GetHistoryByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error)
	GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error)
	GetStatisticsByReminderID(ctx context.Context, reminderID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error)
}

type reminderExecutionUsecase struct {
	repo repository.ReminderExecutionRepository
}

func NewReminderExecutionUsecase(repo repository.ReminderExecutionRepository) ReminderExecutionUsecase {
	return &reminderExecutionUsecase{repo: repo}
}

func (u *reminderExecutionUsecase) RecordSent(ctx context.Context, reminderID, userID uuid.UUID) (*entities.ReminderExecution, error) {
	execution := &entities.ReminderExecution{
		ReminderID: reminderID,
		UserID:     userID,
		Status:     entities.ExecutionStatusSent,
		SentAt:     time.Now(),
	}

	if err := u.repo.Create(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to record sent execution: %w", err)
	}

	return execution, nil
}

func (u *reminderExecutionUsecase) RecordConfirmed(ctx context.Context, executionID uuid.UUID) error {
	if err := u.repo.UpdateStatus(ctx, executionID, entities.ExecutionStatusConfirmed); err != nil {
		return fmt.Errorf("failed to record confirmed execution: %w", err)
	}
	return nil
}

func (u *reminderExecutionUsecase) RecordSkipped(ctx context.Context, executionID uuid.UUID) error {
	if err := u.repo.UpdateStatus(ctx, executionID, entities.ExecutionStatusSkipped); err != nil {
		return fmt.Errorf("failed to record skipped execution: %w", err)
	}
	return nil
}

func (u *reminderExecutionUsecase) GetHistoryByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	if limit <= 0 {
		limit = 50
	}
	executions, err := u.repo.GetByReminderID(ctx, reminderID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution history: %w", err)
	}
	return executions, nil
}

func (u *reminderExecutionUsecase) GetHistoryByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	if limit <= 0 {
		limit = 50
	}
	executions, err := u.repo.GetByUserID(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution history: %w", err)
	}
	return executions, nil
}

func (u *reminderExecutionUsecase) GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error) {
	stats, err := u.repo.GetStatisticsByUserID(ctx, userID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	return stats, nil
}

func (u *reminderExecutionUsecase) GetStatisticsByReminderID(ctx context.Context, reminderID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error) {
	stats, err := u.repo.GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}
	return stats, nil
}
