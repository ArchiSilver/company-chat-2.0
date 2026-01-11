import { create } from 'zustand';
import { api } from '../api/client';

// Унифицированный тип результата поиска
interface SearchResultItem {
    type: 'user' | 'message' | 'task';
    id: string;
    title: string;
    subtitle: string;
    rank: number;
}

interface SearchState {
    query: string;
    results: SearchResultItem[];
    isLoading: boolean;
    error: string | null;

    setQuery: (query: string) => void;
    search: (query: string) => Promise<void>;
    clearResults: () => void;
}

export const useSearchStore = create<SearchState>((set) => ({
    query: '',
    results: [],
    isLoading: false,
    error: null,

    setQuery: (query: string) => set({ query }),

    search: async (query: string) => {
        if (!query.trim()) {
            set({ results: [], isLoading: false });
            return;
        }

        set({ isLoading: true, error: null });

        try {
            const response = await api.get('/search', {
                params: { q: query },
            });

            const data = response.data;

            // Преобразование ответа API в унифицированный формат
            const unified: SearchResultItem[] = [];

            // Пользователи
            if (data.users) {
                data.users.forEach((user: any) => {
                    unified.push({
                        type: 'user',
                        id: user.id,
                        title: user.full_name,
                        subtitle: user.email,
                        rank: user.rank,
                    });
                });
            }

            // Сообщения
            if (data.messages) {
                data.messages.forEach((msg: any) => {
                    unified.push({
                        type: 'message',
                        id: msg.id,
                        title: msg.content.substring(0, 100),
                        subtitle: `Chat: ${msg.chat_id}`,
                        rank: msg.rank,
                    });
                });
            }

            // Задачи
            if (data.tasks) {
                data.tasks.forEach((task: any) => {
                    unified.push({
                        type: 'task',
                        id: task.id,
                        title: task.title,
                        subtitle: task.status,
                        rank: task.rank,
                    });
                });
            }

            // Сортировка по рангу (убывание)
            unified.sort((a, b) => b.rank - a.rank);

            set({ results: unified, isLoading: false });
        } catch (error: any) {
            set({
                error: error.response?.data?.error || 'Search failed',
                isLoading: false,
            });
        }
    },

    clearResults: () => set({ query: '', results: [], error: null }),
}));
