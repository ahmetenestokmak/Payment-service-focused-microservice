package usecase

import (
	"context"
	"fmt"

	//"log"

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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIK YAKALANDI: %v\n", r)
		}
	}()

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
		fmt.Printf("TransactionID boş %v", err)
	}
	result.ID = saveData.ID

	if err != nil {
		err = u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusFailed, result.TransactionID, err.Error())
		return result, err
	}

	if !result.Success {
		return result, u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusFailed, result.TransactionID, result.ErrorMessage)
	}

	err = u.repo.UpdateStatus(ctx, saveData.ID, domain.StatusSuccess, result.TransactionID, result.ErrorMessage)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *paymentUsecase) UpdateStatus(ctx context.Context, request domain.UpdateRequest) error {
	if request.Id == "" {
		return fmt.Errorf("[ERROR] usecase.go:UpdateStatus():Id zorunlu")
	}

	if request.Status == "success" {
		request.Status = domain.StatusSuccess
	}

	if request.Status != "success" {
		request.Status = domain.StatusFailed
	}

	return u.repo.UpdateStatus(ctx, 
		request.Id, 
		request.Status,
		request.TransactionId,
		"",
	)
}
