package grpc

import (
	"context"
	"payment-service/internal/domain"
	payment "payment-service/proto/iyzico"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	payment.UnimplementedPaymentServiceServer
	usecase domain.PaymentUseCase
}

func NewPaymentServer(usecase domain.PaymentUseCase) *PaymentServer {
	return &PaymentServer{
		usecase: usecase,
	}
}

func (p *PaymentServer) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Kullanıcı ID boş olamaz")
	}

	data, err := p.usecase.ProcessPayment(ctx, domain.Payment{
		UserID:        req.GetUserId(),
		ReferenceID:   req.GetReferenceId(),
		ReferenceType: req.GetReferenceType(),
		Amount:        req.GetAmount(),
		Currency:      req.GetCurrency(),
		PaymentMethod: req.GetPaymentMethod(),
		Status:        domain.StatusPending,
		Card:          mapCardToDomain(req.GetCard()),
		Buyer:         mapBuyerToDomain(req.GetBuyer()),
		BasketItems:   mapBasketItemsToDomain(req.GetBasketItems()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ödeme işlenemedi: %v", err)
	}

	return &payment.ProcessPaymentResponse{
		PaymentId:     data.PaymentID,
		TransactionId: data.TransactionID,
		Status:        string(data.Status),
		ThreeDsHtmlContent: data.ThreeDSHTMLContent,
		ErrorMessage:  data.ErrorMessage,
	}, nil
}


func mapCardToDomain(card *payment.Card) domain.Card {
	if card == nil {
		return domain.Card{}
	}
	return domain.Card{
		HolderName:  card.GetHolderName(),
		Number:      card.GetNumber(),
		ExpireYear:  card.GetExpireYear(),
		ExpireMonth: card.GetExpireMonth(),
		CVC:         card.GetCvc(),
	}
}

func mapBuyerToDomain(buyer *payment.Buyer) domain.Buyer {
	if buyer == nil {
		return domain.Buyer{}
	}
	return domain.Buyer{
		ID:               buyer.GetId(),
		Name:             buyer.GetName(),
		Surname:          buyer.GetSurname(),
		IdentityNumber:   buyer.GetIdentityNumber(),
		Email:            buyer.GetEmail(),
		GSMNumber:        buyer.GetGsmNumber(),
		RegistrationDate: buyer.GetRegistrationDate(),
		LastLoginDate:    buyer.GetLastLoginDate(),
		Address:          buyer.GetAddress(),
		City:             buyer.GetCity(),
		Country:          buyer.GetCountry(),
		ZipCode:          buyer.GetZipCode(),
		IP:               buyer.GetIp(),
	}
}

func mapBasketItemsToDomain(items []*payment.BasketItem) []domain.BasketItem {
	if len(items) == 0 {
		return nil
	}
	domainItems := make([]domain.BasketItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		domainItems = append(domainItems, domain.BasketItem{
			ID:        item.GetId(),
			Name:      item.GetName(),
			Category1: item.GetCategory1(),
			ItemType:  item.GetItemType(),
			Price:     item.GetPrice(),
		})
	}
	return domainItems
}

func mapPaymentResultToProto(res *domain.PaymentResult) *payment.ProcessPaymentResponse {
	if res == nil {
		return nil
	}
	return &payment.ProcessPaymentResponse{
		TransactionId:      res.TransactionID,
		Success:            res.Success,
		ThreeDsHtmlContent: res.ThreeDSHTMLContent,
		PaymentId:          res.PaymentID,
		Status:             string(res.Status),
		ErrorMessage:       res.ErrorMessage,
	}
}