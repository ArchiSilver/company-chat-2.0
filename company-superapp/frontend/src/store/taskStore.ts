import { create } from 'zustand';

export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
    id: string;
    title: string;
    description?: string;
    status: TaskStatus;
    creator_id: string;
    assignee_id?: string;
    due_date?: string;
    source_message_id?: number;
    created_at: string;
    updated_at: string;
}

interface TaskState {
    tasks: Task[];
    isLoading: boolean;
    setTasks: (tasks: Task[]) => void;
    addTask: (task: Task) => void;
    updateTask: (id: string, updates: Partial<Task>) => void;
    updateTaskStatus: (id: string, status: TaskStatus) => void;
    deleteTask: (id: string) => void;
    getTasksByStatus: (status: TaskStatus) => Task[];
}

export const useTaskStore = create<TaskState>((set, get) => ({
    tasks: [],
    isLoading: false,
    setTasks: (tasks: Task[]) => set({ tasks }),
    addTask: (task: Task) => set((state) => ({ tasks: [...state.tasks, task] })),
    updateTask: (id: string, updates: Partial<Task>) =>
        set((state) => ({
            tasks: state.tasks.map((task) =>
                task.id === id ? { ...task, ...updates } : task
            ),
        })),
    updateTaskStatus: (id: string, status: TaskStatus) =>
        set((state) => ({
            tasks: state.tasks.map((task) =>
                task.id === id ? { ...task, status } : task
            ),
        })),
    deleteTask: (id: string) =>
        set((state) => ({
            tasks: state.tasks.filter((task) => task.id !== id),
        })),
    getTasksByStatus: (status: TaskStatus) => {
        return get().tasks.filter((task) => task.status === status);
    },
}));
