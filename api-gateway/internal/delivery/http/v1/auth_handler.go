package v1

import (
	"net/http"

	"api-gateway/internal/client"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authClient *client.AuthClient
}

func NewAuthHandler(authClient *client.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authClient.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Mevcut AuthHandler altına eklenecek struct ve fonksiyon:
type RegisterInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=4"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required"`
}
type UpdateInput struct {
	Id        string `json:"id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=4"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json error, " + err.Error()})
		return
	}

	resp, err := h.authClient.Register(c.Request.Context(), input.Email, input.Password, input.FirstName, input.LastName, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kayıt işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kullanıcı başarıyla kaydedildi",
		"id":      resp.Id,
	})
}
func (h *AuthHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bilgisi bulunamadı (id)"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Geçersiz kullanıcı kimliği formatı"})
		return
	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json error, " + err.Error()})
		return
	}

	resp, err := h.authClient.Update(c.Request.Context(), userIDStr, input.Email, input.Password, input.FirstName, input.LastName, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncelleme işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kullanıcı başarıyla güncellendi",
		"id":      resp.Id,
	})
}
