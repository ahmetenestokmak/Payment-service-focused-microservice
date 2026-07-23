package strategy

import (
	"context"
	"fmt"
	"log"
	"payment-service/internal/domain"

	"payment-service/internal/infrastructure/payment/iyzico"
)

type IyzicoStrategy struct {
	client *iyzico.Client
}

func NewIyzicoStrategy(client *iyzico.Client) *IyzicoStrategy {
	return &IyzicoStrategy{
		client: client,
	}
}
func formatAmount(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}
func (i *IyzicoStrategy) Execute(ctx context.Context, payment *domain.Payment) (*domain.PaymentResult, error) {

	request := iyzico.Init3DSRequest{
		Locale:         "tr",
		ConversationID: payment.ConversationId,

		Price:       formatAmount(payment.Amount),
		PaidPrice:   formatAmount(payment.Amount),
		Currency:    payment.Currency,
		Installment: 1,

		PaymentChannel: "WEB",
		PaymentGroup:   "VIRTUAL",
		BasketID:       payment.ID,

		CallbackURL: "https://yourdomain.com/callbacklog.php?payment_id=" + payment.ID,

		PaymentCard: iyzico.PaymentCard{
			CardHolderName: payment.Card.HolderName,
			CardNumber:     payment.Card.Number,
			ExpireYear:     payment.Card.ExpireYear,
			ExpireMonth:    payment.Card.ExpireMonth,
			CVC:            payment.Card.CVC,
			RegisterCard:   0,
		},

		Buyer: iyzico.Buyer{
			ID:                  payment.Buyer.ID,
			Name:                payment.Buyer.Name,
			Surname:             payment.Buyer.Surname,
			IdentityNumber:      payment.Buyer.IdentityNumber,
			Email:               payment.Buyer.Email,
			GSMNumber:           payment.Buyer.GSMNumber,
			RegistrationDate:    payment.Buyer.RegistrationDate,
			LastLoginDate:       payment.Buyer.LastLoginDate,
			RegistrationAddress: payment.Buyer.Address,
			City:                payment.Buyer.City,
			Country:             payment.Buyer.Country,
			ZipCode:             payment.Buyer.ZipCode,
			IP:                  payment.Buyer.IP,
		},

		ShippingAddress: iyzico.Address{
			Address:     payment.ShippingAddress.Address,
			ZipCode:     payment.ShippingAddress.ZipCode,
			ContactName: payment.ShippingAddress.ContactName,
			City:        payment.ShippingAddress.City,
			Country:     payment.ShippingAddress.Country,
		},

		BillingAddress: iyzico.Address{
			Address:     payment.BillingAddress.Address,
			ZipCode:     payment.BillingAddress.ZipCode,
			ContactName: payment.BillingAddress.ContactName,
			City:        payment.BillingAddress.City,
			Country:     payment.BillingAddress.Country,
		},
	}

	for _, item := range payment.BasketItems {
		request.BasketItems = append(
			request.BasketItems,
			iyzico.BasketItem{
				ID:        item.ID,
				Name:      item.Name,
				Category1: item.Category1,
				ItemType:  item.ItemType,
				Price:     formatAmount(item.Price),
			},
		)
	}
	response, err := i.client.Initialize3DS(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] strategy:iyzico_strategy.go:Execute():Initialize3DS:%w", err)
	}
	log.Printf("execute son")
	if response == nil {
		fmt.Printf("[ERROR] strategy:execute(): response nil")
	}

	log.Printf("[DEBUG] Raw Response Data: %+v | Error: %v", response, err)
	if response.Status != "success" {
		return &domain.PaymentResult{
			TransactionID:      "",
			Status:             domain.StatusFailed,
			ErrorMessage:       response.ErrorMessage,
			ID:                 "",
			Success:            false,
			ThreeDSHTMLContent: "",
		}, nil
	}

	return &domain.PaymentResult{
		TransactionID:      response.PaymentID,
		ID:                 "",
		Status:             domain.Status3DS,
		ThreeDSHTMLContent: response.ThreeDSHTMLContent,
		ErrorMessage:       "",
		Success:            true,
	}, nil

}
