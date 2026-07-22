


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY, -- Auth'tan gelen ID doğrudan PK olur (Otomatik artan / serial DEĞİLDİR)
    first_name VARCHAR(100),
    last_name VARCHAR(100)
);