package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TelegramID   int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     *string   `gorm:"size:255" json:"username"`
	FirstName    string    `gorm:"size:255;not null" json:"first_name"`
	LastName     *string   `gorm:"size:255" json:"last_name"`
	LanguageCode *string   `gorm:"size:10" json:"language_code"`
	IsActive     bool      `gorm:"default:true;not null;index" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
