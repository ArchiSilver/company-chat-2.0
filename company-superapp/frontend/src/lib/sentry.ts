// Интеграция с Sentry для отслеживания ошибок
// Примечание: @sentry/react-native должен быть установлен: npx expo install @sentry/react-native

const SENTRY_DSN = process.env.EXPO_PUBLIC_SENTRY_DSN || '';

// Типы Sentry (упрощённые для случая, когда пакет не установлен)
interface SentryEvent {
    event_id?: string;
    message?: string;
    level?: string;
    [key: string]: unknown;
}

interface SentryHint {
    originalException?: Error | string;
    [key: string]: unknown;
}

interface SentryUser {
    id: string;
    email?: string;
    role?: string;
}

// Заглушка Sentry для случая, когда пакет не установлен
const SentryMock = {
    init: (_config: Record<string, unknown>) => {
        console.log('Sentry mock initialized');
    },
    setUser: (_user: SentryUser | null) => { },
    captureException: (error: unknown) => {
        console.error('Sentry capture:', error);
    },
    captureMessage: (message: string) => {
        console.log('Sentry message:', message);
    },
    addBreadcrumb: (_breadcrumb: Record<string, unknown>) => { },
};

// Попытка импорта настоящего Sentry, иначе используем заглушку
let Sentry = SentryMock;
try {
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    Sentry = require('@sentry/react-native');
} catch {
    console.warn('Sentry not installed, using mock');
}

export const initSentry = () => {
    if (!SENTRY_DSN) {
        console.warn('Sentry DSN not configured. Error tracking disabled.');
        return;
    }

    Sentry.init({
        dsn: SENTRY_DSN,

        // Мониторинг производительности
        tracesSampleRate: 1.0, // Захват 100% транзакций в режиме разработки

        // Окружение
        environment: process.env.EXPO_PUBLIC_APP_ENV || 'development',

        // Отслеживание релизов
        release: process.env.EXPO_PUBLIC_APP_VERSION || '1.0.0',

        // Включение автоматической инструментации
        enableAutoSessionTracking: true,
        sessionTrackingIntervalMillis: 30000,

        // Добавление стек-трейсов ко всем сообщениям
        attachStacktrace: true,

        // Игнорирование определённых ошибок
        ignoreErrors: [
            'Network request failed',
            'Failed to fetch',
        ],

        // Хук beforeSend - можно модифицировать или отбросить события
        beforeSend(event: SentryEvent, _hint: SentryHint) {
            return event;
        },
    });
};

// Установка контекста пользователя после входа
export const setUserContext = (userId: string, email?: string, role?: string) => {
    Sentry.setUser({
        id: userId,
        email,
        role,
    });
};

// Очистка контекста пользователя при выходе
export const clearUserContext = () => {
    Sentry.setUser(null);
};

// Add breadcrumb for custom events
export const addBreadcrumb = (message: string, category: string, data?: Record<string, unknown>) => {
    Sentry.addBreadcrumb({
        message,
        category,
        data,
        level: 'info',
    });
};

// Capture custom exception
export const captureException = (error: Error, _context?: Record<string, unknown>) => {
    Sentry.captureException(error);
};

// Capture custom message
export const captureMessage = (message: string) => {
    Sentry.captureMessage(message);
};

// Start a performance transaction (no-op in mock)
export const startTransaction = (_name: string, _op: string) => {
    return {};
};

// Export Sentry for direct access if needed
export { Sentry };
