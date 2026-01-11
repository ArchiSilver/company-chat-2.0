# Company SuperApp

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?logo=typescript&logoColor=white)
![React Native](https://img.shields.io/badge/React_Native-0.74-61DAFB?logo=react&logoColor=white)
![Expo](https://img.shields.io/badge/Expo-51-000020?logo=expo&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?logo=prometheus&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-F46800?logo=grafana&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

Корпоративное мобильное супер-приложение для бизнеса.

**Demo:** _Coming soon_

---

## Стек технологий

### Backend

| Технология | Версия | Назначение |
|------------|--------|------------|
| Go | 1.21+ | Основной язык |
| Gin | 1.9 | HTTP фреймворк |
| PostgreSQL | 15 | Основная БД |
| Redis | 7 | Кэш, сессии, pub/sub |
| MinIO | — | S3-совместимое хранилище |
| WebSocket | Gorilla | Real-time коммуникация |
| JWT | golang-jwt | Аутентификация |
| bcrypt | — | Хеширование паролей |
| AES-256-GCM | — | Шифрование данных |
| gofpdf | — | Генерация PDF |
| slog | — | Структурированное логирование |
| OpenTelemetry | — | Distributed tracing |
| Prometheus | — | Метрики |

### Frontend

| Технология | Версия | Назначение |
|------------|--------|------------|
| React Native | 0.74 | Мобильный фреймворк |
| Expo | 51 | Платформа разработки |
| TypeScript | 5.0 | Типизация |
| Zustand | 4.5 | State management |
| React Navigation | 6 | Навигация |
| Axios | — | HTTP клиент |
| expo-secure-store | — | Безопасное хранилище |
| expo-local-authentication | — | Биометрия |
| expo-file-system | — | Работа с файлами |
| expo-sharing | — | Шаринг файлов |
| expo-notifications | — | Push-уведомления |

### DevOps & Monitoring

| Технология | Назначение |
|------------|------------|
| Docker | Контейнеризация |
| Docker Compose | Оркестрация |
| Prometheus | Сбор метрик |
| Grafana | Визуализация |
| Jaeger | Distributed tracing |
| Sentry | Error tracking |

---

## CDN / Установка

### Docker (рекомендуется)

```bash
git clone https://github.com/your-username/company-superapp.git
cd company-superapp
docker-compose up -d
```

### Ручная установка

```bash
# Backend
cd backend
go mod download
go run ./cmd/api

# Frontend
cd frontend
npm install
npm start
```

---

## Быстрый старт

Добавьте переменные окружения в `.env`:

```env
DATABASE_URL=postgres://admin:superpassword@localhost:5432/company_superapp?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-super-secret-jwt-key
ENCRYPTION_KEY=your-32-byte-secret-key-here!!
```

Запустите:

```bash
make up
```

API доступен на `http://localhost:8080`

---

## Стек технологий

| Компонент | Технология |
|-----------|------------|
| Backend | Go 1.21+, Gin, PostgreSQL 15, Redis 7 |
| Frontend | React Native, Expo 51, TypeScript, Zustand |
| Storage | MinIO (S3-compatible) |
| Monitoring | Prometheus, Grafana, Jaeger, Sentry |

---

## Модули

| Модуль | Описание | Статус |
|--------|----------|--------|
| Auth | JWT + bcrypt + Redis sessions | ✅ |
| Messenger | WebSocket real-time chat | ✅ |
| Tasks | Kanban board | ✅ |
| Finance | AES-256 encrypted salary | ✅ |
| Taxi | S3 receipt upload | ✅ |
| Notifications | FCM + Expo Push | ✅ |
| Search | PostgreSQL FTS | ✅ |
| RBAC | Role-based access | ✅ |
| Reports | PDF generation | ✅ |
| Monitoring | Prometheus + Sentry | ✅ |

---

## API Reference

### Аутентификация

```
POST /api/v1/auth/register    # Регистрация
POST /api/v1/auth/login       # Вход
POST /api/v1/auth/refresh     # Обновление токена
POST /api/v1/auth/logout      # Выход
```

### Мессенджер

```
GET  /api/v1/chats            # Список чатов
GET  /api/v1/chats/:id        # Сообщения чата
POST /api/v1/chats            # Создать чат
WS   /ws/connect              # WebSocket
```

### Задачи

```
GET    /api/v1/tasks          # Список
POST   /api/v1/tasks          # Создать
PUT    /api/v1/tasks/:id      # Обновить
DELETE /api/v1/tasks/:id      # Удалить
```

### Финансы (требуется роль admin/manager)

```
GET /api/v1/finance/salary    # Получить (+ biometric)
PUT /api/v1/finance/salary    # Обновить
```

### Такси

```
POST /api/v1/taxi/generate-upload-url   # Pre-signed URL
POST /api/v1/taxi/confirm-upload        # Подтвердить
GET  /api/v1/taxi/requests              # Список заявок
```

### Поиск и отчёты

```
GET /api/v1/search?q=query                        # Full-text search
GET /api/v1/reports/tasks?from=2026-01-01&to=...  # PDF отчёт
```

### Health & Metrics

```
GET /health         # Полный статус
GET /health/ready   # Kubernetes readiness
GET /health/live    # Kubernetes liveness
GET /metrics        # Prometheus
```

---

## Конфигурация

### Переменные окружения

| Переменная | Описание | Обязательно |
|------------|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string | ✅ |
| `REDIS_URL` | Redis connection string | ✅ |
| `JWT_SECRET` | Секрет для подписи JWT | ✅ |
| `ENCRYPTION_KEY` | 32-байтный ключ AES-256 | ✅ |
| `MINIO_ENDPOINT` | MinIO endpoint | ❌ |
| `MINIO_ACCESS_KEY` | MinIO access key | ❌ |
| `MINIO_SECRET_KEY` | MinIO secret key | ❌ |
| `SENTRY_DSN` | Sentry DSN для error tracking | ❌ |

### Роли пользователей

| Роль | Доступ к Finance | Описание |
|------|------------------|----------|
| `admin` | ✅ | Полный доступ |
| `manager` | ✅ | Доступ к финансам |
| `user` | ❌ | Базовый пользователь |

---

## Makefile

```bash
make up           # Запустить всё
make down         # Остановить
make build        # Пересобрать
make logs         # Все логи
make api-logs     # Логи API
make db-shell     # PostgreSQL CLI
make redis-cli    # Redis CLI
make db-backup    # Бэкап БД
make clean        # Удалить данные
```

---

## Структура проекта

```
company-superapp/
├── backend/
│   ├── cmd/api/main.go           # Entry point
│   ├── internal/
│   │   ├── config/               # Конфигурация
│   │   ├── domain/               # Доменные модели
│   │   ├── repository/postgres/  # Data Access Layer
│   │   ├── service/              # Бизнес-логика
│   │   ├── delivery/http/        # HTTP handlers
│   │   ├── delivery/websocket/   # WebSocket
│   │   ├── infrastructure/       # DB, Redis, S3
│   │   └── pkg/                  # Утилиты
│   └── migrations/               # SQL миграции
├── frontend/
│   ├── src/
│   │   ├── screens/              # Экраны
│   │   ├── components/           # UI компоненты
│   │   ├── store/                # Zustand stores
│   │   ├── hooks/                # Custom hooks
│   │   ├── api/                  # API client
│   │   └── navigation/           # React Navigation
│   └── app.json
├── monitoring/
│   ├── prometheus.yml
│   └── grafana/
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Сервисы

| Сервис | URL | Логин |
|--------|-----|-------|
| API | http://localhost:8080 | — |
| PostgreSQL | localhost:5432 | admin / superpassword |
| Redis | localhost:6379 | — |
| MinIO Console | http://localhost:9001 | minioadmin / minioadminpassword |
| Prometheus | http://localhost:9090 | — |
| Grafana | http://localhost:3000 | admin / admin |
| Jaeger | http://localhost:16686 | — |

---

## Безопасность

- **JWT** — access/refresh tokens (15min / 7d)
- **bcrypt** — хеширование паролей
- **AES-256-GCM** — шифрование финансовых данных
- **RBAC** — контроль доступа по ролям
- **Biometric** — биометрия для просмотра зарплаты
- **Pre-signed URLs** — безопасная загрузка в S3

---

## Тестирование

```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && npm test
```

---

## Roadmap

- [x] Аутентификация (JWT + Redis)
- [x] Мессенджер (WebSocket)
- [x] Таск-трекер (Kanban)
- [x] Финансы (AES-256)
- [x] Такси (S3/MinIO)
- [x] Push-уведомления
- [x] Глобальный поиск (FTS)
- [x] RBAC
- [x] PDF отчёты
- [x] Мониторинг
- [ ] Календарь событий
- [ ] Видеозвонки (WebRTC)
- [ ] Интеграция с 1С

---

## Contributing

1. Fork репозитория
2. Создайте branch (`git checkout -b feature/awesome`)
3. Commit (`git commit -m 'Add awesome feature'`)
4. Push (`git push origin feature/awesome`)
5. Откройте Pull Request

---

## License

MIT

---

**Free as in Open Source.** 🚀
