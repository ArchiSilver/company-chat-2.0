# Changelog

Все значимые изменения проекта документируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/),
проект следует [Semantic Versioning](https://semver.org/lang/ru/).

## [Unreleased]

### Планируется
- Календарь событий
- Видеозвонки (WebRTC)
- Интеграция с 1С

---

## [1.0.0] - 2026-01-11

### Добавлено

#### Аутентификация
- JWT access/refresh tokens (15min / 7 дней)
- Хеширование паролей bcrypt
- Redis сессии с автоматической инвалидацией
- Регистрация, вход, выход, обновление токена

#### Мессенджер
- Real-time чат на WebSocket
- Личные и групповые чаты
- Статусы сообщений (отправлено, доставлено, прочитано)
- Индикатор "печатает..."
- Оптимистичные обновления UI

#### Таск-трекер
- Kanban-доска с drag & drop
- Статусы: todo, in_progress, review, done
- Назначение исполнителей
- Дедлайны и приоритеты

#### Финансы
- Шифрование зарплат AES-256-GCM
- Биометрическая аутентификация для просмотра
- RBAC: доступ только admin/manager

#### Такси
- Загрузка чеков через pre-signed URLs
- Интеграция с MinIO (S3-compatible)
- Статусы заявок: pending, approved, rejected

#### Уведомления
- Push через FCM + Expo Push API
- Категории: message, task, finance, system
- Badges и звуки

#### Поиск
- Full-text search PostgreSQL (tsvector + GIN)
- Поиск по сообщениям, задачам, пользователям
- Ранжирование результатов

#### RBAC
- Роли: admin, manager, user
- Middleware проверки доступа
- Гибкая система разрешений

#### Отчёты
- Генерация PDF отчётов
- Экспорт задач за период
- Статистика активности

#### Мониторинг
- Prometheus метрики
- Grafana дашборды
- Jaeger трейсинг
- Sentry error tracking
- Health checks (ready/live)

### Технологии
- **Backend:** Go 1.21, Gin, PostgreSQL 15, Redis 7
- **Frontend:** React Native 0.74, Expo 51, TypeScript 5.0
- **DevOps:** Docker Compose, Prometheus, Grafana

---

## [0.1.0] - 2025-12-01

### Добавлено
- Инициализация проекта
- Базовая структура backend (Onion Architecture)
- Базовая структура frontend (Expo)
- Docker Compose окружение
- Makefile для команд разработки

---

[Unreleased]: https://github.com/ArchiSilver/company-chat-2.0/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/ArchiSilver/company-chat-2.0/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/ArchiSilver/company-chat-2.0/releases/tag/v0.1.0
