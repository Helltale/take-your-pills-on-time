package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Helltale/take-your-pills-on-time/internal/entities"
	"github.com/Helltale/take-your-pills-on-time/internal/usecases"
)

type BotHandler struct {
	bot      *tgbotapi.BotAPI
	usecases *usecases.Usecases
	logger   *zap.Logger
}

func NewBotHandler(bot *tgbotapi.BotAPI, usecases *usecases.Usecases, logger *zap.Logger) *BotHandler {
	return &BotHandler{
		bot:      bot,
		usecases: usecases,
		logger:   logger,
	}
}

func (h *BotHandler) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(ctx, update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	msg := update.Message
	userID := msg.From.ID

	username := msg.From.UserName
	firstName := msg.From.FirstName
	lastName := msg.From.LastName
	languageCode := msg.From.LanguageCode

	_, err := h.usecases.User.RegisterOrUpdate(ctx, int64(userID), &username, &firstName, &lastName, &languageCode)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err), zap.Int64("user_id", int64(userID)))
	}

	if msg.IsCommand() {
		h.handleCommand(ctx, msg)
		return
	}

	h.handleTextMessage(ctx, msg)
}

func (h *BotHandler) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	command := msg.Command()
	chatID := msg.Chat.ID

	switch command {
	case "start":
		h.handleStart(ctx, chatID, msg.From)
	case "help":
		h.handleHelp(ctx, chatID)
	case "new":
		h.handleNewReminder(ctx, chatID, int64(msg.From.ID))
	case "list":
		h.handleListReminders(ctx, chatID, int64(msg.From.ID))
	case "stats":
		h.handleStats(ctx, chatID, int64(msg.From.ID))
	default:
		h.sendMessage(chatID, "Неизвестная команда. Используйте /help для списка команд.")
	}
}

func (h *BotHandler) handleStart(ctx context.Context, chatID int64, user *tgbotapi.User) {
	text := fmt.Sprintf(
		"Привет, %s! 👋\n\n"+
			"Я бот для напоминаний о приеме лекарств.\n\n"+
			"Доступные команды:\n"+
			"/new - создать новое напоминание\n"+
			"/list - список ваших напоминаний\n"+
			"/stats - статистика выполнения\n"+
			"/help - помощь\n\n"+
			"Начните с команды /new для создания первого напоминания!",
		user.FirstName,
	)
	h.sendMessage(chatID, text)
}

func (h *BotHandler) handleHelp(ctx context.Context, chatID int64) {
	text := `📚 Справка по командам:

/new - Создать новое напоминание
/list - Показать все ваши напоминания
/stats - Показать статистику выполнения напоминаний
/help - Показать эту справку

Для создания напоминания используйте команду /new и следуйте инструкциям.`
	h.sendMessage(chatID, text)
}

func (h *BotHandler) handleNewReminder(ctx context.Context, chatID int64, telegramUserID int64) {
	user, err := h.usecases.User.GetByTelegramID(ctx, telegramUserID)
	if err != nil || user == nil {
		h.sendMessage(chatID, "Ошибка: пользователь не найден. Попробуйте /start")
		return
	}

	text := `Создание нового напоминания 📝

Пожалуйста, отправьте данные в следующем формате:
Название|Тип|Комментарий|Время

Типы напоминаний:
- daily - ежедневно
- weekly - еженедельно
- custom - кастомный интервал (укажите количество часов)
- specific - конкретное время каждый день (формат HH:MM)

Примеры:
Лекарство|daily|Принять после еды|09:00
Витамины|custom|Утром|6
Завтрак|specific|Важно!|08:30

Или используйте упрощенный формат:
Название|daily

Для отмены отправьте /cancel`

	h.sendMessage(chatID, text)
}

func (h *BotHandler) handleListReminders(ctx context.Context, chatID int64, telegramUserID int64) {
	user, err := h.usecases.User.GetByTelegramID(ctx, telegramUserID)
	if err != nil || user == nil {
		h.sendMessage(chatID, "Ошибка: пользователь не найден.")
		return
	}

	reminders, err := h.usecases.Reminder.GetByUserID(ctx, user.ID)
	if err != nil {
		h.logger.Error("failed to get reminders", zap.Error(err))
		h.sendMessage(chatID, "Ошибка при получении списка напоминаний.")
		return
	}

	if len(reminders) == 0 {
		h.sendMessage(chatID, "У вас пока нет напоминаний. Создайте первое с помощью /new")
		return
	}

	var builder strings.Builder
	builder.WriteString("📋 Ваши напоминания:\n\n")

	for i, reminder := range reminders {
		status := "✅ Активно"
		if !reminder.IsActive {
			status = "❌ Неактивно"
		}

		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, reminder.Title))
		builder.WriteString(fmt.Sprintf("   Тип: %s\n", reminder.Type))
		if reminder.Comment != nil {
			builder.WriteString(fmt.Sprintf("   Комментарий: %s\n", *reminder.Comment))
		}
		if reminder.NextSendAt != nil {
			builder.WriteString(fmt.Sprintf("   Следующая отправка: %s\n", reminder.NextSendAt.Format("02.01.2006 15:04")))
		}
		builder.WriteString(fmt.Sprintf("   Статус: %s\n\n", status))
	}

	h.sendMessage(chatID, builder.String())
}

