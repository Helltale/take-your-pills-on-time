package repository

//go:generate mockgen -source=user_repository.go -destination=./mocks/user_repository_mock.go -package=mocks
//go:generate mockgen -source=reminder_repository.go -destination=./mocks/reminder_repository_mock.go -package=mocks
//go:generate mockgen -source=reminder_execution_repository.go -destination=./mocks/reminder_execution_repository_mock.go -package=mocks
