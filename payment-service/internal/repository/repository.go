package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/domain"
)

type paymentRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	
	query := `
    INSERT INTO payments (
        user_id, reference_id, amount, currency, 
        payment_method, status, transaction_id, failure_reason
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)  RETURNING id
`
	var dto domain.Payment
	err := r.db.QueryRowContext(ctx, query,
		payment.UserID,
		payment.ReferenceID,
		payment.Amount,
		payment.Currency,
		payment.PaymentMethod,
		payment.Status,
		payment.TransactionID,
		payment.FailureReason,
	).Scan(&dto.ID)
	if err != nil {
		return nil ,fmt.Errorf("[ERROR] repository.go:Save() Postgres: %w", err)
	}
	return &dto,nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error){
	return nil, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, transactionID string, failureReason string) error {

	query := `
    UPDATE payments 
    SET 
        status = $1, 
        transaction_id = $2, 
        failure_reason = $3, 
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $4
`
	_, err := r.db.ExecContext(ctx, query,
		status,
		transactionID,
		failureReason,
		id,
	)
	if err != nil {
		return fmt.Errorf("[ERROR] repository.go:UpdateStatus() Postgres: %w", err)
	}
	return nil
}
