package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
)

type MockUserRepository struct {
	mock.Mock
}

var _ repository.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error) {
	args := m.Called(ctx, telegramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) SetActive(ctx context.Context, telegramID int64, isActive bool) error {
	args := m.Called(ctx, telegramID, isActive)
	return args.Error(0)
}

func TestUserUsecase_RegisterOrUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("register new user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		firstName := "Test"
		username := "testuser"

		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(nil, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entities.User")).Return(nil)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, &username, &firstName, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, telegramID, user.TelegramID)
		assert.Equal(t, firstName, user.FirstName)
		assert.True(t, user.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("register fails when first name is empty", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		emptyFirstName := ""

		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(nil, nil)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, nil, &emptyFirstName, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "first name is required")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("update existing user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
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

		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(existingUser, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entities.User")).Return(nil)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, &newUsername, &newFirstName, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newFirstName, user.FirstName)
		assert.True(t, user.IsActive)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		firstName := "Test"

		repoError := errors.New("repository error")
		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(nil, repoError)

		user, err := usecase.RegisterOrUpdate(ctx, telegramID, nil, &firstName, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")
	})
}

func TestUserUsecase_GetByTelegramID(t *testing.T) {
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		expectedUser := &entities.User{
			ID:         uuid.New(),
			TelegramID: telegramID,
			FirstName:  "Test",
		}

		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(expectedUser, nil)

		user, err := usecase.GetByTelegramID(ctx, telegramID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		repoError := errors.New("repository error")

		mockRepo.On("GetByTelegramID", ctx, telegramID).Return(nil, repoError)

		user, err := usecase.GetByTelegramID(ctx, telegramID)

		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Deactivate(t *testing.T) {
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		mockRepo.On("SetActive", ctx, telegramID, false).Return(nil)

		err := usecase.Deactivate(ctx, telegramID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		repoError := errors.New("repository error")
		mockRepo.On("SetActive", ctx, telegramID, false).Return(repoError)

		err := usecase.Deactivate(ctx, telegramID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to deactivate")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Activate(t *testing.T) {
	ctx := context.Background()

	t.Run("successful activation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		telegramID := int64(12345)
		mockRepo.On("SetActive", ctx, telegramID, true).Return(nil)

		err := usecase.Activate(ctx, telegramID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
