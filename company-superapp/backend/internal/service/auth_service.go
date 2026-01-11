package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/yourname/company-superapp/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo domain.UserRepository
	redis    *redis.Client
	jwtSecret []byte
}

func NewAuthService(userRepo domain.UserRepository, redisClient *redis.Client, secret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		redis:    redisClient,
		jwtSecret: []byte(secret),
	}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
		Role:         "user",
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthTokens, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateAndStoreTokens(ctx, user.ID, user.Role)
}

func (s *AuthService) generateAndStoreTokens(ctx context.Context, userID uuid.UUID, role string) (*AuthTokens, error) {
	// Generate Access Token (with role claim)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Store Refresh Token in Redis
	err = s.redis.Set(ctx, refreshTokenString, userID.String(), 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
