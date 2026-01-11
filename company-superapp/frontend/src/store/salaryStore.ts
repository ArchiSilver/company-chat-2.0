import { create } from 'zustand';
import { useAuthStore } from './authStore';

interface SalaryState {
    amount: number | null;
    isLoading: boolean;
    error: string | null;
    fetchSalary: () => Promise<void>;
}

const API_URL = 'http://localhost:8080/api/v1';

export const useSalaryStore = create<SalaryState>((set) => ({
    amount: null,
    isLoading: false,
    error: null,

    fetchSalary: async () => {
        set({ isLoading: true, error: null });
        try {
            const token = useAuthStore.getState().accessToken;
            const response = await fetch(`${API_URL}/finance/salary`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                throw new Error('Failed to fetch salary');
            }

            const data = await response.json();
            set({ amount: data.amount, isLoading: false });
        } catch (error) {
            set({ error: (error as Error).message, isLoading: false });
        }
    },
}));
