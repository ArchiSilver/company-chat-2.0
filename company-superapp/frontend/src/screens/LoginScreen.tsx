import { useNavigation } from '@react-navigation/native';
import {
    SafeAreaView,
    StyleSheet,
    Text,
    TextInput,
    TouchableOpacity,
    View,
} from 'react-native';
import { useAuthStore } from '../store/authStore';

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
};

const LoginScreen = () => {
    const navigation = useNavigation();
    const { login } = useAuthStore();

    const handleLogin = async () => {
        // Тестовый вход для демонстрации UI. Заменить на реальный API вызов.
        await login('fake-access-token', 'fake-refresh-token');
    };

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.content}>
                <Text style={styles.title}>Login</Text>
                <TextInput
                    style={styles.input}
                    placeholder="Email"
                    placeholderTextColor={COLORS.textSecondary}
                    keyboardType="email-address"
                    autoCapitalize="none"
                />
                <TextInput
                    style={[styles.input, styles.inputPassword]}
                    placeholder="Password"
                    placeholderTextColor={COLORS.textSecondary}
                    secureTextEntry
                />
                <TouchableOpacity style={styles.button} onPress={handleLogin}>
                    <Text style={styles.buttonText}>Sign In</Text>
                </TouchableOpacity>
                <TouchableOpacity onPress={() => navigation.navigate('Register' as never)}>
                    <Text style={styles.linkText}>Don't have an account? Sign Up</Text>
                </TouchableOpacity>
            </View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
        padding: 16,
    },
    content: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    title: {
        color: COLORS.text,
        fontSize: 30,
        fontWeight: 'bold',
        marginBottom: 32,
    },
    input: {
        width: '100%',
        backgroundColor: COLORS.surface,
        color: COLORS.text,
        borderRadius: 12,
        paddingHorizontal: 16,
        paddingVertical: 12,
        marginBottom: 16,
        fontSize: 16,
    },
    inputPassword: {
        marginBottom: 24,
    },
    button: {
        width: '100%',
        backgroundColor: COLORS.primary,
        borderRadius: 12,
        paddingVertical: 16,
        alignItems: 'center',
    },
    buttonText: {
        color: COLORS.text,
        fontWeight: 'bold',
        fontSize: 18,
    },
    linkText: {
        color: COLORS.textSecondary,
        marginTop: 24,
    },
});

export default LoginScreen;
