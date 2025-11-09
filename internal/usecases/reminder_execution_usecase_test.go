package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
)

type MockReminderExecutionRepository struct {
	mock.Mock
}

var _ repository.ReminderExecutionRepository = (*MockReminderExecutionRepository)(nil)

func (m *MockReminderExecutionRepository) Create(ctx context.Context, execution *entities.ReminderExecution) error {
	args := m.Called(ctx, execution)
	return args.Error(0)
}

func (m *MockReminderExecutionRepository) GetByReminderID(ctx context.Context, reminderID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	args := m.Called(ctx, reminderID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReminderExecution), args.Error(1)
}

func (m *MockReminderExecutionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.ReminderExecution, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReminderExecution), args.Error(1)
}

func (m *MockReminderExecutionRepository) GetStatisticsByUserID(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error) {
	args := m.Called(ctx, userID, fromDate, toDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.ExecutionStatistics), args.Error(1)
}

func (m *MockReminderExecutionRepository) GetStatisticsByReminderID(ctx context.Context, reminderID uuid.UUID, fromDate, toDate time.Time) (*repository.ExecutionStatistics, error) {
	args := m.Called(ctx, reminderID, fromDate, toDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.ExecutionStatistics), args.Error(1)
}

func (m *MockReminderExecutionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ExecutionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestReminderExecutionUsecase_RecordSent(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record sent", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		userID := uuid.New()

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.ReminderExecution")).Return(nil).Run(func(args mock.Arguments) {
			execution := args.Get(1).(*entities.ReminderExecution)
			execution.ID = uuid.New()
		})

		execution, err := usecase.RecordSent(ctx, reminderID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, execution)
		assert.Equal(t, reminderID, execution.ReminderID)
		assert.Equal(t, userID, execution.UserID)
		assert.Equal(t, entities.ExecutionStatusSent, execution.Status)
		assert.False(t, execution.SentAt.IsZero())
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		userID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.ReminderExecution")).Return(repoError)

		execution, err := usecase.RecordSent(ctx, reminderID, userID)

		assert.Error(t, err)
		assert.Nil(t, execution)
		assert.Contains(t, err.Error(), "failed to record sent execution")
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_RecordConfirmed(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record confirmed", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()

		mockRepo.On("UpdateStatus", ctx, executionID, entities.ExecutionStatusConfirmed).Return(nil)

		err := usecase.RecordConfirmed(ctx, executionID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.On("UpdateStatus", ctx, executionID, entities.ExecutionStatusConfirmed).Return(repoError)

		err := usecase.RecordConfirmed(ctx, executionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to record confirmed execution")
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_RecordSkipped(t *testing.T) {
	ctx := context.Background()

	t.Run("successful record skipped", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()

		mockRepo.On("UpdateStatus", ctx, executionID, entities.ExecutionStatusSkipped).Return(nil)

		err := usecase.RecordSkipped(ctx, executionID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		executionID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.On("UpdateStatus", ctx, executionID, entities.ExecutionStatusSkipped).Return(repoError)

		err := usecase.RecordSkipped(ctx, executionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to record skipped execution")
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_GetHistoryByReminderID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get with limit", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		limit := 10
		expectedExecutions := []*entities.ReminderExecution{
			{ID: uuid.New(), ReminderID: reminderID, Status: entities.ExecutionStatusSent},
			{ID: uuid.New(), ReminderID: reminderID, Status: entities.ExecutionStatusConfirmed},
		}

		mockRepo.On("GetByReminderID", ctx, reminderID, limit).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, limit)

		assert.NoError(t, err)
		assert.Len(t, executions, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("uses default limit when limit is 0", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		expectedExecutions := []*entities.ReminderExecution{}

		mockRepo.On("GetByReminderID", ctx, reminderID, 50).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, 0)

		assert.NoError(t, err)
		assert.NotNil(t, executions)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.On("GetByReminderID", ctx, reminderID, 50).Return(nil, repoError)

		executions, err := usecase.GetHistoryByReminderID(ctx, reminderID, 0)

		assert.Error(t, err)
		assert.Nil(t, executions)
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_GetHistoryByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		limit := 20
		expectedExecutions := []*entities.ReminderExecution{
			{ID: uuid.New(), UserID: userID, Status: entities.ExecutionStatusSent},
		}

		mockRepo.On("GetByUserID", ctx, userID, limit).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByUserID(ctx, userID, limit)

		assert.NoError(t, err)
		assert.Len(t, executions, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("uses default limit when limit is negative", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		expectedExecutions := []*entities.ReminderExecution{}

		mockRepo.On("GetByUserID", ctx, userID, 50).Return(expectedExecutions, nil)

		executions, err := usecase.GetHistoryByUserID(ctx, userID, -1)

		assert.NoError(t, err)
		assert.NotNil(t, executions)
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_GetStatisticsByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get statistics", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
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

		mockRepo.On("GetStatisticsByUserID", ctx, userID, fromDate, toDate).Return(expectedStats, nil)

		stats, err := usecase.GetStatisticsByUserID(ctx, userID, fromDate, toDate)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 10, stats.TotalSent)
		assert.Equal(t, 8, stats.TotalConfirmed)
		assert.Equal(t, 2, stats.TotalSkipped)
		assert.Equal(t, 80.0, stats.ConfirmationRate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		userID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -30)
		toDate := time.Now()
		repoError := errors.New("repository error")

		mockRepo.On("GetStatisticsByUserID", ctx, userID, fromDate, toDate).Return(nil, repoError)

		stats, err := usecase.GetStatisticsByUserID(ctx, userID, fromDate, toDate)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get statistics")
		mockRepo.AssertExpectations(t)
	})
}

func TestReminderExecutionUsecase_GetStatisticsByReminderID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get statistics", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
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

		mockRepo.On("GetStatisticsByReminderID", ctx, reminderID, fromDate, toDate).Return(expectedStats, nil)

		stats, err := usecase.GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 7, stats.TotalSent)
		assert.Equal(t, 5, stats.TotalConfirmed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockReminderExecutionRepository)
		usecase := NewReminderExecutionUsecase(mockRepo)

		reminderID := uuid.New()
		fromDate := time.Now().AddDate(0, 0, -7)
		toDate := time.Now()
		repoError := errors.New("repository error")

		mockRepo.On("GetStatisticsByReminderID", ctx, reminderID, fromDate, toDate).Return(nil, repoError)

		stats, err := usecase.GetStatisticsByReminderID(ctx, reminderID, fromDate, toDate)

		assert.Error(t, err)
		assert.Nil(t, stats)
		mockRepo.AssertExpectations(t)
	})
}
