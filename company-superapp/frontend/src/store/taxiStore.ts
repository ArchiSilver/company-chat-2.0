import { create } from 'zustand';
import { useAuthStore } from './authStore';

interface TaxiRequest {
    id: string;
    user_id: string;
    receipt_file_key: string;
    status: 'pending' | 'approved' | 'rejected';
    created_at: string;
    receipt_url?: string;
}

interface TaxiState {
    requests: TaxiRequest[];
    isLoading: boolean;
    error: string | null;
    fetchRequests: () => Promise<void>;
    addRequest: (request: TaxiRequest) => void;
}

const API_URL = 'http://localhost:8080/api/v1';

export const useTaxiStore = create<TaxiState>((set) => ({
    requests: [],
    isLoading: false,
    error: null,

    fetchRequests: async () => {
        set({ isLoading: true, error: null });
        try {
            const token = useAuthStore.getState().accessToken;
            const response = await fetch(`${API_URL}/taxi/requests`, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                throw new Error('Failed to fetch taxi requests');
            }

            const data = await response.json();
            set({ requests: data.requests || [], isLoading: false });
        } catch (error) {
            set({ error: (error as Error).message, isLoading: false });
        }
    },

    addRequest: (request: TaxiRequest) => {
        set((state) => ({
            requests: [request, ...state.requests],
        }));
    },
}));
