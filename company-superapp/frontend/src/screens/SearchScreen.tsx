import { Ionicons } from '@expo/vector-icons';
import React, { useEffect } from 'react';
import {
    ActivityIndicator,
    FlatList,
    StyleSheet,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from 'react-native';
import { useDebounce } from '../hooks/useDebounce';
import { useSearchStore } from '../store/searchStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
};

interface ResultItemProps {
    type: 'user' | 'message' | 'task';
    title: string;
    subtitle: string;
    onPress: () => void;
}

const ResultItem: React.FC<ResultItemProps> = ({ type, title, subtitle, onPress }) => {
    const getIcon = (): keyof typeof Ionicons.glyphMap => {
        switch (type) {
            case 'user':
                return 'person';
            case 'message':
                return 'chatbubble';
            case 'task':
                return 'checkbox';
            default:
                return 'search';
        }
    };

    const getTypeLabel = () => {
        switch (type) {
            case 'user':
                return 'Пользователь';
            case 'message':
                return 'Сообщение';
            case 'task':
                return 'Задача';
            default:
                return '';
        }
    };

    return (
        <TouchableOpacity style={styles.resultItem} onPress={onPress}>
            <View style={styles.iconContainer}>
                <Ionicons name={getIcon()} size={24} color={COLORS.primary} />
            </View>
            <View style={styles.resultContent}>
                <Text style={styles.typeLabel}>{getTypeLabel()}</Text>
                <Text style={styles.resultTitle} numberOfLines={1}>
                    {title}
                </Text>
                <Text style={styles.resultSubtitle} numberOfLines={1}>
                    {subtitle}
                </Text>
            </View>
            <Ionicons name="chevron-forward" size={20} color={COLORS.textSecondary} />
        </TouchableOpacity>
    );
};

export const SearchScreen: React.FC = () => {
    const { query, setQuery, results, isLoading, error, search, clearResults } =
        useSearchStore();

    const debouncedQuery = useDebounce(query, 300);

    useEffect(() => {
        if (debouncedQuery.trim().length >= 2) {
            search(debouncedQuery);
        } else {
            clearResults();
        }
    }, [debouncedQuery, search, clearResults]);

    const handleResultPress = (item: { type: string; id: string }) => {
        // Переход на соответствующий экран в зависимости от типа
        console.log('Navigate to:', item.type, item.id);
    };

    return (
        <View style={styles.container}>
            {/* Заголовок поиска */}
            <View style={styles.header}>
                <View style={styles.searchInputContainer}>
                    <Ionicons
                        name="search"
                        size={20}
                        color={COLORS.textSecondary}
                        style={styles.searchIcon}
                    />
                    <TextInput
                        style={styles.searchInput}
                        placeholder="Поиск по сообщениям, задачам, пользователям..."
                        placeholderTextColor={COLORS.textSecondary}
                        value={query}
                        onChangeText={setQuery}
                        autoFocus
                        returnKeyType="search"
                    />
                    {query.length > 0 && (
                        <TouchableOpacity onPress={() => setQuery('')}>
                            <Ionicons name="close-circle" size={20} color={COLORS.textSecondary} />
                        </TouchableOpacity>
                    )}
                </View>
            </View>

            {/* Индикатор загрузки */}
            {isLoading && (
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                </View>
            )}

            {/* Сообщение об ошибке */}
            {error && (
                <View style={styles.errorContainer}>
                    <Text style={styles.errorText}>{error}</Text>
                </View>
            )}

            {/* Список результатов */}
            {!isLoading && !error && (
                <FlatList
                    data={results}
                    keyExtractor={(item) => `${item.type}-${item.id}`}
                    renderItem={({ item }) => (
                        <ResultItem
                            type={item.type}
                            title={item.title}
                            subtitle={item.subtitle}
                            onPress={() => handleResultPress(item)}
                        />
                    )}
                    contentContainerStyle={styles.listContent}
                    ListEmptyComponent={
                        query.length >= 2 ? (
                            <View style={styles.emptyContainer}>
                                <Ionicons name="search-outline" size={64} color={COLORS.textSecondary} />
                                <Text style={styles.emptyText}>Ничего не найдено</Text>
                                <Text style={styles.emptySubtext}>
                                    Попробуйте изменить поисковый запрос
                                </Text>
                            </View>
                        ) : (
                            <View style={styles.emptyContainer}>
                                <Ionicons name="search" size={64} color={COLORS.textSecondary} />
                                <Text style={styles.emptyText}>Глобальный поиск</Text>
                                <Text style={styles.emptySubtext}>
                                    Введите минимум 2 символа для поиска
                                </Text>
                            </View>
                        )
                    }
                />
            )}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        paddingHorizontal: 16,
        paddingTop: 60,
        paddingBottom: 16,
        backgroundColor: COLORS.surface,
    },
    searchInputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: COLORS.background,
        borderRadius: 12,
        paddingHorizontal: 12,
        height: 48,
    },
    searchIcon: {
        marginRight: 8,
    },
    searchInput: {
        flex: 1,
        color: COLORS.text,
        fontSize: 16,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    errorContainer: {
        padding: 16,
        alignItems: 'center',
    },
    errorText: {
        color: COLORS.primary,
        fontSize: 14,
    },
    listContent: {
        paddingVertical: 8,
    },
    resultItem: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: COLORS.surface,
        marginHorizontal: 16,
        marginVertical: 4,
        padding: 16,
        borderRadius: 12,
    },
    iconContainer: {
        width: 48,
        height: 48,
        borderRadius: 24,
        backgroundColor: COLORS.background,
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: 12,
    },
    resultContent: {
        flex: 1,
    },
    typeLabel: {
        fontSize: 10,
        color: COLORS.primary,
        fontWeight: '600',
        textTransform: 'uppercase',
        marginBottom: 2,
    },
    resultTitle: {
        fontSize: 16,
        color: COLORS.text,
        fontWeight: '600',
        marginBottom: 2,
    },
    resultSubtitle: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        paddingTop: 100,
    },
    emptyText: {
        fontSize: 18,
        color: COLORS.text,
        fontWeight: '600',
        marginTop: 16,
    },
    emptySubtext: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginTop: 8,
        textAlign: 'center',
        paddingHorizontal: 32,
    },
});

export default SearchScreen;
