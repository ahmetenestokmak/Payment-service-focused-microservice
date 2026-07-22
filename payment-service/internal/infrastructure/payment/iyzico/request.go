package iyzico

type Init3DSRequest struct {
	Locale         string      `json:"locale"`
	ConversationID string      `json:"conversationId"`
	Price          string      `json:"price"`
	PaidPrice      string      `json:"paidPrice"`
	Currency       string      `json:"currency"`
	Installment    int         `json:"installment"`
	PaymentChannel string      `json:"paymentChannel"`
	PaymentGroup   string      `json:"paymentGroup"`
	BasketID       string      `json:"basketId"`
	CallbackURL    string      `json:"callbackUrl"`

	PaymentCard PaymentCard `json:"paymentCard"`
	Buyer       Buyer       `json:"buyer"`
	BasketItems []BasketItem `json:"basketItems"`
}

type PaymentCard struct {
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	ExpireYear     string `json:"expireYear"`
	ExpireMonth    string `json:"expireMonth"`
	CVC            string `json:"cvc"`
	RegisterCard   int    `json:"registerCard"`
}

type Buyer struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	IdentityNumber   string `json:"identityNumber"`
	Email            string `json:"email"`
	GSMNumber        string `json:"gsmNumber"`
	RegistrationDate string `json:"registrationDate"`
	LastLoginDate    string `json:"lastLoginDate"`
	RegistrationAddress string `json:"registrationAddress"`
	City             string `json:"city"`
	Country          string `json:"country"`
	ZipCode          string `json:"zipCode"`
	IP               string `json:"ip"`
}

type BasketItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Category1 string `json:"category1"`
	ItemType  string `json:"itemType"`
	Price     string `json:"price"`
}

type Auth3DSRequest struct {
	Locale           string `json:"locale"`
	PaymentId        string `json:"paymentId"`
	ConversationID   string `json:"conversationId"`
	ConversationData string `json:"conversationData"`
}