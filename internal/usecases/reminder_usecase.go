package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
)

type ReminderUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, title string, comment *string, imageURL *string, reminderType entities.ReminderType, intervalHours *int, timeOfDay *string) (*entities.Reminder, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Reminder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error)
	Update(ctx context.Context, id uuid.UUID, title *string, comment *string, imageURL *string, reminderType *entities.ReminderType, intervalHours *int, timeOfDay *string, isActive *bool) (*entities.Reminder, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateNextSendTime(reminder *entities.Reminder) time.Time
}

type reminderUsecase struct {
	repo repository.ReminderRepository
}

func NewReminderUsecase(repo repository.ReminderRepository) ReminderUsecase {
	return &reminderUsecase{repo: repo}
}

func (u *reminderUsecase) Create(ctx context.Context, userID uuid.UUID, title string, comment *string, imageURL *string, reminderType entities.ReminderType, intervalHours *int, timeOfDay *string) (*entities.Reminder, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	switch reminderType {
	case entities.ReminderTypeCustom:
		if intervalHours == nil || *intervalHours <= 0 {
			return nil, fmt.Errorf("interval_hours is required for custom type and must be greater than 0")
		}
	case entities.ReminderTypeSpecific:
		if timeOfDay == nil || *timeOfDay == "" {
			return nil, fmt.Errorf("time_of_day is required for specific type")
		}
		if _, err := time.Parse("15:04", *timeOfDay); err != nil {
			return nil, fmt.Errorf("invalid time_of_day format, expected HH:MM")
		}
	}

	reminder := &entities.Reminder{
		UserID:        userID,
		Title:         title,
		Comment:       comment,
		ImageURL:      imageURL,
		Type:          reminderType,
		IntervalHours: intervalHours,
		TimeOfDay:     timeOfDay,
		IsActive:      true,
	}

	nextTime := u.CalculateNextSendTime(reminder)
	reminder.NextSendAt = &nextTime

	if err := u.repo.Create(ctx, reminder); err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return reminder, nil
}

func (u *reminderUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entities.Reminder, error) {
	reminder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}
	return reminder, nil
}

func (u *reminderUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	reminders, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders: %w", err)
	}
	return reminders, nil
}

func (u *reminderUsecase) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reminder, error) {
	reminders, err := u.repo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reminders: %w", err)
	}
	return reminders, nil
}

func (u *reminderUsecase) Update(ctx context.Context, id uuid.UUID, title *string, comment *string, imageURL *string, reminderType *entities.ReminderType, intervalHours *int, timeOfDay *string, isActive *bool) (*entities.Reminder, error) {
	reminder, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}
	if reminder == nil {
		return nil, fmt.Errorf("reminder not found")
	}

	if title != nil {
		reminder.Title = *title
	}
	if comment != nil {
		reminder.Comment = comment
	}
	if imageURL != nil {
		reminder.ImageURL = imageURL
	}
	if reminderType != nil {
		reminder.Type = *reminderType
		nextTime := u.CalculateNextSendTime(reminder)
		reminder.NextSendAt = &nextTime
	}
	if intervalHours != nil {
		reminder.IntervalHours = intervalHours
		if reminder.Type == entities.ReminderTypeCustom {
			nextTime := u.CalculateNextSendTime(reminder)
			reminder.NextSendAt = &nextTime
		}
	}
	if timeOfDay != nil {
		reminder.TimeOfDay = timeOfDay
		if reminder.Type == entities.ReminderTypeSpecific {
			nextTime := u.CalculateNextSendTime(reminder)
			reminder.NextSendAt = &nextTime
		}
	}
	if isActive != nil {
		reminder.IsActive = *isActive
	}

	if err := u.repo.Update(ctx, reminder); err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	return reminder, nil
}

func (u *reminderUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete reminder: %w", err)
	}
	return nil
}

func (u *reminderUsecase) CalculateNextSendTime(reminder *entities.Reminder) time.Time {
	now := time.Now()

	switch reminder.Type {
	case entities.ReminderTypeDaily:
		next := now.Add(24 * time.Hour)
		return time.Date(next.Year(), next.Month(), next.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	case entities.ReminderTypeWeekly:
		return now.Add(7 * 24 * time.Hour)

	case entities.ReminderTypeCustom:
		if reminder.IntervalHours != nil {
			return now.Add(time.Duration(*reminder.IntervalHours) * time.Hour)
		}
		return now.Add(24 * time.Hour)

	case entities.ReminderTypeSpecific:
		if reminder.TimeOfDay != nil {
			parsedTime, err := time.Parse("15:04", *reminder.TimeOfDay)
			if err == nil {
				hour := parsedTime.Hour()
				minute := parsedTime.Minute()
				next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
				if next.Before(now) || next.Equal(now) {
					next = next.Add(24 * time.Hour)
				}
				return next
			}
		}
		return now.Add(24 * time.Hour)

	default:
		return now.Add(24 * time.Hour)
	}
}
