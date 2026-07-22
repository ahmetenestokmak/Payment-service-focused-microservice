package strategy

import (
	"context"
	//"errors"
	"payment-service/internal/domain"
)

// StripeStrategy Stripe entegrasyonu
type StripeStrategy struct {
	apiKey string
}


func NewStripeStrategy(apiKey string) *StripeStrategy {
	return &StripeStrategy{apiKey: apiKey}
}

func (s *StripeStrategy) Execute(ctx context.Context, payment *domain.Payment) (*domain.PaymentResult, error) {
	// Buraya Stripe API HTTP istek kodları gelecek
	return &domain.PaymentResult{
		TransactionID: "ch_stripe_mock_9988",
		Success:       true,
	}, nil
}
