package usecase

import (
	"context"
	"fmt"
	"user-service/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Save(ctx context.Context, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("[ERROR] usecase.go:Save(): user is nil")
	}
	if user.FirstName == "" {
		return fmt.Errorf("[ERROR] usecase.go:Save(): The first name must not be empty")
	}

	if err := u.repo.Save(ctx, user); err != nil {
		return fmt.Errorf("[ERROR] usecase.go:Save():%w", err)
	}
	return nil
}

func (u *userUsecase) GetProfile(ctx context.Context, id string) (*domain.UserSalt, error) {

	user, err := u.repo.GetProfile(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "[ERROR] usecase.go:GetProfile():  %v", err)
	}
	return user, nil
	

	
}
