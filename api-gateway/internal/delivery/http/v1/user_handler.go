package v1

import (
	"net/http"
	"api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient *client.UserClient
}

func NewUserHandler(userClient *client.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Middleware katmanının token'dan çözüp context'e koyduğu userID'yi alıyoruz
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bilgisi bulunamadı"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geçersiz kullanıcı kimliği formatı"})
		return
	}


	// User mikroservisine gRPC çağrısı yapılıyor
	resp, err := h.userClient.GetProfile(c.Request.Context(), userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profil bilgisi alınamadı: " + err.Error()})
		return
	}

	// gRPC'den gelen protobuf cevabı HTTP/JSON formatına dönüştürülüyor
	c.JSON(http.StatusOK, gin.H{
		"id":    userIDStr,
		"first_name":  resp.FirstName,
		"last_name": resp.LastName,
	})
}

