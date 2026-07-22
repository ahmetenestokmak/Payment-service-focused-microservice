package infrastructure

import (
	"database/sql"
	"fmt"
	"time"
	"payment-service/config"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx sürücüsünü standart sql kütüphanesine kaydeder
)

// NewPostgresDB, optimize edilmiş bir SQL bağlantı havuzu (Connection Pool) döner
func NewPostgresDB(config config.ConfigDB) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)
	// Bağlantıyı aç (pgx sürücüsü ile)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db bağlantı hatası: %w", err)
	}

	// Bağlantı Havuzu (Pool) ayarları:
	db.SetMaxOpenConns(25)                 // Aynı anda açık olabilecek maksimum bağlantı sayısı
	db.SetMaxIdleConns(25)                 // Boşta (idle) hazır bekleyecek maksimum bağlantı sayısı
	db.SetConnMaxLifetime(5 * time.Minute) // Bir bağlantının ömrü (Zombi bağlantıları önler)
	db.SetConnMaxIdleTime(3 * time.Minute) // Boştaki bağlantının maksimum bekleme süresi

	// Canlılık Kontrolü (Ping)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[ERROR] Database: %w", err)
	}

	return db, err
}
