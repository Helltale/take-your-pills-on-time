package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

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
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, telegram_id, username, first_name, last_name, language_code, is_active, created_at, updated_at)
		VALUES (:id, :telegram_id, :username, :first_name, :last_name, :language_code, :is_active, :created_at, :updated_at)
	`

	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE telegram_id = $1`

	err := r.db.GetContext(ctx, &user, query, telegramID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET username = :username, first_name = :first_name, last_name = :last_name, 
		    language_code = :language_code, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id
	`

	user.UpdatedAt = time.Now()
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) SetActive(ctx context.Context, telegramID int64, isActive bool) error {
	query := `UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE telegram_id = $2`
	_, err := r.db.ExecContext(ctx, query, isActive, telegramID)
	return err
}
