package migrations

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
)

type Migrator struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewMigrator(db *gorm.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

func (m *Migrator) Run() error {
	m.logger.Info("Starting database migrations")

	err := m.db.AutoMigrate(
		&entities.User{},
		&entities.Reminder{},
		&entities.ReminderExecution{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Migrations completed successfully")
	return nil
}

func (m *Migrator) Rollback() error {
	m.logger.Info("Rolling back migrations")

	err := m.db.Migrator().DropTable(
		&entities.ReminderExecution{},
		&entities.Reminder{},
		&entities.User{},
	)
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	m.logger.Info("Rollback completed successfully")
	return nil
}
