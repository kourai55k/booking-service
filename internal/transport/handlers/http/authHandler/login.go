package authHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kourai55k/booking-service/internal/domain"
)

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string
}

func (r loginRequest) Validate() error {
	if r.Login == "" || r.Password == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "http.AuthHandler.Login"

	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	var req loginRequest
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
		log.Error("failed to validate request", "error", fmt.Errorf("%s: %w", op, err).Error())
		return
	}

	token, err := h.authService.Login(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "there is no user with this login", http.StatusUnauthorized)
			log.Error("user not found", "err", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		if errors.Is(err, domain.ErrWrongPassword) {
			http.Error(w, "wrong password", http.StatusUnauthorized)
			log.Error("wrong password", "err", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "failed to login user: internal server error", http.StatusInternalServerError)
		log.Error("failed to login user", "err", err.Error())
		return
	}

	var res loginResponse
	res.Token = token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to encode response", "err", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
