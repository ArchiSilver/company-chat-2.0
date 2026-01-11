import { Ionicons } from '@expo/vector-icons';
import React, { useState } from 'react';
import {
    Alert,
    ScrollView,
    StyleSheet,
    Text,
    TouchableOpacity,
    View,
} from 'react-native';
import { Analytics } from '../lib/analytics';
import { useAuthStore } from '../store/authStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    success: '#4CAF50',
    warning: '#FF9800',
    error: '#F44336',
};

interface DebugButtonProps {
    title: string;
    subtitle: string;
    icon: keyof typeof Ionicons.glyphMap;
    color: string;
    onPress: () => void;
}

const DebugButton: React.FC<DebugButtonProps> = ({ title, subtitle, icon, color, onPress }) => (
    <TouchableOpacity style={styles.debugButton} onPress={onPress}>
        <View style={[styles.iconContainer, { backgroundColor: color + '20' }]}>
            <Ionicons name={icon} size={24} color={color} />
        </View>
        <View style={styles.buttonTextContainer}>
            <Text style={styles.buttonTitle}>{title}</Text>
            <Text style={styles.buttonSubtitle}>{subtitle}</Text>
        </View>
        <Ionicons name="chevron-forward" size={20} color={COLORS.textSecondary} />
    </TouchableOpacity>
);

