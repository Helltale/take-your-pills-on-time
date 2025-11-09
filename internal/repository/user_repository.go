package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	SetActive(ctx context.Context, telegramID int64, isActive bool) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) SetActive(ctx context.Context, telegramID int64, isActive bool) error {
	return r.db.WithContext(ctx).Model(&entities.User{}).
		Where("telegram_id = ?", telegramID).
		Updates(map[string]interface{}{
			"is_active":  isActive,
			"updated_at": time.Now(),
		}).Error
}
