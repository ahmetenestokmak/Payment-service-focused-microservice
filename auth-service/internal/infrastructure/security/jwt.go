package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
	jwksJSON   string
}

func NewJWTManager() (*JWTManager, error) {
	// 1. 2048 bitlik RSA Anahtar Çifti Üret
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("RSA anahtarı üretilemedi: %w", err)
	}

	pubKey := &privKey.PublicKey
	kid := "auth-service-key-v1" // Anahtarın benzersiz kimliği

	// 2. Public Key'i JWK formatına dönüştür
	jwkKey, err := jwk.FromRaw(pubKey)
	if err != nil {
		return nil, fmt.Errorf("JWK dönüştürme hatası: %w", err)
	}

	_ = jwkKey.Set(jwk.KeyIDKey, kid)
	_ = jwkKey.Set(jwk.AlgorithmKey, "RS256")
	_ = jwkKey.Set(jwk.KeyUsageKey, "sig")

	// JWK Set oluştur (Gateway bunu okuyacak)
	jwks := jwk.NewSet()
	_ = jwks.AddKey(jwkKey)

	jwksBytes, err := json.Marshal(jwks)
	if err != nil {
		return nil, fmt.Errorf("JWKS JSON marshal hatası: %w", err)
	}

	return &JWTManager{
		privateKey: privKey,
		publicKey:  pubKey,
		keyID:      kid,
		jwksJSON:   string(jwksBytes),
	}, nil
}

// GenerateToken Kullanıcı giriş yaptığında token üretir
func (m *JWTManager) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 1 günlük ömür
		"iat": time.Now().Unix(),
		"iss": "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.keyID // Gateway hangi anahtarla doğrulayacağını bilsin

	return token.SignedString(m.privateKey)
}

// GetJWKSJSON Gateway'in çekeceği public key setini döner
func (m *JWTManager) GetJWKSJSON() string {
	return m.jwksJSON
}