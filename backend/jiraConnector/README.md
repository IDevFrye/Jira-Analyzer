# Инструкция к взаимодействию с jira-connector во время разработки

## Пример файла конфигурации
```yaml
database:
  host: YOUR_POSTGRES_HOST
  port: YOUR_POSTGRES_PORT
  name: DB_NAME
  user: USER_NAME
  password: USER_PWD

jira-connector:
  url: https://issues.apache.org/jira
  thread_count: 10
  issue_in_one_request: 100
  max_sleep: 8000
  min_sleep: 50

server:
  port: ":8080"

log_file: "jiraconnector.log"
env: "local"
```

## Перед запуском

1. Убедитесь, что у вас установлен **Docker** и **Docker Compose**
2. Убедитсь, что у вас установлен Task [опционально] - для лёгкого старта приложения
3. Убедитесь, что добавили файл конфигурации по пути ./configs/сonfig.yaml
4. Клонируйте репозиторий и перейдите в папку микросервиса проекта.

```bash
git clone https://github.com/IDevFrye/Jira-Analyzer.git
cd Jira-Analyzer/jiraConnector
```

## Автоматическая сборка и запуск сервера

> самый простой способ для любой платформы
> убедитесь, что установлен [Task](https://taskfile.dev/#/installation)

```bash
# Собрать образы
task build

# Запустить проект
task up

# Остановить
task down

#Посмотреть все доступные команды task
task --list
```

## Ручной запуск сервера jiraConnector

1. Убедитесь в существовании базы данных
2. Настройте перемнную окружения CONFIG_PATH (путь к вашему файлу конфигурации)
3. Для запуска сервиса выполните команду из ./cmd/service
```bash
go run main.go
```
4. Через браузер\curl\postman проверьте работу выполнив один из доступных запросов

## Запросы
Подробные эндпоинты с описание параметров, тел запросов и ответов можно найти в папке ./docs

1. /api/v1/connector/projects - для получения проектов, доступных для загрузки.
Доступны параметры:
- limit: [int] - количество проектов на одной странице (limit > 0)
- page: [int] - номер страницы, с которой необходимо вернуть проекты (page > 0)
- search: [string] - параметр для фильтрации списка проектов. Будут возвращены только те проекты, имя или ключ которых содержат подстроку заданную в этом параметре без учета регистра.

Возвращает JSON, содержащий массив Projects (проекты на странице под номером page) и структуру PageInfo, которая содержит поле PageCount - общее количество страниц при данном параметре limit и search, CurrentPage - номер текущей страницы, ProjectsCount - общее количество проектов при данном параметре search. Значение limit по умолчанию = 20, значение page по умолчанию = 1

2. /api/v1/connector/updateProject?project=projectKey - Получает (или обновляет) все issues из проекта с ключом 'projectKey' и заносит в базу данных. Что будет происходить - загрузка или
обновление - зависит от того, был ли проект сохранен локально ранее.

*База данных обновляется только при запросе на update.

## Примеры запросов:

- GET http://localhost:8080/api/v1/connector/projects
```json
{
  "projects": [
    {
      "id": "12320120",
      "key": "AAR",
      "name": "aardvark",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12320120"
    },
    {
      "id": "12310505",
      "key": "ABDERA",
      "name": "Abdera",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12310505"
    },
    {
      "id": "12312121",
      "key": "ACCUMULO",
      "name": "Accumulo",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12312121"
    },
    {
      "id": "12310931",
      "key": "ACE",
      "name": "ACE",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12310931"
    },
    {
      "id": "12311200",
      "key": "ACL",
      "name": "ActiveCluster",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311200"
    },
    {
      "id": "12311201",
      "key": "AMQNET",
      "name": "ActiveMQ .Net",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311201"
    },
    {
      "id": "12311310",
      "key": "APLO",
      "name": "ActiveMQ Apollo (Retired)",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311310"
    },
    {
      "id": "12315920",
      "key": "ARTEMIS",
      "name": "ActiveMQ Artemis",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12315920"
    },
    {
      "id": "12311207",
      "key": "AMQCPP",
      "name": "ActiveMQ C++ Client",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311207"
    },
    {
      "id": "12311210",
      "key": "AMQ",
      "name": "ActiveMQ Classic",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311210"
    },
    {
      "id": "12320821",
      "key": "AMQCLI",
      "name": "ActiveMQ CLI Tools",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12320821"
    },
    {
      "id": "12315620",
      "key": "OPENWIRE",
      "name": "ActiveMQ OpenWire",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12315620"
    },
    {
      "id": "12325420",
      "key": "AMQWEBSITE",
      "name": "ActiveMQ Website",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12325420"
    },
    {
      "id": "12311204",
      "key": "BLAZE",
      "name": "ActiveRealTime (Retired)",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311204"
    },
    {
      "id": "12310060",
      "key": "ADDR",
      "name": "Addressing",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12310060"
    },
    {
      "id": "10730",
      "key": "AGILA",
      "name": "Agila",
      "self": "https://issues.apache.org/jira/rest/api/2/project/10730"
    },
    {
      "id": "12311302",
      "key": "AIRAVATA",
      "name": "Airavata",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311302"
    },
    {
      "id": "12311173",
      "key": "ALOIS",
      "name": "ALOIS",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311173"
    },
    {
      "id": "10101",
      "key": "ARMI",
      "name": "AltRMI",
      "self": "https://issues.apache.org/jira/rest/api/2/project/10101"
    },
    {
      "id": "12321521",
      "key": "AMATERASU",
      "name": "AMATERASU",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12321521"
    }
  ],
  "pageInfo": {
    "pageCount": 34,
    "currentPage": 1,
    "projectsCount": 671
  }
}
```
- GET http://localhost:8080/api/v1/connector/projects?limit=5&page=2&search=a
```json
{
  "projects": [
    {
      "id": "12311201",
      "key": "AMQNET",
      "name": "ActiveMQ .Net",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311201"
    },
    {
      "id": "12311310",
      "key": "APLO",
      "name": "ActiveMQ Apollo (Retired)",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311310"
    },
    {
      "id": "12315920",
      "key": "ARTEMIS",
      "name": "ActiveMQ Artemis",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12315920"
    },
    {
      "id": "12311207",
      "key": "AMQCPP",
      "name": "ActiveMQ C++ Client",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311207"
    },
    {
      "id": "12311210",
      "key": "AMQ",
      "name": "ActiveMQ Classic",
      "self": "https://issues.apache.org/jira/rest/api/2/project/12311210"
    }
  ],
  "pageInfo": {
    "pageCount": 89,
    "currentPage": 2,
    "projectsCount": 444
  }
}
```
- POST http://localhost:8080/api/v1/connectorupdateProject?project=AAR
```json
{
  "project": "AAR",
  "status": "updated"
}
```