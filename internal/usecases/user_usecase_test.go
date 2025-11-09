package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository/mocks"
)

func TestUserUsecase_RegisterOrUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("register new user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		firstName := "Test"
		username := "testuser"

		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(nil, nil)
		mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entities.User) error {
			user.ID = uuid.New()
			return nil
		})

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, &username, &firstName, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, telegramID, user.TelegramID)
		assert.Equal(t, firstName, user.FirstName)
		assert.True(t, user.IsActive)
	})

	t.Run("register fails when first name is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		emptyFirstName := ""

		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(nil, nil)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, nil, &emptyFirstName, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "first name is required")
	})

	t.Run("update existing user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		existingUser := &entities.User{
			ID:         uuid.New(),
			TelegramID: telegramID,
			FirstName:  "Old",
			IsActive:   false,
		}
		newFirstName := "New"
		newUsername := "newuser"

		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(existingUser, nil)
		mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entities.User) error {
			assert.Equal(t, newFirstName, user.FirstName)
			assert.True(t, user.IsActive)
			return nil
		})

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, &newUsername, &newFirstName, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newFirstName, user.FirstName)
		assert.True(t, user.IsActive)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		firstName := "Test"

		repoError := errors.New("repository error")
		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(nil, repoError)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, nil, &firstName, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")
	})
}

func TestUserUsecase_GetByTelegramID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		expectedUser := &entities.User{
			ID:         uuid.New(),
			TelegramID: telegramID,
			FirstName:  "Test",
		}

		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(expectedUser, nil)

		user, err := usecase.GetByTelegramID(ctx, telegramID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		repoError := errors.New("repository error")

		mockRepo.EXPECT().GetByTelegramID(ctx, telegramID).Return(nil, repoError)

		user, err := usecase.GetByTelegramID(ctx, telegramID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserUsecase_Deactivate(t *testing.T) {
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		mockRepo.EXPECT().SetActive(ctx, telegramID, false).Return(nil)

		err := usecase.Deactivate(ctx, telegramID)

		assert.NoError(t, err)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		repoError := errors.New("repository error")
		mockRepo.EXPECT().SetActive(ctx, telegramID, false).Return(repoError)

		err := usecase.Deactivate(ctx, telegramID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deactivate")
	})
}

func TestUserUsecase_Activate(t *testing.T) {
	ctx := context.Background()

	t.Run("successful activation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		mockRepo.EXPECT().SetActive(ctx, telegramID, true).Return(nil)

		err := usecase.Activate(ctx, telegramID)

		assert.NoError(t, err)
	})
}
