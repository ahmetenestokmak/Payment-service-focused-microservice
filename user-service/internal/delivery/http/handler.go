package http

import (
	"encoding/json"
	"net/http"
	"user-service/internal/domain"
)



type userHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *userHandler {
	return &userHandler{usecase: usecase}
}

func (h userHandler) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var dto domain.User

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		SendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "failure",
			Message: "Invalid JSON",
		})
		return
	}

	err := h.usecase.Save(r.Context(), &domain.User{
		ID: dto.ID,
		FirstName:	dto.FirstName,
		LastName: 	dto.LastName,
	})
	if err != nil {
		SendJSON(w, http.StatusBadRequest, APIResponse{
			Status:  "failure",
			Message: err.Error(),
		})
		return
	}

	SendJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "",
	})
}

func (h userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	
	userEntity, err := h.usecase.GetProfile(r.Context(), id)
	if err != nil {
		SendJSON(w, http.StatusNotFound, APIResponse{
			Status:  "failure",
			Message: "ID is not found",
		})
		return
	}

	SendJSON(w, http.StatusOK, APIResponse{
		Status: "success",
		Data:   userEntity,
	})
}
