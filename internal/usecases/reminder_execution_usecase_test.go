package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
	"github.com/Helltale/take-your-pills-on-time/internal/repository/mocks"
)

func TestReminderExecutionUsecase_RecordSent(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		userID := uuid.New()

		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, execution *entities.ReminderExecution) error {
			execution.ID = uuid.New()
			return nil
		})

		execution, err := usecase.RecordSent(ctx, reminderID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, execution)
		assert.Equal(t, reminderID, execution.ReminderID)
		assert.Equal(t, userID, execution.UserID)
		assert.Equal(t, entities.ExecutionStatusSent, execution.Status)
		assert.False(t, execution.SentAt.IsZero())
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		userID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoError)

		execution, err := usecase.RecordSent(ctx, reminderID, userID)

		assert.Error(t, err)
		assert.Nil(t, execution)
		assert.Contains(t, err.Error(), "failed to record sent execution")
	})
}

func TestReminderExecutionUsecase_RecordConfirmed(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record confirmed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()

		mockRepo.EXPECT().UpdateStatus(ctx, executionID, entities.ExecutionStatusConfirmed).Return(nil)

		err := usecase.RecordConfirmed(ctx, executionID)

		assert.NoError(t, err)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().UpdateStatus(ctx, executionID, entities.ExecutionStatusConfirmed).Return(repoError)

		err := usecase.RecordConfirmed(ctx, executionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to record confirmed execution")
	})
}

func TestReminderExecutionUsecase_RecordSkipped(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()

		mockRepo.EXPECT().UpdateStatus(ctx, executionID, entities.ExecutionStatusSkipped).Return(nil)

		err := usecase.RecordSkipped(ctx, executionID)

		assert.NoError(t, err)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().UpdateStatus(ctx, executionID, entities.ExecutionStatusSkipped).Return(repoError)

		err := usecase.RecordSkipped(ctx, executionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to record skipped execution")
	})
}

func TestReminderExecutionUsecase_GetHistoryByReminderID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get with limit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		limit := 10
		expectedExecutions := []*entities.ReminderExecution{
			{ID: uuid.New(), ReminderID: reminderID, Status: entities.ExecutionStatusSent},
			{ID: uuid.New(), ReminderID: reminderID, Status: entities.ExecutionStatusConfirmed},
		}

		mockRepo.EXPECT().GetByReminderID(ctx, reminderID, limit).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, limit)

		assert.NoError(t, err)
		assert.Len(t, executions, 2)
	})

	t.Run("uses default limit when limit is 0", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		expectedExecutions := []*entities.ReminderExecution{}

		mockRepo.EXPECT().GetByReminderID(ctx, reminderID, 50).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, 0)

		assert.NoError(t, err)
		assert.NotNil(t, executions)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetByReminderID(ctx, reminderID, 50).Return(nil, repoError)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, 0)

		assert.Error(t, err)
		assert.Nil(t, executions)
	})
}

func TestReminderExecutionUsecase_GetHistoryByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		limit := 20
		expectedExecutions := []*entities.ReminderExecution{
			{ID: uuid.New(), UserID: userID, Status: entities.ExecutionStatusSent},
		}

		mockRepo.EXPECT().GetByUserID(ctx, userID, limit).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByUserID(ctx, userID, limit)

		assert.NoError(t, err)
		assert.Len(t, executions, 1)
	})

	t.Run("uses default limit when limit is negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		expectedExecutions := []*entities.ReminderExecution{}

		mockRepo.EXPECT().GetByUserID(ctx, userID, 50).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByUserID(ctx, userID, -1)

		assert.NoError(t, err)
		assert.NotNil(t, executions)
	})
}

func TestReminderExecutionUsecase_GetStatisticsByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get statistics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -30)
		toDate := time.Now()

		expectedStats := &repository.ExecutionStatistics{
			TotalSent:        10,
			TotalConfirmed:   8,
			TotalSkipped:     2,
			ConfirmationRate: 80.0,
		}

		mockRepo.EXPECT().GetStatisticsByUserID(ctx, userID, fromDate, toDate).Return(expectedStats, nil)

		stats, err := usecase.GetStatisticsByUserID(ctx, userID, fromDate, toDate)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 10, stats.TotalSent)
		assert.Equal(t, 8, stats.TotalConfirmed)
		assert.Equal(t, 2, stats.TotalSkipped)
		assert.Equal(t, 80.0, stats.ConfirmationRate)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -30)
		toDate := time.Now()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetStatisticsByUserID(ctx, userID, fromDate, toDate).Return(nil, repoError)

		stats, err := usecase.GetStatisticsByUserID(ctx, userID, fromDate, toDate)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get statistics")
	})
}

func TestReminderExecutionUsecase_GetStatisticsByReminderID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get statistics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -7)
		toDate := time.Now()

		expectedStats := &repository.ExecutionStatistics{
			TotalSent:        7,
			TotalConfirmed:   5,
			TotalSkipped:     2,
			ConfirmationRate: 71.43,
		}

		mockRepo.EXPECT().GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate).Return(expectedStats, nil)

		stats, err := usecase.GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 7, stats.TotalSent)
		assert.Equal(t, 5, stats.TotalConfirmed)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderExecutionRepository(ctrl)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -7)
		toDate := time.Now()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate).Return(nil, repoError)

		stats, err := usecase.GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate)

		assert.Error(t, err)
		assert.Nil(t, stats)
	})
}
