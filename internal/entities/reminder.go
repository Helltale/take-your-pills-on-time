package entities

import (
	"time"

	"github.com/google/uuid"
)

type ReminderType string

const (
	ReminderTypeDaily    ReminderType = "daily"
	ReminderTypeWeekly   ReminderType = "weekly"
	ReminderTypeCustom   ReminderType = "custom"
	ReminderTypeSpecific ReminderType = "specific"
)

type Reminder struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Title         string       `gorm:"size:255;not null" json:"title"`
	Comment       *string      `gorm:"type:text" json:"comment"`
	ImageURL      *string      `gorm:"type:text" json:"image_url"`
	Type          ReminderType `gorm:"type:varchar(50);not null;index" json:"type"`
	IntervalHours *int         `json:"interval_hours"`
	TimeOfDay     *string      `gorm:"size:5" json:"time_of_day"`
	IsActive      bool         `gorm:"default:true;not null;index" json:"is_active"`
	LastSentAt    *time.Time   `json:"last_sent_at"`
	NextSendAt    *time.Time   `gorm:"index" json:"next_send_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (Reminder) TableName() string {
	return "reminders"
}
