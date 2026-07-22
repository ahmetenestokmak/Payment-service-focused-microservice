package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"payment-service/internal/domain"
)

type paymentUsecase struct {
	repo       domain.PaymentRepository
	strategies map[string]domain.PaymentStrategy
}

func NewUserUsecase(repo domain.PaymentRepository) domain.PaymentUseCase {
	return &paymentUsecase{
		repo:       repo,
		strategies: make(map[string]domain.PaymentStrategy),
	}
}

func (u *paymentUsecase) RegisterStrategy(method string, strategy domain.PaymentStrategy) {
	u.strategies[method] = strategy
}

func (u *paymentUsecase) ProcessPayment(ctx context.Context, payment domain.Payment) (*domain.PaymentResult, error) {
	
	strategy, exists := u.strategies[payment.PaymentMethod]
	if !exists {
		return nil, fmt.Errorf("unsupported payment method: %s", payment.PaymentMethod)
	}
	saveData, err := u.repo.Save(ctx, &payment)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] usecase.go:ProcessPayment():%w", err)
	}

	result, err := strategy.Execute(ctx, &payment)
	if result.TransactionID == "" {
		log.Printf("TransactionID boş")
	}
	// 2 saniye bekle

	time.Sleep(3 * time.Second)
	if err != nil {
		err = u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusFailed, result.TransactionID, err.Error())
		return nil, err
	}

	if !result.Success {
		return nil, u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusFailed, result.TransactionID, result.ErrorMessage)
	}

	err = u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusSuccess, result.TransactionID, result.ErrorMessage)
	if err != nil {
		return nil, err
	}


	return result, nil
}