func (h *BotHandler) handleStats(ctx context.Context, chatID int64, telegramUserID int64) {
	user, err := h.usecases.User.GetByTelegramID(ctx, telegramUserID)
	if err != nil || user == nil {
		h.sendMessage(chatID, "Ошибка: пользователь не найден.")
		return
	}

	toDate := time.Now()
	fromDate := toDate.AddDate(0, 0, -30)

	stats, err := h.usecases.ReminderExecution.GetStatisticsByUserID(ctx, user.ID, fromDate, toDate)
	if err != nil {
		h.logger.Error("failed to get statistics", zap.Error(err))
		h.sendMessage(chatID, "Ошибка при получении статистики.")
		return
	}

	text := fmt.Sprintf(
		"📊 Статистика за последние 30 дней:\n\n"+
			"Отправлено: %d\n"+
			"Подтверждено: %d\n"+
			"Пропущено: %d\n"+
			"Процент выполнения: %.1f%%",
		stats.TotalSent,
		stats.TotalConfirmed,
		stats.TotalSkipped,
		stats.ConfirmationRate,
	)

	h.sendMessage(chatID, text)
}

func (h *BotHandler) handleTextMessage(ctx context.Context, msg *tgbotapi.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID
	telegramUserID := int64(msg.From.ID)

	if text == "/cancel" {
		h.sendMessage(chatID, "Создание напоминания отменено.")
		return
	}

	if !strings.Contains(text, "|") {
		h.sendMessage(chatID, "Неверный формат. Используйте формат: Название|Тип|Комментарий|Время\nИли: Название|Тип\n\nДля отмены отправьте /cancel")
		return
	}

	parts := strings.Split(text, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	if len(parts) < 2 {
		h.sendMessage(chatID, "Неверный формат. Минимум требуется: Название|Тип")
		return
	}

	title := parts[0]
	reminderTypeStr := strings.ToLower(parts[1])

	if title == "" {
		h.sendMessage(chatID, "Ошибка: название не может быть пустым.")
		return
	}

	var reminderType entities.ReminderType
	switch reminderTypeStr {
	case "daily":
		reminderType = entities.ReminderTypeDaily
	case "weekly":
		reminderType = entities.ReminderTypeWeekly
	case "custom":
		reminderType = entities.ReminderTypeCustom
	case "specific":
		reminderType = entities.ReminderTypeSpecific
	default:
		h.sendMessage(chatID, fmt.Sprintf("Неизвестный тип напоминания: %s\nДоступные типы: daily, weekly, custom, specific", reminderTypeStr))
		return
	}

	var comment *string
	var timeOfDay *string
	var intervalHours *int

	if len(parts) >= 3 && parts[2] != "" {
		comment = &parts[2]
	}

	if len(parts) >= 4 && parts[3] != "" {
		if reminderType == entities.ReminderTypeCustom {
			interval, err := strconv.Atoi(parts[3])
			if err != nil || interval <= 0 {
				h.sendMessage(chatID, "Ошибка: для типа 'custom' требуется положительное число часов.")
				return
			}
			intervalHours = &interval
		} else if reminderType == entities.ReminderTypeSpecific {
			if _, err := time.Parse("15:04", parts[3]); err != nil {
				h.sendMessage(chatID, "Ошибка: неверный формат времени. Используйте формат HH:MM (например, 09:00)")
				return
			}
			timeOfDay = &parts[3]
		} else {
			if _, err := time.Parse("15:04", parts[3]); err == nil {
				timeOfDay = &parts[3]
			}
		}
	}

	user, err := h.usecases.User.GetByTelegramID(ctx, telegramUserID)
	if err != nil || user == nil {
		h.sendMessage(chatID, "Ошибка: пользователь не найден. Попробуйте /start")
		return
	}

	reminder, err := h.usecases.Reminder.Create(ctx, user.ID, title, comment, nil, reminderType, intervalHours, timeOfDay)
	if err != nil {
		h.logger.Error("failed to create reminder", zap.Error(err), zap.Int64("user_id", telegramUserID))
		h.sendMessage(chatID, fmt.Sprintf("Ошибка при создании напоминания: %s", err.Error()))
		return
	}

	var responseBuilder strings.Builder
	responseBuilder.WriteString("✅ Напоминание успешно создано!\n\n")
	responseBuilder.WriteString(fmt.Sprintf("📝 Название: %s\n", reminder.Title))
	responseBuilder.WriteString(fmt.Sprintf("🔄 Тип: %s\n", reminder.Type))
	if reminder.Comment != nil {
		responseBuilder.WriteString(fmt.Sprintf("💬 Комментарий: %s\n", *reminder.Comment))
	}
	if reminder.TimeOfDay != nil {
		responseBuilder.WriteString(fmt.Sprintf("⏰ Время: %s\n", *reminder.TimeOfDay))
	}
	if reminder.IntervalHours != nil {
		responseBuilder.WriteString(fmt.Sprintf("⏱ Интервал: %d часов\n", *reminder.IntervalHours))
	}
	if reminder.NextSendAt != nil {
		responseBuilder.WriteString(fmt.Sprintf("📅 Следующая отправка: %s\n", reminder.NextSendAt.Format("02.01.2006 15:04")))
	}

	h.sendMessage(chatID, responseBuilder.String())
}

func (h *BotHandler) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	chatID := callback.Message.Chat.ID

	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		h.answerCallbackQuery(callback.ID, "Ошибка обработки команды")
		return
	}

	action := parts[0]

	switch action {
	case "confirm":
		if len(parts) >= 3 {
			executionID, err := uuid.Parse(parts[2])
			if err == nil {
				if err := h.usecases.ReminderExecution.RecordConfirmed(ctx, executionID); err == nil {
					h.answerCallbackQuery(callback.ID, "✅ Подтверждено!")
					h.sendMessage(chatID, "Спасибо! Напоминание подтверждено.")
				}
			}
		}
	case "skip":
		if len(parts) >= 3 {
			executionID, err := uuid.Parse(parts[2])
			if err == nil {
				if err := h.usecases.ReminderExecution.RecordSkipped(ctx, executionID); err == nil {
					h.answerCallbackQuery(callback.ID, "⏭ Пропущено")
					h.sendMessage(chatID, "Напоминание пропущено.")
				}
			}
		}
	}
}

