import { StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { Task } from '../store/taskStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
    danger: '#EF4444',
    dangerLight: 'rgba(239, 68, 68, 0.2)',
    primaryLight: 'rgba(255, 75, 51, 0.2)',
};

interface TaskCardProps {
    task: Task;
    onPress?: () => void;
    onLongPress?: () => void;
}

const TaskCard = ({ task, onPress, onLongPress }: TaskCardProps) => {
    const formatDate = (dateString?: string) => {
        if (!dateString) return null;
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
    };

    const isOverdue = task.due_date && new Date(task.due_date) < new Date();

    return (
        <TouchableOpacity
            onPress={onPress}
            onLongPress={onLongPress}
            style={styles.card}
            activeOpacity={0.7}
        >
            <Text style={styles.title} numberOfLines={2}>
                {task.title}
            </Text>

            {task.description && (
                <Text style={styles.description} numberOfLines={2}>
                    {task.description}
                </Text>
            )}

            <View style={styles.footer}>
                {task.due_date && (
                    <View style={[styles.dateContainer, isOverdue ? styles.dateOverdue : styles.dateNormal]}>
                        <Text style={[styles.dateText, isOverdue ? styles.dateTextOverdue : styles.dateTextNormal]}>
                            {formatDate(task.due_date)}
                        </Text>
                    </View>
                )}

                {task.assignee_id && (
                    <View style={styles.avatar} />
                )}
            </View>

            {task.source_message_id && (
                <View style={styles.sourceContainer}>
                    <Text style={styles.sourceText}>Из сообщения чата</Text>
                </View>
            )}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    card: {
        backgroundColor: COLORS.surface,
        padding: 16,
        borderRadius: 12,
        marginBottom: 12,
        marginHorizontal: 8,
    },
    title: {
        color: COLORS.text,
        fontSize: 16,
        fontWeight: '500',
        marginBottom: 8,
    },
    description: {
        color: COLORS.textSecondary,
        fontSize: 14,
        marginBottom: 8,
    },
    footer: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginTop: 8,
    },
    dateContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: 8,
        paddingVertical: 4,
        borderRadius: 4,
    },
    dateNormal: {
        backgroundColor: COLORS.primaryLight,
    },
    dateOverdue: {
        backgroundColor: COLORS.dangerLight,
    },
    dateText: {
        fontSize: 12,
    },
    dateTextNormal: {
        color: COLORS.primary,
    },
    dateTextOverdue: {
        color: COLORS.danger,
    },
    avatar: {
        width: 24,
        height: 24,
        backgroundColor: '#4B5563',
        borderRadius: 12,
    },
    sourceContainer: {
        marginTop: 8,
        paddingTop: 8,
        borderTopWidth: 1,
        borderTopColor: COLORS.border,
    },
    sourceText: {
        color: COLORS.textSecondary,
        fontSize: 12,
    },
});

export default TaskCard;
