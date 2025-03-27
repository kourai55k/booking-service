package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
	"github.com/kourai55k/booking-service/pkg/hashing"
)

type updateUserRequest struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Validate checks if ID and at least one other field is provided to update
func (r *updateUserRequest) Validate() error {
	if r.ID == 0 {
		return errors.New("id is required")
	}

	if r.Name == "" && r.Login == "" && r.Password == "" && r.Role == "" {
		return errors.New("at least one field is required")
	}

	return nil
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "http.UserHandler.UpdateUser"

	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	var req updateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent unknown fields
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("failed to decode request body", "error", fmt.Errorf("%s: bad request", op).Error())
		return
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

	req.ID = uint(id)

	if err := req.Validate(); err != nil {
		http.Error(w, "bad request: missing required fields", http.StatusBadRequest)
		log.Error("bad request", "error", fmt.Errorf("%s: %w", op, err).Error())
		return
	}

	hashPass, err := hashing.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to hash password", "error", fmt.Errorf("%s: failed to hash password", op).Error())
		return
	}

	user := &models.User{
		ID:       req.ID,
		Name:     req.Name,
		Login:    req.Login,
		HashPass: hashPass,
		Role:     req.Role,
	}

	if err := h.userService.UpdateUser(user); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Error("user not found", "error", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, "user with this login already exists", http.StatusConflict)
			log.Error("user already exists", "error", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to update user", "error", fmt.Errorf("%s: %w", op, err).Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
