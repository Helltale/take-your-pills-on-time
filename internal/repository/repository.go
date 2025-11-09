package repository

import "gorm.io/gorm"

type Repository struct {
	User              UserRepository
	Reminder          ReminderRepository
	ReminderExecution ReminderExecutionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:              NewUserRepository(db),
		Reminder:          NewReminderRepository(db),
		ReminderExecution: NewReminderExecutionRepository(db),
	}
}
