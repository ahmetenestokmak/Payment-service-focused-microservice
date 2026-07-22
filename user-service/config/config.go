package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LoadConfig, .env dosyasını okur ve ayarları struct olarak geri döner.
func LoadConfigDB() *ConfigDB {
	// .env dosyasını yüklüyoruz. Eğer dosya yoksa sistem patlamasın,
	// işletim sisteminin kendi env değişkenlerine baksın diye hatayı loglayıp geçiyoruz.
	_ = godotenv.Load()

	return &ConfigDB{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		User:     getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Name:     getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
		return os.Getenv(key)
	
}
