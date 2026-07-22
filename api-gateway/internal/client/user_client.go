package client

import (
	"context"
	"log"
	//"time"

	"api-gateway/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client user.UserServiceClient
}

func NewUserClient(addr string) *UserClient {

	// Yeni standart olan NewClient kullanımına geçiyoruz
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: User servisi client'ı oluşturulamadı  %v", err)
	}

	return &UserClient{
		client: user.NewUserServiceClient(conn),
	}

}

func (c *UserClient) GetProfile(ctx context.Context, userID string) (*user.GetProfileResponse, error) {
	user, err := c.client.GetProfile(ctx, &user.GetProfileRequest{UserId: userID})

	if err != nil {
		log.Printf("User servisine GetProfile çağrısı başarısız oldu: %v", err)
		return nil, err
	}

	return user, nil
}
