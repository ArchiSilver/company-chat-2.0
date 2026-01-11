import * as Device from 'expo-device';
import * as Notifications from 'expo-notifications';
import { useEffect, useRef, useState } from 'react';
import { Platform } from 'react-native';
import { useAuthStore } from '../store/authStore';

const API_URL = 'http://localhost:8080/api/v1';

// Configure notification behavior
Notifications.setNotificationHandler({
    handleNotification: async () => ({
        shouldShowAlert: true,
        shouldPlaySound: true,
        shouldSetBadge: true,
        shouldShowBanner: true,
        shouldShowList: true,
    }),
});

interface NotificationData {
    type?: string;
    chat_id?: string;
    task_id?: string;
    [key: string]: string | undefined;
}

export function usePushNotifications(onNotificationTap?: (data: NotificationData) => void) {
    const [expoPushToken, setExpoPushToken] = useState<string | null>(null);
    const [notification, setNotification] = useState<Notifications.Notification | null>(null);
    const notificationListener = useRef<Notifications.EventSubscription | null>(null);
    const responseListener = useRef<Notifications.EventSubscription | null>(null);
    const { accessToken, isLoggedIn } = useAuthStore();

    useEffect(() => {
        if (!isLoggedIn) return;

        registerForPushNotificationsAsync().then((token) => {
            if (token) {
                setExpoPushToken(token);
                registerTokenOnServer(token);
            }
        });

        // Foreground notification listener
        notificationListener.current = Notifications.addNotificationReceivedListener((notification) => {
            setNotification(notification);
            // Update badge count
            updateBadgeCount(1);
        });

        // Background/tap notification listener
        responseListener.current = Notifications.addNotificationResponseReceivedListener((response) => {
            const data = response.notification.request.content.data as NotificationData;
            if (onNotificationTap) {
                onNotificationTap(data);
            }
            // Clear badge on tap
            updateBadgeCount(0);
        });

        return () => {
            if (notificationListener.current) {
                notificationListener.current.remove();
            }
            if (responseListener.current) {
                responseListener.current.remove();
            }
        };
    }, [isLoggedIn, accessToken]);

    const registerTokenOnServer = async (token: string) => {
        if (!accessToken) return;

        try {
            const response = await fetch(`${API_URL}/notifications/register-token`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    token: token,
                    device_info: `${Device.modelName} - ${Platform.OS} ${Platform.Version}`,
                }),
            });

            if (!response.ok) {
                console.error('Failed to register push token on server');
            }
        } catch (error) {
            console.error('Error registering push token:', error);
        }
    };

    const unregisterToken = async () => {
        if (!expoPushToken || !accessToken) return;

        try {
            await fetch(`${API_URL}/notifications/unregister-token`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ token: expoPushToken }),
            });
        } catch (error) {
            console.error('Error unregistering push token:', error);
        }
    };

    return {
        expoPushToken,
        notification,
        unregisterToken,
    };
}

async function registerForPushNotificationsAsync(): Promise<string | null> {
    let token: string | null = null;

    if (Platform.OS === 'android') {
        await Notifications.setNotificationChannelAsync('default', {
            name: 'default',
            importance: Notifications.AndroidImportance.MAX,
            vibrationPattern: [0, 250, 250, 250],
            lightColor: '#FF4B33',
        });
    }

    if (!Device.isDevice) {
        console.log('Must use physical device for Push Notifications');
        return null;
    }

    const { status: existingStatus } = await Notifications.getPermissionsAsync();
    let finalStatus = existingStatus;

    if (existingStatus !== 'granted') {
        const { status } = await Notifications.requestPermissionsAsync();
        finalStatus = status;
    }

    if (finalStatus !== 'granted') {
        console.log('Failed to get push token for push notification!');
        return null;
    }

    try {
        const tokenData = await Notifications.getExpoPushTokenAsync({
            projectId: 'your-project-id', // Replace with your Expo project ID
        });
        token = tokenData.data;
    } catch (error) {
        console.error('Error getting Expo push token:', error);
    }

    return token;
}

export async function updateBadgeCount(count: number): Promise<void> {
    try {
        await Notifications.setBadgeCountAsync(count);
    } catch (error) {
        console.error('Error setting badge count:', error);
    }
}

export async function clearAllNotifications(): Promise<void> {
    try {
        await Notifications.dismissAllNotificationsAsync();
        await updateBadgeCount(0);
    } catch (error) {
        console.error('Error clearing notifications:', error);
    }
}
