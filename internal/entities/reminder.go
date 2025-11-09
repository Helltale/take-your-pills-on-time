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
	ID            uuid.UUID    `db:"id" json:"id"`
	UserID        uuid.UUID    `db:"user_id" json:"user_id"`
	Title         string       `db:"title" json:"title"`
	Comment       *string      `db:"comment" json:"comment"`
	ImageURL      *string      `db:"image_url" json:"image_url"`
	Type          ReminderType `db:"type" json:"type"`
	IntervalHours *int         `db:"interval_hours" json:"interval_hours"`
	TimeOfDay     *string      `db:"time_of_day" json:"time_of_day"`
	IsActive      bool         `db:"is_active" json:"is_active"`
	LastSentAt    *time.Time   `db:"last_sent_at" json:"last_sent_at"`
	NextSendAt    *time.Time   `db:"next_send_at" json:"next_send_at"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
}

func (Reminder) TableName() string {
	return "reminders"
}
