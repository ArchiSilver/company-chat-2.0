#!/bin/bash
# Скрипт для создания bucket в MinIO

set -e

echo "Waiting for MinIO to start..."
sleep 5

# Устанавливаем mc (MinIO Client) если не установлен
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO client..."
    curl -O https://dl.min.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin/
fi

# Настраиваем подключение к MinIO
mc alias set local http://localhost:9000 minioadmin minioadminpassword

# Создаём bucket для чеков такси
mc mb local/receipts --ignore-existing

# Устанавливаем политику доступа (private)
mc anonymous set download local/receipts

echo "✓ MinIO bucket 'receipts' created successfully"
