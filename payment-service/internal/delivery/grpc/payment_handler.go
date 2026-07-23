package grpc

import (
	"context"
	"payment-service/internal/domain"
	paymentIyzico "payment-service/proto/iyzico"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	paymentIyzico.UnimplementedPaymentServiceServer
	usecase domain.PaymentUseCase
}

func NewPaymentServer(usecase domain.PaymentUseCase) *PaymentServer {
	return &PaymentServer{
		usecase: usecase,
	}
}

func (p *PaymentServer) ProcessPayment(ctx context.Context, req *paymentIyzico.ProcessPaymentRequest) (*paymentIyzico.ProcessPaymentResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Kullanıcı ID boş olamaz")
	}

	data, err := p.usecase.ProcessPayment(ctx, domain.Payment{
		UserID:          req.GetUserId(),
		ReferenceID:     req.GetReferenceId(),
		ReferenceType:   req.GetReferenceType(),
		ConversationId:  req.GetConversationId(),
		Amount:          req.GetAmount(),
		Currency:        req.GetCurrency(),
		PaymentMethod:   req.GetPaymentMethod(),
		Status:          domain.StatusPending,
		Card:            mapCardToDomain(req.GetCard()),
		Buyer:           mapBuyerToDomain(req.GetBuyer()),
		ShippingAddress: mapAddressToDomain(req.GetShippingAddress()),
		BillingAddress:  mapAddressToDomain(req.GetBillingAddress()),
		BasketItems:     mapBasketItemsToDomain(req.GetBasketItems()),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ödeme işlenemedi: %v", err)
	}
	if data == nil {
		return &paymentIyzico.ProcessPaymentResponse{
			Id:                 "-",
			TransactionId:      "-",
			Status:             "-",
			Success:            true,
			ThreeDsHtmlContent: "-",
			ErrorMessage:       "-",
		}, nil
	}

	return &paymentIyzico.ProcessPaymentResponse{
		Id:                 data.ID,
		TransactionId:      data.TransactionID,
		Status:             string(data.Status),
		Success:            true,
		ThreeDsHtmlContent: data.ThreeDSHTMLContent,
		ErrorMessage:       data.ErrorMessage,
	}, nil

}

func (p *PaymentServer) UpdateStatus(ctx context.Context, req *paymentIyzico.ProcessUpdateStatusRequest) (*paymentIyzico.ProcessUpdateStatusResponse, error) {
	err := p.usecase.UpdateStatus(ctx,
		domain.UpdateRequest{
			Status:        domain.PaymentStatus(req.Status),
			Id:            req.GetId(),
			MdStatus:      req.GetMdStatus(),
			Signature:     req.GetSignature(),
			TransactionId: req.GetTransactionId(),
		},
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ödeme durumu işlenemedi: %v", err)
	}

	return &paymentIyzico.ProcessUpdateStatusResponse{
		Id:           req.GetId(),
		Status:       string(domain.StatusSuccess),
		ErrorMessage: "",
	}, nil
}

func mapCardToDomain(card *paymentIyzico.Card) domain.Card {
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

func mapBuyerToDomain(buyer *paymentIyzico.Buyer) domain.Buyer {
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
func mapAddressToDomain(address *paymentIyzico.Address) domain.Address {
	if address == nil {
		return domain.Address{}
	}
	return domain.Address{
		Address:     address.Address,
		ZipCode:     address.ZipCode,
		ContactName: address.ContactName,
		City:        address.City,
		Country:     address.Country,
	}
}

func mapBasketItemsToDomain(items []*paymentIyzico.BasketItem) []domain.BasketItem {
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

func mapPaymentResultToProto(res *domain.PaymentResult) *paymentIyzico.ProcessPaymentResponse {
	if res == nil {
		return nil
	}
	return &paymentIyzico.ProcessPaymentResponse{
		TransactionId:      res.TransactionID,
		Success:            res.Success,
		ThreeDsHtmlContent: res.ThreeDSHTMLContent,
		Id:          		res.ID,
		Status:             string(res.Status),
		ErrorMessage:       res.ErrorMessage,
	}
}
