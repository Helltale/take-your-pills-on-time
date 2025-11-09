package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TelegramID   int64     `db:"telegram_id" json:"telegram_id"`
	Username     *string   `db:"username" json:"username"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     *string   `db:"last_name" json:"last_name"`
	LanguageCode *string   `db:"language_code" json:"language_code"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