func (h *BotHandler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Error("failed to send message", zap.Error(err), zap.Int64("chat_id", chatID))
	}
}

func (h *BotHandler) SendReminder(ctx context.Context, reminder *entities.Reminder, executionID uuid.UUID) error {
	user, err := h.usecases.User.GetByID(ctx, reminder.UserID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("🔔 *%s*\n\n", reminder.Title))

	if reminder.Comment != nil {
		builder.WriteString(fmt.Sprintf("%s\n\n", *reminder.Comment))
	}

	confirmBtn := tgbotapi.NewInlineKeyboardButtonData("✅ Выполнено", fmt.Sprintf("confirm:%s:%s", reminder.ID.String(), executionID.String()))
	skipBtn := tgbotapi.NewInlineKeyboardButtonData("⏭ Пропустить", fmt.Sprintf("skip:%s:%s", reminder.ID.String(), executionID.String()))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(confirmBtn, skipBtn),
	)

	if reminder.ImageURL != nil && *reminder.ImageURL != "" {
		photo := tgbotapi.NewPhoto(int64(user.TelegramID), tgbotapi.FileURL(*reminder.ImageURL))
		photo.Caption = builder.String()
		photo.ParseMode = tgbotapi.ModeMarkdown
		photo.ReplyMarkup = keyboard

		if _, err := h.bot.Send(photo); err != nil {
			return fmt.Errorf("failed to send reminder: %w", err)
		}
	} else {
		msg := tgbotapi.NewMessage(int64(user.TelegramID), builder.String())
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = keyboard

		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("failed to send reminder: %w", err)
		}
	}

	return nil
}

func (h *BotHandler) answerCallbackQuery(callbackID string, text string) {
	callback := tgbotapi.NewCallback(callbackID, text)
	if _, err := h.bot.Request(callback); err != nil {
		h.logger.Error("failed to answer callback query", zap.Error(err))
	}
}
