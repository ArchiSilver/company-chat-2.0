package logger

import (
	"context"
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init инициализирует глобальный структурированный логгер
func Init(environment string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	if environment == "production" {
		// JSON-формат для продакшена (удобнее парсить в агрегаторах логов)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Текстовый формат для разработки (более читаемый)
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Get возвращает логгер по умолчанию
func Get() *slog.Logger {
	if defaultLogger == nil {
		Init("development")
	}
	return defaultLogger
}

// WithContext возвращает логгер с контекстом трассировки
func WithContext(ctx context.Context) *slog.Logger {
	logger := Get()

	// Извлекаем trace_id из контекста если присутствует
	if traceID := ctx.Value("trace_id"); traceID != nil {
		logger = logger.With("trace_id", traceID)
	}

	// Извлекаем user_id из контекста если присутствует
	if userID := ctx.Value("user_id"); userID != nil {
		logger = logger.With("user_id", userID)
	}

	return logger
}

// Info логирует информационное сообщение
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Error логирует сообщение об ошибке
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Warn логирует предупреждение
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Debug логирует отладочное сообщение
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// InfoContext логирует информационное сообщение с контекстом
func InfoContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).Info(msg, args...)
}

// ErrorContext logs an error message with context
func ErrorContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).Error(msg, args...)
}
