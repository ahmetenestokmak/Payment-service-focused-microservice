package iyzico

type Init3DSResponse struct {
	Status              string `json:"status"`
	Locale              string `json:"locale"`
	SystemTime          int64  `json:"systemTime"`
	ConversationID      string `json:"conversationId"`
	ThreeDSHTMLContent  string `json:"threeDSHtmlContent"`

	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorGroup   string `json:"errorGroup"`

	PaymentID string `json:"paymentId"`
}

type Auth3DSResponse struct {
	Status              string `json:"status"`
	Locale              string `json:"locale"`
	SystemTime          int64  `json:"systemTime"`
	ConversationID      string `json:"conversationId"`

	PaymentID            string `json:"paymentId"`
	PaymentTransactionID string `json:"paymentTransactionId"`

	Price     string `json:"price"`
	PaidPrice string `json:"paidPrice"`
	Currency  string `json:"currency"`

	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorGroup   string `json:"errorGroup"`
}