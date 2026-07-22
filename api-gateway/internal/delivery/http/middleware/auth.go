package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"api-gateway/internal/client"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type AuthMiddleware struct {
	authClient *client.AuthClient
	keySet     jwk.Set
	mu         sync.RWMutex
	lastUpdate time.Time
}

func NewAuthMiddleware(authClient *client.AuthClient) *AuthMiddleware {
	m := &AuthMiddleware{authClient: authClient}
	m.refreshJWKS() // İlk açılışta anahtarları çekiyoruz
	return m
}

// refreshJWKS Auth servisinden JWKS'i gRPC üzerinden çeker ve parse eder
func (m *AuthMiddleware) refreshJWKS() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 5 dakikada bir güncelleme kontrolü (Throttling)
	if time.Since(m.lastUpdate) < 5*time.Minute && m.keySet != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jwksJSON, err := m.authClient.GetJWKS(ctx)
	if err != nil {
		return // Hata durumunda eski keySet kullanılmaya devam eder
	}

	set, err := jwk.ParseString(jwksJSON)
	if err == nil {
		m.keySet = set
		m.lastUpdate = time.Now()
	}
}

func (m *AuthMiddleware) CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header eksik"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token formatı"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// İmza algoritması kontrolü
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("beklenmeyen imza metodu")
			}

			// JWKS içinden ilgili anahtarı bulma (kid eşleştirmesi)
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("token içinde 'kid' bulunamadı")
			}

			m.mu.RLock()
			key, found := m.keySet.LookupKeyID(kid)
			m.mu.RUnlock()

			if !found {
				// Anahtar bulunamadıysa JWKS güncellenip tekrar denenir
				m.refreshJWKS()
				m.mu.RLock()
				key, found = m.keySet.LookupKeyID(kid)
				m.mu.RUnlock()
				if !found {
					return nil, errors.New("anahtar JWKS içinde bulunamadı")
				}
			}

			var rawKey interface{}
			if err := key.Raw(&rawKey); err != nil {
				return nil, err
			}
			return rawKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		// Token geçerli ise claims bilgilerini context'e ekle
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["sub"])
		}

		c.Next()
	}
}
