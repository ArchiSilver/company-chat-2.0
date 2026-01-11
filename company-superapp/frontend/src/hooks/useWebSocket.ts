import { useCallback, useEffect, useRef } from 'react';
import { Message, useChatStore } from '../store/chatStore';

const WS_URL = 'ws://localhost:8080/api/v1/ws/connect';

interface UseWebSocketOptions {
    chatId: string;
    userId: string;
}

export const useWebSocket = ({ chatId, userId }: UseWebSocketOptions) => {
    const ws = useRef(null as WebSocket | null);
    const addMessage = useChatStore((state) => state.addMessage);

    useEffect(() => {
        const url = `${WS_URL}?user_id=${userId}&chat_id=${chatId}`;
        ws.current = new WebSocket(url);

        ws.current.onopen = () => {
            console.log('WebSocket Connected');
        };

        ws.current.onmessage = (event: MessageEvent) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === 'message') {
                    const message: Message = {
                        id: data.id?.toString() || Date.now().toString(),
                        chat_id: data.chat_id,
                        sender_id: data.sender_id,
                        content: data.content,
                        created_at: data.created_at,
                        sender: data.sender_id === userId ? 'me' : 'other',
                    };
                    addMessage(chatId, message);
                }
            } catch (e) {
                console.error('Failed to parse message:', e);
            }
        };

        ws.current.onerror = (error: Event) => {
            console.error('WebSocket Error:', error);
        };

        ws.current.onclose = () => {
            console.log('WebSocket Disconnected');
        };

        return () => {
            ws.current?.close();
        };
    }, [chatId, userId, addMessage]);

    const sendMessage = useCallback(
        (content: string) => {
            if (ws.current?.readyState === WebSocket.OPEN) {
                const message = {
                    type: 'message',
                    chat_id: chatId,
                    content: content,
                };
                ws.current.send(JSON.stringify(message));
            }
        },
        [chatId]
    );

    return { sendMessage };
};
