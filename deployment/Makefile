.PHONY: build up test-unit test-integration down clean logs help

build:
	docker compose build

up:
	docker compose up -d

unit-test:
	docker compose run --rm backend go test ./... -v

integration-test:
	docker compose run --rm backend go test -tags=integration ./... -v

down:
	docker compose down

clean:
	docker compose down -v --remove-orphans
	docker system prune -f

logs:
	docker compose logs -f

logs-connector:
	docker exec -it deployment-jiraconnector-1 cat log/jiraconnector.log

all:
	make build
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
