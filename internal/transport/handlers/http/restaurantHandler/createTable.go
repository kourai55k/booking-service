package restauranthandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type createTableRequest struct {
	Number   uint `json:"number"`
	Capacity uint `json:"capacity"`

	RestaurantID uint `json:"restaurantID"`
}

type createTableResponse struct {
	ID uint `json:"id"`
}

func (r *createTableRequest) validate() error {
	if r.Number == 0 || r.Capacity == 0 || r.RestaurantID == 0 {
		return errors.New("missing required fields")
	}
	return nil
}

func (h *RestraurantHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	const op = "http.RestaurantHanlder.CreateTabler"

	log := h.logger

	log.Debug("request received", "method", r.Method, "path", r.URL.Path)

	var req createTableRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error("failed to decode request body", "error", fmt.Errorf("%s: bad request", op).Error())
		return
	}

	if err := req.validate(); err != nil {
		http.Error(w, fmt.Sprintf("bad request: %v", err), http.StatusBadRequest)
		log.Error("bad request", "error", fmt.Errorf("%s: %w", op, err).Error())
		return
	}

	// Retrieve userID from context (assuming it has been added by a middleware)
	userID, ok := r.Context().Value(domain.UserIDKey).(uint)
	if !ok {
		http.Error(w, "user ID not found in context", http.StatusUnauthorized)
		log.Error("user ID not found in context", "error", fmt.Errorf("%s: user ID missing", op).Error())
		return
	}

	// Check if the user is the owner of the restaurant
	isOwner, err := h.restaurantService.IsOwnerOfRestaurant(userID, req.RestaurantID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to verify if user is owner", "error", fmt.Errorf("%s: %w", op, err).Error())
		return
	}
	if !isOwner {
		http.Error(w, "forbidden: user is not the owner of the restaurant", http.StatusForbidden)
		log.Error("user is not owner of the restaurant", "error", fmt.Errorf("%s: user is not owner", op).Error())
		return
	}

	table := &models.Table{
		Number:       req.Number,
		Capacity:     req.Capacity,
		RestaurantID: req.RestaurantID,
	}

	id, err := h.restaurantService.CreateTable(table)
	if err != nil {
		if errors.Is(err, domain.ErrTableAlreadyExists) {
			http.Error(w, "table already exists", http.StatusConflict)
			log.Error("table already exists", "error", fmt.Errorf("%s: %w", op, err).Error())
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to create table", "error", fmt.Errorf("%s:%w", op, err).Error())
		return
	}

	var res createTableResponse
	res.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("failed to encode response", "error", fmt.Errorf("%s: failed to encode response", op).Error())
	}
}
