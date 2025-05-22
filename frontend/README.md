# Jira Analyzer — Frontend (React + TypeScript + SCSS)

## Описание проекта
Фронтенд-модуль для системы аналитики Jira-проектов. Позволяет просматривать проекты, собирать и визуализировать статистику, выполнять аналитику и сравнение.

- **Технологии**: React, TypeScript, SCSS

## Функциональность

- Просмотр всех доступных проектов (с пагинацией и фильтрацией)
- Добавление проекта себе и их удаление
- Просмотр "Моих проектов" с сухой статистикой (все метрики)
- Отображение графиков (Chart.js)
- Сравнение 2-3 проектов по метрикам и аналитике (в том числе графикам)
- Адаптивное меню: Проекты, Мои проекты, Сравнение

## Структура проекта

```
frontend/
├── public/
│   └── index.html
│   └── favicon.png
├── src/
│   ├── assets/             # Медиафайлы
│   ├── components/         # Все UI-компоненты
│   ├── config/             # Описание эндпоинтов и конфиг диаграмм
│   ├── pages/              # Основные страницы (Projects, MyProjects, Compare)
│   ├── styles/             # Базовые SCSS-стили
│   ├── types/              # Базовые модели данных (интерфейсы)
│   ├── router.tsx          # React Router DOM
│   ├── index.tsx
│   └── JiraAnalyzerApp.tsx
├── .eslintrc.json
├── .prettierrc
├── tsconfig.json
├── webpack.config.js
├── package.json
└── README.md
```

## Как запустить

```bash
npm install         # или npm install
npm run start       # запускает dev-сервер на localhost:3000
```

## Сборка и деплой

```bash
npm run build       # => Создаётся папка dist/
```

## Подключение к бэкенду

1. API находится на `http://localhost:8080`
2. Все fetch-запросы делают вызовы к REST API, например:
   ```ts
   await fetch(`${apiHost}/api/v1/projects`)
   ```
3. Настройка эндпоинтов осуществляется через ```src/config/config.ts```

---

Разработка под React 18 + TypeScript. Все компоненты типизированы, оптимизированы под расширение. Собирается в SPA, развертывается на любой статический сервер или в Tomcat через WAR-like структуру (dist).
