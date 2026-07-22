package grpc

import (
	"context"

	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/broker"
	"auth-service/internal/infrastructure/security" // Eklendi
	"auth-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	publisher  *broker.RabbitMQPublisher
	jwtManager *security.JWTManager // Eklendi
	usecase  domain.AuthUsecase
}

func NewAuthServer(pub *broker.RabbitMQPublisher, jwtMng *security.JWTManager, usecase domain.AuthUsecase) *AuthServer {
	return &AuthServer{
		publisher:  pub,
		jwtManager: jwtMng,
		usecase: usecase,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email ve parola boş olamaz")
	}

	data,err := s.usecase.Register(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Kayıt başarısız: "+err.Error())
	}

	event := domain.UserCreatedEvent{
		ID:        data.ID,
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	}

	err = s.publisher.PublishUserCreated(event)
	if err != nil {
		return nil, status.Error(codes.Internal, "Kayıt kuyruğa işlenemedi")
	}

	return &auth.RegisterResponse{
		Id:      data.ID,
		Message: "Kullanıcı kaydı alındı.",
	}, nil
}
func (s *AuthServer) Update(ctx context.Context, req *auth.UpdateRequest) (*auth.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email ve parola boş olamaz")
	}

	err := s.usecase.Update(ctx, req.Id, req.Email, req.Password, req.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "Güncelleme başarısız: "+err.Error())
	}

	event := domain.UserCreatedEvent{
		ID:        req.Id,
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	}

	err = s.publisher.PublishUserCreated(event)
	if err != nil {
		return nil, status.Error(codes.Internal, "Kayıt kuyruğa işlenemedi")
	}


	return &auth.RegisterResponse{
		Id:      req.Id,
		Message: "Kullanıcı güncellendi.",
	}, nil
}

// Login Giriş isteğini işler ve RS256 imzalı JWT döner
func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	
	login, err := s.usecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	// 2. ADIM: JWT Token üret
	token, err := s.jwtManager.GenerateToken(login.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Token üretilemedi")
	}

	return &auth.LoginResponse{Token: token}, nil
}

// GetJWKS API Gateway'in public key'leri çekebilmesi için endpoint
func (s *AuthServer) GetJWKS(ctx context.Context, req *auth.JWKSRequest) (*auth.JWKSResponse, error) {
	return &auth.JWKSResponse{
		JwksJson: s.jwtManager.GetJWKSJSON(),
	}, nil
}