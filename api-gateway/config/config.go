package config

import (
	"os"
)

type Config struct {
	Port           string
	AuthService    string
	UserService    string
	PaymentService string
	// Diğer servis adresleri (UserService, BankService vb.) buraya eklenebilir.
}

func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		AuthService:    getEnv("AUTH_SERVICE_ADDR", "[IP_ADDRESS]:50051"),
		UserService:    getEnv("USER_SERVICE_ADDR", "[IP_ADDRESS]:50053"),
		PaymentService: getEnv("PAYMENT_SERVICE_ADDR", "[IP_ADDRESS]:50055"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
