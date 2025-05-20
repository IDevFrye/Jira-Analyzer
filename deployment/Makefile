.PHONY: build up unit-test integration-test down clean logs help

build:
	# Сборка приложения
	docker compose build

up:
	# Запуск приложения
	docker compose up -d

unit-test:
	# Запуск юнит-тестов эндпоинтов
	# docker compose run --rm backend go test ./... -v

	# Запуск юнит-тестов коннектора
	# docker compose run --rm jiraconnector go test ./... -v

integration-test:
	# Запуск интеграционных тестов для коннектора
	docker compose run --rm jiraconnector go test ./tests/integration/... -v -tags=integration

	# Запуск интеграционных тестов для всей системы
	docker compose run --rm backend-test

down:
	# Остановка приложения
	docker compose down

clean:
	# Удаление контейнеров и томов приложения
	docker compose down -v --remove-orphans

logs:
	# Просмотр логов сервисов (кроме коннетора)
	docker compose logs -f

logs-connector:
	# Просмотр логов коннектора
	docker exec -it deployment-jiraconnector-1 cat log/jiraconnector.log

db:
	# Подключение к просмотру базы данных в консоли
	docker exec -it deployment-postgres-1 psql -U postgres -d testdb

all:
	# Запуск полного пайплайна
	make build
	make unit-test
	make integration-test
	make up

help:
	@echo "Makefile команды:"
	@echo "  make build             — Сборка контейнеров"
	@echo "  make up                — Запуск приложения"
	@echo "  make unit-test         — Юнит тесты"
	@echo "  make integration-test  — Интеграционные тесты"
	@echo "  make down              — Остановка приложения"
	@echo "  make clean             — Удаление контейнеров и мусора"
	@echo "  make logs              — Логи всех сервисов кроме коннектора"
	@echo "  make logs-connector    — Логи всех коннектора"
	@echo "  make all              	— Запуск полного пайплайна"
	@echo "  make db             	— Консоль базы данных"
