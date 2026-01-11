import { BlurView } from 'expo-blur';
import * as LocalAuthentication from 'expo-local-authentication';
import { useEffect, useState } from 'react';
import {
    ActivityIndicator,
    StyleSheet,
    Text,
    TouchableOpacity,
    View
} from 'react-native';
import { useSalaryStore } from '../store/salaryStore';
import { colors } from '../theme/colors';

export default function SalaryWidget() {
    const { amount, isLoading, error, fetchSalary } = useSalaryStore();
    const [isBlurred, setIsBlurred] = useState(true);
    const [isBiometricSupported, setIsBiometricSupported] = useState(false);

    useEffect(() => {
        checkBiometricSupport();
        fetchSalary();
    }, []);

    const checkBiometricSupport = async () => {
        const compatible = await LocalAuthentication.hasHardwareAsync();
        const enrolled = await LocalAuthentication.isEnrolledAsync();
        setIsBiometricSupported(compatible && enrolled);
    };

    const handleReveal = async () => {
        if (!isBiometricSupported) {
            setIsBlurred(false);
            setTimeout(() => setIsBlurred(true), 15000);
            return;
        }

        const result = await LocalAuthentication.authenticateAsync({
            promptMessage: 'Подтвердите личность для просмотра баланса',
            fallbackLabel: 'Использовать пароль',
            cancelLabel: 'Отмена',
        });

        if (result.success) {
            setIsBlurred(false);
            setTimeout(() => setIsBlurred(true), 15000);
        }
    };

    const formatCurrency = (value: number) => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'KGS',
            minimumFractionDigits: 2,
        }).format(value);
    };

    if (isLoading) {
        return (
            <View style={styles.container}>
                <ActivityIndicator size="small" color={colors.primary} />
            </View>
        );
    }

    if (error) {
        return (
            <View style={styles.container}>
                <Text style={styles.errorText}>Ошибка загрузки</Text>
            </View>
        );
    }

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Зарплата</Text>
                <TouchableOpacity onPress={handleReveal} style={styles.eyeButton}>
                    <Text style={styles.eyeIcon}>{isBlurred ? '👁️' : '🙈'}</Text>
                </TouchableOpacity>
            </View>

            <View style={styles.amountContainer}>
                <Text style={styles.amount}>
                    {amount !== null ? formatCurrency(amount) : '—'}
                </Text>
                {isBlurred && (
                    <View style={styles.blur}>
                        <BlurView intensity={50} tint="dark" style={StyleSheet.absoluteFill} />
                    </View>
                )}
            </View>

            <Text style={styles.hint}>
                {isBlurred
                    ? 'Нажмите 👁️ для просмотра'
                    : 'Скроется через 15 секунд'}
            </Text>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        backgroundColor: colors.surface,
        borderRadius: 16,
        padding: 20,
        marginHorizontal: 16,
        marginVertical: 8,
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: 12,
    },
    title: {
        color: colors.textSecondary,
        fontSize: 14,
        fontWeight: '500',
    },
    eyeButton: {
        padding: 8,
    },
    eyeIcon: {
        fontSize: 20,
    },
    amountContainer: {
        position: 'relative',
        overflow: 'hidden',
        borderRadius: 8,
    },
    amount: {
        color: colors.text,
        fontSize: 32,
        fontWeight: '700',
        textAlign: 'center',
        paddingVertical: 8,
    },
    blur: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
    },
    hint: {
        color: colors.textSecondary,
        fontSize: 12,
        textAlign: 'center',
        marginTop: 12,
    },
    errorText: {
        color: colors.primary,
        fontSize: 14,
        textAlign: 'center',
    },
});
