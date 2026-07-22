package repository

import (
	"context"
	"database/sql"
	"fmt"
	"user-service/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {

	query := "INSERT INTO users (id,first_name,last_name) VALUES ($1,$2,$3)"

	_, err := r.db.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName)

	if err != nil {
		return fmt.Errorf("[ERROR] repository.go:Save() Postgres: %w", err)
	}
	return nil
}

func (r *userRepository) GetProfile(ctx context.Context, id string) (*domain.UserSalt, error) {
	
	query := "SELECT first_name,last_name FROM users WHERE id=$1"

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "[ERROR] repository.go:GetByID() Postgres:  %v", err)
	} 
	defer rows.Close()

	var user domain.UserSalt

	rows.Next()
	if err := rows.Scan(&user.FirstName, &user.LastName); err != nil {
		return nil, status.Errorf(codes.Internal, "Veri tabanından kullanıcı profili alınamadı: %v", err)
	}

	return &user, nil
	
}
