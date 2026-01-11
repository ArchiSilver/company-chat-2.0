# 🏢 Company SuperApp

Корпоративное мобильное супер-приложение с мессенджером, таск-трекером, финансами и такси.

## 🏗️ Архитектура

### Backend (Go)
- **Onion Architecture**: Domain → Repository → Service → Handler
- **Go 1.21+** с Gin framework
- **PostgreSQL 15** с полнотекстовым поиском
- **Redis 7** для кэширования и сессий
- **MinIO** для S3-совместимого хранилища
- **WebSocket** для real-time мессенджера

### Frontend (React Native)
- **Expo SDK 51+** с TypeScript
- **Atomic Design** структура компонентов
- **Zustand** для state management
- **NativeWind** (или StyleSheet) для стилей

### Мониторинг
- **Prometheus** — метрики
- **Grafana** — дашборды
- **Jaeger** — distributed tracing
- **Sentry** — error tracking

---

## 🚀 Быстрый старт

### Требования
- Docker & Docker Compose
- Go 1.21+ (для локальной разработки)
- Node.js 18+ (для frontend)
- Expo CLI

### 1. Запуск всех сервисов

\`\`\`bash
cd company-superapp

# Запуск через Docker Compose
make up

# Или напрямую:
docker-compose up -d --build
\`\`\`

### 2. Проверка сервисов

| Сервис     | URL                          | Credentials           |
|------------|------------------------------|-----------------------|
| API        | http://localhost:8080        | -                     |
| PostgreSQL | localhost:5432               | admin / superpassword |
| Redis      | localhost:6379               | -                     |
| MinIO      | http://localhost:9001        | minioadmin / minioadminpassword |
| Prometheus | http://localhost:9090        | -                     |
| Grafana    | http://localhost:3000        | admin / admin         |
| Jaeger UI  | http://localhost:16686       | -                     |

### 3. Запуск Frontend

\`\`\`bash
cd frontend
npm install
npm start
\`\`\`

Отсканируйте QR-код в Expo Go или нажмите \`i\` для iOS / \`a\` для Android.

---

## 📁 Структура проекта

\`\`\`
company-superapp/
├── backend/
│   ├── cmd/api/              # Entry point
│   ├── internal/
│   │   ├── config/           # Configuration
│   │   ├── domain/           # Domain entities
│   │   ├── repository/       # Data access layer
│   │   ├── service/          # Business logic
│   │   ├── delivery/         # HTTP handlers
│   │   └── infrastructure/   # DB, Redis, S3
│   ├── migrations/           # SQL migrations
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/              # API client
│   │   ├── components/       # UI components
│   │   ├── screens/          # Screens
│   │   ├── navigation/       # Navigation
│   │   ├── store/            # Zustand stores
│   │   ├── hooks/            # Custom hooks
│   │   └── lib/              # Utilities
│   └── app.json
├── monitoring/
│   ├── prometheus.yml
│   └── grafana/
├── docker-compose.yml
├── Makefile
└── README.md
\`\`\`

---

## 🔌 API Endpoints

### Auth
\`\`\`
POST /api/v1/auth/register   - Регистрация
POST /api/v1/auth/login      - Вход
POST /api/v1/auth/refresh    - Обновление токена
POST /api/v1/auth/logout     - Выход
\`\`\`

### Chats
\`\`\`
GET  /api/v1/chats           - Список чатов
GET  /api/v1/chats/:id       - Чат с сообщениями
POST /api/v1/chats           - Создать чат
WS   /ws/connect             - WebSocket соединение
\`\`\`

### Tasks
\`\`\`
GET    /api/v1/tasks         - Список задач
POST   /api/v1/tasks         - Создать задачу
PUT    /api/v1/tasks/:id     - Обновить задачу
DELETE /api/v1/tasks/:id     - Удалить задачу
\`\`\`

### Finance (RBAC protected)
\`\`\`
GET /api/v1/finance/salary   - Получить зарплату (biometric)
PUT /api/v1/finance/salary   - Обновить зарплату (admin/manager)
\`\`\`

### Taxi
\`\`\`
POST /api/v1/taxi/generate-upload-url  - Получить URL для загрузки чека
POST /api/v1/taxi/confirm-upload       - Подтвердить загрузку
GET  /api/v1/taxi/requests             - Список заявок
\`\`\`

### Search
\`\`\`
GET /api/v1/search?q=query   - Глобальный поиск (FTS)
\`\`\`

### Reports
\`\`\`
GET /api/v1/reports/tasks?from=2026-01-01&to=2026-01-31 - PDF отчёт
\`\`\`

### Health
\`\`\`
GET /health        - Health check
GET /health/ready  - Readiness probe
GET /health/live   - Liveness probe
GET /metrics       - Prometheus metrics
\`\`\`

---

## 🛠️ Makefile команды

\`\`\`bash
make up          # Запустить все сервисы
make down        # Остановить сервисы
make build       # Пересобрать и запустить
make logs        # Логи всех сервисов
make api-logs    # Логи API
make db-shell    # PostgreSQL shell
make redis-cli   # Redis CLI
make clean       # Удалить все данные
make db-backup   # Бэкап базы данных
\`\`\`

---

## 🔐 Безопасность

- **JWT** с access/refresh токенами
- **bcrypt** для хешинга паролей
- **AES-256** для шифрования финансовых данных
- **RBAC** (Role-Based Access Control)
- **Biometric auth** для просмотра зарплаты
- **Pre-signed URLs** для S3 загрузок

---

## 📊 Мониторинг

### Prometheus метрики
- `http_requests_total` — количество запросов
- `http_request_duration_seconds` — latency
- `go_goroutines` — активные горутины

### Grafana дашборды
После первого запуска:
1. Откройте http://localhost:3000
2. Войдите (admin/admin)
3. Добавьте Prometheus datasource: http://prometheus:9090

---

## 🧪 Тестирование

\`\`\`bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test
\`\`\`

---

## 📝 License

MIT
