package client

import (
	"context"
	"log"

	"api-gateway/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client auth.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Auth servisine bağlanılamadı: %v", err)
	}

	return &AuthClient{
		client: auth.NewAuthServiceClient(conn),
	}
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (string, error) {
	resp, err := c.client.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *AuthClient) GetJWKS(ctx context.Context) (string, error) {
	resp, err := c.client.GetJWKS(ctx, &auth.JWKSRequest{})
	if err != nil {
		return "", err
	}
	return resp.JwksJson, nil
}

func (c *AuthClient) Register(ctx context.Context, email, password, firstName, lastName, role string) (*auth.RegisterResponse, error) {
	return c.client.Register(ctx, &auth.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	})
}

func (c *AuthClient) Update(ctx context.Context, id, email, password, firstName, lastName, role string) (*auth.RegisterResponse, error) {
	return c.client.Update(ctx, &auth.UpdateRequest{
		Id:        id,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	})
}
