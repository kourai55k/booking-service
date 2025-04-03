package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type GetUsersResponse struct {
	Users []*models.User
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	const op = "http.userHandler.GetUsers"
	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	users, err := h.userService.GetUsers()
	if err != nil {
		if errors.Is(err, domain.ErrUsersNotFound) {
			http.Error(w, "users not found", http.StatusNotFound)
			log.Error("users not found", "err", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		log.Error("failed to get users", "err", err.Error())
		return
	}

	var res GetUsersResponse
	res.Users = users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to encode response", "err", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
