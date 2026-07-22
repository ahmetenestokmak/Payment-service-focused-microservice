-- Eğer enum tipi önceden tanımlanmadıysa oluşturuyoruz
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
        CREATE TYPE payment_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED', 'REFUNDED');
    END IF;
END $$;

-- Tablo oluşturma
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    reference_id VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL, -- Kuruş / Cents cinsinden
    currency VARCHAR(3) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- STRIPE, IYZICO vb.
    status payment_status NOT NULL DEFAULT 'PENDING',
    transaction_id VARCHAR(255), -- Bankadan dönen makbuz no
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- GORM modelinde belirtilen indekslerin oluşturulması
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_reference_id ON payments(reference_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);