FROM golang:1.26-alpine AS builder

WORKDIR /app

# Bağımlılıkları kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Tüm proje kodunu (yerelde ürettiğiniz .pb.go dosyaları dahil) kopyala
COPY . .

# Uygulamayı derle (protoc adımlarını sildik)
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/api/main.go


FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]