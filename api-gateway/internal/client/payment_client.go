package client

import (
	"context"
	"log"
	//"time"

	paymentIyzico "api-gateway/proto/payment/iyzico"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client paymentIyzico.PaymentServiceClient
}

func NewPaymentClient(addr string) *PaymentClient {

	// Yeni standart olan NewClient kullanımına geçiyoruz
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Payment servisi client'ı oluşturulamadı  %v", err)
	}

	return &PaymentClient{
		client: paymentIyzico.NewPaymentServiceClient(conn),
	}

}

func (c *PaymentClient) Create(ctx context.Context, payy *paymentIyzico.ProcessPaymentRequest) (*paymentIyzico.ProcessPaymentResponse, error) {

	data, err := c.client.ProcessPayment(ctx, payy)

	if err != nil {
		log.Printf("payment servisine ProcessPayment çağrısı başarısız oldu: %v", err)
		return nil, err
	}

	return data, nil
}
func (c *PaymentClient) Update(ctx context.Context, update *paymentIyzico.ProcessUpdateStatusRequest) (*paymentIyzico.ProcessUpdateStatusResponse, error) {

	data, err := c.client.UpdateStatus(ctx, update)

	if err != nil {
		log.Printf("payment servisine UpdateStatus çağrısı başarısız oldu: %v", err)
		return nil, err
	}

	return data, nil
}
