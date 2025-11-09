package entities

import (
	"time"

	"github.com/google/uuid"
)

type ExecutionStatus string

const (
	ExecutionStatusSent      ExecutionStatus = "sent"
	ExecutionStatusConfirmed ExecutionStatus = "confirmed"
	ExecutionStatusSkipped   ExecutionStatus = "skipped"
)

type ReminderExecution struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReminderID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"reminder_id"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	Status      ExecutionStatus `gorm:"type:varchar(50);not null;index" json:"status"`
	SentAt      time.Time       `gorm:"not null;index" json:"sent_at"`
	ConfirmedAt *time.Time      `json:"confirmed_at"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (ReminderExecution) TableName() string {
	return "reminder_executions"
}
