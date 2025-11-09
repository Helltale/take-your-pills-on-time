package scheduler

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/handlers"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
	"github.com/Helltale/take-your-pills-on-time/internal/usecases"
)

type Scheduler struct {
	reminderRepo     repository.ReminderRepository
	executionUsecase usecases.ReminderExecutionUsecase
	reminderUsecase  usecases.ReminderUsecase
	handler          *handlers.BotHandler
	logger           *zap.Logger
	ticker           *time.Ticker
	stopChan         chan struct{}
}

func NewScheduler(
	reminderRepo repository.ReminderRepository,
	executionUsecase usecases.ReminderExecutionUsecase,
	reminderUsecase usecases.ReminderUsecase,
	handler *handlers.BotHandler,
	logger *zap.Logger,
) *Scheduler {
	return &Scheduler{
		reminderRepo:     reminderRepo,
		executionUsecase: executionUsecase,
		reminderUsecase:  reminderUsecase,
		handler:          handler,
		logger:           logger,
		stopChan:         make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.processReminders(ctx)
			case <-s.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	s.logger.Info("Scheduler started")
}

func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
	s.logger.Info("Scheduler stopped")
}

func (s *Scheduler) processReminders(ctx context.Context) {
	now := time.Now()
	reminders, err := s.reminderRepo.GetDueReminders(ctx)
	if err != nil {
		s.logger.Error("failed to get due reminders", zap.Error(err))
		return
	}

	if len(reminders) > 0 {
		s.logger.Info("processing reminders",
			zap.Time("current_time", now),
			zap.Int("found_count", len(reminders)),
		)
	}

	for _, reminder := range reminders {
		s.logger.Info("found due reminder",
			zap.String("reminder_id", reminder.ID.String()),
			zap.String("title", reminder.Title),
			zap.String("type", string(reminder.Type)),
			zap.Time("next_send_at", *reminder.NextSendAt),
			zap.Time("current_time", now),
		)
		if err := s.sendReminder(ctx, reminder); err != nil {
			s.logger.Error("failed to send reminder",
				zap.Error(err),
				zap.String("reminder_id", reminder.ID.String()),
			)
			continue
		}

		nextSendTime := s.reminderUsecase.CalculateNextSendTime(reminder)
		if err := s.reminderRepo.UpdateNextSendAt(ctx, reminder.ID, nextSendTime); err != nil {
			s.logger.Error("failed to update next send time",
				zap.Error(err),
				zap.String("reminder_id", reminder.ID.String()),
			)
		}

		now := time.Now()
		if err := s.reminderRepo.UpdateLastSentAt(ctx, reminder.ID, now); err != nil {
			s.logger.Error("failed to update last sent time",
				zap.Error(err),
				zap.String("reminder_id", reminder.ID.String()),
			)
		}
	}
}

func (s *Scheduler) sendReminder(ctx context.Context, reminder *entities.Reminder) error {
	execution, err := s.executionUsecase.RecordSent(ctx, reminder.ID, reminder.UserID)
	if err != nil {
		return fmt.Errorf("failed to record sent execution: %w", err)
	}

	if err := s.handler.SendReminder(ctx, reminder, execution.ID); err != nil {
		return fmt.Errorf("failed to send reminder: %w", err)
	}

	s.logger.Info("reminder sent",
		zap.String("reminder_id", reminder.ID.String()),
		zap.String("execution_id", execution.ID.String()),
	)

	return nil
}
