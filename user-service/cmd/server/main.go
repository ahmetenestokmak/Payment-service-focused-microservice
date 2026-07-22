package main

import (
	"log"
	"net"
	"os"
	"user-service/config"
	"user-service/internal/infrastructure/broker"
	"user-service/internal/repository"
	"user-service/internal/usecase"

	deliveryGrpc "user-service/internal/delivery/grpc"
	user "user-service/proto"

	"google.golang.org/grpc"
	postgre "user-service/internal/infrastructure"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Port dinlenemedi :50053: %v", err)
	}

	db, _ := postgre.NewPostgresDB(*config.LoadConfigDB())
	repo := repository.NewUserRepository(db)
	usecase := usecase.NewUserUsecase(repo)

	rmqURL := os.Getenv("RABBITMQ_URL")
	if rmqURL == "" {
		rmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	publisher := broker.NewRabbitMQConsumer(rmqURL, usecase)
	publisher.StartConsume()
	defer publisher.Close()

	grpcServer := grpc.NewServer()
	userServer := deliveryGrpc.NewUserServer(usecase)
	user.RegisterUserServiceServer(grpcServer, userServer)

	log.Println("[INFO] User Servisi :50053 portunda gRPC dinliyor...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC sunucusu başlatılamadı: %v", err)
	}
}
