import { useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { useTaxiStore } from '../store/taxiStore';

const API_URL = 'http://localhost:8080/api/v1';

interface UploadState {
    isUploading: boolean;
    progress: string;
    error: string | null;
}

export function useReceiptUpload() {
    const [state, setState] = useState<UploadState>({
        isUploading: false,
        progress: '',
        error: null,
    });

    const { accessToken } = useAuthStore();
    const { addRequest } = useTaxiStore();

    const uploadReceipt = async (imageUri: string, filename: string) => {
        setState({ isUploading: true, progress: 'Получение URL для загрузки...', error: null });

        try {
            // Шаг 1: Получение presigned URL от бэкенда
            const urlResponse = await fetch(`${API_URL}/taxi/generate-upload-url`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    filename: filename,
                    content_type: 'image/jpeg',
                }),
            });

            if (!urlResponse.ok) {
                throw new Error('Не удалось получить URL для загрузки');
            }

            const { upload_url, file_key } = await urlResponse.json();

            setState((prev) => ({ ...prev, progress: 'Загрузка изображения...' }));

            // Шаг 2: Загрузка изображения напрямую в MinIO через presigned URL
            const imageResponse = await fetch(imageUri);
            const blob = await imageResponse.blob();

            const uploadResponse = await fetch(upload_url, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'image/jpeg',
                },
                body: blob,
            });

            if (!uploadResponse.ok) {
                throw new Error('Не удалось загрузить изображение');
            }

            setState((prev) => ({ ...prev, progress: 'Подтверждение загрузки...' }));

            // Шаг 3: Подтверждение загрузки на бэкенде
            const confirmResponse = await fetch(`${API_URL}/taxi/confirm-upload`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    file_key: file_key,
                }),
            });

            if (!confirmResponse.ok) {
                throw new Error('Не удалось подтвердить загрузку');
            }

            const taxiRequest = await confirmResponse.json();
            addRequest(taxiRequest);

            setState({ isUploading: false, progress: '', error: null });
            return true;
        } catch (error) {
            setState({
                isUploading: false,
                progress: '',
                error: (error as Error).message,
            });
            return false;
        }
    };

    return {
        ...state,
        uploadReceipt,
    };
}
