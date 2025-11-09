package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	User              UserRepository
	Reminder          ReminderRepository
	ReminderExecution ReminderExecutionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:              NewUserRepository(db),
		Reminder:          NewReminderRepository(db),
		ReminderExecution: NewReminderExecutionRepository(db),
	}
}
