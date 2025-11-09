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
	ID          uuid.UUID       `db:"id" json:"id"`
	ReminderID  uuid.UUID       `db:"reminder_id" json:"reminder_id"`
	UserID      uuid.UUID       `db:"user_id" json:"user_id"`
	Status      ExecutionStatus `db:"status" json:"status"`
	SentAt      time.Time       `db:"sent_at" json:"sent_at"`
	ConfirmedAt *time.Time      `db:"confirmed_at" json:"confirmed_at"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}

func (ReminderExecution) TableName() string {
	return "reminder_executions"
}
