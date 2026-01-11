import AsyncStorage from '@react-native-async-storage/async-storage';
import { jwtDecode } from 'jwt-decode';
import { create } from 'zustand';
import { Analytics } from '../lib/analytics';
import { clearUserContext, setUserContext } from '../lib/sentry';

interface JWTPayload {
    sub: string;
    role: string;
    exp: number;
}

export type UserRole = 'admin' | 'manager' | 'user';

interface AuthState {
    isLoggedIn: boolean;
    accessToken: string | null;
    userId: string | null;
    role: UserRole | null;
    login: (accessToken: string, refreshToken: string) => Promise<void>;
    logout: () => Promise<void>;
    checkAuth: () => Promise<void>;
    hasRole: (allowedRoles: UserRole[]) => boolean;
}

const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';

const decodeToken = (token: string): { userId: string; role: UserRole } => {
    try {
        const decoded = jwtDecode<JWTPayload>(token);
        return {
            userId: decoded.sub,
            role: (decoded.role as UserRole) || 'user',
        };
    } catch {
        return { userId: '', role: 'user' };
    }
};

export const useAuthStore = create<AuthState>((set, get) => ({
    isLoggedIn: false,
    accessToken: null,
    userId: null,
    role: null,

    login: async (accessToken: string, refreshToken: string) => {
        await AsyncStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
        await AsyncStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
        const { userId, role } = decodeToken(accessToken);

        // Установка контекста пользователя в Sentry
        setUserContext(userId, undefined, role);

        // Отслеживание события входа
        Analytics.identify(userId, { role });
        Analytics.events.userLoggedIn(userId);

        set({ isLoggedIn: true, accessToken, userId, role });
    },

    logout: async () => {
        // Отслеживание события выхода
        Analytics.events.userLoggedOut();
        Analytics.reset();

        // Очистка контекста пользователя в Sentry
        clearUserContext();

        await AsyncStorage.removeItem(ACCESS_TOKEN_KEY);
        await AsyncStorage.removeItem(REFRESH_TOKEN_KEY);
        set({ isLoggedIn: false, accessToken: null, userId: null, role: null });
    },

    checkAuth: async () => {
        const accessToken = await AsyncStorage.getItem(ACCESS_TOKEN_KEY);
        if (accessToken) {
            const { userId, role } = decodeToken(accessToken);

            // Установка контекста пользователя в Sentry
            setUserContext(userId, undefined, role);

            // Идентификация пользователя в аналитике
            Analytics.identify(userId, { role });

            set({ isLoggedIn: true, accessToken, userId, role });
        } else {
            set({ isLoggedIn: false, accessToken: null, userId: null, role: null });
        }
    },

    hasRole: (allowedRoles: UserRole[]) => {
        const { role } = get();
        if (!role) return false;
        return allowedRoles.includes(role);
    },
}));
