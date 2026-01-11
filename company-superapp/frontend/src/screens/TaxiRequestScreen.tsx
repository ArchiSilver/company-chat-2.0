import * as ImagePicker from 'expo-image-picker';
import { useEffect } from 'react';
import {
    ActivityIndicator,
    Alert,
    FlatList,
    Image,
    SafeAreaView,
    StyleSheet,
    Text,
    TouchableOpacity,
    View,
} from 'react-native';
import { useReceiptUpload } from '../hooks/useReceiptUpload';
import { useTaxiStore } from '../store/taxiStore';
import { colors } from '../theme/colors';

export default function TaxiRequestScreen() {
    const { requests, isLoading, fetchRequests } = useTaxiStore();
    const { isUploading, progress, error, uploadReceipt } = useReceiptUpload();

    useEffect(() => {
        fetchRequests();
    }, []);

    const handlePickImage = async () => {
        const permissionResult = await ImagePicker.requestMediaLibraryPermissionsAsync();

        if (!permissionResult.granted) {
            Alert.alert('Ошибка', 'Необходимо разрешение на доступ к галерее');
            return;
        }

        const result = await ImagePicker.launchImageLibraryAsync({
            mediaTypes: ImagePicker.MediaTypeOptions.Images,
            allowsEditing: true,
            quality: 0.8,
        });

        if (!result.canceled && result.assets[0]) {
            const asset = result.assets[0];
            const filename = asset.fileName || `receipt_${Date.now()}.jpg`;
            await uploadReceipt(asset.uri, filename);
        }
    };

    const handleTakePhoto = async () => {
        const permissionResult = await ImagePicker.requestCameraPermissionsAsync();

        if (!permissionResult.granted) {
            Alert.alert('Ошибка', 'Необходимо разрешение на доступ к камере');
            return;
        }

        const result = await ImagePicker.launchCameraAsync({
            allowsEditing: true,
            quality: 0.8,
        });

        if (!result.canceled && result.assets[0]) {
            const asset = result.assets[0];
            const filename = `receipt_${Date.now()}.jpg`;
            await uploadReceipt(asset.uri, filename);
        }
    };

    const showUploadOptions = () => {
        Alert.alert(
            'Загрузить чек',
            'Выберите способ загрузки',
            [
                { text: 'Камера', onPress: handleTakePhoto },
                { text: 'Галерея', onPress: handlePickImage },
                { text: 'Отмена', style: 'cancel' },
            ]
        );
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'approved':
                return '#4CAF50';
            case 'rejected':
                return '#F44336';
            default:
                return '#FFC107';
        }
    };

    const getStatusText = (status: string) => {
        switch (status) {
            case 'approved':
                return 'Одобрено';
            case 'rejected':
                return 'Отклонено';
            default:
                return 'На рассмотрении';
        }
    };

    const renderRequestItem = ({ item }: { item: typeof requests[0] }) => (
        <View style={styles.requestCard}>
            {item.receipt_url && (
                <Image
                    source={{ uri: item.receipt_url }}
                    style={styles.receiptImage}
                    resizeMode="cover"
                />
            )}
            <View style={styles.requestInfo}>
                <View style={[styles.statusBadge, { backgroundColor: getStatusColor(item.status) }]}>
                    <Text style={styles.statusText}>{getStatusText(item.status)}</Text>
                </View>
                <Text style={styles.dateText}>
                    {new Date(item.created_at).toLocaleDateString('ru-RU', {
                        day: 'numeric',
                        month: 'long',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                    })}
                </Text>
            </View>
        </View>
    );

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Такси</Text>
                <Text style={styles.subtitle}>Загрузите чек для компенсации</Text>
            </View>

            <TouchableOpacity
                style={[styles.uploadButton, isUploading && styles.uploadButtonDisabled]}
                onPress={showUploadOptions}
                disabled={isUploading}
            >
                {isUploading ? (
                    <View style={styles.uploadingContainer}>
                        <ActivityIndicator size="small" color={colors.text} />
                        <Text style={styles.uploadButtonText}>{progress}</Text>
                    </View>
                ) : (
                    <>
                        <Text style={styles.uploadIcon}>📷</Text>
                        <Text style={styles.uploadButtonText}>Загрузить чек</Text>
                    </>
                )}
            </TouchableOpacity>

            {error && (
                <View style={styles.errorContainer}>
                    <Text style={styles.errorText}>{error}</Text>
                </View>
            )}

            <View style={styles.listHeader}>
                <Text style={styles.listTitle}>История заявок</Text>
            </View>

            {isLoading ? (
                <ActivityIndicator size="large" color={colors.primary} style={styles.loader} />
            ) : requests.length === 0 ? (
                <View style={styles.emptyContainer}>
                    <Text style={styles.emptyText}>Нет заявок</Text>
                    <Text style={styles.emptySubtext}>Загрузите чек, чтобы создать заявку</Text>
                </View>
            ) : (
                <FlatList
                    data={requests}
                    keyExtractor={(item) => item.id}
                    renderItem={renderRequestItem}
                    contentContainerStyle={styles.listContent}
                    showsVerticalScrollIndicator={false}
                />
            )}
        </SafeAreaView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    header: {
        paddingHorizontal: 16,
        paddingTop: 20,
        paddingBottom: 16,
    },
    title: {
        color: colors.text,
        fontSize: 28,
        fontWeight: '700',
    },
    subtitle: {
        color: colors.textSecondary,
        fontSize: 14,
        marginTop: 4,
    },
    uploadButton: {
        backgroundColor: colors.primary,
        marginHorizontal: 16,
        paddingVertical: 16,
        borderRadius: 12,
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
    },
    uploadButtonDisabled: {
        opacity: 0.7,
    },
    uploadingContainer: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    uploadIcon: {
        fontSize: 20,
        marginRight: 8,
    },
    uploadButtonText: {
        color: colors.text,
        fontSize: 16,
        fontWeight: '600',
        marginLeft: 8,
    },
    errorContainer: {
        marginHorizontal: 16,
        marginTop: 12,
        padding: 12,
        backgroundColor: 'rgba(244, 67, 54, 0.1)',
        borderRadius: 8,
    },
    errorText: {
        color: '#F44336',
        fontSize: 14,
        textAlign: 'center',
    },
    listHeader: {
        paddingHorizontal: 16,
        paddingTop: 24,
        paddingBottom: 12,
    },
    listTitle: {
        color: colors.text,
        fontSize: 18,
        fontWeight: '600',
    },
    loader: {
        marginTop: 40,
    },
    emptyContainer: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
        paddingBottom: 100,
    },
    emptyText: {
        color: colors.textSecondary,
        fontSize: 16,
    },
    emptySubtext: {
        color: colors.textSecondary,
        fontSize: 14,
        marginTop: 4,
    },
    listContent: {
        paddingHorizontal: 16,
        paddingBottom: 20,
    },
    requestCard: {
        backgroundColor: colors.surface,
        borderRadius: 12,
        marginBottom: 12,
        overflow: 'hidden',
    },
    receiptImage: {
        width: '100%',
        height: 150,
    },
    requestInfo: {
        padding: 12,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    statusBadge: {
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 16,
    },
    statusText: {
        color: colors.text,
        fontSize: 12,
        fontWeight: '600',
    },
    dateText: {
        color: colors.textSecondary,
        fontSize: 12,
    },
});
