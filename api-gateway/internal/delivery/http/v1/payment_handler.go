package v1

import (
	"api-gateway/internal/client"
	"net/http"

	paymentIyzico "api-gateway/proto/payment/iyzico"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentClient *client.PaymentClient
}

func NewPaymentHandler(paymentClient *client.PaymentClient) *PaymentHandler {
	return &PaymentHandler{paymentClient: paymentClient}
}

// CardInput kart bilgilerini taşır
type CardInput struct {
	HolderName  string `json:"holder_name" binding:"required"`
	Number      string `json:"number" binding:"required,credit_card"`
	ExpireYear  string `json:"expire_year" binding:"required,len=4"`
	ExpireMonth string `json:"expire_month" binding:"required,len=2"`
	CVC         string `json:"cvc" binding:"required,min=3,max=4"`
}

// BuyerInput müşteri/alıcı bilgilerini taşır
type BuyerInput struct {
	ID               string `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Surname          string `json:"surname" binding:"required"`
	IdentityNumber   string `json:"identity_number" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	GSMNumber        string `json:"gsm_number" binding:"required"`
	RegistrationDate string `json:"registration_date"`
	LastLoginDate    string `json:"last_login_date"`
	Address          string `json:"address" binding:"required"`
	City             string `json:"city" binding:"required"`
	Country          string `json:"country" binding:"required"`
	ZipCode          string `json:"zip_code" binding:"required"`
	IP               string `json:"ip" binding:"required,ip"`
}

// BasketItemInput sepet öğelerini taşır
type BasketItemInput struct {
	ID        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Category1 string `json:"category1" binding:"required"`
	ItemType  string `json:"item_type" binding:"required"`
	Price     int64  `json:"price" binding:"required,gt=0"`
}

type Address struct {
	Address     string `json:"address" binding:"required"`
	ZipCode     string `json:"zipCode"`
	ContactName string `json:"contactName" binding:"required"`
	City        string `json:"city" binding:"required"`
	Country     string `json:"country" binding:"required"`
}

// ProcessPaymentInput HTTP Gateway isteğini karşılayan ana DTO
type ProcessPaymentInput struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	UserID         string `json:"user_id" binding:"required"`

	ReferenceID    string `json:"reference_id" binding:"required"`   // 'a' harfi hatası düzeltildi
	ReferenceType  string `json:"reference_type" binding:"required"` // ORDER, SUBSCRIPTION vb.
	ConversationId string `json:"conversation_id"`

	Amount        int64  `json:"amount" binding:"required,gt=0"`    // Kuruş cinsinden
	Currency      string `json:"currency" binding:"required,len=3"` // TRY, USD vb.
	PaymentMethod string `json:"payment_method" binding:"required"` // STRIPE, IYZICO vb.

	Card *CardInput `json:"card" binding:"required_without=CardToken"`

	Buyer       BuyerInput        `json:"buyer" binding:"required"`
	BasketItems []BasketItemInput `json:"basket_items" binding:"required,gt=0,dive"`

	ShippingAddress Address `json:"shippingAddress" binding:"required"`
	BillingAddress  Address `json:"billingAddress" binding:"required"`
}

type UpdateRequest struct {
	Status        string `json:"status"`
	Id            string `json:"id"`
	MdStatus      string `json:"mdStatus"`
	Signature     string `json:"signature"`
	TransactionId string `json:"transaction_id"`
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var input ProcessPaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json error, " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bilgisi bulunamadı"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geçersiz kullanıcı kimliği formatı"})
		return
	}

	var paymentCard *paymentIyzico.Card
	if input.Card != nil {
		paymentCard = &paymentIyzico.Card{
			HolderName:  input.Card.HolderName,
			Number:      input.Card.Number,
			ExpireYear:  input.Card.ExpireYear,
			ExpireMonth: input.Card.ExpireMonth,
			Cvc:         input.Card.CVC,
		}
	}

	paymentBuyer := &paymentIyzico.Buyer{
		Id:               input.Buyer.ID,
		Name:             input.Buyer.Name,
		Surname:          input.Buyer.Surname,
		IdentityNumber:   input.Buyer.IdentityNumber,
		Email:            input.Buyer.Email,
		GsmNumber:        input.Buyer.GSMNumber,
		RegistrationDate: input.Buyer.RegistrationDate,
		LastLoginDate:    input.Buyer.LastLoginDate,
		Address:          input.Buyer.Address,
		City:             input.Buyer.City,
		Country:          input.Buyer.Country,
		ZipCode:          input.Buyer.ZipCode,
		Ip:               input.Buyer.IP,
	}

	paymentShippingAddress := &paymentIyzico.Address{
		Address:     input.ShippingAddress.Address,
		ZipCode:     input.ShippingAddress.ZipCode,
		ContactName: input.ShippingAddress.ContactName,
		City:        input.ShippingAddress.City,
		Country:     input.ShippingAddress.Country,
	}

	paymentBillingAddress := &paymentIyzico.Address{
		Address:     input.BillingAddress.Address,
		ZipCode:     input.BillingAddress.ZipCode,
		ContactName: input.BillingAddress.ContactName,
		City:        input.BillingAddress.City,
		Country:     input.BillingAddress.Country,
	}

	paymentBasketItems := make([]*paymentIyzico.BasketItem, 0, len(input.BasketItems))
	for _, item := range input.BasketItems {
		paymentBasketItems = append(paymentBasketItems, &paymentIyzico.BasketItem{
			Id:        item.ID,
			Name:      item.Name,
			Category1: item.Category1,
			ItemType:  item.ItemType,
			Price:     item.Price,
		})
	}

	paymentReq := &paymentIyzico.ProcessPaymentRequest{
		UserId:          userIDStr,
		Amount:          input.Amount,
		Currency:        input.Currency,
		ReferenceId:     input.ReferenceID,
		ReferenceType:   input.ReferenceType,
		ConversationId:  input.ConversationId,
		PaymentMethod:   input.PaymentMethod,
		Card:            paymentCard,
		Buyer:           paymentBuyer,
		BasketItems:     paymentBasketItems,
		BillingAddress:  paymentBillingAddress,
		ShippingAddress: paymentShippingAddress,
	}

	resp, err := h.paymentClient.Create(c.Request.Context(), paymentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error_message":         resp.ErrorMessage,
		"payment_id":            resp.Id,
		"status":                resp.Status,
		"transaction_id":        resp.TransactionId,
		"three_ds_html_content": resp.ThreeDsHtmlContent,
	})

}

func (h *PaymentHandler) Update(c *gin.Context) {
	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json error, " + err.Error()})
		return
	}

	statusReq := &paymentIyzico.ProcessUpdateStatusRequest{
		Status:        input.Status,
		Id:            input.Id,
		MdStatus:      input.MdStatus,
		Signature:     input.Signature,
		TransactionId: input.TransactionId,
	}

	resp, err := h.paymentClient.Update(c.Request.Context(), statusReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error_message": resp.ErrorMessage,
		"payment_id":    resp.Id,
		"status":        resp.Status,
	})
}
