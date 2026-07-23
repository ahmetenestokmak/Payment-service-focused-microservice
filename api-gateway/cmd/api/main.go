package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/config"
	"api-gateway/internal/client"
	"api-gateway/internal/delivery/http/middleware"
	"api-gateway/internal/delivery/http/v1"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env dosyası bulunamadı, sistem ortam değişkenleri kullanılacak.")
	}

	cfg := config.LoadConfig()

	// 1. gRPC Client'larını initialize et
	authClient := client.NewAuthClient(cfg.AuthService)
	userClient := client.NewUserClient(cfg.UserService)
	paymentClient := client.NewPaymentClient(cfg.PaymentService)

	// 2. Middleware ve Handler'ları oluştur
	authMiddleware := middleware.NewAuthMiddleware(authClient)
	authHandler := v1.NewAuthHandler(authClient)
	userHandler := v1.NewUserHandler(userClient)
	paymentHandler := v1.NewPaymentHandler(paymentClient)

	// 3. HTTP Router ayarla
	r := gin.Default()

	// Ana API Grubu
	v1Group := r.Group("/api/v1")
	{
		// Public Routes
		v1Group.POST("/auth/login", authHandler.Login)
		v1Group.POST("/auth/register", authHandler.Register)

		// Protected Routes
		protectedGroup := v1Group.Group("/")
		protectedGroup.Use(authMiddleware.CheckJWT())
		{
			protectedGroup.POST("/payment/create", paymentHandler.Create)
			protectedGroup.POST("/payment/update", paymentHandler.Update)


			protectedGroup.PUT("/auth/update", authHandler.Update)

			protectedGroup.GET("/users/profile", userHandler.GetUserProfile)
		}
	}

	// 4. Graceful Shutdown ve Sunucu Başlatma Kurulumu
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Sunucuyu ayrı bir goroutine'de başlatıyoruz ki main thread kilitlenmesin
	go func() {
		log.Printf("API Gateway %s portunda çalışıyor...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server başlatılamadı: %v", err)
		}
	}()

	// Kapatma sinyallerini bekleme (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("API Gateway kapatılıyor...")

	// Kapatma işlemi için 5 saniyelik tolerans süresi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server zorla kapatıldı: %v", err)
	}

	log.Println("API Gateway tamamen durduruldu.")
}