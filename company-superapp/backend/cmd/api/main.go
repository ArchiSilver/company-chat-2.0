package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/yourname/company-superapp/internal/config"
	"github.com/yourname/company-superapp/internal/delivery/http"
	"github.com/yourname/company-superapp/internal/delivery/websocket"
	"github.com/yourname/company-superapp/internal/infrastructure/migrations"
	"github.com/yourname/company-superapp/internal/pkg/encryption"
	"github.com/yourname/company-superapp/internal/pkg/fcm"
	"github.com/yourname/company-superapp/internal/pkg/s3"
	"github.com/yourname/company-superapp/internal/repository/postgres"
	"github.com/yourname/company-superapp/internal/service"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем структурированное логирование
	initLogger(cfg.Server.Environment)
	slog.Info("Запуск Company SuperApp API",
		"environment", cfg.Server.Environment,
		"port", cfg.Server.Port,
	)

	// Подключение к базе данных
	db, err := connectDB(cfg.Database)
	if err != nil {
		slog.Error("Не удалось подключиться к базе данных", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Подключено к PostgreSQL", "host", cfg.Database.Host)

	// Выполняем миграции
	migrationRunner := migrations.NewRunner(db.DB, "./migrations")
	if err := migrationRunner.Run(); err != nil {
		slog.Error("Не удалось выполнить миграции", "error", err)
		os.Exit(1)
	}
	slog.Info("Миграции базы данных применены успешно")

	// Подключение к Redis
	redisClient, err := connectRedis(cfg.Redis)
	if err != nil {
		slog.Error("Не удалось подключиться к Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	slog.Info("Подключено к Redis", "host", cfg.Redis.Host)

	// Сервис шифрования
	encryptionService, err := encryption.NewEncryptionService()
	if err != nil {
		slog.Error("Не удалось создать сервис шифрования", "error", err)
		os.Exit(1)
	}

	// Клиент MinIO
	minioClient, err := s3.NewMinioClient()
	if err != nil {
		slog.Warn("Не удалось создать MinIO клиент, S3 функции отключены", "error", err)
	}

	// Клиент FCM
	fcmClient, err := fcm.NewFCMClient()
	if err != nil {
		slog.Warn("Не удалось создать FCM клиент, push-уведомления отключены", "error", err)
	}

	// Настройка Onion Architecture — Репозитории
	userRepo := postgres.NewUserRepository(db)
	chatRepo := postgres.NewChatRepository(db)
	messageRepo := postgres.NewMessageRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	salaryRepo := postgres.NewSalaryRepository(db)
	taxiRequestRepo := postgres.NewTaxiRequestRepository(db)
	pushTokenRepo := postgres.NewPushTokenRepository(db)
	searchRepo := postgres.NewSearchRepository(db)

	// Настройка Onion Architecture — Сервисы
	authService := service.NewAuthService(userRepo, redisClient, cfg.JWT.Secret)
	notificationService := service.NewNotificationService(pushTokenRepo, fcmClient)
	chatService := service.NewChatService(chatRepo, messageRepo)
	taskService := service.NewTaskService(taskRepo, messageRepo)
	salaryService := service.NewSalaryService(salaryRepo, encryptionService)
	taxiService := service.NewTaxiService(taxiRequestRepo, minioClient)
	searchService := service.NewGlobalSearchService(searchRepo)
	reportService := service.NewReportService(taskRepo)

	// WebSocket Hub для real-time соединений
	hub := websocket.NewHub(redisClient)
	go hub.Run()

	// Настройка HTTP обработчиков
	authHandler := http.NewAuthHandler(authService)
	chatHandler := http.NewChatHandler(chatService, hub)
	taskHandler := http.NewTaskHandler(taskService)
	financeHandler := http.NewFinanceHandler(salaryService)
	taxiHandler := http.NewTaxiHandler(taxiService)
	notificationHandler := http.NewNotificationHandler(notificationService)
	searchHandler := http.NewSearchHandler(searchService)
	reportHandler := http.NewReportHandler(reportService)
	healthHandler := http.NewHealthHandler(db, redisClient)

	// Настройка Gin Router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Применяем middleware для мониторинга
	router.Use(http.TracingMiddleware())
	router.Use(http.PrometheusMiddleware())

	// Health и метрики (без авторизации)
	healthHandler.RegisterRoutes(router)

	// Маршруты API v1
	apiV1 := router.Group("/api/v1")
	authHandler.RegisterRoutes(apiV1)
	chatHandler.RegisterRoutes(apiV1)
	taskHandler.RegisterRoutes(apiV1)
	financeHandler.RegisterRoutes(apiV1)
	taxiHandler.RegisterRoutes(apiV1)
	notificationHandler.RegisterRoutes(apiV1)
	searchHandler.RegisterRoutes(apiV1)
	reportHandler.RegisterRoutes(apiV1)

	// Graceful shutdown — плавное завершение
	go func() {
		slog.Info("Запуск HTTP сервера", "port", cfg.Server.Port)
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			slog.Error("Не удалось запустить сервер", "error", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Завершение работы сервера...")
	time.Sleep(1 * time.Second) // Даём время завершить текущие запросы
	slog.Info("Сервер остановлен")
}

// connectDB подключается к PostgreSQL
func connectDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// connectRedis подключается к Redis
func connectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return client, nil
}

// initLogger инициализирует структурированный логгер
func initLogger(env string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: env != "production",
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
