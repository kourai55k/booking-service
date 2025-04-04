package authHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
	"github.com/kourai55k/booking-service/pkg/hashing"
)

type registerRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r registerRequest) Validate() error {
	if r.Name == "" || r.Login == "" || r.Password == "" {
		return errors.New("missing required fields")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

type registerResponse struct {
	ID uint `json:"id"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "http.AuthHandler.Register"

	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	var req registerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent unknown fields
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("failed to decode request body", "error", fmt.Errorf("%s: bad request", op).Error())
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("bad request: %v", err), http.StatusBadRequest)
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
		Name:     req.Name,
		Login:    req.Login,
		HashPass: hashPass,
		Role:     "user",
	}

	id, err := h.authService.Register(user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, "user already exists", http.StatusConflict)
			log.Error("user already exists", "error", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to create user", "error", fmt.Errorf("%s:%w", op, err).Error())
		return
	}

	var res registerResponse
	res.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to encode response", "error", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
