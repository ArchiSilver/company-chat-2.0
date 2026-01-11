import { useRoute } from '@react-navigation/native';
import { useCallback, useState } from 'react';
import {
    Alert,
    FlatList,
    KeyboardAvoidingView,
    Platform,
    StyleSheet,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useWebSocket } from '../hooks/useWebSocket';
import { Message, useChatStore } from '../store/chatStore';
import { useTaskStore } from '../store/taskStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
};

// ID текущего пользователя — в реальном приложении берётся из auth store
const CURRENT_USER_ID = 'a1b2c3d4-e5f6-7890-abcd-ef1234567890';

const ChatScreen = () => {
    const route = useRoute();
    const { chatId, name } = route.params as { chatId: string; name: string };
    const [input, setInput] = useState('');

    const messages = useChatStore((state) => state.messages[chatId] || []);
    const { sendMessage } = useWebSocket({ chatId, userId: CURRENT_USER_ID });
    const addTask = useTaskStore((state) => state.addTask);

    const onSend = useCallback(() => {
        if (input.trim().length === 0) return;
        sendMessage(input.trim());
        setInput('');
    }, [input, sendMessage]);

    const handleCreateTaskFromMessage = useCallback((message: Message) => {
        Alert.alert(
            'Create Task',
            `Create task from message: "${message.content}"?`,
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Create',
                    onPress: () => {
                        const newTask = {
                            id: Date.now().toString(),
                            title: message.content,
                            description: 'Created from chat message',
                            status: 'todo' as const,
                            creator_id: CURRENT_USER_ID,
                            source_message_id: parseInt(message.id) || undefined,
                            created_at: new Date().toISOString(),
                            updated_at: new Date().toISOString(),
                        };
                        addTask(newTask);
                        Alert.alert('Success', 'Task created!');
                    },
                },
            ]
        );
    }, [addTask]);

    const renderItem = ({ item }: { item: Message }) => {
        const isMyMessage = item.sender === 'me';
        return (
            <TouchableOpacity
                onLongPress={() => handleCreateTaskFromMessage(item)}
                delayLongPress={500}
            >
                <View style={[styles.messageRow, isMyMessage ? styles.messageRowRight : styles.messageRowLeft]}>
                    <View style={[styles.messageBubble, isMyMessage ? styles.myBubble : styles.otherBubble]}>
                        <Text style={styles.messageText}>{item.content}</Text>
                    </View>
                </View>
            </TouchableOpacity>
        );
    };

    return (
        <SafeAreaView style={styles.container} edges={['bottom']}>
            <KeyboardAvoidingView
                behavior={Platform.OS === "ios" ? "padding" : "height"}
                style={styles.flex}
                keyboardVerticalOffset={90}
            >
                <View style={styles.header}>
                    <Text style={styles.headerTitle}>{name}</Text>
                </View>
                <FlatList
                    data={messages}
                    renderItem={renderItem}
                    keyExtractor={(item: Message) => item.id}
                    inverted
                />
                <View style={styles.inputContainer}>
                    <TextInput
                        style={styles.input}
                        placeholder="Message..."
                        placeholderTextColor={COLORS.textSecondary}
                        value={input}
                        onChangeText={setInput}
                    />
                    <TouchableOpacity onPress={onSend} style={styles.sendButton}>
                        <Text style={styles.sendText}>Send</Text>
                    </TouchableOpacity>
                </View>
            </KeyboardAvoidingView>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    flex: {
        flex: 1,
    },
    header: {
        padding: 16,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.border,
    },
    headerTitle: {
        color: COLORS.text,
        fontSize: 20,
        textAlign: 'center',
        fontWeight: 'bold',
    },
    messageRow: {
        flexDirection: 'row',
        marginVertical: 4,
        paddingHorizontal: 16,
    },
    messageRowRight: {
        justifyContent: 'flex-end',
    },
    messageRowLeft: {
        justifyContent: 'flex-start',
    },
    messageBubble: {
        padding: 12,
        borderRadius: 16,
        maxWidth: '80%',
    },
    myBubble: {
        backgroundColor: COLORS.primary,
    },
    otherBubble: {
        backgroundColor: COLORS.surface,
    },
    messageText: {
        color: COLORS.text,
        fontSize: 16,
    },
    inputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 8,
        borderTopWidth: 1,
        borderTopColor: COLORS.border,
    },
    input: {
        flex: 1,
        backgroundColor: COLORS.surface,
        color: COLORS.text,
        borderRadius: 20,
        paddingHorizontal: 16,
        paddingVertical: 8,
        marginRight: 8,
        fontSize: 16,
    },
    sendButton: {
        backgroundColor: COLORS.primary,
        padding: 12,
        borderRadius: 20,
    },
    sendText: {
        color: COLORS.text,
        fontWeight: 'bold',
    },
});

export default ChatScreen;
