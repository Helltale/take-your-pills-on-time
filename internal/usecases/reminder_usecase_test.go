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
	"github.com/Helltale/take-your-pills-on-time/internal/repository/mocks"
)

func TestReminderUsecase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("successful creation daily", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test Reminder"
		reminderType := entities.ReminderTypeDaily

		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			reminder.ID = uuid.New()
			return nil
		})

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, reminder)
		assert.Equal(t, title, reminder.Title)
		assert.Equal(t, reminderType, reminder.Type)
		assert.True(t, reminder.IsActive)
		assert.NotNil(t, reminder.NextSendAt)
	})

	t.Run("successful creation with comment and image", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Medicine"
		comment := "Take after meal"
		imageURL := "https://example.com/image.jpg"
		reminderType := entities.ReminderTypeDaily

		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			reminder.ID = uuid.New()
			return nil
		})

		reminder, err := usecase.Create(ctx, userID, title, &comment, &imageURL, reminderType, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, reminder)
		assert.Equal(t, comment, *reminder.Comment)
		assert.Equal(t, imageURL, *reminder.ImageURL)
	})

	t.Run("error when title is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		reminderType := entities.ReminderTypeDaily

		reminder, err := usecase.Create(ctx, userID, "", nil, nil, reminderType, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("error when custom type without interval_hours", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeCustom

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "interval_hours is required")
	})

	t.Run("error when custom type with invalid interval_hours", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeCustom
		invalidInterval := 0

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, &invalidInterval, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "interval_hours")
	})

	t.Run("successful creation custom type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeCustom
		intervalHours := 6

		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			reminder.ID = uuid.New()
			return nil
		})

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, &intervalHours, nil)

		assert.NoError(t, err)
		assert.NotNil(t, reminder)
		assert.Equal(t, intervalHours, *reminder.IntervalHours)
	})

	t.Run("error when specific type without time_of_day", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeSpecific

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "time_of_day is required")
	})

	t.Run("error when specific type with invalid time format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeSpecific
		invalidTime := "25:00"

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, &invalidTime)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "invalid time_of_day format")
	})

	t.Run("successful creation specific type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeSpecific
		timeOfDay := "09:00"

		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			reminder.ID = uuid.New()
			return nil
		})

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, &timeOfDay)

		assert.NoError(t, err)
		assert.NotNil(t, reminder)
		assert.Equal(t, timeOfDay, *reminder.TimeOfDay)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		title := "Test"
		reminderType := entities.ReminderTypeDaily
		repoError := errors.New("repository error")

		mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(repoError)

		reminder, err := usecase.Create(ctx, userID, title, nil, nil, reminderType, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "failed to create reminder")
	})
}

func TestReminderUsecase_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()
		expectedReminder := &entities.Reminder{
			ID:     reminderID,
			Title:  "Test",
			Type:   entities.ReminderTypeDaily,
			UserID: uuid.New(),
		}

		mockRepo.EXPECT().GetByID(ctx, reminderID).Return(expectedReminder, nil)

		reminder, err := usecase.GetByID(ctx, reminderID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReminder, reminder)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetByID(ctx, reminderID).Return(nil, repoError)

		reminder, err := usecase.GetByID(ctx, reminderID)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "failed to get reminder")
	})
}

func TestReminderUsecase_GetByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		expectedReminders := []*entities.Reminder{
			{ID: uuid.New(), Title: "Reminder 1", UserID: userID},
			{ID: uuid.New(), Title: "Reminder 2", UserID: userID},
		}

		mockRepo.EXPECT().GetByUserID(ctx, userID).Return(expectedReminders, nil)

		reminders, err := usecase.GetByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, reminders, 2)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetByUserID(ctx, userID).Return(nil, repoError)

		reminders, err := usecase.GetByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, reminders)
	})
}

func TestReminderUsecase_GetActiveByUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		userID := uuid.New()
		expectedReminders := []*entities.Reminder{
			{ID: uuid.New(), Title: "Active Reminder", UserID: userID, IsActive: true},
		}

		mockRepo.EXPECT().GetActiveByUserID(ctx, userID).Return(expectedReminders, nil)

		reminders, err := usecase.GetActiveByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, reminders, 1)
	})
}

