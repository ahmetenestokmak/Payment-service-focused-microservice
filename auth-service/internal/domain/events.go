package domain

// UserCreatedEvent RabbitMQ'ya göndereceğimiz mesajın şeması
type UserCreatedEvent struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"` // Güncellendi
    LastName  string `json:"last_name"`  // Güncellendi
}