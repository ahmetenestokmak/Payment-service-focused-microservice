package client

import (
	"context"
	"log"
	//"time"

	payment "api-gateway/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client payment.PaymentServiceClient
}

func NewPaymentClient(addr string) *PaymentClient {

	// Yeni standart olan NewClient kullanımına geçiyoruz
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Payment servisi client'ı oluşturulamadı  %v", err)
	}

	return &PaymentClient{
		client: payment.NewPaymentServiceClient(conn),
	}

}

func (c *PaymentClient) Create(ctx context.Context, payy *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {

	user, err := c.client.ProcessPayment(ctx, payy)

	if err != nil {
		log.Printf("payment servisine ProcessPayment çağrısı başarısız oldu: %v", err)
		return nil, err
	}

	return user, nil
}
