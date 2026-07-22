package repository

import (
	"auth-service/internal/domain"
	"auth-service/pkg/security"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Login(ctx context.Context, email string, password string) (*domain.UserSalt, error) {

	query := "SELECT id, email, password, role FROM auth_users WHERE email=$1"

	var hashedPassword string

	var user domain.UserSalt
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("User not found: %w", err)
		}
		return nil, fmt.Errorf("[ERROR] repository.go:Login() Postgres: %w", err)
	}
	if user.Role != "admin" {
		return nil, fmt.Errorf("Unauthorized access")
	}
	er := security.CheckPasswordHash(password, hashedPassword)
	if !er {
		return nil, fmt.Errorf("Incorrect password.")
	}

	return &user, nil
}

func (r *authRepository) Register(ctx context.Context, email string, password string, role string) (*domain.User, error) {

	query := "INSERT INTO auth_users (email,password,role) VALUES ($1,$2,$3) RETURNING id"
	var dto domain.User
	err := r.db.QueryRowContext(ctx, query, email, password, role).Scan(&dto.ID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] repository.go:Register() Postgres: %s", err.Error())
	}
	return &dto, nil
}

func (r *authRepository) Update(ctx context.Context, id string, email string, password string, role string) error {

	query := "UPDATE auth_users SET email = $1, password = $2, role = $3 WHERE id = $4"

	_, err := r.db.ExecContext(ctx, query, email, password, role, id)
	if err != nil {
		return fmt.Errorf("[ERROR] repository.go:Update() Postgres: %s", err.Error())
	}
	return nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*domain.UserSalt, error) {
	query := `SELECT id, email, role FROM auth_users WHERE email = $1`

	var user domain.UserSalt
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
