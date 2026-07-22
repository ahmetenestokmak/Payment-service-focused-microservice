package usecase

import (
	"auth-service/internal/domain"
	"context"
	"errors"
	

	

	"auth-service/pkg/security"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	// userClient pb.UserServiceClient // gRPC üzerinden User Service istemcisi
}

// NewAuthUsecase, iş mantığı katmanına hem veritabanını (Redis) hem de dış servis bağını (gRPC) enjekte eder.
func NewAuthUsecase(repo domain.AuthRepository /*, grpcClient pb.UserServiceClient*/) domain.AuthUsecase {
	return &authUsecase{
		authRepo: repo,
		// userClient: grpcClient,
	}
}

func (u authUsecase) Login(ctx context.Context, email string, password string) (*domain.UserSalt, error) {
	if email == "" {
		return nil, errors.New("Invalid mail.")
	}


	user, err := u.authRepo.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u authUsecase) Register(ctx context.Context, email string, password string, role string) (*domain.User, error) {
	if role == "" {
		return  nil,errors.New("Invalid role.")
	}
	if email == "" {
		return  nil,errors.New("Invalid mail.")
	}
	existingUser, err := u.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil,err
	}
	if existingUser != nil {
		return  nil,errors.New("Duplicated mail.")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return  nil,err
	}

	data, err := u.authRepo.Register(ctx, email, hashedPassword, role)
	if err != nil {
		return  nil, errors.New("Failed register: " + err.Error())
	}

	return data,nil
}



func (u authUsecase) Update(ctx context.Context,id string, email string, password string, role string) error {
	if role == "" {
		return  errors.New("Invalid role.")
	}
	if email == "" {
		return  errors.New("Invalid mail.")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return  err
	}

	err = u.authRepo.Update(ctx, id,email, hashedPassword, role)
	if err != nil {
		return  errors.New("Failed update: " + err.Error())
	}

	return nil
}
