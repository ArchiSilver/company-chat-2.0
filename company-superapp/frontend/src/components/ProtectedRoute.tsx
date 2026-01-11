import { Ionicons } from '@expo/vector-icons';
import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useAuthStore, UserRole } from '../store/authStore';

interface ProtectedRouteProps {
    allowedRoles: UserRole[];
    children: React.ReactNode;
}

const COLORS = {
    primary: '#FF4B33',
    background: '#0F0F0F',
    surface: '#1E1E1E',
    text: '#FFFFFF',
    textSecondary: '#A0A0A0',
};

/**
 * Protected route wrapper that shows access denied screen for unauthorized users.
 */
export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
    allowedRoles,
    children,
}) => {
    const { hasRole, role } = useAuthStore();

    if (hasRole(allowedRoles)) {
        return <>{children}</>;
    }

    return (
        <View style={styles.container}>
            <View style={styles.content}>
                <View style={styles.iconContainer}>
                    <Ionicons name="lock-closed" size={64} color={COLORS.primary} />
                </View>
                <Text style={styles.title}>Доступ запрещён</Text>
                <Text style={styles.subtitle}>
                    У вас недостаточно прав для просмотра этой страницы.
                </Text>
                <View style={styles.roleInfo}>
                    <Text style={styles.roleLabel}>Ваша роль:</Text>
                    <View style={styles.roleBadge}>
                        <Text style={styles.roleText}>{role || 'unknown'}</Text>
                    </View>
                </View>
                <Text style={styles.requiredText}>
                    Требуемые роли: {allowedRoles.join(', ')}
                </Text>
            </View>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
        justifyContent: 'center',
        alignItems: 'center',
    },
    content: {
        alignItems: 'center',
        paddingHorizontal: 32,
    },
    iconContainer: {
        width: 120,
        height: 120,
        borderRadius: 60,
        backgroundColor: COLORS.surface,
        justifyContent: 'center',
        alignItems: 'center',
        marginBottom: 24,
    },
    title: {
        fontSize: 24,
        fontWeight: '700',
        color: COLORS.text,
        marginBottom: 12,
    },
    subtitle: {
        fontSize: 16,
        color: COLORS.textSecondary,
        textAlign: 'center',
        marginBottom: 32,
        lineHeight: 24,
    },
    roleInfo: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 16,
    },
    roleLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginRight: 8,
    },
    roleBadge: {
        backgroundColor: COLORS.surface,
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 16,
        borderWidth: 1,
        borderColor: COLORS.primary,
    },
    roleText: {
        fontSize: 14,
        color: COLORS.primary,
        fontWeight: '600',
        textTransform: 'uppercase',
    },
    requiredText: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
});

export default ProtectedRoute;
