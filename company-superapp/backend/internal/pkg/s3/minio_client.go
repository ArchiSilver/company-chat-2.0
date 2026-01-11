package s3

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient() (*MinioClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadminpassword"
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "taxi-receipts"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Проверяем существование bucket
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", bucketName)
	}

	return &MinioClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (m *MinioClient) GeneratePresignedUploadURL(ctx context.Context, objectKey string, contentType string) (string, error) {
	expiry := 15 * time.Minute

	reqParams := make(url.Values)
	reqParams.Set("Content-Type", contentType)

	presignedURL, err := m.client.PresignedPutObject(ctx, m.bucketName, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

func (m *MinioClient) GeneratePresignedDownloadURL(ctx context.Context, objectKey string) (string, error) {
	expiry := 1 * time.Hour

	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	return presignedURL.String(), nil
}

func (m *MinioClient) DeleteObject(ctx context.Context, objectKey string) error {
	return m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
}
