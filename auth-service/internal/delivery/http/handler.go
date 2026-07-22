package handler

import (
	response "auth-service/internal/delivery"
	"auth-service/internal/domain"
	//"auth-service/pkg/security"
	"encoding/json"
	"fmt"
	"net/http"
	// "auth-service/pb" // gRPC proto paketiniz
)

type authHandler struct {
	// pb.UnimplementedAuthServiceServer
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(usecase domain.AuthUsecase) *authHandler {
	return &authHandler{
		authUsecase: usecase,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var dto domain.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendJSON(w, http.StatusBadRequest, response.APIResponse{
			Status:  "failure",
			Message: "Invalid JSON",
		})
		return
	}

	login, err := h.authUsecase.Login(r.Context(), dto.Email, dto.Password)
	if err != nil {
		response.SendJSON(w, http.StatusNotFound, response.APIResponse{
			Status:  "failure",
			Message: fmt.Sprintf("%s", err.Error()),
		})
		return
	}
	//token, _, _ := security.GenerateJWT(login.ID, login.Email, "secret_key")
	//login.Token = token
	response.SendJSON(w, http.StatusOK, response.APIResponse{
		Status: "success",
		Data:   login,
	})
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var dto domain.User

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendJSON(w, http.StatusBadRequest, response.APIResponse{
			Status:  "failure",
			Message: "Invalid JSON",
		})
		return
	}

	_,err := h.authUsecase.Register(r.Context(), dto.Email, dto.Password, dto.Role)
	if err != nil {
		response.SendJSON(w, http.StatusNotFound, response.APIResponse{
			Status:  "failure",
			Message: fmt.Sprintf("Failed Register: %v", err),
		})
		return
	}

	response.SendJSON(w, http.StatusOK, response.APIResponse{
		Status: "success",
	})
}
