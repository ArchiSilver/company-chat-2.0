import { Ionicons } from '@expo/vector-icons';
import * as Sharing from 'expo-sharing';
import React, { useState } from 'react';
import {
    ActivityIndicator,
    Alert,
    Platform,
    StyleSheet,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from 'react-native';
import { useAuthStore } from '../store/authStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
    success: '#4CAF50',
};

const API_BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const ReportsScreen: React.FC = () => {
    const { accessToken } = useAuthStore();

    // Состояние дат (формат: YYYY-MM-DD)
    const [fromDate, setFromDate] = useState(() => {
        const date = new Date();
        date.setDate(date.getDate() - 7); // По умолчанию: последние 7 дней
        return date.toISOString().split('T')[0];
    });

    const [toDate, setToDate] = useState(() => {
        return new Date().toISOString().split('T')[0];
    });

    const [isLoading, setIsLoading] = useState(false);
    const [lastGeneratedFile, setLastGeneratedFile] = useState<string | null>(null);

    const generateReport = async () => {
        if (!fromDate || !toDate) {
            Alert.alert('Ошибка', 'Выберите период для отчёта');
            return;
        }

        if (new Date(fromDate) > new Date(toDate)) {
            Alert.alert('Ошибка', 'Начальная дата не может быть позже конечной');
            return;
        }

        setIsLoading(true);
        setLastGeneratedFile(null);

        try {
            // Запрос PDF из API
            const response = await fetch(
                `${API_BASE_URL}/reports/tasks?from=${fromDate}&to=${toDate}`,
                {
                    method: 'GET',
                    headers: {
                        Authorization: `Bearer ${accessToken}`,
                    },
                }
            );

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.error || 'Failed to generate report');
            }

            // Получение PDF как blob и создание локального URL
            const blob = await response.blob();
            const blobUrl = URL.createObjectURL(blob);

            setLastGeneratedFile(blobUrl);
            Alert.alert('Успех', 'Отчёт успешно сгенерирован!');
        } catch (error: unknown) {
            console.error('Report generation error:', error);
            const message = error instanceof Error ? error.message : 'Не удалось сгенерировать отчёт';
            Alert.alert('Ошибка', message);
        } finally {
            setIsLoading(false);
        }
    };

    const shareReport = async () => {
        if (!lastGeneratedFile) {
            Alert.alert('Ошибка', 'Сначала сгенерируйте отчёт');
            return;
        }

        try {
            const isAvailable = await Sharing.isAvailableAsync();
            if (!isAvailable) {
                Alert.alert('Ошибка', 'Функция "Поделиться" недоступна на этом устройстве');
                return;
            }

            await Sharing.shareAsync(lastGeneratedFile, {
                mimeType: 'application/pdf',
                dialogTitle: 'Поделиться отчётом',
                UTI: 'com.adobe.pdf',
            });
        } catch (error: unknown) {
            console.error('Share error:', error);
            Alert.alert('Ошибка', 'Не удалось поделиться файлом');
        }
    };

    return (
        <View style={styles.container}>
            {/* Заголовок */}
            <View style={styles.header}>
                <Text style={styles.headerTitle}>Отчёты</Text>
                <Text style={styles.headerSubtitle}>Генерация PDF отчётов по задачам</Text>
            </View>

            {/* Карточка типа отчёта */}
            <View style={styles.card}>
                <View style={styles.cardHeader}>
                    <View style={styles.iconContainer}>
                        <Ionicons name="document-text" size={28} color={COLORS.primary} />
                    </View>
                    <View style={styles.cardTitleContainer}>
                        <Text style={styles.cardTitle}>Отчёт по задачам</Text>
                        <Text style={styles.cardSubtitle}>Список задач за выбранный период</Text>
                    </View>
                </View>

                {/* Поля выбора дат */}
                <View style={styles.dateContainer}>
                    <View style={styles.dateField}>
                        <Text style={styles.dateLabel}>С</Text>
                        <TextInput
                            style={styles.dateInput}
                            value={fromDate}
                            onChangeText={setFromDate}
                            placeholder="YYYY-MM-DD"
                            placeholderTextColor={COLORS.textSecondary}
                            keyboardType={Platform.OS === 'ios' ? 'default' : 'default'}
                        />
                    </View>
                    <View style={styles.dateSeparator}>
                        <Ionicons name="arrow-forward" size={20} color={COLORS.textSecondary} />
                    </View>
                    <View style={styles.dateField}>
                        <Text style={styles.dateLabel}>По</Text>
                        <TextInput
                            style={styles.dateInput}
                            value={toDate}
                            onChangeText={setToDate}
                            placeholder="YYYY-MM-DD"
                            placeholderTextColor={COLORS.textSecondary}
                            keyboardType={Platform.OS === 'ios' ? 'default' : 'default'}
                        />
                    </View>
                </View>

                {/* Кнопка генерации */}
                <TouchableOpacity
                    style={[styles.generateButton, isLoading && styles.generateButtonDisabled]}
                    onPress={generateReport}
                    disabled={isLoading}
                >
                    {isLoading ? (
                        <ActivityIndicator size="small" color={COLORS.text} />
                    ) : (
                        <>
                            <Ionicons name="download" size={20} color={COLORS.text} />
                            <Text style={styles.generateButtonText}>Сгенерировать отчёт</Text>
                        </>
                    )}
                </TouchableOpacity>

                {/* Кнопка "Поделиться" (видна только после генерации) */}
                {lastGeneratedFile && (
                    <TouchableOpacity style={styles.shareButton} onPress={shareReport}>
                        <Ionicons name="share-outline" size={20} color={COLORS.primary} />
                        <Text style={styles.shareButtonText}>Поделиться</Text>
                    </TouchableOpacity>
                )}
            </View>

            {/* Информационная карточка */}
            <View style={styles.infoCard}>
                <Ionicons name="information-circle" size={24} color={COLORS.textSecondary} />
                <Text style={styles.infoText}>
                    Отчёт включает все ваши задачи (созданные и назначенные) за выбранный период с
                    разбивкой по статусам.
                </Text>
            </View>
        </View>
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
        borderRadius: 16,
        padding: 20,
    },
    cardHeader: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 20,
    },
    iconContainer: {
        width: 56,
        height: 56,
        borderRadius: 28,
        backgroundColor: COLORS.background,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: 16,
    },
    cardTitleContainer: {
        flex: 1,
    },
    cardTitle: {
        fontSize: 18,
        fontWeight: '600',
        color: COLORS.text,
        marginBottom: 4,
    },
    cardSubtitle: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    dateContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 20,
    },
    dateField: {
        flex: 1,
    },
    dateLabel: {
        fontSize: 12,
        color: COLORS.textSecondary,
        marginBottom: 6,
        textTransform: 'uppercase',
        fontWeight: '600',
    },
    dateInput: {
        backgroundColor: COLORS.background,
        borderRadius: 12,
        paddingHorizontal: 16,
        paddingVertical: 14,
        color: COLORS.text,
        fontSize: 16,
        borderWidth: 1,
        borderColor: COLORS.border,
    },
    dateSeparator: {
        paddingHorizontal: 12,
        paddingTop: 20,
    },
    generateButton: {
        backgroundColor: COLORS.primary,
        borderRadius: 12,
        paddingVertical: 16,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        gap: 8,
    },
    generateButtonDisabled: {
        opacity: 0.6,
    },
    generateButtonText: {
        color: COLORS.text,
        fontSize: 16,
        fontWeight: '600',
    },
    shareButton: {
        marginTop: 12,
        borderRadius: 12,
        paddingVertical: 16,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        gap: 8,
        borderWidth: 1,
        borderColor: COLORS.primary,
    },
    shareButtonText: {
        color: COLORS.primary,
        fontSize: 16,
        fontWeight: '600',
    },
    infoCard: {
        flexDirection: 'row',
        backgroundColor: COLORS.surface,
        marginHorizontal: 16,
        borderRadius: 12,
        padding: 16,
        gap: 12,
    },
    infoText: {
        flex: 1,
        fontSize: 13,
        color: COLORS.textSecondary,
        lineHeight: 20,
    },
});

export default ReportsScreen;
