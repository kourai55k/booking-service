package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type getUserByLoginResponse struct {
	User *models.User `json:"user"`
}

func (h *UserHandler) GetUserByLogin(w http.ResponseWriter, r *http.Request) {
	const op = "http.userHandler.GetUserByID"
	log := h.logger

	if r.URL.Path != "/favicon.ico" {
		log.Debug("request received", "method", r.Method, "path", r.URL.Path)
	}

	login := r.URL.Query().Get("login")
	if login == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("bad request", "err", fmt.Errorf("%s: bad request", op).Error())
		return
	}

	user, err := h.userService.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Error("user not found", "err", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "failed to get user by id", http.StatusInternalServerError)
		log.Error("failed to get user by id", "err", err.Error())
		return
	}

	var res getUserByLoginResponse
	res.User = user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to encode response", "err", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
