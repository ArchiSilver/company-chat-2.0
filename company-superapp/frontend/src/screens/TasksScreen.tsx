import { useCallback, useEffect } from 'react';
import { Alert, ScrollView, StyleSheet, Text, TouchableOpacity, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import TaskColumn from '../components/TaskColumn';
import { Task, TaskStatus, useTaskStore } from '../store/taskStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
    border: '#2E2E2E',
};

// Тестовые данные для демонстрации
const mockTasks: Task[] = [
    {
        id: '1',
        title: 'Design new login screen',
        description: 'Create mockups for the new login flow',
        status: 'todo',
        creator_id: 'user1',
        due_date: '2026-01-15T00:00:00Z',
        created_at: '2026-01-10T00:00:00Z',
        updated_at: '2026-01-10T00:00:00Z',
    },
    {
        id: '2',
        title: 'Implement WebSocket connection',
        description: 'Set up real-time messaging',
        status: 'in_progress',
        creator_id: 'user1',
        assignee_id: 'user2',
        created_at: '2026-01-09T00:00:00Z',
        updated_at: '2026-01-10T00:00:00Z',
    },
    {
        id: '3',
        title: 'Fix navigation bug',
        status: 'done',
        creator_id: 'user2',
        created_at: '2026-01-08T00:00:00Z',
        updated_at: '2026-01-10T00:00:00Z',
    },
    {
        id: '4',
        title: 'Review PR #42',
        description: 'Check the new feature implementation',
        status: 'todo',
        creator_id: 'user1',
        due_date: '2026-01-12T00:00:00Z',
        created_at: '2026-01-10T00:00:00Z',
        updated_at: '2026-01-10T00:00:00Z',
    },
];

const TasksScreen = () => {
    const { tasks, setTasks, updateTaskStatus } = useTaskStore();

    useEffect(() => {
        // В реальном приложении здесь загрузка задач из API
        // Пока используем тестовые данные
        if (tasks.length === 0) {
            setTasks(mockTasks);
        }
    }, [setTasks, tasks.length]);

    const todoTasks = tasks.filter((t) => t.status === 'todo');
    const inProgressTasks = tasks.filter((t) => t.status === 'in_progress');
    const doneTasks = tasks.filter((t) => t.status === 'done');

    const handleTaskPress = useCallback((task: Task) => {
        // Открыть модальное окно с деталями задачи
        Alert.alert(task.title, task.description || 'No description');
    }, []);

    const handleTaskLongPress = useCallback((task: Task) => {
        // Показать меню действий
        const nextStatus: Record<TaskStatus, TaskStatus> = {
            todo: 'in_progress',
            in_progress: 'done',
            done: 'todo',
        };

        Alert.alert(
            'Task Actions',
            `Move "${task.title}" to next status?`,
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: `Move to ${nextStatus[task.status].replace('_', ' ')}`,
                    onPress: () => {
                        updateTaskStatus(task.id, nextStatus[task.status]);
                        // В реальном приложении также вызов API:
                        // fetch(`/api/v1/tasks/${task.id}/status`, { method: 'PUT', body: ... })
                    },
                },
            ]
        );
    }, [updateTaskStatus]);

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Tasks</Text>
            </View>

            <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                contentContainerStyle={styles.scrollContent}
            >
                <TaskColumn
                    title="To Do"
                    status="todo"
                    tasks={todoTasks}
                    onTaskPress={handleTaskPress}
                    onTaskLongPress={handleTaskLongPress}
                />
                <TaskColumn
                    title="In Progress"
                    status="in_progress"
                    tasks={inProgressTasks}
                    onTaskPress={handleTaskPress}
                    onTaskLongPress={handleTaskLongPress}
                />
                <TaskColumn
                    title="Done"
                    status="done"
                    tasks={doneTasks}
                    onTaskPress={handleTaskPress}
                    onTaskLongPress={handleTaskLongPress}
                />
            </ScrollView>

            <TouchableOpacity
                style={styles.fab}
                onPress={() => Alert.alert('Create Task', 'Task creation modal would open here')}
            >
                <Text style={styles.fabText}>+</Text>
            </TouchableOpacity>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        padding: 16,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.border,
    },
    title: {
        color: COLORS.text,
        fontSize: 30,
        fontWeight: 'bold',
    },
    scrollContent: {
        paddingHorizontal: 16,
        paddingVertical: 16,
    },
    fab: {
        position: 'absolute',
        bottom: 24,
        right: 24,
        backgroundColor: COLORS.primary,
        width: 56,
        height: 56,
        borderRadius: 28,
        alignItems: 'center',
        justifyContent: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.3,
        shadowRadius: 4,
        elevation: 8,
    },
    fabText: {
        color: COLORS.text,
        fontSize: 30,
        fontWeight: '300',
    },
});

export default TasksScreen;
