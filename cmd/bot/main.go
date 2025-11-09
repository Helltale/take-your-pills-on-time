package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Helltale/take-your-pills-on-time/internal/config"
	"github.com/Helltale/take-your-pills-on-time/internal/handlers"
	"github.com/Helltale/take-your-pills-on-time/internal/migrations"
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

	appLogger := initLogger(cfg.App.LogLevel)
	defer appLogger.Sync()

	appLogger.Info("Starting application", zap.String("env", cfg.App.Env))

	db, err := connectDatabase(cfg.Database.DSN(), cfg.App.LogLevel, appLogger)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Fatal("Failed to get sql.DB", zap.Error(err))
	}
	defer sqlDB.Close()

	migrator := migrations.NewMigrator(db, appLogger)
	if err := migrator.Run(); err != nil {
		appLogger.Fatal("Failed to run migrations", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		appLogger.Fatal("Failed to initialize bot", zap.Error(err))
	}

	appLogger.Info("Bot authorized", zap.String("username", bot.Self.UserName))

	repo := repository.NewRepository(db)

	usecases := usecases.NewUsecases(repo)

	handler := handlers.NewBotHandler(bot, usecases, appLogger)

	sched := scheduler.NewScheduler(repo.Reminder, usecases.ReminderExecution, usecases.Reminder, handler, appLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched.Start(ctx)
	defer sched.Stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	appLogger.Info("Bot is running. Press Ctrl+C to stop.")

	for {
		select {
		case update := <-updates:
			if update.CallbackQuery != nil {
				handler.HandleUpdate(ctx, update)
			} else if update.Message != nil {
				handler.HandleUpdate(ctx, update)
			}
		case <-sigChan:
			appLogger.Info("Shutting down...")
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

func connectDatabase(dsn string, logLevel string, appLogger *zap.Logger) (*gorm.DB, error) {
	var gormLogger logger.Interface
	if logLevel == "development" || logLevel == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	appLogger.Info("Database connection established")
	return db, nil
}
