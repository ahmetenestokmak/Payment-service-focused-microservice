# 1. Aşama: Derleme (Build Stage)
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/server/main.go

# 2. Aşama: Çalıştırma (Run Stage)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 50053
CMD ["./user-service"]