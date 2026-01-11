import { useNavigation } from '@react-navigation/native';
import { useEffect } from 'react';
import { FlatList, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Chat, useChatStore } from '../store/chatStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
};

// Тестовые данные — в реальном приложении загружаются с API
const mockChats: Chat[] = [
    { id: 'a1b2c3d4-0001-0000-0000-000000000001', name: 'Argen', lastMessage: 'See you tomorrow!', unread: 2 },
    { id: 'a1b2c3d4-0001-0000-0000-000000000002', name: 'Design Team', lastMessage: 'The new mockups are ready.', unread: 0 },
    { id: 'a1b2c3d4-0001-0000-0000-000000000003', name: 'John Doe', lastMessage: 'Okay, sounds good.', unread: 0 },
];

const ChatListScreen = () => {
    const navigation = useNavigation();
    const { chats, setChats } = useChatStore();

    useEffect(() => {
        // В реальном приложении здесь загружаем чаты с API
        // Пока используем тестовые данные
        setChats(mockChats);
    }, [setChats]);

    const renderItem = ({ item }: { item: Chat }) => (
        <TouchableOpacity
            style={styles.chatItem}
            onPress={() => (navigation as any).navigate('Chat', { chatId: item.id, name: item.name })}
        >
            {/* Заглушка аватара */}
            <View style={styles.avatar} />
            <View style={styles.chatInfo}>
                <Text style={styles.chatName}>{item.name}</Text>
                <Text style={styles.lastMessage}>{item.lastMessage}</Text>
            </View>
            {item.unread > 0 && (
                <View style={styles.unreadBadge}>
                    <Text style={styles.unreadText}>{item.unread}</Text>
                </View>
            )}
        </TouchableOpacity>
    );

    return (
        <SafeAreaView style={styles.container}>
            <Text style={styles.title}>Chats</Text>
            <FlatList
                data={chats.length > 0 ? chats : mockChats}
                renderItem={renderItem}
                keyExtractor={(item: Chat) => item.id}
            />
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    title: {
        color: COLORS.text,
        fontSize: 30,
        fontWeight: 'bold',
        padding: 16,
    },
    chatItem: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 16,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.border,
    },
    avatar: {
        width: 48,
        height: 48,
        backgroundColor: COLORS.surface,
        borderRadius: 24,
        marginRight: 16,
    },
    chatInfo: {
        flex: 1,
    },
    chatName: {
        color: COLORS.text,
        fontSize: 18,
    },
    lastMessage: {
        color: COLORS.textSecondary,
        fontSize: 14,
        marginTop: 4,
    },
    unreadBadge: {
        backgroundColor: COLORS.primary,
        borderRadius: 12,
        width: 24,
        height: 24,
        justifyContent: 'center',
        alignItems: 'center',
    },
    unreadText: {
        color: COLORS.text,
        fontWeight: 'bold',
        fontSize: 12,
    },
});

export default ChatListScreen;
