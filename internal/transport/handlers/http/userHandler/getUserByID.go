package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type GetUserByIDResponse struct {
	User *models.User `json:"user"`
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	const op = "userHandler.GetUserByID"
	log := h.logger

	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path != "/favicon.ico" {
		log.Debug("GET /user was called")
	}

	// Extract the 'id' path parameter using Go 1.22's PathValue
	idStr := r.PathValue("id")

	// Convert 'id' to uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || idStr == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("bad request", "err", fmt.Errorf("%s: bad request", op).Error())
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Error("user not found", "err", err.Error())
			return
		}
		http.Error(w, "failed to get user by id", http.StatusInternalServerError)
		log.Error("failed to get user by id", "err", err.Error())
		return
	}

	var res GetUserByIDResponse
	res.User = user
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error("failed to encode response", "err", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
