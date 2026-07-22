package main

import (
	"log"
	"net"
	"os"
"auth-service/internal/repository"
"auth-service/internal/usecase"
"auth-service/config"

	deliveryGrpc "auth-service/internal/delivery/grpc"
	"auth-service/internal/infrastructure/broker"
	"auth-service/internal/infrastructure/security"
	postgre "auth-service/internal/infrastructure"
	auth "auth-service/proto"

	"google.golang.org/grpc"
)

func main() {
	// 1. JWT Manager başlat (Asimetrik anahtarlar RAM'de üretilir)
	jwtMng, err := security.NewJWTManager()
	if err != nil {
		log.Fatalf("JWT Manager başlatılamadı: %v", err)
	}

	// 2. RabbitMQ URL bilgisini al ve Publisher'ı başlat
	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		rmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	publisher := broker.NewRabbitMQPublisher(rmqURL)
	defer publisher.Close()


	// 3. gRPC Sunucusunu Yapılandır
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Port dinlenemedi :50051: %v", err)
	}


	db,_ := postgre.NewPostgresDB(*config.LoadConfigDB())
	repo := repository.NewAuthRepository(db)
	usecase := usecase.NewAuthUsecase(repo)

	grpcServer := grpc.NewServer()
	authServer := deliveryGrpc.NewAuthServer(publisher, jwtMng, usecase)
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	log.Println("[INFO] Auth Servisi :50051 portunda gRPC dinliyor...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC sunucusu başlatılamadı: %v", err)
	}
}