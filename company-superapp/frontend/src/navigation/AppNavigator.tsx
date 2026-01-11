import { NavigationContainer, useNavigationContainerRef } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { useEffect, useState } from 'react';
import { usePushNotifications } from '../hooks/usePushNotifications';
import { useAuthStore } from '../store/authStore';

import { ActivityIndicator, View } from 'react-native';
import LoginScreen from '../screens/LoginScreen';
import RegisterScreen from '../screens/RegisterScreen';

import { Ionicons } from '@expo/vector-icons';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { ProtectedRoute } from '../components/ProtectedRoute';
import ChatListScreen from '../screens/ChatListScreen';
import ChatScreen from '../screens/ChatScreen';
import DebugScreen from '../screens/DebugScreen';
import FinanceScreen from '../screens/FinanceScreen';
import ReportsScreen from '../screens/ReportsScreen';
import SearchScreen from '../screens/SearchScreen';
import TasksScreen from '../screens/TasksScreen';
import TaxiRequestScreen from '../screens/TaxiRequestScreen';

// Protected Finance Screen Wrapper
const ProtectedFinanceScreen = () => (
    <ProtectedRoute allowedRoles={['admin', 'manager']}>
        <FinanceScreen />
    </ProtectedRoute>
);

// Stubs for the main app stacks
const MainStack = createNativeStackNavigator();
const MainAppStack = () => (
    <MainStack.Navigator screenOptions={{ headerShown: false }}>
        <MainStack.Screen name="HomeTabs" component={HomeTabs} />
        <MainStack.Screen name="Chat" component={ChatScreen} />
    </MainStack.Navigator>
);

const Tab = createBottomTabNavigator();
const HomeTabs = () => {
    const { hasRole } = useAuthStore();
    const canAccessFinance = hasRole(['admin', 'manager']);

    return (
        <Tab.Navigator
            screenOptions={({ route }) => ({
                headerShown: false,
                tabBarStyle: { backgroundColor: '#1E1E1E' },
                tabBarActiveTintColor: '#FF4B33',
                tabBarInactiveTintColor: '#A0A0A0',
                tabBarIcon: ({ color, size }) => {
                    let iconName: keyof typeof Ionicons.glyphMap = 'ellipse';

                    if (route.name === 'Chats') iconName = 'chatbubbles';
                    else if (route.name === 'Tasks') iconName = 'checkbox';
                    else if (route.name === 'Finance') iconName = 'wallet';
                    else if (route.name === 'Taxi') iconName = 'car';
                    else if (route.name === 'Search') iconName = 'search';
                    else if (route.name === 'Reports') iconName = 'document-text';
                    else if (route.name === 'Debug') iconName = 'bug';

                    return <Ionicons name={iconName} size={size} color={color} />;
                },
            })}
        >
            <Tab.Screen name="Chats" component={ChatListScreen} />
            <Tab.Screen name="Tasks" component={TasksScreen} />
            {canAccessFinance && (
                <Tab.Screen name="Finance" component={ProtectedFinanceScreen} />
            )}
            <Tab.Screen name="Taxi" component={TaxiRequestScreen} />
            <Tab.Screen name="Reports" component={ReportsScreen} />
            <Tab.Screen name="Search" component={SearchScreen} />
            {__DEV__ && <Tab.Screen name="Debug" component={DebugScreen} />}
        </Tab.Navigator>
    );
};

const AuthStack = createNativeStackNavigator();
const AuthNavigator = () => (
    <AuthStack.Navigator screenOptions={{ headerShown: false }}>
        <AuthStack.Screen name="Login" component={LoginScreen} />
        <AuthStack.Screen name="Register" component={RegisterScreen} />
    </AuthStack.Navigator>
);

const AppNavigator = () => {
    const { isLoggedIn, checkAuth } = useAuthStore();
    const [isLoading, setIsLoading] = useState(true);
    const navigationRef = useNavigationContainerRef();

    // Handle notification tap navigation
    const handleNotificationTap = (data: { type?: string; chat_id?: string; task_id?: string }) => {
        if (data.type === 'chat' && data.chat_id && navigationRef.isReady()) {
            (navigationRef as any).navigate('Chat', { chatId: data.chat_id, name: 'Chat' });
        } else if (data.type === 'task' && navigationRef.isReady()) {
            (navigationRef as any).navigate('HomeTabs', { screen: 'Tasks' });
        }
    };

    // Initialize push notifications
    usePushNotifications(handleNotificationTap);

    useEffect(() => {
        const bootstrapAsync = async () => {
            await checkAuth();
            setIsLoading(false);
        };
        bootstrapAsync();
    }, [checkAuth]);

    if (isLoading) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#0F0F0F' }}>
                <ActivityIndicator size="large" color="#FF4B33" />
            </View>
        );
    }

    return (
        <NavigationContainer ref={navigationRef}>
            {isLoggedIn ? <MainAppStack /> : <AuthNavigator />}
        </NavigationContainer>
    );
};

export default AppNavigator;
