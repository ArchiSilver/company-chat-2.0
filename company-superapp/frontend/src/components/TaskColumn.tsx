import { ScrollView, StyleSheet, Text, View, ViewStyle } from 'react-native';
import { Task, TaskStatus } from '../store/taskStore';
import TaskCard from './TaskCard';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
    blue: '#3B82F6',
    green: '#22C55E',
    gray: '#6B7280',
};

interface TaskColumnProps {
    title: string;
    status: TaskStatus;
    tasks: Task[];
    onTaskPress?: (task: Task) => void;
    onTaskLongPress?: (task: Task) => void;
}

const getColumnBorderColor = (status: TaskStatus): ViewStyle => {
    switch (status) {
        case 'todo':
            return { borderTopColor: COLORS.gray };
        case 'in_progress':
            return { borderTopColor: COLORS.blue };
        case 'done':
            return { borderTopColor: COLORS.green };
        default:
            return { borderTopColor: COLORS.gray };
    }
};

const TaskColumn = ({
    title,
    status,
    tasks,
    onTaskPress,
    onTaskLongPress,
}: TaskColumnProps) => {
    return (
        <View style={[styles.column, getColumnBorderColor(status)]}>
            <View style={styles.header}>
                <View style={styles.headerRow}>
                    <Text style={styles.title}>{title}</Text>
                    <View style={styles.badge}>
                        <Text style={styles.badgeText}>{tasks.length}</Text>
                    </View>
                </View>
            </View>

            <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
                {tasks.map((task: Task) => (
                    <View key={task.id}>
                        <TaskCard
                            task={task}
                            onPress={() => onTaskPress?.(task)}
                            onLongPress={() => onTaskLongPress?.(task)}
                        />
                    </View>
                ))}

                {tasks.length === 0 && (
                    <View style={styles.emptyContainer}>
                        <Text style={styles.emptyText}>No tasks</Text>
                    </View>
                )}
            </ScrollView>
        </View>
    );
};

const styles = StyleSheet.create({
    column: {
        width: 288,
        marginRight: 16,
        backgroundColor: COLORS.background,
        borderRadius: 12,
        borderTopWidth: 4,
    },
    header: {
        padding: 16,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.border,
    },
    headerRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    title: {
        color: COLORS.text,
        fontSize: 18,
        fontWeight: 'bold',
    },
    badge: {
        backgroundColor: COLORS.surface,
        paddingHorizontal: 8,
        paddingVertical: 4,
        borderRadius: 12,
    },
    badgeText: {
        color: COLORS.text,
        fontSize: 12,
        fontWeight: '500',
    },
    content: {
        flex: 1,
        padding: 8,
    },
    emptyContainer: {
        padding: 16,
        alignItems: 'center',
    },
    emptyText: {
        color: COLORS.textSecondary,
        fontSize: 14,
    },
});

export default TaskColumn;