export const DebugScreen: React.FC = () => {
    const { role, accessToken } = useAuthStore();
    const [testResults, setTestResults] = useState<string[]>([]);

    const addResult = (message: string) => {
        setTestResults(prev => [...prev, `${new Date().toLocaleTimeString()}: ${message}`]);
    };

    const triggerTestError = () => {
        try {
            throw new Error('Sentry Test Error - This is a test exception');
        } catch (error) {
            // В продакшене с Sentry: Sentry.captureException(error);
            addResult('Test error triggered (would be sent to Sentry)');
            Alert.alert('Test Error', 'Error captured! Check Sentry dashboard.');
        }
    };

    const triggerUnhandledError = () => {
        Alert.alert(
            'Warning',
            'This will trigger an unhandled error and may crash the app. Continue?',
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Continue',
                    style: 'destructive',
                    onPress: () => {
                        // Это вызовет необработанную ошибку
                        // @ts-ignore
                        undefined.foo();
                    },
                },
            ]
        );
    };

    const testAnalyticsEvent = () => {
        Analytics.track('debug_test_event', {
            test_property: 'test_value',
            timestamp: Date.now(),
        });
        addResult('Analytics event sent: debug_test_event');
    };

    const testPerformanceSpan = () => {
        const startTime = Date.now();
        // Имитация выполнения работы
        setTimeout(() => {
            const duration = Date.now() - startTime;
            addResult(`Performance span completed: ${duration}ms`);
        }, 500);
    };

    const checkHealthEndpoint = async () => {
        try {
            const API_BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';
            const response = await fetch(`${API_BASE_URL}/health`);
            const data = await response.json();
            addResult(`Health check: ${data.status} (uptime: ${data.uptime})`);
            Alert.alert('Health Check', JSON.stringify(data, null, 2));
        } catch (error: any) {
            addResult(`Health check failed: ${error.message}`);
            Alert.alert('Error', 'Failed to check health endpoint');
        }
    };

    const clearResults = () => {
        setTestResults([]);
    };

    return (
        <ScrollView style={styles.container}>
            {/* Заголовок */}
            <View style={styles.header}>
                <Text style={styles.headerTitle}>Debug & Monitoring</Text>
                <Text style={styles.headerSubtitle}>Test error tracking and analytics</Text>
            </View>

            {/* Карточка информации о пользователе */}
            <View style={styles.card}>
                <Text style={styles.cardTitle}>Current Session</Text>
                <View style={styles.infoRow}>
                    <Text style={styles.infoLabel}>Role:</Text>
                    <Text style={styles.infoValue}>{role || 'Not logged in'}</Text>
                </View>
                <View style={styles.infoRow}>
                    <Text style={styles.infoLabel}>Token:</Text>
                    <Text style={styles.infoValue} numberOfLines={1}>
                        {accessToken ? `${accessToken.substring(0, 20)}...` : 'None'}
                    </Text>
                </View>
            </View>

            {/* Тестирование ошибок */}
            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Error Tracking</Text>
                <DebugButton
                    title="Trigger Test Error"
                    subtitle="Capture a handled exception"
                    icon="bug"
                    color={COLORS.warning}
                    onPress={triggerTestError}
                />
                <DebugButton
                    title="Trigger Unhandled Error"
                    subtitle="Simulate app crash (careful!)"
                    icon="alert-circle"
                    color={COLORS.error}
                    onPress={triggerUnhandledError}
                />
            </View>

            {/* Тестирование аналитики */}
            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Analytics</Text>
                <DebugButton
                    title="Send Test Event"
                    subtitle="Track a custom analytics event"
                    icon="analytics"
                    color={COLORS.primary}
                    onPress={testAnalyticsEvent}
                />
                <DebugButton
                    title="Test Performance Span"
                    subtitle="Measure a timed operation"
                    icon="speedometer"
                    color={COLORS.success}
                    onPress={testPerformanceSpan}
                />
            </View>

            {/* Проверка бэкенда */}
            <View style={styles.section}>
                <Text style={styles.sectionTitle}>Backend Health</Text>
                <DebugButton
                    title="Check Health Endpoint"
                    subtitle="GET /health"
                    icon="pulse"
                    color={COLORS.success}
                    onPress={checkHealthEndpoint}
                />
            </View>

            {/* Лог результатов тестирования */}
            {testResults.length > 0 && (
                <View style={styles.resultsCard}>
                    <View style={styles.resultsHeader}>
                        <Text style={styles.cardTitle}>Test Log</Text>
                        <TouchableOpacity onPress={clearResults}>
                            <Text style={styles.clearButton}>Clear</Text>
                        </TouchableOpacity>
                    </View>
                    {testResults.map((result, index) => (
                        <Text key={index} style={styles.resultText}>
                            {result}
                        </Text>
                    ))}
                </View>
            )}

            <View style={styles.bottomPadding} />
        </ScrollView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        paddingHorizontal: 20,
        paddingTop: 60,
        paddingBottom: 20,
        backgroundColor: COLORS.surface,
    },
    headerTitle: {
        fontSize: 28,
        fontWeight: '700',
        color: COLORS.text,
        marginBottom: 4,
    },
    headerSubtitle: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    card: {
        backgroundColor: COLORS.surface,
        margin: 16,
        borderRadius: 12,
        padding: 16,
    },
    cardTitle: {
        fontSize: 16,
        fontWeight: '600',
        color: COLORS.text,
        marginBottom: 12,
    },
    infoRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        paddingVertical: 8,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.background,
    },
    infoLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    infoValue: {
        fontSize: 14,
        color: COLORS.text,
        fontWeight: '500',
        flex: 1,
        textAlign: 'right',
    },
    section: {
        marginTop: 8,
        marginHorizontal: 16,
    },
    sectionTitle: {
        fontSize: 14,
        fontWeight: '600',
        color: COLORS.textSecondary,
        textTransform: 'uppercase',
        marginBottom: 8,
        marginLeft: 4,
    },
    debugButton: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: COLORS.surface,
        borderRadius: 12,
        padding: 16,
        marginBottom: 8,
    },
    iconContainer: {
        width: 44,
        height: 44,
        borderRadius: 22,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: 12,
    },
    buttonTextContainer: {
        flex: 1,
    },
    buttonTitle: {
        fontSize: 16,
        fontWeight: '600',
        color: COLORS.text,
        marginBottom: 2,
    },
    buttonSubtitle: {
        fontSize: 13,
        color: COLORS.textSecondary,
    },
    resultsCard: {
        backgroundColor: COLORS.surface,
        margin: 16,
        borderRadius: 12,
        padding: 16,
    },
    resultsHeader: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: 12,
    },
    clearButton: {
        color: COLORS.primary,
        fontSize: 14,
        fontWeight: '600',
    },
    resultText: {
        fontSize: 12,
        color: COLORS.textSecondary,
        fontFamily: 'monospace',
        paddingVertical: 4,
    },
    bottomPadding: {
        height: 40,
    },
});

export default DebugScreen;
