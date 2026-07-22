package v1

import (
	"api-gateway/internal/client"
	"net/http"

	"api-gateway/proto/payment"

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

// ProcessPaymentInput HTTP Gateway isteğini karşılayan ana DTO
type ProcessPaymentInput struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	UserID         string `json:"user_id" binding:"required"`

	ReferenceID   string `json:"reference_id" binding:"required"`   // 'a' harfi hatası düzeltildi
	ReferenceType string `json:"reference_type" binding:"required"` // ORDER, SUBSCRIPTION vb.

	Amount        int64  `json:"amount" binding:"required,gt=0"`    // Kuruş cinsinden
	Currency      string `json:"currency" binding:"required,len=3"` // TRY, USD vb.
	PaymentMethod string `json:"payment_method" binding:"required"` // STRIPE, IYZICO vb.

	Card *CardInput `json:"card" binding:"required_without=CardToken"`

	Buyer       BuyerInput        `json:"buyer" binding:"required"`
	BasketItems []BasketItemInput `json:"basket_items" binding:"required,gt=0,dive"`
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

	var paymentCard *payment.Card
	if input.Card != nil {
		paymentCard = &payment.Card{
			HolderName:  input.Card.HolderName,
			Number:      input.Card.Number,
			ExpireYear:  input.Card.ExpireYear,
			ExpireMonth: input.Card.ExpireMonth,
			Cvc:         input.Card.CVC,
		}
	}

	paymentBuyer := &payment.Buyer{
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

	paymentBasketItems := make([]*payment.BasketItem, 0, len(input.BasketItems))
	for _, item := range input.BasketItems {
		paymentBasketItems = append(paymentBasketItems, &payment.BasketItem{
			Id:        item.ID,
			Name:      item.Name,
			Category1: item.Category1,
			ItemType:  item.ItemType,
			Price:     item.Price,
		})
	}

	paymentReq := &payment.ProcessPaymentRequest{
		UserId:        userIDStr,
		Amount:        input.Amount,
		Currency:      input.Currency,
		ReferenceId:   input.ReferenceID,
		ReferenceType: input.ReferenceType,
		PaymentMethod: input.PaymentMethod,
		Card:          paymentCard,
		Buyer:         paymentBuyer,
		BasketItems:   paymentBasketItems,
	}

	resp, err := h.paymentClient.Create(c.Request.Context(), paymentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        resp.ErrorMessage,
		"payment_id":     resp.PaymentId,
		"status":         resp.Status,
		"transaction_id": resp.TransactionId,
		"three_ds_html_content": resp.ThreeDsHtmlContent,
	})

}
