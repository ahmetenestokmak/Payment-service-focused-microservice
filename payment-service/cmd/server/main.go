package main

import (
	"context"
	"log"
	"net"
	//"os"
	"payment-service/config"
	"time"
	//"payment-service/internal/infrastructure/broker"
	"payment-service/internal/repository"
	"payment-service/internal/usecase"

	deliveryGrpc "payment-service/internal/delivery/grpc"
	pbPayment "payment-service/proto/iyzico"

	postgre "payment-service/internal/infrastructure"
	"payment-service/internal/infrastructure/payment/iyzico"
	"payment-service/internal/infrastructure/strategy"
	"payment-service/internal/interceptor"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Port dinlenemedi :50055: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       1,
	})

	// Redis Bağlantı Testi
	redisCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(redisCtx).Err(); err != nil {
		log.Printf("Redis bağlantısı kurulamadı: %v", err)
	}

	/*

		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
		fmt.Println("Successfully connected to Redis!")

		// Set
		err = rdb.Set(ctx, "framework", "Go", 1*time.Hour).Err()
		if err != nil {
			panic(err)
		}

		// Get
		val, err := rdb.Get(ctx, "framework").Result()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Key 'framework' value: %s\n", val)
	*/

	iyzi := iyzico.NewClient(
		"https://sandbox-api.iyzipay.com",
		"sandbox-",
		"",
	)

	idempotencyInterceptor := interceptor.NewIdempotencyInterceptor(rdb)

	db, _ := postgre.NewPostgresDB(*config.LoadConfigDB())
	repo := repository.NewUserRepository(db)

	usecase := usecase.NewUserUsecase(repo)

	// STRATEGY PATTERN: Ödeme yöntemlerini usecase'e kaydediyoruz
	stripeStrat := strategy.NewStripeStrategy("sk_test_stripe_key_123")
	iyzicoStrat := strategy.NewIyzicoStrategy(iyzi)

	usecase.RegisterStrategy("STRIPE", stripeStrat)
	usecase.RegisterStrategy("IYZICO", iyzicoStrat)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(idempotencyInterceptor.UnaryServerInterceptor()),
	)

	paymentServer := deliveryGrpc.NewPaymentServer(usecase)
	pbPayment.RegisterPaymentServiceServer(grpcServer, paymentServer)

	log.Println("[INFO] Payment Servisi :50055 portunda gRPC dinliyor...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC sunucusu başlatılamadı: %v", err)
	}
}
