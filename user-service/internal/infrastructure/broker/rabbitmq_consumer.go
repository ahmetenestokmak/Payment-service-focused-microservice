package broker

import (
	"context"
	"encoding/json"
	"log"

	"user-service/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	usecase domain.UserUsecase
}

func NewRabbitMQConsumer(url string, uc domain.UserUsecase) *RabbitMQConsumer {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Consumer için RabbitMQ bağlantısı kurulamadı: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Consumer için RabbitMQ kanalı açılamadı: %v", err)
	}
	
	
	// 1. Exchange Tanımla (Eğer Auth servisi henüz ayağa kalkmadıysa crash olmamak için)
	err = ch.ExchangeDeclare(
		"user.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Exchange declare hatası: %v", err)
	}

	// 2. Sadece bu servise özel bir Kuyruk (Queue) oluştur
	q, err := ch.QueueDeclare(
		"user_service_queue", // kuyruk adı
		true,                 // durable (kalıcı)
		false,                // auto-delete
		false,                // exclusive
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Queue declare hatası: %v", err)
	}

	// 3. Kuyruğu, "user.created" etiketiyle Exchange'e bağla (Bind)
	err = ch.QueueBind(
		q.Name,
		"user.created", // routing key
		"user.events",  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Queue bind hatası: %v", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		usecase: uc,
	}
}

// StartConsume Arka planda mesajları dinlemeye başlar (Sonsuz döngü)
func (c *RabbitMQConsumer) StartConsume() {
	msgs, err := c.channel.Consume(
		"user_service_queue", // dinlenecek kuyruk
		"",                   // consumer tag
		false,                // auto-ack (Manuel onaylama yapacağız - daha güvenli)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Consume başlatma hatası: %v", err)
	}

	// Mesajları asenkron olarak dinlemek için goroutine başlatıyoruz
	go func() {
		for d := range msgs {
			var event domain.UserCreatedEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("[ERROR] Mesaj JSON'a dönüştürülemedi, dropping: %v", err)
				d.Ack(false) // Hatalı formatı kuyruktan at
				continue
			}
			/*	ID:        data.ID,
				FirstName: req.GetFirstName(),
				LastName:  req.GetLastName(),


				(
					context.Background(),
					event.ID,
					event.FirstName,
					event.LastName,
				)
			*/
			// 4. ADIM: Veri tabanına profil kaydetme mantığı burada çalışacak
			log.Printf("[CONSUMER] Yeni kullanıcı event'i yakalandı! ID: %s, İsim: %s %s",
				event.ID, event.FirstName, event.LastName)

			err = c.usecase.Save(context.Background(), &domain.User{
				ID:        event.ID,
				FirstName: event.FirstName,
				LastName:  event.LastName,
			})
			if err != nil {
				log.Printf("[ERROR] Kullanıcı kaydedilemedi: %v", err)
				continue
			}

			d.Ack(false) // Başarılı işleme sonrası mesajı onayla
		}
	}()

	log.Println("[INFO] RabbitMQ Consumer başarıyla başlatıldı, 'user.created' eventleri bekleniyor...")
}

func (c *RabbitMQConsumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
