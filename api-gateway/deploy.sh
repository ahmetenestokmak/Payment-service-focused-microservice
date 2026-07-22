# chmod +x deploy.sh

SERVER_IP="[IP_ADDRESS]"
USER="developer"
TARGET_DIR="/home/developer//microservices/api-gateway"

echo "1. Sunucuda hedef klasör kontrol ediliyor..."
ssh $USER@$SERVER_IP "mkdir -p $TARGET_DIR"

echo "2. Proje dosyaları sunucuya yükleniyor..."
rsync -avz --exclude='.git' --exclude='node_modules' ./ $USER@$SERVER_IP:$TARGET_DIR/

echo "3. Sunucuda Docker üzerinde konteyner derleniyor ve arka planda ayağa kaldırılıyor..."
# Sunucudaki klasöre girip, docker-compose'u tetikliyoruz. --build ile yeni kodu zorunlu derletiyoruz.
ssh $USER@$SERVER_IP "cd $TARGET_DIR && docker compose up -d --build"

echo "4. Konteyner durumu kontrol ediliyor..."
ssh $USER@$SERVER_IP "docker ps --filter name=gateway_microservice"

echo "İşlem Başarılı! api-gateway sorunsuz şekilde Docker üzerinde canlıya alındı."