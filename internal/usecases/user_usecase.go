package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
)

type UserUsecase interface {
	RegisterOrUpdate(ctx context.Context, telegramID int64, username, firstName, lastName, languageCode *string) (*entities.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Deactivate(ctx context.Context, telegramID int64) error
	Activate(ctx context.Context, telegramID int64) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) RegisterOrUpdate(ctx context.Context, telegramID int64, username, firstName, lastName, languageCode *string) (*entities.User, error) {
	user, err := u.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		if firstName == nil || *firstName == "" {
			return nil, fmt.Errorf("first name is required")
		}

		user = &entities.User{
			TelegramID:   telegramID,
			Username:     username,
			FirstName:    *firstName,
			LastName:     lastName,
			LanguageCode: languageCode,
			IsActive:     true,
		}

		if err := u.repo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		if username != nil {
			user.Username = username
		}
		if firstName != nil {
			user.FirstName = *firstName
		}
		if lastName != nil {
			user.LastName = lastName
		}
		if languageCode != nil {
			user.LanguageCode = languageCode
		}
		user.IsActive = true

		if err := u.repo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return user, nil
}

func (u *userUsecase) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error) {
	user, err := u.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *userUsecase) Deactivate(ctx context.Context, telegramID int64) error {
	if err := u.repo.SetActive(ctx, telegramID, false); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}

func (u *userUsecase) Activate(ctx context.Context, telegramID int64) error {
	if err := u.repo.SetActive(ctx, telegramID, true); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	return nil
}
