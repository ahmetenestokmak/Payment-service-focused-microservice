package domain

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	StatusPending PaymentStatus = "PENDING"
	Status3DS     PaymentStatus = "3DS"
	StatusSuccess PaymentStatus = "SUCCESS"
	StatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`

	ReferenceID    string `json:"reference_id" db:"reference_id"`     // Sipariş ID vb.
	ReferenceType  string `json:"reference_type" db:"reference_type"` // ORDER, SUBSCRIPTION
	ConversationId string `json:"conversation_id" db:"conversation_id"`

	Amount   int64  `json:"amount" db:"amount"`     // Kuruş cinsinden (Cents)
	Currency string `json:"currency" db:"currency"` // ISO 4217 formatında

	Card Card

	Buyer Buyer

	BasketItems []BasketItem

	ShippingAddress Address `json:"shippingAddress" db:"shippingAddress"`
	BillingAddress  Address `json:"billingAddress" db:"billingAddress"`

	PaymentMethod string        `json:"payment_method" db:"payment_method"` // STRIPE, IYZICO vb.
	Status        PaymentStatus `json:"status" db:"status"`

	TransactionID string `json:"transaction_id" db:"transaction_id"`
	FailureReason string `json:"failure_reason" db:"failure_reason"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Card struct {
	HolderName  string
	Number      string
	ExpireYear  string
	ExpireMonth string
	CVC         string
}
type Address struct {
	Address     string `json:"address"`
	ZipCode     string `json:"zipCode"`
	ContactName string `json:"contactName"`
	City        string `json:"city"`
	Country     string `json:"country"`
}
type Buyer struct {
	ID               string
	Name             string
	Surname          string
	IdentityNumber   string
	Email            string
	GSMNumber        string
	RegistrationDate string
	LastLoginDate    string
	Address          string
	City             string
	Country          string
	ZipCode          string
	IP               string
}

type BasketItem struct {
	ID        string
	Name      string
	Category1 string
	ItemType  string
	Price     int64
}

type PaymentResult struct {
	TransactionID string `json:"transaction_id"`
	Success       bool   `json:"success"`

	ThreeDSHTMLContent string `json:"three_ds_html_content"`
	ID                 string `json:"id"`

	Status       PaymentStatus `json:"status"`
	ErrorMessage string        `json:"error_message"`
}

type UpdateRequest struct {
	Status        PaymentStatus `json:"status"`
	Id            string        `json:"id"`
	MdStatus      string        `json:"mdStatus"`
	Signature     string        `json:"signature"`
	TransactionId string        `json:"transaction_id"`
}

type PaymentStrategy interface {
	Execute(ctx context.Context, payment *Payment) (*PaymentResult, error)
}

type PaymentRepository interface {
	Save(ctx context.Context, payment *Payment) (*Payment, error)
	GetByID(ctx context.Context, id string) (*Payment, error)
	UpdateStatus(ctx context.Context, id string, status PaymentStatus, transactionID string, failureReason string) error
}

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, payment Payment) (*PaymentResult, error)
	RegisterStrategy(method string, strategy PaymentStrategy)
	UpdateStatus(ctx context.Context, request UpdateRequest) error
}
