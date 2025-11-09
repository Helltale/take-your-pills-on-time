package usecases

import "github.com/Helltale/take-your-pills-on-time/internal/repository"

type Usecases struct {
	User              UserUsecase
	Reminder          ReminderUsecase
	ReminderExecution ReminderExecutionUsecase
}

func NewUsecases(repo *repository.Repository) *Usecases {
	return &Usecases{
		User:              NewUserUsecase(repo.User),
		Reminder:          NewReminderUsecase(repo.Reminder),
		ReminderExecution: NewReminderExecutionUsecase(repo.ReminderExecution),
	}
}
