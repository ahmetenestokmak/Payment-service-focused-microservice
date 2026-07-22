package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTClaims token içerisine gömeceğimiz özel veriler
type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT kullanıcı için 15 dakika ömürlü bir Access Token üretir (RS256)
func GenerateJWT(id string, email string, privateKey *rsa.PrivateKey) (string, int64, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &JWTClaims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   id, // Standart id alanı
		},
	}

	// RS256 algoritması seçildi
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

// GenerateRefreshToken 32 byte'lık rastgele bir string üretir.
// Bu token'ı, user_id ve expire_time ile birlikte Redis veya DB'ye kaydetmelisin.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

