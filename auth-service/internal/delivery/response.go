package delivery

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`           // "success" veya "error"
	Message string      `json:"message"`          // Genel açıklama mesajı
	Data    interface{} `json:"data,omitempty"`   // Başarılıysa dönecek veri
	Errors  interface{} `json:"errors,omitempty"` // Validation hataları detayları
}

func SendJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
