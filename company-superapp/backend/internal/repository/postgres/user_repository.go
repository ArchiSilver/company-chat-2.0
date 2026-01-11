package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/yourname/company-superapp/internal/domain"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO system.users (email, password_hash, full_name, role) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	
	err := r.db.QueryRowxContext(ctx, query, user.Email, user.PasswordHash, user.FullName, user.Role).
		Scan(&user.ID, &user.CreatedAt)

	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, full_name, role, created_at FROM system.users WHERE email=$1`
	
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil // Or a custom not found error
	}

	return &user, err
}
