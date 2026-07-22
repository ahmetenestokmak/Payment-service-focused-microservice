package domain

// UserCreatedEvent Auth servisinin fırlattığı şema ile birebir aynı olmalı
type UserCreatedEvent struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}