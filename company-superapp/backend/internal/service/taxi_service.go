package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
	"github.com/yourname/company-superapp/internal/pkg/s3"
)

type TaxiService struct {
	taxiRepo    domain.TaxiRequestRepository
	minioClient *s3.MinioClient
}

func NewTaxiService(taxiRepo domain.TaxiRequestRepository, minioClient *s3.MinioClient) *TaxiService {
	return &TaxiService{
		taxiRepo:    taxiRepo,
		minioClient: minioClient,
	}
}

type GenerateUploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type GenerateUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileKey   string `json:"file_key"`
}

func (s *TaxiService) GenerateUploadURL(ctx context.Context, userID uuid.UUID, req GenerateUploadURLRequest) (*GenerateUploadURLResponse, error) {
	ext := filepath.Ext(req.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	fileKey := fmt.Sprintf("receipts/%s/%s%s", userID.String(), uuid.New().String(), ext)

	uploadURL, err := s.minioClient.GeneratePresignedUploadURL(ctx, fileKey, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &GenerateUploadURLResponse{
		UploadURL: uploadURL,
		FileKey:   fileKey,
	}, nil
}

type ConfirmUploadRequest struct {
	FileKey string `json:"file_key"`
}

func (s *TaxiService) ConfirmUpload(ctx context.Context, userID uuid.UUID, req ConfirmUploadRequest) (*domain.TaxiRequest, error) {
	taxiRequest := &domain.TaxiRequest{
		ID:             uuid.New(),
		UserID:         userID,
		ReceiptFileKey: req.FileKey,
		Status:         domain.TaxiStatusPending,
		CreatedAt:      time.Now(),
	}

	if err := s.taxiRepo.Create(ctx, taxiRequest); err != nil {
		return nil, err
	}

	return taxiRequest, nil
}

func (s *TaxiService) GetUserRequests(ctx context.Context, userID uuid.UUID) ([]domain.TaxiRequest, error) {
	return s.taxiRepo.GetByUserID(ctx, userID)
}

type TaxiRequestWithURL struct {
	domain.TaxiRequest
	ReceiptURL string `json:"receipt_url"`
}

func (s *TaxiService) GetUserRequestsWithURLs(ctx context.Context, userID uuid.UUID) ([]TaxiRequestWithURL, error) {
	requests, err := s.taxiRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]TaxiRequestWithURL, len(requests))
	for i, req := range requests {
		receiptURL, err := s.minioClient.GeneratePresignedDownloadURL(ctx, req.ReceiptFileKey)
		if err != nil {
			receiptURL = ""
		}
		result[i] = TaxiRequestWithURL{
			TaxiRequest: req,
			ReceiptURL:  receiptURL,
		}
	}

	return result, nil
}

func (s *TaxiService) UpdateStatus(ctx context.Context, requestID uuid.UUID, status domain.TaxiRequestStatus) error {
	return s.taxiRepo.UpdateStatus(ctx, requestID, status)
}
