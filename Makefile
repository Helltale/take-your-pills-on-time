.PHONY: test test-cover test-all build docker-build docker-up docker-down docker-restart docker-logs run clean help generate-mocks

# Переменные
DOCKER_COMPOSE = docker-compose
GO_TEST = go test
GO_TEST_COVER = go test -cover

# Цвета для вывода
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

help: ## Показать справку по командам
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

test: ## Запустить все юнит-тесты
	@echo "$(YELLOW)Запуск юнит-тестов...$(NC)"
	@$(GO_TEST) ./...

test-cover: ## Запустить тесты с покрытием
	@echo "$(YELLOW)Запуск тестов с покрытием...$(NC)"
	@$(GO_TEST_COVER) ./...

test-all: test-cover ## Запустить все тесты с покрытием (алиас для test-cover)

generate-mocks: ## Сгенерировать моки для репозиториев
	@echo "$(YELLOW)Генерация моков...$(NC)"
	@export PATH=$$PATH:$$HOME/go/bin && go generate ./internal/repository/...
	@echo "$(GREEN)Моки успешно сгенерированы$(NC)"

docker-build: ## Собрать Docker образ
	@echo "$(YELLOW)Сборка Docker образа...$(NC)"
	@$(DOCKER_COMPOSE) build

docker-up: ## Запустить контейнеры
	@echo "$(YELLOW)Запуск контейнеров...$(NC)"
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Остановить контейнеры
	@echo "$(YELLOW)Остановка контейнеров...$(NC)"
	@$(DOCKER_COMPOSE) down

docker-restart: docker-down docker-up ## Перезапустить контейнеры

docker-logs: ## Показать логи контейнеров
	@$(DOCKER_COMPOSE) logs -f bot

check-env: ## Проверить наличие .env файла
	@if [ ! -f .env ]; then \
		echo "$(RED)Ошибка: файл .env не найден!$(NC)"; \
		echo "$(YELLOW)Скопируйте .env.example в .env и заполните необходимые переменные:$(NC)"; \
		echo "  cp .env.example .env"; \
		exit 1; \
	fi
	@if grep -q "^TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here" .env 2>/dev/null || \
	   ! grep -q "^TELEGRAM_BOT_TOKEN=" .env 2>/dev/null || \
	   grep -q "^TELEGRAM_BOT_TOKEN=$$" .env 2>/dev/null; then \
		echo "$(RED)Ошибка: TELEGRAM_BOT_TOKEN не установлен в .env файле!$(NC)"; \
		echo "$(YELLOW)Установите реальный токен бота от @BotFather$(NC)"; \
		exit 1; \
	fi

docker-build-up: check-env docker-build docker-up ## Собрать и запустить контейнеры

run: test check-env docker-build-up ## Запустить проект: сначала тесты, потом сборка и запуск в Docker
	@echo "$(GREEN)✓ Все тесты пройдены успешно$(NC)"
	@echo "$(GREEN)✓ Docker образ собран$(NC)"
	@echo "$(GREEN)✓ Контейнеры запущены$(NC)"
	@echo "$(GREEN)Проект успешно запущен!$(NC)"
	@echo "$(YELLOW)Для просмотра логов выполните: make docker-logs$(NC)"

clean: ## Очистить временные файлы и остановить контейнеры
	@echo "$(YELLOW)Очистка...$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@go clean -cache
	@echo "$(GREEN)Очистка завершена$(NC)"

