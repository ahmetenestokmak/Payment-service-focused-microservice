package interceptor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// Kısa süreli koruma (Double-click / Debounce kilidi)
	FingerprintTTL = 5 * time.Second
	FingerprintPrefix = "fingerprint:"

	// Uzun süreli koruma (Network timeout / Retry kilidi)
	IdempotencyTTL = 48 * time.Hour
	IdempotencyPrefix = "idempotency:"

	// Durumlar
	StatusProcessing = "PROCESSING"
)

type IdempotencyInterceptor struct {
	rdb *redis.Client
}

func NewIdempotencyInterceptor(rdb *redis.Client) *IdempotencyInterceptor {
	return &IdempotencyInterceptor{rdb: rdb}
}

// UnaryServerInterceptor gRPC isteklerini sarmalayan ana fonksiyondur.
func (i *IdempotencyInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		
		// Metadata'yı (gRPC Header'larını) oku
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "Metadata okunamadı")
		}

		var idempotencyKey string
		if keys := md.Get("x-idempotency-key"); len(keys) > 0 {
			idempotencyKey = keys[0]
		}

		// 1. KATMAN: IDEMPOTENCY KEY KONTROLÜ (Uzun Süreli)
		if idempotencyKey != "" {
			redisKey := IdempotencyPrefix + idempotencyKey
			val, err := i.rdb.Get(ctx, redisKey).Result()
			if err == nil {
				if val == StatusProcessing {
					return nil, status.Errorf(codes.AlreadyExists, "İşleminiz şu an yürütülüyor, lütfen bekleyin")
				}

				// Eğer daha önce başarıyla tamamlandıysa, kaydedilmiş response JSON'ını çözüp dön
				// handler parametresindeki dönüş tipine unmarshal etmemiz gerekir.
				// Bu interceptor genel olduğu için handler'ın döneceği boş bir instance yaratıp dolduruyoruz.
				res, err := handler(ctx, req) // Şablon yapısını almak için geçici çağrı yerine doğrudan unmarshal adımı:
				if err == nil {
					if unmarshalErr := json.Unmarshal([]byte(val), &res); unmarshalErr == nil {
						return res, nil
					}
				}
			}
		}

		// 2. KATMAN: İSTEK PARMAK İZİ (HASH) KONTROLÜ (Kısa Süreli)
		// İstek parametrelerini (protobuf struct) JSON'a çevirerek hash'liyoruz
		reqBytes, err := json.Marshal(req)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "İstek içeriği analiz edilemedi: %v", err)
		}

		// Kullanıcı bilgisini metadata içindeki "user-id" veya "authorization" bilgisinden çekebilirsiniz
		var userID string
		if uids := md.Get("user-id"); len(uids) > 0 {
			userID = uids[0]
		} else {
			userID = "anonymous" // Fallback fallback (Güvenlik için gerçek projede JWT'den çözülmelidir)
		}

		// SHA-256 ile benzersiz istek hash'i üret (UserID + gRPC Method + Request Payload)
		hash := sha256.New()
		hash.Write([]byte(userID + info.FullMethod))
		hash.Write(reqBytes)
		fingerprintHash := hex.EncodeToString(hash.Sum(nil))
		fingerprintKey := FingerprintPrefix + fingerprintHash

		// Redis SETNX ile atomik debounce kilidi (5 saniye geçerli)
		success, err := i.rdb.SetNX(ctx, fingerprintKey, StatusProcessing, FingerprintTTL).Result()
		if err != nil || !success {
			return nil, status.Errorf(codes.ResourceExhausted, "Çok fazla istek gönderildi, lütfen bekleyin (Debounce Lock)")
		}

		// Idempotency Key mevcutsa PROCESSING durumuna çekerek uzun vadeli kilidi başlat
		if idempotencyKey != "" {
			_ = i.rdb.Set(ctx, IdempotencyPrefix+idempotencyKey, StatusProcessing, IdempotencyTTL).Err()
		}

		// Esas gRPC Ödeme Handler logic'ini çalıştır
		resp, err := handler(ctx, req)

		if err != nil {
			// İşlem başarısız veya hata kodu döndüyse kilitleri temizle ki kullanıcı tekrar deneyebilsin
			if idempotencyKey != "" {
				_ = i.rdb.Del(ctx, IdempotencyPrefix+idempotencyKey).Err()
			}
			_ = i.rdb.Del(ctx, fingerprintKey).Err()
			return nil, err
		}

		// İşlem başarılı! Yanıtı JSON formatına serileştirip Redis'e kaydet
		respBytes, marshalErr := json.Marshal(resp)
		if marshalErr == nil {
			if idempotencyKey != "" {
				_ = i.rdb.Set(ctx, IdempotencyPrefix+idempotencyKey, string(respBytes), IdempotencyTTL).Err()
			}
		}

		return resp, nil
	}
}