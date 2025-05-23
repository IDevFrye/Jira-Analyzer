# Запуск приложения

Описание запуска приложения через task, make, docker compose.
Перейдите в папку Jira-Analyzer/deployment и все дальнейшие шаги выполняйте оттуда.

## Taskfile
> Удобно, кроссплатформенно, простой синтаксис (требует предварительную установку).
```bash
task up                 # Запустить всё в Docker
task unit-test          # Юнит тесты
task integration-test   # Интеграционные тесты
task down               # Остановить контейнеры
task clean              # Полная очистка
task logs               # Просмотр логов
task logs-connector     # Просмотр логов коннетора
task db             	# Консоль базы данных
task all                # Полный пайплайн сборки-тестирования-запуска
```
---
## Makefile
> Для Unix-систем. На Windows работает через WSL или Git Bash
```bash
make up                 # Запустить всё
make unit-test          # Юнит тесты
make integration-test   # Интеграционные тесты
make down               # Остановить
make clean              # Удалить контейнеры и тома
make logs               # Логи всех сервисов
make logs-connector     # Просмотр логов коннетора
make db             	# Консоль базы данных
make all                # Полный пайплайн сборки-тестирования-запуска
```
---

## Команды docker compose
> Прямой и универсальный способ — просто команды Docker.

```bash
# Сборка и запуск
docker compose up --build -d

# Запуск юнит-тестов эндпоинтов
docker compose run --rm backend go test ./... -v
# Запуск юнит-тестов коннектора
docker compose run --rm jiraconnector go test ./... -v

# Запуск интеграционных тестов для коннектора
docker compose run --rm jiraconnector go test ./tests/integration/... -v -tags=integration
# Запуск интеграционных тестов для всей системы
docker compose run --rm backend-test

# Остановка и удаление
docker compose down

# Полная очистка
docker compose down -v --remove-orphans

# Логи
docker compose logs -f
docker exec -it deployment-jiraconnector-1 cat log/jiraconnector.log

# Состояние базы данных
docker exec -it deployment-postgres-1 psql -U postgres -d testdb
```
---