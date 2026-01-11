// Модуль аналитики для отслеживания событий пользователя
// Может быть подключён к PostHog, Amplitude или любому провайдеру аналитики

type EventProperties = Record<string, string | number | boolean | null>;

interface AnalyticsProvider {
    track: (event: string, properties?: EventProperties) => void;
    identify: (userId: string, traits?: EventProperties) => void;
    reset: () => void;
}

// Простая консольная аналитика для разработки
const consoleAnalytics: AnalyticsProvider = {
    track: (event, properties) => {
        if (__DEV__) {
            console.log('[Analytics] Track:', event, properties);
        }
    },
    identify: (userId, traits) => {
        if (__DEV__) {
            console.log('[Analytics] Identify:', userId, traits);
        }
    },
    reset: () => {
        if (__DEV__) {
            console.log('[Analytics] Reset');
        }
    },
};

// Текущий провайдер (можно заменить на PostHog, Amplitude и т.д.)
let provider: AnalyticsProvider = consoleAnalytics;

export const Analytics = {
    // Инициализация аналитики с провайдером
    init: (customProvider?: AnalyticsProvider) => {
        if (customProvider) {
            provider = customProvider;
        }
    },

    // Отслеживание пользовательского события
    track: (event: string, properties?: EventProperties) => {
        provider.track(event, {
            ...properties,
            timestamp: Date.now(),
        });
    },

    // Идентификация пользователя
    identify: (userId: string, traits?: EventProperties) => {
        provider.identify(userId, traits);
    },

    // Сброс (при выходе)
    reset: () => {
        provider.reset();
    },

    // Предопределённые события
    events: {
        // События авторизации
        userRegistered: (userId: string) => {
            Analytics.track('user_registered', { user_id: userId });
        },
        userLoggedIn: (userId: string) => {
            Analytics.track('user_logged_in', { user_id: userId });
        },
        userLoggedOut: () => {
            Analytics.track('user_logged_out');
        },

        // События мессенджера
        messageSent: (chatId: string) => {
            Analytics.track('message_sent', { chat_id: chatId });
        },
        chatOpened: (chatId: string) => {
            Analytics.track('chat_opened', { chat_id: chatId });
        },

        // События задач
        taskCreated: (taskId: string) => {
            Analytics.track('task_created', { task_id: taskId });
        },
        taskStatusChanged: (taskId: string, status: string) => {
            Analytics.track('task_status_changed', { task_id: taskId, status });
        },

        // События финансов
        salaryViewed: () => {
            Analytics.track('salary_viewed');
        },

        // События такси
        taxiRequestCreated: (requestId: string) => {
            Analytics.track('taxi_request_created', { request_id: requestId });
        },

        // События отчётов
        reportGenerated: (reportType: string, dateRange: string) => {
            Analytics.track('report_generated', { report_type: reportType, date_range: dateRange });
        },

        // События поиска
        searchPerformed: (query: string, resultsCount: number) => {
            Analytics.track('search_performed', { query, results_count: resultsCount });
        },

        // События навигации
        screenViewed: (screenName: string) => {
            Analytics.track('screen_viewed', { screen_name: screenName });
        },
    },
};

export default Analytics;
