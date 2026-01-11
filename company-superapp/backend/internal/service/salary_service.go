package service

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/yourname/company-superapp/internal/domain"
	"github.com/yourname/company-superapp/internal/pkg/encryption"
)

type SalaryService struct {
	salaryRepo        domain.SalaryRepository
	encryptionService *encryption.EncryptionService
}

func NewSalaryService(salaryRepo domain.SalaryRepository, encryptionService *encryption.EncryptionService) *SalaryService {
	return &SalaryService{
		salaryRepo:        salaryRepo,
		encryptionService: encryptionService,
	}
}

type SalaryResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	UpdatedAt string    `json:"updated_at"`
}

func (s *SalaryService) GetSalary(ctx context.Context, userID uuid.UUID) (*SalaryResponse, error) {
	salary, err := s.salaryRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if salary == nil {
		return &SalaryResponse{
			UserID: userID,
			Amount: 0,
		}, nil
	}

	decrypted, err := s.encryptionService.Decrypt(salary.AmountEncrypted)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(string(decrypted), 64)
	if err != nil {
		return nil, err
	}

	return &SalaryResponse{
		UserID:    salary.UserID,
		Amount:    amount,
		UpdatedAt: salary.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *SalaryService) UpdateSalary(ctx context.Context, userID uuid.UUID, amount float64) error {
	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)

	encrypted, err := s.encryptionService.Encrypt([]byte(amountStr))
	if err != nil {
		return err
	}

	salary := &domain.Salary{
		UserID:          userID,
		AmountEncrypted: encrypted,
	}

	return s.salaryRepo.Upsert(ctx, salary)
}