func TestReminderUsecase_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("successful update title", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()
		existingReminder := &entities.Reminder{
			ID:     reminderID,
			Title:  "Old Title",
			Type:   entities.ReminderTypeDaily,
			UserID: uuid.New(),
		}
		newTitle := "New Title"

		mockRepo.EXPECT().GetByID(ctx, reminderID).Return(existingReminder, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			assert.Equal(t, newTitle, reminder.Title)
			return nil
		})

		reminder, err := usecase.Update(ctx, reminderID, &newTitle, nil, nil, nil, nil, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, reminder)
		assert.Equal(t, newTitle, reminder.Title)
	})

	t.Run("error when reminder not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()

		mockRepo.EXPECT().GetByID(ctx, reminderID).Return(nil, nil)

		reminder, err := usecase.Update(ctx, reminderID, nil, nil, nil, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, reminder)
		assert.Contains(t, err.Error(), "reminder not found")
	})

	t.Run("update reminder type recalculates next_send_at", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()
		existingReminder := &entities.Reminder{
			ID:     reminderID,
			Title:  "Test",
			Type:   entities.ReminderTypeDaily,
			UserID: uuid.New(),
		}
		newType := entities.ReminderTypeWeekly

		mockRepo.EXPECT().GetByID(ctx, reminderID).Return(existingReminder, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, reminder *entities.Reminder) error {
			assert.Equal(t, newType, reminder.Type)
			assert.NotNil(t, reminder.NextSendAt)
			return nil
		})

		reminder, err := usecase.Update(ctx, reminderID, nil, nil, nil, &newType, nil, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, newType, reminder.Type)
	})
}

func TestReminderUsecase_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()

		mockRepo.EXPECT().Delete(ctx, reminderID).Return(nil)

		err := usecase.Delete(ctx, reminderID)

		assert.NoError(t, err)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReminderRepository(ctrl)
		usecase := NewReminderUsecase(mockRepo)

		reminderID := uuid.New()
		repoError := errors.New("repository error")

		mockRepo.EXPECT().Delete(ctx, reminderID).Return(repoError)

		err := usecase.Delete(ctx, reminderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete reminder")
	})
}

func TestReminderUsecase_CalculateNextSendTime(t *testing.T) {
	usecase := &reminderUsecase{}

	t.Run("daily type", func(t *testing.T) {
		reminder := &entities.Reminder{
			Type: entities.ReminderTypeDaily,
		}

		nextTime := usecase.CalculateNextSendTime(reminder)
		now := time.Now()
		expected := now.Add(24 * time.Hour)

		assert.WithinDuration(t, expected, nextTime, 1*time.Minute)
	})

	t.Run("weekly type", func(t *testing.T) {
		reminder := &entities.Reminder{
			Type: entities.ReminderTypeWeekly,
		}

		nextTime := usecase.CalculateNextSendTime(reminder)
		now := time.Now()
		expected := now.Add(7 * 24 * time.Hour)

		assert.WithinDuration(t, expected, nextTime, 1*time.Minute)
	})

	t.Run("custom type", func(t *testing.T) {
		intervalHours := 6
		reminder := &entities.Reminder{
			Type:          entities.ReminderTypeCustom,
			IntervalHours: &intervalHours,
		}

		nextTime := usecase.CalculateNextSendTime(reminder)
		now := time.Now()
		expected := now.Add(time.Duration(intervalHours) * time.Hour)

		assert.WithinDuration(t, expected, nextTime, 1*time.Minute)
	})

	t.Run("specific type - future time today", func(t *testing.T) {
		timeOfDay := "15:00"
		reminder := &entities.Reminder{
			Type:      entities.ReminderTypeSpecific,
			TimeOfDay: &timeOfDay,
		}

		nextTime := usecase.CalculateNextSendTime(reminder)
		now := time.Now()

		parsedTime, _ := time.Parse("15:04", timeOfDay)
		expected := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())

		if expected.Before(now) || expected.Equal(now) {
			expected = expected.Add(24 * time.Hour)
		}

		assert.WithinDuration(t, expected, nextTime, 1*time.Minute)
	})

	t.Run("specific type - past time today", func(t *testing.T) {
		now := time.Now()
		hour := now.Hour() - 1
		if hour < 0 {
			hour = 23
		}
		timeOfDay := time.Date(2000, 1, 1, hour, now.Minute(), 0, 0, time.UTC).Format("15:04")

		reminder := &entities.Reminder{
			Type:      entities.ReminderTypeSpecific,
			TimeOfDay: &timeOfDay,
		}

		nextTime := usecase.CalculateNextSendTime(reminder)
		expected := now.Add(24 * time.Hour)

		assert.True(t, nextTime.After(now))
		assert.WithinDuration(t, expected, nextTime, 2*time.Hour)
	})
}
