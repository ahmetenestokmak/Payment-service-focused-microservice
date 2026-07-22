package broker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"auth-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQPublisher(url string) *RabbitMQPublisher {
	// RabbitMQ Sunucusuna bağlan
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
	}

	// İletişim kanalı aç
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ kanalı açılamadı: %v", err)
	}

	// "user.events" adında bir Exchange (dağıtıcı) tanımla
	err = ch.ExchangeDeclare(
		"user.events", // exchange adı
		"direct",      // tipi
		true,          // kalıcı mı (durable)
		false,         // otomatik silinsin mi
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Exchange declare hatası: %v", err)
	}

	return &RabbitMQPublisher{conn: conn, channel: ch}
}

func (p *RabbitMQPublisher) PublishUserCreated(userEvent domain.UserCreatedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Veriyi JSON formatına dönüştür
	body, err := json.Marshal(userEvent)
	if err != nil {
		return err
	}

	// Mesajı "user.created" routing key'i ile exchange'e fırlat
	err = p.channel.PublishWithContext(ctx,
		"user.events",  // exchange
		"user.created", // routing key (hedef etiket)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Mesaj diskte saklansın (güvenli)
		},
	)
	if err != nil {
		log.Fatalf("Exchange declare hatası: %v", err)
		return err
	}

	return nil

}

func (p *RabbitMQPublisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
