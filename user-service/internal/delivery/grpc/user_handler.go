package grpc

import (
	"context"
	"user-service/internal/domain"
	user "user-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	usecase domain.UserUsecase
}

func NewUserServer(usecase domain.UserUsecase) *UserServer {
	return &UserServer{
		usecase: usecase,
	}
}


func (s *UserServer) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Kullanıcı ID boş olamaz")
	}

	profile, err := s.usecase.GetProfile(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Veri tabanından kullanıcı profili alınamadı: %v", err)
	}

	return &user.GetProfileResponse{
		Id:        req.GetUserId(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}, nil

}
