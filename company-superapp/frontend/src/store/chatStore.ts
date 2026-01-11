import { create } from 'zustand';

export interface Message {
    id: string;
    chat_id: string;
    sender_id: string;
    content: string;
    created_at: string;
    sender: 'me' | 'other';
}

export interface Chat {
    id: string;
    name: string;
    lastMessage: string;
    unread: number;
}

interface ChatState {
    chats: Chat[];
    messages: Record<string, Message[]>;
    setChats: (chats: Chat[]) => void;
    addMessage: (chatId: string, message: Message) => void;
    setMessages: (chatId: string, messages: Message[]) => void;
}

export const useChatStore = create<ChatState>((set) => ({
    chats: [],
    messages: {},
    setChats: (chats) => set({ chats }),
    addMessage: (chatId, message) =>
        set((state) => ({
            messages: {
                ...state.messages,
                [chatId]: [message, ...(state.messages[chatId] || [])],
            },
        })),
    setMessages: (chatId, messages) =>
        set((state) => ({
            messages: {
                ...state.messages,
                [chatId]: messages,
            },
        })),
}));
