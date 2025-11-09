package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/Helltale/take-your-pills-on-time/internal/config"
	"github.com/Helltale/take-your-pills-on-time/internal/handlers"
	"github.com/Helltale/take-your-pills-on-time/internal/repository"
	"github.com/Helltale/take-your-pills-on-time/internal/scheduler"
	"github.com/Helltale/take-your-pills-on-time/internal/usecases"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := initLogger(cfg.App.LogLevel)
	defer logger.Sync()

	logger.Info("Starting application", zap.String("env", cfg.App.Env))

	db, err := connectDatabase(cfg.Database.DSN(), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := runMigrations(cfg.Database.URL(), logger); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		logger.Fatal("Failed to initialize bot", zap.Error(err))
	}

	logger.Info("Bot authorized", zap.String("username", bot.Self.UserName))

	repo := repository.NewRepository(db)

	usecases := usecases.NewUsecases(repo)

	handler := handlers.NewBotHandler(bot, usecases, logger)

	sched := scheduler.NewScheduler(repo.Reminder, usecases.ReminderExecution, usecases.Reminder, handler, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched.Start(ctx)
	defer sched.Stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("Bot is running. Press Ctrl+C to stop.")

	for {
		select {
		case update := <-updates:
			if update.CallbackQuery != nil {
				handler.HandleUpdate(ctx, update)
			} else if update.Message != nil {
				handler.HandleUpdate(ctx, update)
			}
		case <-sigChan:
			logger.Info("Shutting down...")
			cancel()
			time.Sleep(2 * time.Second)
			return
		}
	}
}

func initLogger(level string) *zap.Logger {
	var zapConfig zap.Config

	if level == "development" || level == "debug" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	logger, err := zapConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}

func connectDatabase(dsn string, logger *zap.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")
	return db, nil
}

func runMigrations(url string, logger *zap.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		url,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Migrations completed")
	return nil
}
