package security

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type TokenValidator struct {
	jwks jwk.Set
}

// NewTokenValidator Auth servisinden gelen JWKS JSON string'ini parse ederek validator oluşturur
func NewTokenValidator(jwksJSON string) (*TokenValidator, error) {
	set, err := jwk.ParseString(jwksJSON)
	if err != nil {
		return nil, fmt.Errorf("JWKS parse hatası: %w", err)
	}
	return &TokenValidator{jwks: set}, nil
}

// ValidateToken Gelen JWT token'ı JWKS içindeki public key'ler ile doğrular ve userID (sub) döner
func (v *TokenValidator) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Algoritma kontrolü
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("beklenmeyen imza yöntemi: %v", token.Header["alg"])
		}

		// Token başlığından anahtar kimliğini (kid) al
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token başlığında 'kid' bulunamadı")
		}

		// JWKS içinde bu kid'e ait public key'i bul
		key, found := v.jwks.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("JWKS içinde anahtar bulunamadı: %s", kid)
		}

		var rawKey interface{}
		if err := key.Raw(&rawKey); err != nil {
			return nil, fmt.Errorf("raw key alınamadı: %w", err)
		}

		return rawKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("geçersiz token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token claims okunamadı")
	}

	// Kullanıcı ID'sini (subject) döndür
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("token içinde 'sub' alanı bulunamadı")
	}

	return userID, nil
}